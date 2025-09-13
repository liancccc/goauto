package xscan_spider

import (
	"fmt"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/xscan"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
	*xscan.XscanBeaseModule
}

func (m *ModuleStruct) Name() string {
	return "xscan spider"
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var command string
	if params.IsFileTarget() {
		command = fmt.Sprintf("%s --config %s spider --no-md --file %s --root-scope --xss-json %s --gau", m.GetBin(), m.GetConfigPath(), params.Target, params.Output)
	} else {
		command = fmt.Sprintf("%s --config %s spider --no-md  --url %s --root-scope --xss-json %s --gau", m.GetBin(), m.GetConfigPath(), params.Target, params.Output)
	}

	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}

	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
