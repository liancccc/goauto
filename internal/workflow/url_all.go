package workflow

import (
	"fmt"
	"path/filepath"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/httpx"
	"github.com/liancccc/goauto/internal/modules/notify"
	"github.com/liancccc/goauto/internal/modules/nuclei"
	"github.com/liancccc/goauto/internal/modules/wih"
	"github.com/liancccc/goauto/internal/modules/xray"
	"github.com/liancccc/goauto/internal/modules/xscan"
	xscan_spider "github.com/liancccc/goauto/internal/modules/xscan/spider"
)

func init() {
	RegisterWorkflow(&urlAllFlow{})
}

type urlAllFlow struct {
}

func (d *urlAllFlow) Name() string {
	return "urlAll"
}

func (d *urlAllFlow) Description() string {
	return "url -> httpx -> vulscan alive url -> spider -> vulscan spider link"
}

func (d *urlAllFlow) Run(params *workflowParams) {
	// httpx web 验活
	var httpxOutDir = filepath.Join(params.workSpace, "httpx")
	new(httpx.ModuleStruct).Run(modules.BaseParams{
		Target:          params.target,
		Output:          filepath.Join(httpxOutDir, "all.json"),
		CustomizeParams: "-hash simhash -json",
		Proxy:           params.opt.Proxy,
	})
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
