package subfinder

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
		Target: "monzo.com",
		Output: "test_output/monzo.com.txt",
	}
	new(ModuleStruct).Run(params)
}

func TestDomainFile(t *testing.T) {
	params := modules.BaseParams{
		Target: "targets.txt",
		Output: "test_output/subdomains.txt",
	}
	new(ModuleStruct).Run(params)
}
