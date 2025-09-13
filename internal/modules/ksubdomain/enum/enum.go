package ksubdomain_enum

import (
	"fmt"
	"path/filepath"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/ksubdomain"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
	*ksubdomain.KBaseModule
}

func (m *ModuleStruct) Name() string {
	return "ksubdomain enum"
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}
	_ = params.MkOutDir()
	var toolOut = filepath.Join(filepath.Dir(params.Output), "ksubdomain-enum.json")
	fileutil.Remove(toolOut)
	var command string
	if params.IsFileTarget() {
		command = fmt.Sprintf("ksubdomain enum -ds %s --wild-filter-mode basic --output-type json -o %s", params.Target, toolOut)
	} else {
		command = fmt.Sprintf("ksubdomain enum -d %s --wild-filter-mode basic --output-type json -o %s", params.Target, toolOut)
	}

	if params.Dict != "" {
		command = fmt.Sprintf("%s -f %s", command, params.Dict)
	}

	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}

	err = ksubdomain.CleanResult(toolOut, params.Output)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	fileutil.Remove(toolOut)
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)

}
