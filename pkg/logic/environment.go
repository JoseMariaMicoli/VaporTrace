package logic

func SenseEnvironment() {
    // Keep the logic
    if CurrentSession.Threads <= 0 {
        CurrentSession.Threads = 10
    }

    // SILENCE THE OUTPUT
    // Comment out or delete the entire pterm block below:
    /*
    pterm.Info.WithPrefix(pterm.Prefix{Text: "SENSE", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgCyan)}).
        Println("Tactical environment synchronized. Industrialized engines standing by.")
    */
}