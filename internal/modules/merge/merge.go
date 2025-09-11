package merge

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
	getContentHelpCommand = `powershell -Command "Get-Help Get-Content -Detailed"`
	outFileHelpCommand    = `powershell -Command "Get-Help Out-File -Detailed"`
	catHelpCommand        = "cat --help"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type Params struct {
	*modules.BaseParams
	Targets []string
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "merge"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	var installed bool
	if runtime.GOOS == "windows" {
		var (
			getContentInstalled, outFileInstalled bool
		)
		if output, _ := executil.RunCommandSteamOutput(getContentHelpCommand); strings.Contains(output, "Get-Content") {
			getContentInstalled = true
		}
		if output, _ := executil.RunCommandSteamOutput(outFileHelpCommand); strings.Contains(output, "Out-File") {
			outFileInstalled = true
		}
		installed = getContentInstalled && outFileInstalled
	} else {
		if output, _ := executil.RunCommandSteamOutput(catHelpCommand); strings.Contains(output, "Usage:") {
			installed = true
		}
	}
	return installed
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(Params)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var commands []string
	var command string

	for _, target := range params.Targets {
		if !fileutil.FileExists(target) {
			continue
		}
		if runtime.GOOS == "windows" {
			command = fmt.Sprintf(`powershell -Command "Get-Content -Path '%s' -Encoding UTF8 | Out-File -FilePath '%s' -Encoding UTF8 -Append"`, target, params.Output)
		} else {
			command = fmt.Sprintf("cat %s >> %s", target, params.Output)
		}
		commands = append(commands, command)
	}
	for _, cmd := range commands {
		_, err := executil.RunCommandSteamOutput(cmd)
		if err != nil {
			gologger.Error().Str("module", m.Name()).Msg(err.Error())
			return
		}
	}

	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
