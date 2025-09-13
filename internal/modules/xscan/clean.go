package xscan

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/netutil"
	"github.com/liancccc/goauto/internal/reportutil"
)

type JsonData struct {
	Desc              string `json:"desc"`
	IsHiddenParameter bool   `json:"isHiddenParameter"`
	Key               string `json:"key"`
	Line              string `json:"line"`
	Payload           string `json:"payload"`
	Position          string `json:"position"`
	Req               string `json:"req"`
	Response          string `json:"response"`
	SuggestPayload    string `json:"suggest-payload"`
	URL               string `json:"url"`
	XSSType           string `json:"xssType"`
}

func Clean(toolOutput, output string) {
	lines := fileutil.ReadingLines(toolOutput)
	var jsonData JsonData
	var report reportutil.XrayReportVulItem
	var reports []reportutil.XrayReportVulItem
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
			continue
		}
		report.CreateTime = time.Now().Unix()
		report.Plugin = jsonData.Desc
		report.Target.URL = netutil.GetBaseUrl(jsonData.URL)
		report.Detail.Addr = genSuggestUrl(jsonData.URL, jsonData.Payload, jsonData.SuggestPayload)
		report.Detail.Payload = jsonData.SuggestPayload
		report.Detail.Snapshot = [][]string{
			{jsonData.Req, jsonData.Response},
		}
		report.Detail.Extra.Param.Key = jsonData.Key
		report.Detail.Extra.Param.Position = jsonData.Position
		reports = append(reports, report)
		jsonData = JsonData{}
		report = reportutil.XrayReportVulItem{}
	}
	reportutil.GenHtmlReport(reports, output)
}

func genSuggestUrl(target, rawParams, suggestParams string) string {
	unescape, err := url.QueryUnescape(target)
	if err != nil {
		return target
	}
	if strings.Contains(unescape, rawParams) {
		return strings.ReplaceAll(unescape, rawParams, suggestParams)
	}
	return unescape
}
