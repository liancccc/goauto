package processutil

import (
	"fmt"
	"runtime"

	"github.com/liancccc/goauto/internal/executil"
)

func Kill(pid int) error {
	var command string
	if runtime.GOOS == "windows" {
		command = fmt.Sprintf("taskkill /F /PID %d ", pid)
	} else {
		command = fmt.Sprintf("kill -9 %d ", pid)
	}
	_, err := executil.RunCommandSteamOutput(command)
	if err != nil {
		return err
	}
	return nil
}
