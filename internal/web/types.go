package web

import (
	"github.com/gin-gonic/gin"
	"github.com/liancccc/goauto/internal/db"
	"github.com/liancccc/goauto/internal/sysutil"
)

func Response(success bool, msg string, data interface{}) map[string]any {
	return gin.H{
		"success": success,
		"message": msg,
		"data":    data,
	}
}

type ExecHelp struct {
	Flows []string `json:"flows"`
	Help  string   `json:"help"`
}

type TaskDetail struct {
	*db.DB
	ChildProcess []*sysutil.ProcessItem       `json:"child_process"`
	Reports      map[string][]*TaskReportItem `json:"reports"`
}

type TaskReportItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Link  string `json:"link"`
}
