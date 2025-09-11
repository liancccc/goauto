package gospider

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/liancccc/goauto/internal/netutil"
)

type JsonData struct {
	Input  string `json:"input"`
	Source string `json:"source"`
	Type   string `json:"type"`
	Output string `json:"output"`
	Status int    `json:"status"`
	Length int    `json:"length"`
}

func Clean(jsonSrc, output string) error {
	srcFile, err := os.Open(jsonSrc)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	outFile, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)

	var jsonData JsonData
	scanner := bufio.NewScanner(srcFile)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		if val == "" {
			continue
		}
		if err := json.Unmarshal([]byte(val), &jsonData); err != nil {
			continue
		}

		if strings.Contains(jsonData.Output, netutil.GetUrlMainDomain(jsonData.Input)) && strings.Contains(jsonData.Output, "http") {
			writer.WriteString(jsonData.Output + "\n")
		}
		jsonData = JsonData{}
	}
	writer.Flush()
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
