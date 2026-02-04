package ui

import (
	"fmt"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	// reportFlex holds the layout (Editor or Preview)
	reportFlex    *tview.Flex
	reportEditor  *tview.TextArea
	reportPreview *tview.TextView

	// State
	isPreviewMode bool
)

// InitReportTab creates the component for the F7 page.
// Returns a Flex container that switches between Editor and Preview
func InitReportTab() *tview.Flex {
	reportFlex = tview.NewFlex().SetDirection(tview.FlexRow)

	// 1. Initialize Editor
	reportEditor = tview.NewTextArea()
	reportEditor.SetTextStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	reportEditor.SetTitle(" [white:red] REPORT EDITOR (EDIT MODE) [white] - [yellow]Ctrl+W/S: Save[-] | [cyan]Ctrl+P: Preview (Color)[-] | [red]Ctrl+X: Delete[-] ").
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	// 2. Initialize Previewer (Read Only, renders colors)
	reportPreview = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)
	reportPreview.SetTitle(" [white:blue] REPORT PREVIEW (READ ONLY) [white] - [cyan]Ctrl+P: Edit Mode[-] ").
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	// 3. Editor Input Capture
	reportEditor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlW, tcell.KeyCtrlS:
			SaveReport()
			return nil
		case tcell.KeyCtrlX:
			DeleteSession(app, pages, ClearReport)
			return nil
		case tcell.KeyCtrlP:
			ToggleReportMode()
			return nil
		}
		return event
	})

	// 4. Preview Input Capture
	reportPreview.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlP:
			ToggleReportMode()
			return nil
		case tcell.KeyCtrlW, tcell.KeyCtrlS:
			// Allow save from preview too
			SaveReport()
			return nil
		}
		return event
	})

	// Start in Editor Mode
	reportFlex.AddItem(reportEditor, 0, 1, true)
	isPreviewMode = false

	return reportFlex
}

// ToggleReportMode switches the view in the Flex container
func ToggleReportMode() {
	reportFlex.Clear()
	if isPreviewMode {
		// Switch to Edit
		reportFlex.AddItem(reportEditor, 0, 1, true)
		app.SetFocus(reportEditor)
		isPreviewMode = false
	} else {
		// Switch to Preview
		// Sync content: Render markdown colors roughly
		raw := reportEditor.GetText()
		rendered := renderMarkdownToTview(raw)
		reportPreview.SetText(rendered)

		reportFlex.AddItem(reportPreview, 0, 1, true)
		app.SetFocus(reportPreview)
		isPreviewMode = true
	}
}

// Simple heuristic renderer for tview colors
func renderMarkdownToTview(md string) string {
	// Replace headers with Blue/Bold
	md = strings.ReplaceAll(md, "## ", "[blue::b]## ")
	md = strings.ReplaceAll(md, "# ", "[blue::b]# ")
	// Replace bold with white/bold
	md = strings.ReplaceAll(md, "**", "[white::b]")
	// Highlight specific keywords
	md = strings.ReplaceAll(md, "CRITICAL", "[red::b]CRITICAL[white]")
	md = strings.ReplaceAll(md, "HIGH", "[orange::b]HIGH[white]")
	md = strings.ReplaceAll(md, "MEDIUM", "[yellow]MEDIUM[white]")

	// Ensure resets exist at newlines for headers (simplified)
	lines := strings.Split(md, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "[blue::b]") {
			lines[i] = line + "[-]"
		}
	}
	return strings.Join(lines, "\n")
}

// LoadFindings fetches data from DB into the Editor
func LoadFindings() {
	if IsReportDirty {
		return
	}
	markdown := GenerateBaseTemplate()
	if reportEditor != nil {
		reportEditor.SetText(markdown, false)
	}
}

// SaveReport exports the buffer to a file
func SaveReport() {
	if reportEditor == nil {
		return
	}
	// Always save the raw editor text, even if in preview mode
	content := reportEditor.GetText()
	filename, err := SaveReportDisk(content)

	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]REPORT ERROR: Could not save file: %v[-]", err))
	} else {
		utils.TacticalLog(fmt.Sprintf("[green]REPORT SAVED: %s[-]", filename))
	}
}

// ClearReport wipes the buffer
func ClearReport() {
	if reportEditor != nil {
		reportEditor.SetText("", false)
		reportEditor.SetPlaceholder("Session Wiped.")
	}
	if reportPreview != nil {
		reportPreview.SetText("")
	}
	utils.TacticalLog("[yellow]REPORT: Buffer wiped via Session Purge.[-]")
}
