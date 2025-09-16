package workflow

import (
	"path/filepath"
	"sync"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/alterx"
	ksubdomain_enum "github.com/liancccc/goauto/internal/modules/ksubdomain/enum"
	ksubdomain_verify "github.com/liancccc/goauto/internal/modules/ksubdomain/verify"
	"github.com/liancccc/goauto/internal/modules/oneforall"
	"github.com/liancccc/goauto/internal/modules/subfinder"
)

func init() {
	RegisterWorkflow(&subdomainFlow{})
}

type subdomainFlow struct {
	finalOut string
}

func (s *subdomainFlow) Name() string {
	return "subdomain"
}

func (s *subdomainFlow) Description() string {
	return "subfinder + oneforall + ksubdomain enum -> alterx -> ksubdomain verify"
}

func (s *subdomainFlow) Run(params *workflowParams) {
	var subdomainOutDir = filepath.Join(params.workSpace, "subdomain")
	s.finalOut = filepath.Join(subdomainOutDir, "final.txt")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		new(subfinder.ModuleStruct).Run(modules.BaseParams{
			Target: params.target,
			Output: filepath.Join(subdomainOutDir, "subfinder.txt"),
		})
	}()
	go func() {
		defer wg.Done()
		new(oneforall.ModuleStruct).Run(modules.BaseParams{
			Target: params.target,
			Output: filepath.Join(subdomainOutDir, "oneforall.txt"),
		})
	}()
	wg.Wait()
	new(ksubdomain_enum.ModuleStruct).Run(modules.BaseParams{
		Target:  params.target,
		Output:  filepath.Join(subdomainOutDir, "ksubdomain.txt"),
		Timeout: "5h",
	})
	MergeAndUnique(
		[]string{
			filepath.Join(subdomainOutDir, "subfinder.txt"),
			filepath.Join(subdomainOutDir, "oneforall.txt"),
			filepath.Join(subdomainOutDir, "ksubdomain.txt"),
		},
		filepath.Join(subdomainOutDir, "merge.txt"),
	)

	new(alterx.ModuleStruct).Run(modules.BaseParams{
		Target:  filepath.Join(subdomainOutDir, "merge.txt"),
		Output:  filepath.Join(subdomainOutDir, "alterx.txt"),
		Timeout: "1h",
	})
	new(ksubdomain_verify.ModuleStruct).Run(modules.BaseParams{
		Target: filepath.Join(subdomainOutDir, "alterx.txt"),
		Output: filepath.Join(subdomainOutDir, "alterx-alive.txt"),
	})
	MergeAndUnique(
		[]string{
			filepath.Join(subdomainOutDir, "merge.txt"),
			filepath.Join(subdomainOutDir, "alterx-alive.txt"),
		},
		filepath.Join(subdomainOutDir, s.finalOut),
	)
	fileutil.Cleaning(subdomainOutDir, []string{s.finalOut})
}
