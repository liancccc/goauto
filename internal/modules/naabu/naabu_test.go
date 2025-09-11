package naabu

import (
	"testing"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestRunHostNaabu(t *testing.T) {
	var target = "molimer.com"
	new(ModuleStruct).Run(modules.BaseParams{
		Target: target,
		Output: "test_output/molimer.com.services.txt",
	})
}

func TestRunHostNaabuFile(t *testing.T) {
	var target = "targets.txt"
	new(ModuleStruct).Run(modules.BaseParams{
		Target: target,
		Output: "test_output/services.txt",
	})
}
