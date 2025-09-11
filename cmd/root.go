package cmd

import (
	"os"

	"github.com/liancccc/goauto/internal/banner"
	"github.com/liancccc/goauto/internal/options"
	"github.com/spf13/cobra"
)

var opt = options.New()
var debug bool

var rootCmd = &cobra.Command{
	Use:   "goauto",
	Short: "",
	Long:  banner.Banner,
	Run: func(cmd *cobra.Command, args []string) {
		banner.Print()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "")
}
