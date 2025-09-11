package httpx_unique

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/projectdiscovery/gologger"
)

type HttpJsonData struct {
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

func CleanHttpxInvalidTargets(hashJsonPath string, validOut string) error {
	gologger.Info().Msgf("Cleaning Httpx Json result: %s", hashJsonPath)

	file, err := os.Open(hashJsonPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var seen = make(map[string]struct{})
	var httpResult HttpJsonData
	var hash string
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		if val == "" {
			continue
		}
		if err := json.Unmarshal([]byte(val), &httpResult); err != nil {
			continue
		}
		hash = fmt.Sprintf("%s-%s", httpResult.Hash.BodySimhash, httpResult.Host)
		if _, ok := seen[hash]; !ok {
			seen[hash] = struct{}{}
			fileutil.AppendToContent(validOut, httpResult.URL)
		}
		httpResult = HttpJsonData{}
		hash = ""
	}
	gologger.Info().Msgf("Found %d valid targets in %s", len(seen), validOut)
	return nil
}
