package cmd

import (
	"github.com/liancccc/goauto/internal/banner"
	"github.com/liancccc/goauto/internal/initializer"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "初始化",
	Long:  banner.Banner,
	Run: func(cmd *cobra.Command, args []string) {
		banner.Print()
		if debug {
			gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
		}
		gologger.Info().Msgf("start initializer env ...")
		initializer.VenvInit()
		initializer.ModuleInstall()
		initializer.DictInit()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
