package httpx_info

import (
	"testing"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestUrls(t *testing.T) {
	params := modules.BaseParams{
		Target: "urls.txt",
		Output: "test_output/web.txt",
	}
	new(ModuleStruct).Run(params)
}
