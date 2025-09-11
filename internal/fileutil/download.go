package fileutil

import (
	"fmt"
	"path/filepath"

	"github.com/liancccc/goauto/internal/executil"
)

// Download 下载
// 系统工具大部分时间都胜过自己写
func Download(link, output string) bool {
	MakeDir(filepath.Dir(output))
	var command = fmt.Sprintf("curl %s -o %s", link, output)
	executil.RunCommandSteamOutput(command)
	return IsFile(output)
}
