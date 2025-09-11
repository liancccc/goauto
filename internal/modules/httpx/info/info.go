package httpx_info

import (
	"fmt"
	"path/filepath"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/httpx"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
	*httpx.HttpxBeaseModule
}

func (m *ModuleStruct) Name() string {
	return "httpx info"
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
		command = fmt.Sprintf("httpx -list %s -title -screenshot -location -no-color -sc -server -td -srd %s -o %s", params.Target, filepath.Dir(params.Output), params.Output)
	} else {
		command = fmt.Sprintf("httpx -target %s -title -screenshot -location -no-color -sc -server -td -srd %s -o %s", params.Target, filepath.Dir(params.Output), params.Output)
	}
	if params.CustomizeParams != "" {
		command = fmt.Sprintf("%s %s", command, params.CustomizeParams)
	}
	if params.Proxy != "" {
		command = fmt.Sprintf("%s -proxy %s", command, params.Proxy)
	}

	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
