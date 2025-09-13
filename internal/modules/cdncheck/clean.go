package cdncheck

import (
	"encoding/json"
	"strings"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/merge"
	"github.com/projectdiscovery/gologger"
)

type JsonData struct {
	Input string `json:"input"`
}

func CleanCdnCheckResult(target, toolResultPath, cdnPath, noCdnPath string) error {
	gologger.Info().Msgf("Cleaning Cdncheck result: %s", toolResultPath)

	if fileutil.CountLines(toolResultPath) == 0 {
		gologger.Info().Msgf("Cdn check not found in %s, all no cdn", toolResultPath)
		if fileutil.IsFile(target) {
			new(merge.ModuleStruct).Run(merge.Params{
				BaseParams: &modules.BaseParams{
					Output: noCdnPath,
				},
				Targets: []string{target},
			})
		} else {
			fileutil.WriteSliceToFile(noCdnPath, []string{target})
		}
		return nil
	}

	var targets []string
	if fileutil.IsFile(target) {
		targets = fileutil.ReadingLines(target)
	} else {
		targets = []string{target}
	}
	var cdncheckResult JsonData
	var cdnTargetMap = make(map[string]struct{}, len(targets))
	var cdnTargets []string
	cdnCheckResults := fileutil.ReadingLines(toolResultPath)
	for _, line := range cdnCheckResults {
		if strings.Contains(strings.ToLower(line), "cdn") {
			if err := json.Unmarshal([]byte(line), &cdncheckResult); err != nil {
				continue
			}
			gologger.Debug().Msgf("Found cdn domain: %s", cdncheckResult.Input)
			cdnTargetMap[cdncheckResult.Input] = struct{}{}
			cdnTargets = append(cdnTargets, cdncheckResult.Input)
			cdncheckResult = JsonData{}
		}
	}

	var noCdnTargets []string

	for _, value := range targets {
		if _, ok := cdnTargetMap[value]; !ok {
			noCdnTargets = append(noCdnTargets, value)
		}
	}

	gologger.Debug().Msgf("Found %d CDN Target", len(cdnTargetMap))
	gologger.Debug().Msgf("Found %d No CDN Target", len(noCdnTargets))

	fileutil.WriteSliceToFile(cdnPath, cdnTargets)
	fileutil.WriteSliceToFile(noCdnPath, noCdnTargets)
	return nil
}
