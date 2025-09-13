package notify

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

type Params struct {
	Msg  string
	File string
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "notify"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install -v github.com/projectdiscovery/notify/cmd/notify@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("notify")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(Params)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}
	if params.Msg != "" && params.File == "" {
		params.File, _ = fileutil.WriteTempFile(params.Msg)
		defer fileutil.Remove(params.File)
	}

	var command = fmt.Sprintf("notify -data %s -bulk", params.File)
	_, err := executil.RunCommandSteamOutput(command)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	var msg string
	if params.Msg != "" {
		msg = fmt.Sprintf("Notify Msg %s Success", params.Msg)
	} else {
		msg = fmt.Sprintf("Notify File %s Success", params.File)
	}
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
