package wih

import (
	"testing"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestInstall(t *testing.T) {
	new(ModuleStruct).Install()
}

func TestDomainFile(t *testing.T) {
	//params := modules.BaseParams{
	//	Target: "all.txt",
	//	Output: "test_output/wih.json",
	//}
	//new(ModuleStruct).Run(params)
	Clean("test_output/wih.json")
}
