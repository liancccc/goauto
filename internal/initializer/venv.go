package initializer

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/liancccc/goauto/internal/venv/python"
	"github.com/projectdiscovery/gologger"
)

func VenvInit() {
	var status = make(map[string]bool)
	pythonEnv := python.New()
	if !pythonEnv.CheckInited() {
		gologger.Info().Msg("need to initializer python environment")
		pythonEnv.Init()
	}
	status[python.VenvName] = pythonEnv.CheckInited()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedRowLength(300)
	t.SetTitle("Env Report")
	t.Style().Title.Align = text.AlignCenter
	t.AppendHeader(table.Row{"ENV", "STATUS"})
	for name, installed := range status {
		t.AppendRow(table.Row{name, installed})
	}
	fmt.Println()
	t.Render()
	fmt.Println()
}
