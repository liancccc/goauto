package sysutil

import (
	"fmt"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	info.Hostname = hostname

	hostInfo, err := host.Info()
	if err == nil {
		info.OS = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
	} else {
		info.OS = runtime.GOOS
	}
	info.Arch = runtime.GOARCH
	info.GoVersion = runtime.Version()

	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		info.CPUUsage = cpuPercent[0]
	} else {
		info.CPUUsage = 0
	}

	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.MemoryTotal = float64(memInfo.Total) / (1024 * 1024 * 1024) // 转换为GB
		info.MemoryUsed = float64(memInfo.Used) / (1024 * 1024 * 1024)   // 转换为GB
		info.MemoryUsage = memInfo.UsedPercent
	} else {
		info.MemoryTotal = 0
		info.MemoryUsed = 0
		info.MemoryUsage = 0
	}

	return info, nil
}
