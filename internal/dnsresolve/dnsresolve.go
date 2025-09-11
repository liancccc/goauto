package dnsresolve

import (
	"bufio"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/panjf2000/ants/v2"
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
		resolve, err := resolveDomain(i.(string))
		if err == nil && resolve != nil {
			results <- resolve
		}
		wg.Done()
	})
	defer pool.Release()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
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
