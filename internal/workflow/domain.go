package workflow

import (
	"path/filepath"
	"sync"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/alterx"
	"github.com/liancccc/goauto/internal/modules/cdncheck"
	"github.com/liancccc/goauto/internal/modules/gospider"
	httpx_info "github.com/liancccc/goauto/internal/modules/httpx/info"
	httpx_unique "github.com/liancccc/goauto/internal/modules/httpx/unique"
	"github.com/liancccc/goauto/internal/modules/katana"
	"github.com/liancccc/goauto/internal/modules/ksubdomain/enum"
	"github.com/liancccc/goauto/internal/modules/ksubdomain/verify"
	"github.com/liancccc/goauto/internal/modules/merge"
	"github.com/liancccc/goauto/internal/modules/naabu"
	"github.com/liancccc/goauto/internal/modules/nuclei"
	"github.com/liancccc/goauto/internal/modules/oneforall"
	"github.com/liancccc/goauto/internal/modules/subfinder"
	"github.com/liancccc/goauto/internal/modules/unique"
	"github.com/liancccc/goauto/internal/modules/urlfinder"
	"github.com/liancccc/goauto/internal/modules/xray"
	xscan_spider "github.com/liancccc/goauto/internal/modules/xscan/spider"
	"github.com/projectdiscovery/gologger"
)

type DomainALLFlow struct {
}

func init() {
	RegisterWorkflow(&DomainALLFlow{})
}

func (f *DomainALLFlow) Name() string {
	return "DomainALL"
}

func (f *DomainALLFlow) Description() string {
	return "子域名收集 -> CDN识别 -> 端口扫描[TOP-1000] -> 验活去重 -> WEB获取信息和截图 -> 爬虫 -> 漏洞扫描"
}

func (f *DomainALLFlow) Run(runner *Runner) {
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
		Timeout: "5m",
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
		new(naabu.ModuleStruct).Run(modules.BaseParams{
			Target: filepath.Join(cdncheckOutDir, "noCdn.txt"),
			Output: filepath.Join(portscanOutDir, "noCdn-services.txt"),
		})
	}

	// URL 去重验活 + 获取截图等信息
	var httpxOutDir = filepath.Join(runner.workSpace, "httpx")
	if fileutil.CountLines(filepath.Join(cdncheckOutDir, "cdn.txt")) > 0 {
		new(httpx_unique.ModuleStruct).Run(modules.BaseParams{
			Target:          filepath.Join(cdncheckOutDir, "cdn.txt"),
			Output:          filepath.Join(httpxOutDir, "cdn-alive.txt"),
			CustomizeParams: "-mc 200,302 -p 80,443,8080,8000,8888,4848,7070,8089,8181,9080,9443,5000,8443,5001,81,8081,50805,3000,88,7547",
			Proxy:           runner.opt.Proxy,
		})
	}
	if fileutil.CountLines(filepath.Join(portscanOutDir, "noCdn-services.txt")) > 0 {
		new(httpx_unique.ModuleStruct).Run(modules.BaseParams{
			Target: filepath.Join(portscanOutDir, "noCdn-services.txt"),
			Output: filepath.Join(httpxOutDir, "noCdn-alive.txt"),
			Proxy:  runner.opt.Proxy,
		})
	}

	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(httpxOutDir, "merge.txt"),
		},
		Targets: []string{filepath.Join(httpxOutDir, "noCdn-alive.txt"), filepath.Join(httpxOutDir, "cdn-alive.txt")},
	})

	new(unique.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "merge.txt"),
		Output: filepath.Join(httpxOutDir, "all.txt"),
	})

	new(httpx_info.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(httpxOutDir, "web.txt"),
	})

	// 爬虫
	var spiderOutDir = filepath.Join(runner.workSpace, "spider")
	new(gospider.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(spiderOutDir, "gospider.txt"),
		Proxy:  runner.opt.Proxy,
	})
	new(katana.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(spiderOutDir, "katana.txt"),
		Proxy:  runner.opt.Proxy,
	})
	new(urlfinder.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(spiderOutDir, "urlfinder.txt"),
		Proxy:  runner.opt.Proxy,
	})
	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: filepath.Join(spiderOutDir, "all.txt"),
		},
		Targets: []string{filepath.Join(spiderOutDir, "gospider.txt"), filepath.Join(spiderOutDir, "katana.txt"), filepath.Join(spiderOutDir, "urlfinder.txt")},
	})

	new(unique.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(spiderOutDir, "all.txt"),
		Output: filepath.Join(spiderOutDir, "links.txt"),
	})

	// 漏洞扫描
	var vulscanOutDir = filepath.Join(runner.workSpace, "vulscan")
	new(nuclei.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "all.txt"),
		Output: filepath.Join(vulscanOutDir, "nuclei.txt"),
	})
	new(xscan_spider.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(spiderOutDir, "links.txt"),
		Output: filepath.Join(vulscanOutDir, "xscan.json"),
	})
	new(xray.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(spiderOutDir, "links.txt"),
		Output: filepath.Join(vulscanOutDir, "xray.txt"),
	})
}
