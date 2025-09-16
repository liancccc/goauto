package workflow

import (
	"fmt"
	"path/filepath"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/notify"
	"github.com/liancccc/goauto/internal/modules/nuclei"
	"github.com/liancccc/goauto/internal/modules/xray"
	"github.com/liancccc/goauto/internal/modules/xscan"
	xscan_spider "github.com/liancccc/goauto/internal/modules/xscan/spider"
)

type vulscanFlow struct {
}

func (v *vulscanFlow) Name() string {
	return "vulscan"
}

func (v *vulscanFlow) Description() string {
	return "vulscan: wih + xscan + xray"
}

func (v *vulscanFlow) Run(params *workflowParams) {
	var vulscanOutDir = filepath.Join(params.workSpace, "vulscan")
	new(xscan_spider.ModuleStruct).Run(modules.BaseParams{
		Target: params.target,
		Output: filepath.Join(vulscanOutDir, "xscan-spider.json"),
	})
	if fileutil.IsFile(filepath.Join(vulscanOutDir, "xscan-spider.json")) {
		xscan.Clean(filepath.Join(vulscanOutDir, "xscan-spider.json"), filepath.Join(vulscanOutDir, "xscan-spider.html"))
		new(notify.ModuleStruct).Run(notify.Params{
			Msg: fmt.Sprintf("Task Name: %s, xscan-spider.json Count: %d", params.opt.TaskName, fileutil.CountLines(filepath.Join(vulscanOutDir, "xscan-spider.json"))),
		})
	}

	new(xray.ModuleStruct).Run(xray.Params{
		BaseParams: &modules.BaseParams{
			Target: params.target,
			Output: filepath.Join(vulscanOutDir, "xray.html"),
		},
	})
	if fileutil.IsFile(filepath.Join(vulscanOutDir, "xray.html")) {
		new(notify.ModuleStruct).Run(notify.Params{
			Msg: fmt.Sprintf("Task Name: %s, xray.html Count: %d", params.opt.TaskName, fileutil.CountLines(filepath.Join(vulscanOutDir, "xray-links.html"))),
		})
	}

	new(nuclei.ModuleStruct).Run(modules.BaseParams{
		Target: params.target,
		Output: filepath.Join(vulscanOutDir, "nuclei.txt"),
		Proxy:  params.opt.Proxy,
	})
	if fileutil.IsFile(filepath.Join(vulscanOutDir, "nuclei.txt")) {
		new(notify.ModuleStruct).Run(notify.Params{
			Msg: fmt.Sprintf("Task Name: %s, nuclei.txt Count: %d", params.opt.TaskName, fileutil.CountLines(filepath.Join(vulscanOutDir, "nuclei.txt"))),
		})
	}
}
