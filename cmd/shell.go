/*
Copyright (c) 2026 José María Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

package cmd

import (
	"github.com/JoseMariaMicoli/VaporTrace/pkg/ui"
	"github.com/spf13/cobra"
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
