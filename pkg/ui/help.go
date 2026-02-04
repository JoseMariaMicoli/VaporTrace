package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowHelpModal displays the global keybinding reference
func ShowHelpModal(app *tview.Application, pages *tview.Pages) {
	modal := tview.NewTable().
		SetBorders(true).
		SetBordersColor(tcell.ColorDarkCyan).
		SetSelectable(true, false)

	headers := []string{"KEY COMBINATION", "SCOPE", "FUNCTION"}
	for i, h := range headers {
		modal.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorBlack).
			SetBackgroundColor(tcell.ColorDarkCyan).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	}

	data := [][]string{
		{"Ctrl + I", "Global", "Toggle Interceptor (On/Off)"},
		{"Ctrl + F", "Modal", "Forward packet to network"},
		{"Ctrl + D", "Modal", "Drop packet"},
		{"Ctrl + B", "Modal", "Neuro Brute: Gen payloads for current field"},
		{"Ctrl + S", "Modal", "Sync: Save to Loot DB"},
		{"Ctrl + A", "F4 Tab", "Analyze: Send snapshot to AI Brain"},
		{"F1 - F6", "Global", "Switch Tabs (Logs...Neural)"},
		{"F7", "Global", "Report Tab (Read/Edit)"},
		{"Ctrl + W", "F7 Tab", "Save Report to Disk"},
		{"Ctrl + X", "F7 Tab", "Delete Session / Clear Report"},
	}

	for i, row := range data {
		for j, col := range row {
			modal.SetCell(i+1, j, tview.NewTableCell(col).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))
		}
	}

	// Wrapper for centering
	frame := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modal, 14, 1, true).
			AddItem(nil, 0, 1, false), 80, 1, true).
		AddItem(nil, 0, 1, false)

	// Close on Input
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyCtrlH {
			pages.RemovePage("help_modal")
			if cmdInput != nil {
				app.SetFocus(cmdInput)
			}
			return nil
		}
		return event
	})

	pages.AddPage("help_modal", frame, true, true)
	app.SetFocus(modal)
}
