package dnsx

import (
	"testing"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestDomainFile(t *testing.T) {
	params := modules.BaseParams{
		Target: "noCdn.txt",
		Output: "test_output/dnsx.json",
	}
	new(ModuleStruct).Run(params)
}

func TestClean(t *testing.T) {
	t.Log(CleanAndGenCustomizeFormat("test_output/dnsx.json", "test_output/ips.txt"))
}
