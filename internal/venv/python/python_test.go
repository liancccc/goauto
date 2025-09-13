package python

import (
	"fmt"
	"testing"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestName(t *testing.T) {
	executil.RunCommandSteamOutput(fmt.Sprintf("%s install uro", New().PipxBin))
}
