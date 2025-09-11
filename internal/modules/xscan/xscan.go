package xscan

import (
	"path/filepath"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/paths"
)

type XscanBeaseModule struct {
}

func (m *XscanBeaseModule) Install() error {
	return nil
}

func (m *XscanBeaseModule) GetBin() string {
	return filepath.Join(paths.ToolsDir, "xscan/xscan")
}

func (m *XscanBeaseModule) GetConfigPath() string {
	return filepath.Join(paths.ToolsDir, "xscan/config.yaml")
}

func (m *XscanBeaseModule) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput(m.GetBin())
	return strings.Contains(commandSteamOutput, "w8ay")
}
