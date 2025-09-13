package uncover

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
	params := Params{
		BaseParams: &modules.BaseParams{
			Output: "test_output/services.txt",
		},
		Targets: []Target{
			{
				Query:  `domain: "kucoin.com"`,
				Engine: "quake",
			},
		},
	}
	new(ModuleStruct).Run(params)
}
