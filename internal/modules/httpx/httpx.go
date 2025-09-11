package httpx

import (
	"strings"

	"github.com/liancccc/goauto/internal/executil"
)

type HttpxBeaseModule struct {
}

func (m *HttpxBeaseModule) Install() error {
	var installCmd = "go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *HttpxBeaseModule) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("httpx")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}
