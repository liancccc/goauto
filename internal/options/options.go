package options

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/liancccc/goauto/internal/fileutil"
)

type Options struct {
	Target     string
	ConfigPath string
	TaskName   string
	Proxy      string

	WorkFlow string
	LogFile  bool
	*WebServer

	MaxTaskNum int
}

func New() *Options {
	o := new(Options)
	o.WebServer = new(WebServer)
	return o
}

type WebServer struct {
	Addr string
	User string
	Pass string
}

func (o *Options) Complete() error {
	var err error
	if err = o.validate(); err != nil {
		return err
	}
	if err = o.setTaskName(); err != nil {
		return err
	}
	if fileutil.IsFile(o.Target) {
		absPath, err := filepath.Abs(o.Target)
		if err != nil {
			return err
		}
		o.Target = absPath
	}
	return nil
}

func (o *Options) validate() error {
	if fileutil.IsFile(o.Target) && o.TaskName == "" {
		return errors.New("file target must have task params")
	}
	return nil
}

func (o *Options) setTaskName() error {
	if o.TaskName != "" {
		return nil
	}
	if strings.Contains(o.Target, "://") {
		parse, err := url.Parse(o.Target)
		if err != nil {
			return err
		}
		o.TaskName = strings.ReplaceAll(parse.Hostname(), ":", "_")
	} else {
		o.TaskName = o.Target
	}
	return nil
}
