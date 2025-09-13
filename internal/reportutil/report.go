package reportutil

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/liancccc/goauto/internal/fileutil"
)

//go:embed xray-template.html
var xrayReportTemplate string

type XrayReportVulItem struct {
	CreateTime int64 `json:"create_time"`
	Detail     struct {
		Addr     string     `json:"addr"`
		Payload  string     `json:"payload"`
		Snapshot [][]string `json:"snapshot"`
		Extra    struct {
			Param struct {
				Key      string `json:"key"`
				Position string `json:"position"`
				Value    string `json:"value"`
			} `json:"param"`
		} `json:"extra"`
	} `json:"detail"`
	Plugin string `json:"plugin"`
	Target struct {
		URL    string `json:"url"`
		Params []struct {
			Position string   `json:"position"`
			Path     []string `json:"path"`
		} `json:"params"`
	} `json:"target"`
}

func GenHtmlReport(items []XrayReportVulItem, output string) {
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString(xrayReportTemplate)
	for _, item := range items {
		jsonData, err := json.Marshal(item)
		if err != nil {
			continue
		}
		htmlBuilder.WriteString("\n")
		htmlBuilder.WriteString(fmt.Sprintf("<script class='web-vulns'>webVulns.push(%s)</script>", string(jsonData)))
	}
	fileutil.WriteToFile(output, htmlBuilder.String())
}
