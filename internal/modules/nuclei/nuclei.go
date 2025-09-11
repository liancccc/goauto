package nuclei

import (
	"fmt"
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
	return "nuclei"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("nuclei")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var command string
	if params.IsFileTarget() {
		command = fmt.Sprintf("nuclei -l %s -o %s", params.Target, params.Output)
	} else {
		command = fmt.Sprintf("nuclei -u %s -o %s", params.Target, params.Output)
	}
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
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
