package xscan_spider

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
		Output: "SRC-20250907/xscan.txt",
	}
	new(ModuleStruct).Run(params)
}
