package unique

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "unique"
}

func (m *ModuleStruct) Install() error {
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	var command string
	var flag string
	if runtime.GOOS == "windows" {
		command = `powershell -Command "Get-Help Sort-Object -Detailed"`
		flag = "-Unique"
	} else {
		command = `sort --help`
		flag = "Usage:"
	}
	commandSteamOutput, _ := executil.RunCommandSteamOutput(command)
	return strings.Contains(commandSteamOutput, flag)
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var command string
	if runtime.GOOS == "windows" {
		command = fmt.Sprintf(`powershell -Command "Get-Content '%s' -Encoding UTF8 | Sort-Object -Unique | Out-File '%s' -Encoding UTF8"`, params.Target, params.Output)
	} else {
		command = fmt.Sprintf(`sort -u %s > %s`, params.Target, params.Output)
	}
	_, err := executil.RunCommandSteamOutput(command)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}

	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
