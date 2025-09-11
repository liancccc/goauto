package fileutil

import (
	"fmt"
	"runtime"

	"github.com/liancccc/goauto/internal/executil"
)

func Unzip(src, outDir string) {
	var command string
	MakeDir(outDir)
	if runtime.GOOS == "windows" {
		command = fmt.Sprintf(`Expand-Archive -Path "%s" -DestinationPath "%s"`, src, outDir)
	} else {
		command = fmt.Sprintf(`unzip %s -d %s`, src, outDir)
	}
	executil.RunCommandSteamOutput(command)
}
