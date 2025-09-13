package cdncheck

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
)

const (
	installCommand = "go install -v github.com/liancccc/cdncheck/cmd/cdncheck@latest"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type Params struct {
	*modules.BaseParams
	CDNPath   string
	NoCDNPath string
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "cdncheck"
}

func (m *ModuleStruct) Install() error {
	_, err := executil.RunCommandSteamOutput(installCommand)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("cdncheck")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(Params)
	if !ok || !params.IsFileTarget() {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	fileutil.MakeDir(filepath.Dir(params.CDNPath))
	fileutil.MakeDir(filepath.Dir(params.NoCDNPath))
	var toolOut = filepath.Join(filepath.Dir(params.CDNPath), "cdncheck.json")
	var command = fmt.Sprintf("cdncheck -i %s -jsonl -o %s", params.Target, toolOut)

	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	err = CleanCdnCheckResult(params.Target, toolOut, params.CDNPath, params.NoCDNPath)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
	}
	var msg = fmt.Sprintf("CDNOutput: %s, Count: %d", params.CDNPath, fileutil.CountLines(params.CDNPath))
	gologger.Info().Str("module", m.Name()).Msg(msg)
	msg = fmt.Sprintf("NoCDNOutput: %s, Count: %d", params.NoCDNPath, fileutil.CountLines(params.NoCDNPath))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
