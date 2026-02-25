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