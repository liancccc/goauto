package db

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/options"
	"github.com/liancccc/goauto/internal/paths"
	"github.com/projectdiscovery/gologger"
)

const (
	RunningStatus = "running"
	DoneStatus    = "done"
	ExitStatus    = "exit"
)

type DB struct {
	TaskName string `json:"task"`
	WorkFlow string `json:"flow"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
	PID      int    `json:"pid"`
	Status   string `json:"status"`
	Path     string `json:"-"`
	Args     string `json:"command"`
}

func Create(opt *options.Options) *DB {
	r := new(DB)
	r.TaskName = opt.TaskName
	r.WorkFlow = opt.WorkFlow
	r.StartAt = time.Now().Format("2006-01-02 15:04:05")
	r.PID = os.Getpid()
	r.Args = strings.Join(os.Args, " ")
	r.Status = RunningStatus
	r.Path = filepath.Join(paths.WorkspaceDir, opt.TaskName, "DB")
	fileutil.MakeDir(filepath.Dir(r.Path))
	r.Save()
	gologger.Info().Msgf("%s DB Create", r.Path)
	return r
}

func (r *DB) Done() {
	r.EndAt = time.Now().Format("2006-01-02 15:04:05")
	r.Status = DoneStatus
	r.Save()
	gologger.Info().Msgf("%s DB Done", r.Path)
}

func (r *DB) Exit() {
	r.EndAt = time.Now().Format("2006-01-02 15:04:05")
	r.Status = ExitStatus
	r.Save()
	gologger.Info().Msgf("%s DB Exit", r.Path)
}

func (r *DB) Save() {
	jsonData, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	file, err := os.OpenFile(r.Path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return
	}
}

func Load(filename string) (*DB, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var d = new(DB)
	err = json.Unmarshal(bytes, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
