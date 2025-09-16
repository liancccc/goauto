package wih

import (
	"encoding/json"
	"path/filepath"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/reportutil"
	"github.com/projectdiscovery/gologger"
)

type JsonData struct {
	Records []struct {
		ID        string `json:"id"`
		Content   string `json:"content"`
		Source    string `json:"source"`
		SourceTag string `json:"source_tag"`
		Count     int    `json:"count"`
		Hash      int64  `json:"-"`
		Tag       string `json:"tag"`
	} `json:"records"`
	Target string `json:"target"`
}

func Clean(toolOutput string) {
	gologger.Debug().Msgf("Cleaning %s", toolOutput)
	var outDir = filepath.Dir(toolOutput)
	var jsonData JsonData
	var seen = make(map[string]bool)
	lines := fileutil.ReadingLines(toolOutput)
	var reports []reportutil.XrayReportVulItem
	for _, line := range lines {
		if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
			gologger.Fatal().Msgf("Error parsing %s: %s", line, err)
			continue
		}
		for _, record := range jsonData.Records {
			if !seen[record.Content] {
				seen[record.Content] = true
				report := new(reportutil.XrayReportVulItem)
				report.Target.URL = record.Content
				report.Plugin = record.ID
				report.Detail.Addr = record.Source
				report.Detail.Extra.Param.Position = record.SourceTag
				reports = append(reports, *report)
			}
		}
		jsonData = JsonData{}
	}
	reportutil.GenHtmlReport(reports, filepath.Join(outDir, "wih_unique.html"))
}
