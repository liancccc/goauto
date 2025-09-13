package notify

import (
	"testing"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestName(t *testing.T) {
	params := Params{
		Msg:  "发送测试",
		File: "targets.txt",
	}
	new(ModuleStruct).Run(params)
}
