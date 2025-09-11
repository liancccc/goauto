package oneforall

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
	pythonBin = python.New().PythonBin
	pipBin    = python.New().PipBin

	oenforallDir           = filepath.Join(paths.ToolsDir, "OneForAll")
	oneforallBin           = filepath.Join(oenforallDir, "oneforall.py")
	oneforallRequirements  = filepath.Join(oenforallDir, "requirements.txt")
	oneforallExec          = fmt.Sprintf("%s %s", pythonBin, oneforallBin)
	installRequireCommands = []string{
		fmt.Sprintf("%s -m pip install -U pip setuptools wheel", pythonBin),
		fmt.Sprintf("%s install -r %s", pipBin, oneforallRequirements),
		fmt.Sprintf("%s install fire", pipBin),
	}
	cloneCommand = fmt.Sprintf("git clone https://github.com/liancccc/OneForAll.git %s", oenforallDir)
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "oneforall"
}

func (m *ModuleStruct) Install() error {
	fileutil.MakeDir(filepath.Dir(oneforallBin))
	if !fileutil.IsFile(oneforallBin) {
		executil.RunCommandSteamOutput(cloneCommand)
	}
	for _, cmd := range installRequireCommands {
		executil.RunCommandSteamOutput(cmd)
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput(oneforallExec)
	return strings.Contains(commandSteamOutput, "python3 oneforall.py")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()
	var csvFilePath = filepath.Join(filepath.Dir(params.Output), "oneforall.csv")
	fileutil.Remove(csvFilePath)

	var command string
	if params.IsFileTarget() {
		command = fmt.Sprintf("%s --target %s --dns False --brute False --req False --path %s run", oneforallExec, params.Target, csvFilePath)
	} else {
		command = fmt.Sprintf("%s --target %s --dns False --brute False --req False --path %s run", oneforallExec, params.Target, csvFilePath)
	}
	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	if !params.IsFileTarget() && fileutil.IsFile(csvFilePath) {
		subdomains := fileutil.GetCsvColumn(csvFilePath, 6)
		fileutil.WriteSliceToFile(params.Output, subdomains)
		fileutil.Remove(csvFilePath)
	}
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
