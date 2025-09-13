package cdncheck

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
	params := Params{
		BaseParams: &modules.BaseParams{
			Target: "all.txt",
		},
		CDNPath:   "test_output/cdn.txt",
		NoCDNPath: "test_output/noCdn.txt",
	}
	new(ModuleStruct).Run(params)
}
