package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/liancccc/goauto/internal/banner"
	"github.com/liancccc/goauto/internal/workflow"
	"github.com/spf13/cobra"
)

var flowsCmd = &cobra.Command{
	Use:   "flows",
	Short: "列出工作流",
	Long:  banner.Banner,
	Run: func(cmd *cobra.Command, args []string) {
		banner.Print()
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetAllowedRowLength(500)
		t.SetTitle("Flow Report")
		t.Style().Title.Align = text.AlignCenter
		t.AppendHeader(table.Row{"Flow", "Description"})
		for _, flow := range workflow.WorkFlows {
			t.AppendRow(table.Row{flow.Name(), flow.Description()})
		}
		fmt.Println()
		t.Render()
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(flowsCmd)
}
