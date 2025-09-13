package uro

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/paths"
	"github.com/liancccc/goauto/internal/venv/python"
	"github.com/projectdiscovery/gologger"
)

var (
	uroDir       = filepath.Join(paths.ToolsDir, "uro")
	urlPy        = filepath.Join(uroDir, "uro.py")
	cloneCommand = fmt.Sprintf("git clone https://github.com/liancccc/uro.git %s", uroDir)
	pythonBin    = python.New().PythonBin
	uroExec      = fmt.Sprintf("%s %s", pythonBin, urlPy)
)

var uroBin = filepath.Join(fileutil.GetHomeDir(), ".local", "bin", "uro")

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "uro"
}

func (m *ModuleStruct) Install() error {
	_, err := executil.RunCommandSteamOutput(cloneCommand)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput(fmt.Sprintf("%s -h", uroExec))
	return strings.Contains(commandSteamOutput, "-i")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok || !params.IsFileTarget() {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var command = fmt.Sprintf("%s -i %s -o %s", uroExec, params.Target, params.Output)
	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
