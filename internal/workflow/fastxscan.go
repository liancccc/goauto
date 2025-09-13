package workflow

import (
	"fmt"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/alterx"
	"github.com/liancccc/goauto/internal/modules/cdncheck"
	"github.com/liancccc/goauto/internal/modules/dnsx"
	httpx_info "github.com/liancccc/goauto/internal/modules/httpx/info"
	httpx_unique "github.com/liancccc/goauto/internal/modules/httpx/unique"
	ksubdomain_enum "github.com/liancccc/goauto/internal/modules/ksubdomain/enum"
	ksubdomain_verify "github.com/liancccc/goauto/internal/modules/ksubdomain/verify"
	"github.com/liancccc/goauto/internal/modules/merge"
	"github.com/liancccc/goauto/internal/modules/naabu"
	"github.com/liancccc/goauto/internal/modules/notify"
	"github.com/liancccc/goauto/internal/modules/oneforall"
	"github.com/liancccc/goauto/internal/modules/subfinder"
	"github.com/liancccc/goauto/internal/modules/uncover"
	"github.com/liancccc/goauto/internal/modules/unique"
	"github.com/liancccc/goauto/internal/modules/xscan"
	xscan_spider "github.com/liancccc/goauto/internal/modules/xscan/spider"
	"github.com/projectdiscovery/gologger"
)

type FastXscanFlow struct {
}

func init() {
	RegisterWorkflow(&FastXscanFlow{})
}

func (f *FastXscanFlow) Name() string {
	return "DomainFastXscan"
}

func (f *FastXscanFlow) Description() string {
	return "subdomain -> cdncheck -> portscan-top-1000 -> quake -> httpx -> xscan"
}

func (f *FastXscanFlow) Run(runner *Runner) {
	var subdomainOutDir = filepath.Join(runner.workSpace, "subdomain")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		new(subfinder.ModuleStruct).Run(modules.BaseParams{
			Target: runner.opt.Target,
			Output: filepath.Join(subdomainOutDir, "subfinder.txt"),
		})
	}()
	go func() {
		defer wg.Done()
		new(oneforall.ModuleStruct).Run(modules.BaseParams{
			Target: runner.opt.Target,
			Output: filepath.Join(subdomainOutDir, "oneforall.txt"),
		})
	}()
	wg.Wait()
	new(ksubdomain_enum.ModuleStruct).Run(modules.BaseParams{
		Target:  runner.opt.Target,
		Output:  filepath.Join(subdomainOutDir, "ksubdomain.txt"),
		Timeout: "5h",
	})
	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(subdomainOutDir, "merge.txt"),
		},
		Targets: []string{filepath.Join(subdomainOutDir, "subfinder.txt"), filepath.Join(subdomainOutDir, "oneforall.txt"), filepath.Join(subdomainOutDir, "ksubdomain.txt")},
	})

	new(alterx.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(subdomainOutDir, "merge.txt"),
		Output: filepath.Join(subdomainOutDir, "alterx.txt"),
	})

	new(ksubdomain_verify.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(subdomainOutDir, "alterx.txt"),
		Output: filepath.Join(subdomainOutDir, "alterx-alive.txt"),
	})
	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(subdomainOutDir, "merge-alterx.txt"),
		},
		Targets: []string{filepath.Join(subdomainOutDir, "merge.txt"), filepath.Join(subdomainOutDir, "alterx-alive.txt")},
	})

	new(unique.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(subdomainOutDir, "merge-alterx.txt"),
		Output: filepath.Join(subdomainOutDir, "all.txt"),
	})

	if fileutil.CountLines(filepath.Join(subdomainOutDir, "all.txt")) == 0 {
		gologger.Error().Msgf("%s Get 0 Subdomains Exit", runner.opt.Target)
		return
	}

	fileutil.Cleaning(subdomainOutDir, []string{
		filepath.Join(subdomainOutDir, "all.txt"),
	})

	// CDN 识别
	var cdncheckOutDir = filepath.Join(runner.workSpace, "cdncheck")
	new(cdncheck.ModuleStruct).Run(cdncheck.Params{
		BaseParams: &modules.BaseParams{
			Target: filepath.Join(subdomainOutDir, "all.txt"),
		},
		CDNPath:   filepath.Join(cdncheckOutDir, "cdn.txt"),
		NoCDNPath: filepath.Join(cdncheckOutDir, "noCdn.txt"),
	})

	// 对非 CDN 目标进行端口扫描
	var portscanOutDir = filepath.Join(runner.workSpace, "portscan")
	if fileutil.CountLines(filepath.Join(cdncheckOutDir, "noCdn.txt")) > 0 {
		// 解析域名获取 IP 地址
		new(dnsx.ModuleStruct).Run(modules.BaseParams{
			Target: filepath.Join(cdncheckOutDir, "noCdn.txt"),
			Output: filepath.Join(portscanOutDir, "dnsx.json"),
		})
		// 解析 IP 地址, 获取 IP 和域名列表映射
		ipDomainsMap, _ := dnsx.CleanAndGenCustomizeFormat(filepath.Join(portscanOutDir, "dnsx.json"), filepath.Join(portscanOutDir, "ips.txt"))
		// 端口扫描, 获取如 ssh://ip:port, http://ip:port 的链接
		new(naabu.ModuleStruct).Run(modules.BaseParams{
			Target: filepath.Join(portscanOutDir, "ips.txt"),
			Output: filepath.Join(portscanOutDir, "noCdn-services.txt"),
		})
		// 给 http 服务把 IP 对应的域名都添加上
		if fileutil.CountLines(filepath.Join(portscanOutDir, "noCdn-services.txt")) > 0 {
			var services []string
			serviceUrls := fileutil.ReadingLines(filepath.Join(portscanOutDir, "noCdn-services.txt"))
			for _, serviceUrl := range serviceUrls {
				parseUrl, err := url.Parse(serviceUrl)
				if err != nil {
					services = append(services, serviceUrl)
					continue
				}
				if _, exists := ipDomainsMap[parseUrl.Hostname()]; !exists {
					continue
				}
				if parseUrl.Scheme == "http" || parseUrl.Scheme == "https" {
					for _, domain := range ipDomainsMap[parseUrl.Hostname()] {
						services = append(services, fmt.Sprintf("%s://%s:%s", parseUrl.Scheme, domain, parseUrl.Port()))
					}
				} else {
					services = append(services, serviceUrl)
				}
			}
			fileutil.WriteSliceToFile(filepath.Join(portscanOutDir, "services.txt"), services)
		}
	}
	// 从测绘平台拉取服务信息 quake, 目前每其他的
	var uncoverOutDir = filepath.Join(runner.workSpace, "uncover")
	var uncoverTargets []uncover.Target
	if fileutil.IsFile(runner.opt.Target) {
		var domains = fileutil.ReadingLines(runner.opt.Target)
		for _, domain := range domains {
			uncoverTargets = append(uncoverTargets, uncover.Target{
				Query:  fmt.Sprintf(`domain:"%s"`, domain),
				Engine: "quake",
			})
		}
	} else {
		uncoverTargets = append(uncoverTargets, uncover.Target{
			Query:  fmt.Sprintf(`domain:"%s"`, runner.opt.Target),
			Engine: "quake",
		})
	}
	new(uncover.ModuleStruct).Run(uncover.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(uncoverOutDir, "services.txt"),
		},
		Targets: uncoverTargets,
	})

	// URL 去重验活 + 获取截图等信息
	var httpxOutDir = filepath.Join(runner.workSpace, "httpx")
	if fileutil.CountLines(filepath.Join(cdncheckOutDir, "cdn.txt")) > 0 {
		new(httpx_unique.ModuleStruct).Run(modules.BaseParams{
			Target:          filepath.Join(cdncheckOutDir, "cdn.txt"),
			Output:          filepath.Join(httpxOutDir, "cdn-alive.txt"),
			CustomizeParams: "-p 80,443,8080,8000,8888,4848,7070,8089,8181,9080,9443,5000,8443,5001,81,8081,50805,3000,88,7547",
			Proxy:           runner.opt.Proxy,
		})
	}
	if fileutil.CountLines(filepath.Join(portscanOutDir, "services.txt")) > 0 {
		new(httpx_unique.ModuleStruct).Run(modules.BaseParams{
			Target: filepath.Join(portscanOutDir, "services.txt"),
			Output: filepath.Join(httpxOutDir, "noCdn-alive.txt"),
			Proxy:  runner.opt.Proxy,
		})
	}
	if fileutil.CountLines(filepath.Join(uncoverOutDir, "services.txt")) > 0 {
		new(httpx_unique.ModuleStruct).Run(modules.BaseParams{
			Target: filepath.Join(uncoverOutDir, "services.txt"),
			Output: filepath.Join(httpxOutDir, "uncover-alive.txt"),
			Proxy:  runner.opt.Proxy,
		})
	}

	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(httpxOutDir, "merge.txt"),
		},
		Targets: []string{filepath.Join(httpxOutDir, "noCdn-alive.txt"), filepath.Join(httpxOutDir, "cdn-alive.txt"), filepath.Join(httpxOutDir, "uncover-alive.txt")},
	})

	new(unique.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "merge.txt"),
		Output: filepath.Join(httpxOutDir, "all.txt"),
	})

	new(httpx_info.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(httpxOutDir, "web.txt"),
	})

	// xscan 扫描
	var vulscanOutDir = filepath.Join(runner.workSpace, "vulscan")
	new(xscan_spider.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(vulscanOutDir, "xscan-spider.json"),
	})
	if fileutil.IsFile(filepath.Join(vulscanOutDir, "xscan-spider.json")) {
		xscan.Clean(filepath.Join(vulscanOutDir, "xscan-spider.json"), filepath.Join(vulscanOutDir, "xscan-spider.html"))
		new(notify.ModuleStruct).Run(notify.Params{
			Msg: fmt.Sprintf("Task Name: %s, xscan-spider.json Count: %d", runner.opt.TaskName, fileutil.CountLines(filepath.Join(vulscanOutDir, "xscan-spider.json"))),
		})
	}
}
