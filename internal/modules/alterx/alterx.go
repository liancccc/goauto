package alterx

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
)

const (
	installCommand = "go install github.com/projectdiscovery/alterx/cmd/alterx@latest"
	toolBin        = "alterx"
	linuxCat       = "cat"
	windowsCat     = "type"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "alterx"
}

func (m *ModuleStruct) Install() error {
	_, err := executil.RunCommandSteamOutput(installCommand)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("alterx")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok || !params.IsFileTarget() {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var catBin string

	if runtime.GOOS == "windows" {
		catBin = windowsCat
	} else {
		catBin = linuxCat
	}

	var command = fmt.Sprintf("%s %s | alterx -enrich -o %s", catBin, params.Target, params.Output)

	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}

	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
