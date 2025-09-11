package katana

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
	params := modules.BaseParams{
		Target: "https://wiki.xazlsec.com:443",
		Output: "test_output/wiki.xazlsec.com.txt",
	}
	new(ModuleStruct).Run(params)
}
