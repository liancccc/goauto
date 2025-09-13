package cmd

import (
	"github.com/liancccc/goauto/internal/banner"
	"github.com/liancccc/goauto/internal/web"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "WEB 模式",
	Long:  banner.Banner,
	Run: func(cmd *cobra.Command, args []string) {
		banner.Print()
		web.StartWebServer(opt)
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.PersistentFlags().StringVarP(&opt.WebServer.Addr, "addr", "", "0.0.0.0:8848", "")
	webCmd.PersistentFlags().StringVarP(&opt.WebServer.User, "user", "", "", "")
	webCmd.PersistentFlags().StringVarP(&opt.WebServer.Pass, "pass", "", "", "")
	webCmd.PersistentFlags().IntVarP(&opt.MaxTaskNum, "max-task", "", 3, "")
}
