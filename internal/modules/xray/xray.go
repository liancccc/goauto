package xray

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

type Params struct {
	*modules.BaseParams
	Listen string
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "xray"
}

func (m *ModuleStruct) Install() error {
	var downloadUrl = fmt.Sprintf("https://github.com/chaitin/xray/releases/download/1.9.11/%s", m.GetDownloadName())
	var zipOut = filepath.Join(paths.ToolsDir, "xray.zip")
	if !fileutil.Download(downloadUrl, zipOut) {
		return fmt.Errorf("xray.zip does not exist, download fail")
	}
	var toolDir = filepath.Join(paths.ToolsDir, "xray")
	fileutil.Unzip(zipOut, toolDir)
	var rawBinName = strings.Replace(m.GetDownloadName(), ".zip", "", 1)
	if !fileutil.IsFile(filepath.Join(toolDir, rawBinName)) {
		return fmt.Errorf("unzip xray.zip fail")
	}
	fileutil.Move(filepath.Join(toolDir, rawBinName), m.GetBin())
	executil.RunCommandSteamOutput(fmt.Sprintf("%s genca", m.GetBin()))
	fileutil.Remove(zipOut)
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput(m.GetBin())
	return strings.Contains(commandSteamOutput, "USAGE:")
}

func (m *ModuleStruct) GetDownloadName() string {
	var ext string
	if runtime.GOOS == "windows" {
		ext = ".exe.zip"
	} else {
		ext = ".zip"
	}
	return fmt.Sprintf("xray_%s_%s%s", runtime.GOOS, runtime.GOARCH, ext)
}

func (m *ModuleStruct) GetBin() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(paths.ToolsDir, "xray/xray.exe")
	} else {
		return filepath.Join(paths.ToolsDir, "xray/xray")
	}
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(Params)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var command string

	if params.Listen != "" {
		command = fmt.Sprintf("%s webscan --listen %s --html-output %s", m.GetBin(), params.Listen, params.Output)
	} else if params.IsFileTarget() {
		command = fmt.Sprintf("%s webscan --url-file %s --html-output %s", m.GetBin(), params.Target, params.Output)
	} else {
		command = fmt.Sprintf("%s webscan --url %s --html-output %s", m.GetBin(), params.Target, params.Output)
	}
	if params.CustomizeParams != "" {
		command = fmt.Sprintf("%s %s", command, params.CustomizeParams)
	}

	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
