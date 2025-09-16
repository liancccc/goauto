package workflow

import (
	"fmt"
	"path/filepath"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/modules"
	"github.com/liancccc/goauto/internal/modules/merge"
	"github.com/liancccc/goauto/internal/modules/unique"
)

// MergeAndUnique 合并去重
func MergeAndUnique(targets []string, output string) {
	var tempOutput = filepath.Join(filepath.Dir(output), fmt.Sprintf("%s.txt", fileutil.GetUnixNmae()))
	new(merge.ModuleStruct).Run(merge.Params{
		BaseParams: &modules.BaseParams{
			Output: tempOutput,
		},
		Targets: targets,
	})
	new(unique.ModuleStruct).Run(modules.BaseParams{
		Target: tempOutput,
		Output: output,
	})
	fileutil.Remove(tempOutput)
}
