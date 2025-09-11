package xray

import (
	"testing"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestName(t *testing.T) {
	//new(ModuleStruct).Install()
	//new(ModuleStruct).CheckInstalled()
	params := Params{
		BaseParams: &modules.BaseParams{
			Target: "http://testphp.vulnweb.com",
			Output: "xray.html",
		},
	}
	new(ModuleStruct).Run(params)
}
