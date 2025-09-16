package httpx

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/liancccc/goauto/internal/fileutil"
)

type JsonData struct {
	Timestamp time.Time `json:"timestamp"`
	Hash      struct {
		BodySimhash   string `json:"body_simhash"`
		HeaderSimhash string `json:"header_simhash"`
	} `json:"hash"`
	Port          string   `json:"port"`
	URL           string   `json:"url"`
	Input         string   `json:"input"`
	Title         string   `json:"title"`
	Scheme        string   `json:"scheme"`
	Webserver     string   `json:"webserver"`
	ContentType   string   `json:"content_type"`
	Method        string   `json:"method"`
	Host          string   `json:"host"`
	Path          string   `json:"path"`
	Time          string   `json:"time"`
	A             []string `json:"a"`
	Tech          []string `json:"tech"`
	Words         int      `json:"words"`
	Lines         int      `json:"lines"`
	StatusCode    int      `json:"status_code"`
	ContentLength int      `json:"content_length"`
	Failed        bool     `json:"failed"`
	Knowledgebase struct {
		PageType string `json:"PageType"`
		PHash    int    `json:"pHash"`
	} `json:"knowledgebase"`
}

// ParseAndUnique 解析 httpx json 文件根据 body simhash 去重
func ParseAndUnique(jsonOut string) []JsonData {
	var results []JsonData
	var seen = make(map[string]struct{})
	var result JsonData
	var hash string
	var lines = fileutil.ReadingLines(jsonOut)
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}
		results = append(results, result)
		hash = fmt.Sprintf("%s-%s", result.Hash.BodySimhash, result.Host)
		if _, ok := seen[hash]; !ok {
			seen[hash] = struct{}{}
		}
		result = JsonData{}
		hash = ""
	}
	return results
}
