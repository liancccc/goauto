package gospider

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
	return "gospider"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install github.com/jaeles-project/gospider@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("gospider -h")
	return strings.Contains(commandSteamOutput, "Usage:")
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}

	_ = params.MkOutDir()

	var toolOutDir = filepath.Join(filepath.Dir(params.Output), "gospider")
	defer fileutil.Remove(toolOutDir)
	fileutil.MakeDir(toolOutDir)
	var command string
	if params.IsFileTarget() {
		command = fmt.Sprintf("gospider -S %s -o %s --json -c 50 -d 2 --other-source --subs --sitemap --robots", params.Target, toolOutDir)
	} else {
		command = fmt.Sprintf("gospider -s %s -o %s --json -c 50 -d 2 --other-source --subs --sitemap --robots", params.Target, toolOutDir)
	}
	if params.CustomizeParams != "" {
		command = fmt.Sprintf("%s %s", command, params.CustomizeParams)
	}
	if params.Proxy != "" {
		command = fmt.Sprintf("%s --proxy %s", command, params.Proxy)
	}
	_, err := executil.RunCommandSteamOutput(command, params.Timeout)
	if err != nil {
		gologger.Error().Str("module", m.Name()).Msg(err.Error())
		return
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
		var allJsonOut = filepath.Join(filepath.Dir(params.Output), "gospider.json")
		new(merge.ModuleStruct).Run(merge.Params{
			BaseParams: &modules.BaseParams{Output: allJsonOut},
			Targets:    toolOutputFiles,
		})
		// 只抓主域名和 http 的链接
		Clean(allJsonOut, params.Output)
	}
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}
