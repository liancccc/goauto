package paths

import (
	"path/filepath"

	"github.com/liancccc/goauto/internal/fileutil"
)

var homeDir = fileutil.GetHomeDir()
var BaseDir = filepath.Join(homeDir, "goauto-base")
var VenvDir = filepath.Join(BaseDir, "venv")
var ToolsDir = filepath.Join(BaseDir, "tools")
var DictDir = filepath.Join(BaseDir, "wordlist")
var TargetDir = filepath.Join(BaseDir, "targets")
var WorkspaceDir = filepath.Join(homeDir, "goauto-workspace")

func init() {
	fileutil.MakeDir(BaseDir)
	fileutil.MakeDir(VenvDir)
	fileutil.MakeDir(ToolsDir)
	fileutil.MakeDir(DictDir)
	fileutil.MakeDir(TargetDir)
	fileutil.MakeDir(WorkspaceDir)
}
