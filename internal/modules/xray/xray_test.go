package xray

import (
	"testing"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestName(t *testing.T) {
	new(ModuleStruct).CheckInstalled()
}
