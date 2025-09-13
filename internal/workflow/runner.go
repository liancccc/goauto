package workflow

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/liancccc/goauto/internal/db"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/options"
	"github.com/liancccc/goauto/internal/paths"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/gologger/writer"
)

type Runner struct {
	opt       *options.Options
	db        *db.DB
	workSpace string
}

func New(opt *options.Options) *Runner {
	w := new(Runner)
	w.opt = opt
	_ = w.opt.Complete()
	_ = w.setWorkSpace()
	_ = w.seLoggerFile()
	w.db = db.Create(opt)

	return w
}

func (runner *Runner) setWorkSpace() error {
	runner.workSpace = filepath.Join(paths.WorkspaceDir, runner.opt.TaskName)
	return fileutil.MakeDir(runner.workSpace)
}

func (runner *Runner) seLoggerFile() error {
	if runner.opt.LogFile {
		gologger.Info().Msgf("Log echo %s", filepath.Join(runner.workSpace, "goauto.log"))
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
		gologger.DefaultLogger.SetFormatter(formatter.NewCLI(true))
		fileWithRotation, _ := writer.NewFileWithRotation(&writer.FileWithRotationOptions{
			Location: runner.workSpace,
			FileName: "goauto.log",
		})
		gologger.DefaultLogger.SetWriter(fileWithRotation)
	}
	return nil
}

func (runner *Runner) Run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	go func() {
		for range c {
			gologger.Info().Msg("Received an interrupt, stopping...")
			runner.db.Exit()
			os.Exit(1)
		}
	}()
	flow, exists := WorkFlows[runner.opt.WorkFlow]
	if !exists {
		return fmt.Errorf("work flow %s does not exist", runner.opt.WorkFlow)
	}
	flow.Run(runner)
	return nil
}

func (runner *Runner) Close() error {
	runner.db.Done()
	return nil
}
