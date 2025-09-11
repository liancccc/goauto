package ksubdomain_verify

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
	return "ksubdomain verify"
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok || !params.IsFileTarget() {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}
	_ = params.MkOutDir()
	var toolOut = filepath.Join(filepath.Dir(params.Output), "ksubdomain-verify.txt")
	fileutil.Remove(toolOut)
	var command = fmt.Sprintf("ksubdomain verify -f %s -o %s --wild-filter-mode advanced", params.Target, toolOut)
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
