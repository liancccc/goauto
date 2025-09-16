package workflow

import (
	"path/filepath"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/cdncheck"
)

func init() {
	RegisterWorkflow(&cdncheckFlow{})
}

type cdncheckFlow struct {
	cdnOut   string
	noCdnOut string
}

func (s *cdncheckFlow) Name() string {
	return "cdncheck"
}

func (s *cdncheckFlow) Description() string {
	return "cdncheck"
}

func (s *cdncheckFlow) Run(params *workflowParams) {
	var cdncheckOutDir = filepath.Join(params.workSpace, "cdncheck")
	s.cdnOut = filepath.Join(cdncheckOutDir, "cdn.txt")
	s.noCdnOut = filepath.Join(cdncheckOutDir, "noCdn.txt")
	new(cdncheck.ModuleStruct).Run(cdncheck.Params{
		BaseParams: &modules.BaseParams{
			Target: params.target,
		},
		CDNPath:   s.cdnOut,
		NoCDNPath: s.noCdnOut,
	})
}
