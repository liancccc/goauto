package wih

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/paths"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "wih"
}

func (m *ModuleStruct) Install() error {
	var downloadUrl = "https://raw.githubusercontent.com/liancccc/arl_files/master/wih/%s"
	if runtime.GOOS == "windows" {
		downloadUrl = fmt.Sprintf(downloadUrl, "wih_amd64.exe")
	} else {
		downloadUrl = fmt.Sprintf(downloadUrl, fmt.Sprintf("wih_%s_%s", runtime.GOOS, runtime.GOARCH))
	}
	fileutil.Download(downloadUrl, m.GetBin())
	executil.RunCommandSteamOutput(fmt.Sprintf("%s -G", m.GetBin()))
	return nil
}

func (m *ModuleStruct) GetBin() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(paths.ToolsDir, "wih/wih.exe")
	} else {
		return filepath.Join(paths.ToolsDir, "wih/wih")
	}
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput(fmt.Sprintf("%s -h", m.GetBin()))
	executil.RunCommandSteamOutput(fmt.Sprintf("%s -G", m.GetBin()))
	return strings.Contains(commandSteamOutput, "WebInfoHunter")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var command = fmt.Sprintf("%s --target %s --output %s --disable-ak-sk-output --output-json --max-collect 5000", m.GetBin(), params.Target, params.Output)
	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	Clean(params.Output)
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
