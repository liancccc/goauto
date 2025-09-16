package workflow

import (
	"path/filepath"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/gospider"
	"github.com/liancccc/goauto/internal/modules/katana"
	"github.com/liancccc/goauto/internal/modules/urlfinder"
	"github.com/liancccc/goauto/internal/modules/uro"
)

type spiderFlow struct {
	finalOut string
}

func (s *spiderFlow) Name() string {
	return "spider"
}

func (s *spiderFlow) Description() string {
	return "gospider + katana + urlfinder"
}

func (s *spiderFlow) Run(params *workflowParams) {
	var outDir = filepath.Join(params.workSpace, "spider")
	s.finalOut = filepath.Join(outDir, "links.txt")
	new(gospider.ModuleStruct).Run(modules.BaseParams{
		Target: params.target,
		Output: filepath.Join(outDir, "gospider.txt"),
		Proxy:  params.opt.Proxy,
	})
	new(katana.ModuleStruct).Run(modules.BaseParams{
		Target: params.target,
		Output: filepath.Join(outDir, "katana.txt"),
		Proxy:  params.opt.Proxy,
	})
	new(urlfinder.ModuleStruct).Run(modules.BaseParams{
		Target: params.target,
		Output: filepath.Join(outDir, "urlfinder.txt"),
		Proxy:  params.opt.Proxy,
	})
	var mergeFiles = []string{
		filepath.Join(outDir, "gospider.txt"),
		filepath.Join(outDir, "katana.txt"),
		filepath.Join(outDir, "urlfinder.txt"),
	}
	MergeAndUnique(mergeFiles, filepath.Join(outDir, "unique.txt"))
	new(uro.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(outDir, "unique.txt"),
		Output: s.finalOut,
	})
}
