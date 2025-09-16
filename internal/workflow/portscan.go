package workflow

import (
	"path/filepath"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/naabu"
)

type portscanFlow struct {
	output string
}

func (p *portscanFlow) Name() string {
	return "portscan"
}

func (p *portscanFlow) Description() string {
	return "naabu + nmap portscan"
}

func (p *portscanFlow) Run(params *workflowParams) {
	var outDir = filepath.Join(params.workSpace, "portscan")
	new(naabu.ModuleStruct).Run(modules.BaseParams{
		Target: params.target,
		Output: filepath.Join(outDir, "portscan.txt"),
	})
}
