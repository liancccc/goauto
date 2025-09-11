package urlfinder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/merge"
	"github.com/liancccc/goauto/internal/netutil"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "urlfinder"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install -v github.com/projectdiscovery/urlfinder/cmd/urlfinder@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("urlfinder")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var toolOutDir = filepath.Join(filepath.Dir(params.Output), "urlfinder")
	fileutil.MakeDir(toolOutDir)
	var commands []string
	var hostname string
	if params.IsFileTarget() {
		lines := fileutil.ReadingLines(params.Target)
		for _, target := range lines {
			hostname = netutil.GetUrlHostname(target)
			if hostname == "" {
				continue
			}
			commands = append(commands, fmt.Sprintf("urlfinder -d %s -all -o %s", hostname, filepath.Join(toolOutDir, fileutil.GetUrlFileName(target))))
			hostname = ""
		}
	} else {
		hostname = netutil.GetUrlHostname(params.Target)
		if hostname == "" {
			gologger.Error().Str("target", params.Target).Msg("invalid target")
			return
		}
		commands = append(commands, fmt.Sprintf("urlfinder -d %s -all -o %s", hostname, filepath.Join(toolOutDir, fileutil.GetUrlFileName(params.Target))))
	}

	for _, command := range commands {
		if params.CustomizeParams != "" {
			command = fmt.Sprintf("%s %s", command, params.CustomizeParams)
		}
		if params.Proxy != "" {
			command = fmt.Sprintf("%s -proxy %s", command, params.Proxy)
		}
		_, err := executil.RunCommandSteamOutput(command, params.Timeout)
		if err != nil {
			gologger.Error().Str("module", m.Name()).Msg(err.Error())
			continue
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
	}
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
