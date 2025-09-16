package workflow

import (
	"testing"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

func init() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
}

func TestSubdomainFlow(t *testing.T) {
	subFlow := new(subdomainFlow)
	subFlow.Run(&workflowParams{
		target:    "nba.com",
		workSpace: "test_workspace",
	})
	t.Log(subFlow)
}
