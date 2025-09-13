package fileutil

// osm 摘出来

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/projectdiscovery/gologger"
	"github.com/thoas/go-funk"
)

func WriteTempFile(data string) (string, error) {
	tmpFile, err := os.CreateTemp("", "goauto-*.temp")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	_, err = tmpFile.WriteString(data)
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

// WriteToFile write string to a file
func WriteToFile(filename string, data string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.WriteString(file, data+"\n")
	if err != nil {
		return "", err
	}
	return filename, file.Sync()
}

// GetFileContent Reading file and return content of it
func GetFileContent(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Cleaning(folder string, reports []string) {
	gologger.Info().Msgf("Cleaning result: %v %v", folder, reports)
	// list all the file
	items, err := filepath.Glob(fmt.Sprintf("%v/*", folder))
	if err != nil {
		return
	}

	for _, item := range items {
		gologger.Debug().Msgf("Check Cleaning: %v", item)
		if funk.Contains(reports, item) {
			gologger.Debug().Msgf("Skip cleaning file: %v", item)
			continue
		}
		Remove(item)
	}
}

func ReadingLines(filename string) []string {
	var result []string

	file, err := os.Open(filename)
	if err != nil {
		return result
	}
	defer file.Close()

	// increase the buffer size
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		if val == "" {
			continue
		}
		result = append(result, val)
	}

	if err := scanner.Err(); err != nil {
		return result
	}
	return result
}

// AppendToContent 效率不好 但是方便, 不用在乎那点
func AppendToContent(filename string, data string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	if _, err := writer.WriteString(data + "\n"); err != nil {
		return err
	}
	return writer.Flush()
}

func GetHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func GetCsvColumn(filename string, column int) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return []string{}
	}
	defer file.Close()
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
		return []string{}
	}
	seen := make(map[string]bool)
	var cols []string
	for i := 1; i < len(data); i++ {
		val := strings.TrimSpace(data[i][column-1])
		if val == "" {
			continue
		}
		if seen[val] {
			continue
		}
		seen[val] = true
		cols = append(cols, val)
	}
	return cols
}

func IsFile(src string) bool {
	fi, err := os.Stat(src)
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return false
	case mode.IsRegular():
		if FileLength(src) > 0 {
			return true
		}
		return false
	}
	return false
}

// FileLength count len of file
func FileLength(filename string) int {
	if !FileExists(filename) {
		return 0
	}
	return CountLines(filename)
}

// FileExists check if file is exist or not
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CountLines Return the lines amount of the file
func CountLines(filename string) int {
	var amount int
	file, err := os.Open(filename)
	if err != nil {
		return amount
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		if val == "" {
			continue
		}
		amount++
	}
	if err := scanner.Err(); err != nil {
		return amount
	}
	return amount
}

// MakeDir just make a folder
func MakeDir(folder string) error {
	gologger.Debug().Msgf("Creating folder: %s", folder)
	return os.MkdirAll(folder, 0750)
}

func Move(src string, dest string) error {
	gologger.Debug().Msgf("Moving %s to %s", src, dest)
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	return os.Rename(src, dest)
}

func Remove(src string) error {
	gologger.Debug().Msgf("Remove %s ...", src)
	err := os.RemoveAll(src)
	if err != nil {
		gologger.Error().Msgf("Error removing %s: %s", src, err)
		return err
	}
	return nil
}

func WriteSliceToFile(path string, lines []string) error {
	gologger.Debug().Msgf("Write File to %s", path)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
