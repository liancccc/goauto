package workflow

import (
	"path/filepath"

	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/dnsx"
)

type dnsxFlow struct {
	output       string              // dnsx 结果
	ipDomainsMap map[string][]string // IP => Domains
	ipsFile      string              // ip 文件
}

func (s *dnsxFlow) Name() string {
	return "dnsx"
}

func (s *dnsxFlow) Description() string {
	return "domain -> dnsx get a ip address"
}

func (s *dnsxFlow) Run(params *workflowParams) {
	var outDir = filepath.Join(params.workSpace, "dnsx")
	s.output = filepath.Join(outDir, "dnsx.json")
	s.ipsFile = filepath.Join(outDir, "ips.txt")
	new(dnsx.ModuleStruct).Run(modules.BaseParams{
		Target: params.target,
		Output: s.output,
	})
	ipDomainsMap, err := dnsx.CleanAndGenCustomizeFormat(s.output, s.ipsFile)
	if err == nil {
		s.ipDomainsMap = ipDomainsMap
	} else {
		s.ipDomainsMap = make(map[string][]string)
	}
}
