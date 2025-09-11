package httpx_unique

import (
	"fmt"
	"path/filepath"
	"time"

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
	return "httpx unique"
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok || !params.IsFileTarget() {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()
	var toolOut = filepath.Join(filepath.Dir(params.Output), fmt.Sprintf("%v.json", time.Now().Unix()))
	var command = fmt.Sprintf("httpx -list %s -hash simhash -json -o %s", params.Target, toolOut)
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
	if !fileutil.FileExists(toolOut) {
		gologger.Error().Str("module", m.Name()).Msgf("%s does not exist", toolOut)
		return
	}
	defer fileutil.Remove(toolOut)
	err = CleanHttpxInvalidTargets(toolOut, params.Output)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
	}
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
