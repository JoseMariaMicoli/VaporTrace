package main

import (
	"os"
	"github.com/JoseMariaMicoli/VaporTrace/cmd"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/ui"
)

func main() {
	// 1. EXECUTION MODE DETECTION (Immediate exit for CLI)
	if len(os.Args) > 1 && os.Args[1] == "shell" {
		logic.SenseEnvironment() // Only run pterm-heavy sense if using the shell
		cmd.Execute()
		return
	}

	// 2. DEFAULT PATH: DASHBOARD
	// We run logic checks AFTER or DURING UI initialization to prevent TTY lock.
	// For now, let's bypass SenseEnvironment for the dashboard test.
	ui.InitTacticalDashboard()
}