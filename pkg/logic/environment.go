package logic

import (
	"github.com/pterm/pterm"
)

// SenseEnvironment performs pre-flight checks to ensure the tactical 
// environment and session defaults are correctly initialized before the shell starts.
func SenseEnvironment() {
	// Ensure threading defaults are set if the session was not fully initialized
	if CurrentSession.Threads <= 0 {
		CurrentSession.Threads = 10
	}

	// Tactical status message for CLI startup
	pterm.Info.WithPrefix(pterm.Prefix{Text: "SENSE", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgCyan)}).
		Println("Tactical environment synchronized. Industrialized engines standing by.")
}