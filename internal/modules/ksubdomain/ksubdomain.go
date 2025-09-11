package ksubdomain

import (
	"strings"

	"github.com/liancccc/goauto/internal/executil"
)

type KBaseModule struct {
}

func (m *KBaseModule) Install() error {
	var installCmd = "go install -v github.com/boy-hack/ksubdomain_enum/v2/cmd/ksubdomain_enum@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *KBaseModule) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("ksubdomain")
	return strings.Contains(commandSteamOutput, "enum")
}
