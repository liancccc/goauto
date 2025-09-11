package dict

import (
	"path/filepath"

	"github.com/liancccc/goauto/internal/paths"
)

type Dict struct {
	Name string
	Path string
	Link string
}

var Dicts = map[string]Dict{
	"subdomain-all": {
		Name: "Subdomain-ALL",
		Path: filepath.Join(paths.DictDir, "subdomain-all.txt"),
		Link: "https://gist.githubusercontent.com/jhaddix/f64c97d0863a78454e44c2f7119c2a6a/raw",
	},
}
