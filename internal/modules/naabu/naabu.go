package naabu

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liancccc/goauto/internal/dnsresolve"
	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/panjf2000/ants/v2"
	"github.com/projectdiscovery/gologger"
)

func init() {
	modules.RegisterModule(&ModuleStruct{})
}

type ModuleStruct struct {
}

func (m *ModuleStruct) Name() string {
	return "naabu"
}

func (m *ModuleStruct) Install() error {
	var installCmd = "go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest"
	_, err := executil.RunCommandSteamOutput(installCmd)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModuleStruct) CheckInstalled() bool {
	commandSteamOutput, _ := executil.RunCommandSteamOutput("naabu")
	return strings.Contains(commandSteamOutput, "projectdiscovery.io")
}

type naabuHostScanParams struct {
	Host       string
	BaseParams *modules.BaseParams
}

func (m *ModuleStruct) Run(funcParams any) {
	params, ok := funcParams.(modules.BaseParams)
	if !ok {
		gologger.Error().Str("module", m.Name()).Msg("invalid params")
		return
	}
	if params.CustomizeParams == "" {
		params.CustomizeParams = "-top-ports 1000 -Pn -timeout 30 -warm-up-time 5 -scan-type CONNECT"
	}
	_ = params.MkOutDir()

	var processWg sync.WaitGroup
	processWg.Add(1)
	serviceUrls := make(chan string, 100)

	go func() {
		defer processWg.Done()
		outFile, err := os.Create(params.Output)
		if err != nil {
			return
		}
		defer outFile.Close()
		writer := bufio.NewWriter(outFile)
		for serviceUrl := range serviceUrls {
			writer.WriteString(serviceUrl + "\n")
		}
		writer.Flush()
	}()
	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(5, func(i interface{}) {
		defer wg.Done()
		scanParams := i.(naabuHostScanParams)
		hostServices, _ := runNaabuHostScan(&scanParams)
		for _, service := range hostServices {
			serviceUrls <- service
		}
	})
	defer pool.Release()
	var hosts []string
	if params.IsFileTarget() {
		hosts = fileutil.ReadingLines(params.Target)
	} else {
		hosts = []string{params.Target}
	}

	for _, host := range hosts {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}
		wg.Add(1)
		_ = pool.Invoke(naabuHostScanParams{
			Host:       host,
			BaseParams: &params,
		})
	}
	wg.Wait()
	close(serviceUrls)
	processWg.Wait()
	var msg = fmt.Sprintf("Output: %s, Count: %d", params.Output, fileutil.CountLines(params.Output))
	gologger.Info().Str("module", m.Name()).Msg(msg)
}

func runNaabuHostScan(params *naabuHostScanParams) ([]string, error) {
	var services []string
	var toolOutput = filepath.Join(filepath.Dir(params.BaseParams.Output), fmt.Sprintf("%s.xml", dnsresolve.GetIPFormatFileName(params.Host)))
	var command = fmt.Sprintf(`naabu -host %s -nmap-cli "nmap -sV -Pn --open -oX %s"`, params.Host, toolOutput)
	if params.BaseParams.CustomizeParams != "" {
		command = fmt.Sprintf("%s %s", command, params.BaseParams.CustomizeParams)
	}
	_, err := executil.RunCommandSteamOutput(command, params.BaseParams.Timeout)
	if err != nil {
		gologger.Error().Str("module", "naabu").Str("host", params.Host).Msg(err.Error())
		return services, err
	}
	if !fileutil.FileExists(toolOutput) {
		return services, errors.New(fmt.Sprintf("%s not found", toolOutput))
	}
	nmapServices, err := CleanNmapResult(toolOutput)
	if err != nil {
		return services, err
	}
	for _, serviceUrl := range nmapServices {
		services = append(services, serviceUrl)
	}
	fileutil.Remove(toolOutput)
	return services, nil
}
