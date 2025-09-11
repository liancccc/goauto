package sysutil

import (
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func GetChidrenByPID(pid int32) ([]*ProcessItem, error) {
	var currentProcess = new(process.Process)
	processes, err := process.Processes()
	for _, item := range processes {
		if item.Pid == pid {
			currentProcess = item
			break
		}
	}
	var childens []*ProcessItem
	children, err := currentProcess.Children()
	if err != nil {
		return nil, err
	}
	var createAt string
	for _, child := range children {
		cmdline, err := child.Cmdline()
		if err != nil || cmdline == "" {
			continue
		}
		createTime, err := child.CreateTime()
		if err == nil {
			createAt = time.UnixMilli(createTime).Local().Format("2006-01-02 15:04:05")
		}
		childens = append(childens, &ProcessItem{
			PID:      child.Pid,
			Command:  cmdline,
			CreateAt: createAt,
		})
		createAt = ""
	}
	return childens, nil
}
