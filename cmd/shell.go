package cmd

import (
	"github.com/spf13/cobra"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/ui"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Launch interactive tactical UI",
	Run: func(cmd *cobra.Command, args []string) {
		shell := ui.NewShell()
		shell.Start()
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}