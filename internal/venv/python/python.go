package python

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/paths"
)

const VenvName = "python"

type PythonEnv struct {
	VenvPath  string
	PythonBin string
	PipBin    string
	PipxBin   string
}

func New() *PythonEnv {
	p := new(PythonEnv)
	p.VenvPath = filepath.Join(paths.VenvDir, "python")
	if runtime.GOOS == "windows" {
		p.PythonBin = filepath.Join(p.VenvPath, "Scripts", "python.exe")
		p.PipBin = filepath.Join(p.VenvPath, "Scripts", "pip.exe")
		p.PipxBin = filepath.Join(p.VenvPath, "Scripts", "pipx.exe")
	} else {
		p.PythonBin = filepath.Join(p.VenvPath, "bin", "python")
		p.PipBin = filepath.Join(p.VenvPath, "bin", "pip")
		p.PipxBin = filepath.Join(p.VenvPath, "bin", "pipx")
	}
	return p
}

func (p *PythonEnv) Init() error {
	fileutil.MakeDir(p.VenvPath)
	if output, _ := executil.RunCommandSteamOutput("python3 -h"); strings.Contains(output, "usage:") {
		executil.RunCommandSteamOutput(fmt.Sprintf("python3 -m venv %s", p.VenvPath))
	} else {
		executil.RunCommandSteamOutput(fmt.Sprintf("python -m venv %s", p.VenvPath))
	}
	fileutil.Download("https://bootstrap.pypa.io/get-pip.py", "get-pip.py")
	executil.RunCommandSteamOutput(fmt.Sprintf("%s get-pip.py", p.PythonBin))
	fileutil.Remove("get-pip.py")
	executil.RunCommandSteamOutput(fmt.Sprintf("%s install pipx", p.PipBin))
	return nil
}

func (p *PythonEnv) CheckInited() bool {
	command := fmt.Sprintf("%s -h", p.PythonBin)
	if output, _ := executil.RunCommandSteamOutput(command); strings.Contains(output, "usage:") {
		return true
	}
	return false
}
