package dnsresolve

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/panjf2000/ants/v2"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"github.com/projectdiscovery/gologger"
)

type Resolve struct {
	Domain string
	IP     string
}

// resolveDomain
func resolveDomain(target string) (*Resolve, error) {
	if ip := net.ParseIP(target); ip != nil {
		return &Resolve{IP: target}, nil
	}
	ips, err := net.LookupIP(target)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, errors.New("no ips found")
	}
	for _, ip := range ips {
		if !ip.IsLoopback() {
			return &Resolve{Domain: target, IP: ip.String()}, nil
		}
	}
	return nil, errors.New("no valid ip found")
}

func DoResolve(target, ipDomainDir, ipOutput string) error {
	fileutil.Remove(ipDomainDir)
	fileutil.Remove(ipOutput)
	fileutil.MakeDir(ipDomainDir)
	fileutil.MakeDir(filepath.Dir(ipOutput))
	if fileutil.IsFile(target) {
		return doFileResolve(target, ipDomainDir, ipOutput)
	} else {
		resolve, err := resolveDomain(target)
		if err != nil {
			return err
		}
		fileutil.AppendToContent(filepath.Join(ipDomainDir, GetIPFormatFileName(resolve.IP)), resolve.Domain)
		fileutil.AppendToContent(ipOutput, resolve.IP)
	}
	return nil
}

func GetIPFormatFileName(ip string) string {
	return strings.ReplaceAll(ip, ":", "_")
}

func doFileResolve(targetFile, ipDomainDir, ipOutput string) error {
	file, err := os.Open(targetFile)
	if err != nil {
		return err
	}
	defer file.Close()

	results := make(chan *Resolve, 100)
	var processWg sync.WaitGroup
	processWg.Add(1)
	go func() {
		defer processWg.Done()
		seen := make(map[string]struct{})
		for result := range results {
			gologger.Debug().Msgf("IP: %s Domain: %s\n", result.IP, result.Domain)
			ip := net.ParseIP(result.IP)
			if ip == nil {
				continue
			}
			out := filepath.Join(ipDomainDir, GetIPFormatFileName(result.IP))
			fileutil.AppendToContent(out, result.Domain)
			if _, exists := seen[result.IP]; !exists {
				fileutil.AppendToContent(ipOutput, result.IP)
				seen[result.IP] = struct{}{}
			}
		}
	}()

	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(5, func(i interface{}) {
		dnsClient, err := dnsx.New(dnsx.DefaultOptions)
		if err != nil {
			gologger.Debug().Msgf("Error resolving domain: %s %s\n", err.Error(), i.(string))
			return
		}
		result, err := dnsClient.Lookup(i.(string))
		println(i.(string))
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		println(strings.Join(result, " "))
		results <- &Resolve{
			Domain: i.(string),
			IP:     result[0],
		}
		wg.Done()
	})
	defer pool.Release()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := strings.TrimSpace(scanner.Text())
		if val == "" {
			continue
		}
		wg.Add(1)
		_ = pool.Invoke(val)
	}

	wg.Wait()
	close(results)
	processWg.Wait()
	return nil
}
