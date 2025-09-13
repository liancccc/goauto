package uncover

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/merge"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type Target struct {
	Query  string
	Engine string
}

type Params struct {
	*modules.BaseParams
	Targets []Target
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "uncover"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install -v github.com/projectdiscovery/uncover/cmd/uncover@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("uncover")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(Params)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()
	var toolOutDir = filepath.Join(filepath.Dir(params.Output), "uncover")
	fileutil.MakeDir(toolOutDir)

	var commands []string
	for _, target := range params.Targets {
		command := fmt.Sprintf("uncover -q '%s' -e %s -o %s.txt", target.Query, target.Engine, filepath.Join(toolOutDir, fileutil.GetUnixNmae()))
		commands = append(commands, command)
	}

	for _, command := range commands {
		_, err := executil.RunCommandSteamOutput(command, params.Timeout)
		if err != nil {
			gologger.Error().Str("module", m.Name()).Msg(err.Error())
		}
	}

	var toolOutputFiles []string
	filepath.Walk(toolOutDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			toolOutputFiles = append(toolOutputFiles, path)
		}
		return nil
	})
	if len(toolOutputFiles) > 0 {
		new(merge.ModuleStruct).Run(merge.Params{
			BaseParams: &modules.BaseParams{Output: params.Output},
			Targets:    toolOutputFiles,
		})
		fileutil.Remove(toolOutDir)
	}

	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
