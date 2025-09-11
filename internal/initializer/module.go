package initializer

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"

	_ "github.com/liancccc/goauto/internal/modules/alterx"
	_ "github.com/liancccc/goauto/internal/modules/cdncheck"
	_ "github.com/liancccc/goauto/internal/modules/gospider"
	_ "github.com/liancccc/goauto/internal/modules/httpx/info"
	_ "github.com/liancccc/goauto/internal/modules/httpx/unique"
	_ "github.com/liancccc/goauto/internal/modules/katana"
	_ "github.com/liancccc/goauto/internal/modules/ksubdomain/enum"
	_ "github.com/liancccc/goauto/internal/modules/ksubdomain/verify"
	_ "github.com/liancccc/goauto/internal/modules/merge"
	_ "github.com/liancccc/goauto/internal/modules/naabu"
	_ "github.com/liancccc/goauto/internal/modules/nuclei"
	_ "github.com/liancccc/goauto/internal/modules/oneforall"
	_ "github.com/liancccc/goauto/internal/modules/subfinder"
	_ "github.com/liancccc/goauto/internal/modules/unique"
	_ "github.com/liancccc/goauto/internal/modules/urlfinder"
	_ "github.com/liancccc/goauto/internal/modules/xray"
	_ "github.com/liancccc/goauto/internal/modules/xscan/spider"
)

func ModuleInstall() {
	var status = make(map[string]bool, len(modules.Modules))
	for _, module := range modules.Modules {
		if module.CheckInstalled() {
			status[module.Name()] = true
		} else {
			if err := module.Install(); err != nil {
				gologger.Error().Msg(err.Error())
			}
			status[module.Name()] = module.CheckInstalled()
		}
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedRowLength(300)
	t.SetTitle("Module Report")
	t.Style().Title.Align = text.AlignCenter
	t.AppendHeader(table.Row{"Module", "Status"})
	for name, installed := range status {
		t.AppendRow(table.Row{name, installed})
	}
	fmt.Println()
	t.Render()
	fmt.Println()
}
