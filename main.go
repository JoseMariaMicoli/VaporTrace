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

package main

import (
	"os"
	"github.com/JoseMariaMicoli/VaporTrace/cmd"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/ui"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

func main() {
	// 1. EXECUTION MODE DETECTION
	if len(os.Args) > 1 && os.Args[1] == "shell" {
		utils.SetLoggerMode("CLI") // Standard pterm output
		logic.SenseEnvironment()
		cmd.Execute()
		return
	}

	// 2. DASHBOARD MODE (Default)
	utils.SetLoggerMode("TUI") // Route logs to Dashboard Aggregator
	ui.InitTacticalDashboard()
}