package uro

import (
	"testing"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestDomain(t *testing.T) {
	t.Log(new(ModuleStruct).Install())
	t.Log(new(ModuleStruct).CheckInstalled())
	params := modules.BaseParams{
		Target: "C:\\Users\\admin\\Downloads\\spiders-clean.txt",
		Output: "test_output\\urls.txt",
	}
	new(ModuleStruct).Run(params)
}
