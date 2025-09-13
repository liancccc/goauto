package dnsx

import (
	"encoding/json"
	"strings"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/projectdiscovery/gologger"
)

type JSONData struct {
	Host     string   `json:"host"`
	TTL      int      `json:"ttl"`
	Resolver []string `json:"resolver"`
	A        []string `json:"a"`
	All      []string `json:"all"`
}

// CleanAndGenCustomizeFormat 解析 Dnsx 的输出并保持为自定义的格式
// 目前这样会漏掉 CNAME 和 IPV6 的情况, 但是我的 VPS 好像也扫描不了 IPV6 CNAME 大部分 cndcheck 应该都可以检测到了
// 暂时就只管 IPV4
func CleanAndGenCustomizeFormat(dnsxOut, ipsOut string) (map[string][]string, error) {
	lines := fileutil.ReadingLines(dnsxOut)
	var ipDomains = make(map[string][]string)
	var ips []string
	var jsonData JSONData
	var err error
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if err = json.Unmarshal([]byte(line), &jsonData); err != nil {
			continue
		}
		if len(jsonData.A) == 0 {
			continue
		}
		if _, exists := ipDomains[jsonData.A[0]]; !exists {
			ips = append(ips, jsonData.A[0])
		}
		ipDomains[jsonData.A[0]] = append(ipDomains[jsonData.A[0]], jsonData.Host)
		jsonData = JSONData{}
	}
	gologger.Info().Msgf("Found %d ips", len(ips))
	fileutil.WriteSliceToFile(ipsOut, ips)
	return ipDomains, nil
}
