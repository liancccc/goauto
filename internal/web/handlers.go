package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liancccc/goauto/internal/db"
	"github.com/liancccc/goauto/internal/executil"
	"github.com/liancccc/goauto/internal/fileutil"
	"github.com/liancccc/goauto/internal/paths"
	"github.com/liancccc/goauto/internal/sysutil"
	"github.com/liancccc/goauto/internal/workflow"
	"github.com/projectdiscovery/gologger"
)

func execCommandHandler(c *gin.Context) {
	command := c.Query("command")
	if command == "" {
		c.JSON(http.StatusOK, Response(false, "command is nil", nil))
		return
	}
	split := strings.Split(command, " ")
	if split[0] == "goauto" {
		split[0] = os.Args[0]
	}
	command = strings.Join(split, " ")
	go executil.RunCommandSteamOutput(command)
	c.JSON(http.StatusOK, Response(true, "success", nil))
}

func execHelpHandler(c *gin.Context) {
	execHelp := new(ExecHelp)
	for flowNmae := range workflow.WorkFlows {
		execHelp.Flows = append(execHelp.Flows, flowNmae)
	}
	help, _ := executil.RunCommandSteamOutput(fmt.Sprintf("%s scan -h", os.Args[0]))
	if help != "" {
		execHelp.Help = help
	}
	c.JSON(http.StatusOK, Response(true, "success", execHelp))
}

func getALLTaskDBHandler(c *gin.Context) {
	var dbs []*db.DB
	files, err := filepath.Glob(filepath.Join(paths.WorkspaceDir, "*", "DB"))
	if err != nil {
		c.JSON(http.StatusOK, Response(false, err.Error(), nil))
		return
	}
	for _, file := range files {
		taskDB, err := db.Load(file)
		if err != nil {
			gologger.Error().Msgf("Load %s error: %v", file, err)
			continue
		}
		dbs = append(dbs, taskDB)
	}
	c.JSON(http.StatusOK, Response(true, "success", dbs))
}

func getTaskReports(taskName string) (map[string][]*TaskReportItem, error) {
	var workspace = filepath.Join(paths.WorkspaceDir, taskName)
	var reports = make(map[string][]*TaskReportItem)
	var files []string
	var exts = []string{".txt", ".json", ".html", ".md"}
	var blacklist = []string{filepath.Join("httpx", "response"), "index_screenshot.txt", "-ksubdomain-"}
	var dirFunc = func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		for _, noPaht := range blacklist {
			if strings.Contains(path, noPaht) {
				return nil
			}
		}
		if d.IsDir() {
			return nil
		}
		for _, ext := range exts {
			if filepath.Ext(path) == ext {
				files = append(files, path)
			}
		}
		return nil
	}
	if err := filepath.WalkDir(workspace, dirFunc); err != nil {
		return nil, err
	}

	for _, filePath := range files {
		count := fileutil.CountLines(filePath)
		if count == 0 {
			continue
		}
		filePath = strings.Replace(strings.Replace(filePath, workspace, "", -1), "\\", "/", -1)
		filePath = strings.TrimPrefix(filePath, "/")
		split := strings.Split(filePath, "/")
		mode := split[0]
		name := strings.Join(split[1:], "/")
		link, _ := url.JoinPath("workspace", taskName, filePath)
		reports[mode] = append(reports[mode], &TaskReportItem{
			Name:  name,
			Count: count,
			Link:  link,
		})
	}
	return reports, nil
}

func getTaskDetail(c *gin.Context) {
	taskName := c.Query("task")
	if taskName == "" {
		c.JSON(http.StatusOK, Response(false, "task name is nil", nil))
		return
	}
	var workspace = filepath.Join(paths.WorkspaceDir, taskName)
	detail := new(TaskDetail)
	taskDB, err := db.Load(filepath.Join(workspace, "DB"))
	if err != nil {
		c.JSON(http.StatusOK, Response(false, err.Error(), nil))
		return
	}
	detail.DB = taskDB
	processItems, err := sysutil.GetChidrenByPID(int32(detail.DB.PID))
	if err != nil {
		detail.ChildProcess = make([]*sysutil.ProcessItem, 0)
	}
	detail.ChildProcess = processItems

	reports, err := getTaskReports(taskName)
	if err != nil {
		detail.Reports = make(map[string][]*TaskReportItem)
	}
	detail.Reports = reports
	c.JSON(http.StatusOK, Response(true, "success", detail))
}

func getSysInfoHandler(c *gin.Context) {
	systemInfo, err := sysutil.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusOK, Response(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, Response(true, "success", systemInfo))
}

func uploadTargets(c *gin.Context) {
	targets := c.PostForm("targets")
	if targets == "" {
		c.JSON(http.StatusOK, Response(false, "targets is empty", nil))
		return
	}

	fileutil.MakeDir(paths.TargetDir)
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s.txt", timestamp)
	filePath := filepath.Join(paths.TargetDir, fileName)
	if err := os.WriteFile(filePath, []byte(targets), 0644); err != nil {
		c.JSON(http.StatusOK, Response(false, "upload failed", nil))
		return
	}
	c.JSON(http.StatusOK, Response(true, "upload success", gin.H{"filePath": filePath}))
}

func deleteTaskDir(c *gin.Context) {
	taskName := c.Query("task")
	if taskName == "" {
		c.JSON(http.StatusOK, Response(false, "task name is nil", nil))
		return
	}
	if err := fileutil.Remove(filepath.Join(paths.WorkspaceDir, taskName)); err != nil {
		c.JSON(http.StatusOK, Response(false, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, Response(true, "success", nil))
}
