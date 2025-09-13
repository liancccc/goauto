package ksubdomain

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/projectdiscovery/gologger"
)

type JsonData struct {
	Subdomain string `json:"subdomain"`
}

func CleanResult(src string, dest string) error {
	gologger.Info().Msgf("Cleaning Ksubdomain result: %s => %s", src, dest)

	if fileutil.CountLines(src) == 0 {
		return errors.New(fmt.Sprintf("no results found in %s", src))
	}

	var subdomains []string
	var jsonData []JsonData
	fileContent, err := fileutil.GetFileContent(src)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(fileContent), &jsonData); err != nil {
		return err
	}
	for _, item := range jsonData {
		if item.Subdomain != "" {
			subdomains = append(subdomains, item.Subdomain)
		}
	}

	fileutil.WriteSliceToFile(dest, subdomains)
	gologger.Info().Msgf("Clean %s Ksubdomain Complete", src)
	return nil
}
