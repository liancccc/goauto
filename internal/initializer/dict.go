package initializer

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/liancccc/goauto/internal/dict"
	"github.com/liancccc/goauto/internal/fileutil"
)

func DictInit() {
	var status = make(map[string]bool)
	for name, dictStruct := range dict.Dicts {
		if !fileutil.IsFile(dictStruct.Path) {
			fileutil.Download(dictStruct.Link, dictStruct.Path)
		}
		status[name] = fileutil.IsFile(dictStruct.Path)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedRowLength(300)
	t.SetTitle("Dict Report")
	t.Style().Title.Align = text.AlignCenter
	t.AppendHeader(table.Row{"Dict", "Inited"})
	for name, installed := range status {
		t.AppendRow(table.Row{name, installed})
	}
	fmt.Println()
	t.Render()
	fmt.Println()
}
