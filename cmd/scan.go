package cmd

import (
	"fmt"
	"strings"

	"github.com/liancccc/goauto/internal/banner"
	"github.com/liancccc/goauto/internal/workflow"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "扫描模式",
	Long:  banner.Banner,
	Run: func(cmd *cobra.Command, args []string) {
		banner.Print()
		if debug {
			gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
		}
		flow := workflow.New(opt)
		defer flow.Close()
		flow.Run()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.PersistentFlags().StringVarP(&opt.Target, "target", "", "", "目标 ( 域名 | IP | URL )")
	scanCmd.PersistentFlags().StringVarP(&opt.TaskName, "task-name", "", "", "任务名称 ( 目标是文件的时候必须 )")
	scanCmd.PersistentFlags().StringVarP(&opt.WorkFlow, "flow", "", "", fmt.Sprintf("工作流名称 ( %s )", strings.Join(workflow.GetFlowNames(), " | ")))
	scanCmd.PersistentFlags().StringVarP(&opt.Proxy, "proxy", "", "", "代理 ( 给支持代理的工具添加代理 )")
	scanCmd.PersistentFlags().BoolVarP(&opt.LogFile, "log-file", "", false, "")
}
