package cmd

import (
	"github.com/flashguru-git/node-monitor/app"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts notification monitor",
	Long:  `Launch monitoring jobs and send emails`,
	Run: func(cmd *cobra.Command, args []string) {
		app.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
