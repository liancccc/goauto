package sysutil

type ProcessItem struct {
	PID      int32  `json:"pid"`
	Command  string `json:"command"`
	CreateAt string `json:"create_at"`
}

type SystemInfo struct {
	Hostname    string  `json:"hostname"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	MemoryTotal float64 `json:"memory_total"`
	MemoryUsed  float64 `json:"memory_used"`
	OS          string  `json:"os"`
	Arch        string  `json:"arch"`
	GoVersion   string  `json:"go_version"`
}
