package httpx_unique

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
		Output: "test_output/alive.txt",
	}
	new(ModuleStruct).Run(params)
}

func TestDomains(t *testing.T) {
	params := modules.BaseParams{
		Target:          "urls.txt",
		Output:          "test_output/alive.txt",
		CustomizeParams: "-mc 200,302 -p 80,443,8080,8000,8888,4848,7070,8089,8181,9080,9443,5000,8443,5001,81,8081,50805,3000,88,7547",
	}
	new(ModuleStruct).Run(params)
}
