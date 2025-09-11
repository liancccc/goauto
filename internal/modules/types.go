package modules

import (
	"path/filepath"

	"github.com/liancccc/goauto/internal/fileutil"
)

type Module interface {
	Name() string
	Install() error
	CheckInstalled() bool
	Run(funcParams any)
}

type BaseParams struct {
	Target string
	Output string

	Proxy   string
	Dict    string
	Timeout string

	CustomizeParams string
}

func (p BaseParams) IsFileTarget() bool {
	return fileutil.IsFile(p.Target)
}

func (p BaseParams) MkOutDir() error {
	fileutil.Remove(p.Output)
	return fileutil.MakeDir(filepath.Dir(p.Output))
}
