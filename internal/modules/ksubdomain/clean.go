package ksubdomain

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/projectdiscovery/gologger"
)

func CleanResult(src string, dest string) error {
	gologger.Info().Msgf("Cleaning Ksubdomain result: %s => %s", src, dest)

	if fileutil.CountLines(src) == 0 {
		return errors.New(fmt.Sprintf("no results found in %s", src))
	}

	destFile, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(destFile)

	defer func() {
		writer.Flush()
		destFile.Close()
	}()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	scanner := bufio.NewScanner(srcFile)

	var seen = make(map[string][]string)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=>")
		domain := parts[0]
		hash := strings.Join(parts[1:], ",")
		if len(seen[hash]) < 15 {
			seen[hash] = append(seen[hash], domain)
			writer.WriteString(domain + "\n")
		}
	}
	gologger.Info().Msgf("Clean %s Ksubdomain Complete", src)

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
