package katana

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

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "katana"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install github.com/projectdiscovery/katana/cmd/katana@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("katana")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var toolOutDir = filepath.Join(filepath.Dir(params.Output), "katana")
	fileutil.MakeDir(toolOutDir)
	var commands []string
	if params.IsFileTarget() {
		lines := fileutil.ReadingLines(params.Target)
		for _, target := range lines {
			commands = append(commands, fmt.Sprintf("katana -u %s -sc -kf robotstxt,sitemapxml -jc -c 50 -o %s", target, filepath.Join(toolOutDir, fileutil.GetUrlFileName(target))))
		}
	} else {
		commands = append(commands, fmt.Sprintf("katana -u %s -sc -kf robotstxt,sitemapxml -jc -c 50 -o %s", params.Target, filepath.Join(toolOutDir, fileutil.GetUrlFileName(params.Target))))
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
