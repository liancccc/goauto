package workflow

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/httpx"
	"github.com/liancccc/goauto/internal/modules/naabu"
	"github.com/liancccc/goauto/internal/modules/notify"
	"github.com/liancccc/goauto/internal/modules/nuclei"
	"github.com/liancccc/goauto/internal/modules/uncover"
	"github.com/liancccc/goauto/internal/modules/wih"
	"github.com/liancccc/goauto/internal/modules/xray"
	"github.com/liancccc/goauto/internal/modules/xscan"
	xscan_spider "github.com/liancccc/goauto/internal/modules/xscan/spider"
)

func init() {
	RegisterWorkflow(&subdomainAllFlow{})
}

type subdomainAllFlow struct {
}

func (d *subdomainAllFlow) Name() string {
	return "subdomainAll"
}

func (d *subdomainAllFlow) Description() string {
	return "subdomain -> dnsx -> portscan -> httpx -> vulscan alive url -> spider -> vulscan spider link"
}

func (d *subdomainAllFlow) Run(params *workflowParams) {
	// CDN 识别
	cdnFlow := new(cdncheckFlow)
	cdnFlow.Run(&workflowParams{
		target:    params.target,
		workSpace: params.workSpace,
		opt:       params.opt,
	})

	// 非 CDN 域名解析
	dnsxFlow := new(dnsxFlow)
	dnsxFlow.Run(&workflowParams{
		target:    cdnFlow.noCdnOut,
		workSpace: params.workSpace,
		opt:       params.opt,
	})

	// 非 CDN IP 端口扫描
	var portscanOutDir = filepath.Join(params.workSpace, "portscan")
	var portscanOut = filepath.Join(portscanOutDir, "portscan.txt")
	var portscanWebOut = filepath.Join(portscanOutDir, "http-urls.txt")
	new(naabu.ModuleStruct).Run(modules.BaseParams{
		Target: dnsxFlow.ipsFile,
		Output: portscanOut,
	})
	// 给 web 服务添加对应 ip 的域名
	if fileutil.CountLines(portscanOut) > 0 {
		var webUrls []string
		serviceUrls := fileutil.ReadingLines(portscanOut)
		for _, serviceUrl := range serviceUrls {
			parseUrl, err := url.Parse(serviceUrl)
			if err != nil {
				continue
			}
			if parseUrl.Scheme != "https" && parseUrl.Scheme != "http" {
				continue
			}
			webUrls = append(webUrls, serviceUrl)
			for _, domain := range dnsxFlow.ipDomainsMap[parseUrl.Hostname()] {
				webUrls = append(webUrls, fmt.Sprintf("%s://%s:%s", parseUrl.Scheme, domain, parseUrl.Port()))
			}
		}
		fileutil.WriteSliceToFile(portscanWebOut, webUrls)
	}

	// 从 Quake 拉取资产
	var uncoverOutDir = filepath.Join(params.workSpace, "uncover")
	var uncoverTargets []uncover.Target
	var uncoverSvrOut = filepath.Join(uncoverOutDir, "services.txt")
	if fileutil.IsFile(params.opt.Target) {
		var domains = fileutil.ReadingLines(params.opt.Target)
		for _, domain := range domains {
			uncoverTargets = append(uncoverTargets, uncover.Target{
				Query:  fmt.Sprintf(`domain:"%s"`, domain),
				Engine: "quake",
			})
		}
	} else {
		uncoverTargets = append(uncoverTargets, uncover.Target{
			Query:  fmt.Sprintf(`domain:"%s"`, params.opt.Target),
			Engine: "quake",
		})
	}
	new(uncover.ModuleStruct).Run(uncover.Params{
		BaseParams: &modules.BaseParams{
			Output: uncoverSvrOut,
		},
		Targets: uncoverTargets,
	})

	// httpx web 验活
	var httpxOutDir = filepath.Join(params.workSpace, "httpx")
	// CDN web 就拼接常见的 WEB 端口
	new(httpx.ModuleStruct).Run(modules.BaseParams{
		Target:          cdnFlow.cdnOut,
		Output:          filepath.Join(httpxOutDir, "cdn.json"),
		CustomizeParams: "-hash simhash -json -p 80,443,8080,8000,8888,4848,7070,8089,8181,9080,9443,5000,8443,5001,81,8081,50805,3000,88,7547",
		Proxy:           params.opt.Proxy,
	})
	new(httpx.ModuleStruct).Run(modules.BaseParams{
		Target:          portscanWebOut,
		Output:          filepath.Join(httpxOutDir, "noCdn.json"),
		CustomizeParams: "-hash simhash -json",
		Proxy:           params.opt.Proxy,
	})
	new(httpx.ModuleStruct).Run(modules.BaseParams{
		Target:          uncoverSvrOut,
		Output:          filepath.Join(httpxOutDir, "uncover.json"),
		CustomizeParams: "-hash simhash -json",
		Proxy:           params.opt.Proxy,
	})
	MergeAndUnique(
		[]string{
			filepath.Join(httpxOutDir, "cdn.json"),
			filepath.Join(httpxOutDir, "noCdn.json"),
			filepath.Join(httpxOutDir, "uncover.json"),
		},
		filepath.Join(httpxOutDir, "all.json"),
	)
	httpxUniqueResults := httpx.ParseAndUnique(filepath.Join(httpxOutDir, "all.json"))
	// 收集 URL , 存活的和其他状态码的, 其他的也有可能出东西的
	for _, result := range httpxUniqueResults {
		fileutil.AppendToContent(filepath.Join(httpxOutDir, "all.txt"), result.URL)
		if 200 <= result.StatusCode && result.StatusCode < 400 {
			fileutil.AppendToContent(filepath.Join(httpxOutDir, "alive.txt"), result.URL)
		} else {
			fileutil.AppendToContent(filepath.Join(httpxOutDir, fmt.Sprintf("%d.txt", result.StatusCode)), result.URL)
		}
	}

	// 漏洞扫描
	var vulscanOutDir = filepath.Join(params.workSpace, "vulscan")
	new(xscan_spider.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "alive.txt"),
		Output: filepath.Join(vulscanOutDir, "xscan-spider.json"),
	})
	if fileutil.IsFile(filepath.Join(vulscanOutDir, "xscan-spider.json")) {
		xscan.Clean(filepath.Join(vulscanOutDir, "xscan-spider.json"), filepath.Join(vulscanOutDir, "xscan-spider.html"))
		new(notify.ModuleStruct).Run(notify.Params{
			Msg: fmt.Sprintf("Task Name: %s, xscan-spider.json Count: %d", params.opt.TaskName, fileutil.CountLines(filepath.Join(vulscanOutDir, "xscan-spider.json"))),
		})
	}
	new(nuclei.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(httpxOutDir, "alive.txt"),
		Output: filepath.Join(vulscanOutDir, "nuclei.txt"),
		Proxy:  params.opt.Proxy,
	})
	if fileutil.IsFile(filepath.Join(vulscanOutDir, "nuclei.txt")) {
		new(notify.ModuleStruct).Run(notify.Params{
			Msg: fmt.Sprintf("Task Name: %s, nuclei.txt Count: %d", params.opt.TaskName, fileutil.CountLines(filepath.Join(vulscanOutDir, "nuclei.txt"))),
		})
	}

	// 爬虫
	spidFlow := new(spiderFlow)
	spidFlow.Run(&workflowParams{
		target:    filepath.Join(httpxOutDir, "alive.txt"),
		workSpace: params.workSpace,
		opt:       params.opt,
	})
	// wih
	var wihOutDir = filepath.Join(params.workSpace, "wih")
	new(wih.ModuleStruct).Run(modules.BaseParams{
		Target: spidFlow.finalOut,
		Output: filepath.Join(wihOutDir, "wih.txt"),
		Proxy:  params.opt.Proxy,
	})
	// 漏洞扫描 -> 爬的链接
	new(xscan_spider.ModuleStruct).Run(modules.BaseParams{
		Target: spidFlow.finalOut,
		Output: filepath.Join(vulscanOutDir, "xscan-links.json"),
	})
	if fileutil.IsFile(filepath.Join(vulscanOutDir, "xscan-links.json")) {
		xscan.Clean(filepath.Join(vulscanOutDir, "xscan-links.json"), filepath.Join(vulscanOutDir, "xscan-links.html"))
		new(notify.ModuleStruct).Run(notify.Params{
			Msg: fmt.Sprintf("Task Name: %s, xscan-links.json Count: %d", params.opt.TaskName, fileutil.CountLines(filepath.Join(vulscanOutDir, "xscan-links.json"))),
		})
	}

	new(xray.ModuleStruct).Run(xray.Params{
		BaseParams: &modules.BaseParams{
			Target: spidFlow.finalOut,
			Output: filepath.Join(vulscanOutDir, "xray-links.html"),
		},
	})
	if fileutil.IsFile(filepath.Join(vulscanOutDir, "xray-links.html")) {
		new(notify.ModuleStruct).Run(notify.Params{
			Msg: fmt.Sprintf("Task Name: %s, xray-links.html Count: %d", params.opt.TaskName, fileutil.CountLines(filepath.Join(vulscanOutDir, "xray-links.html"))),
		})
	}
}
