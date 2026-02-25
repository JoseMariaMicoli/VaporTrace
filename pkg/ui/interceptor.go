package ui

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowInterceptorModal builds and displays the Interceptor UI
// Optimized for strict blocking/resuming of the Logic Thread.
func ShowInterceptorModal(app *tview.Application, pages *tview.Pages, payload *logic.InterceptorPayload) {
	req := payload.Request

	// 1. Safe Body Reading
	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err == nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}
	bodyStr := string(bodyBytes)

	// 2. Normalization
	methodStr := req.Method
	if methodStr == "" {
		methodStr = "GET"
	}
	urlStr := req.URL.String()

	// 3. Format Headers
	var headerBuilder strings.Builder
	for k, v := range req.Header {
		headerBuilder.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, ",")))
	}
	headerStr := headerBuilder.String()

	// --- UI COMPOSITION ---

	// Container Frame
	formFrame := tview.NewFlex().SetDirection(tview.FlexRow)
	formFrame.SetBorder(true).
		SetTitle(" [red::b]TACTICAL INTERCEPTOR [white]").
		SetTitleAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorBlack)

	// --- Input Fields ---

	methodField := tview.NewInputField().
		SetLabel("Method: ").
		SetLabelColor(tcell.ColorAqua).
		SetText(methodStr).
		SetFieldWidth(10).
		SetFieldBackgroundColor(tcell.ColorDarkSlateGray).
		SetFieldTextColor(tcell.ColorWhite)

	urlField := tview.NewInputField().
		SetLabel("URL: ").
		SetLabelColor(tcell.ColorAqua).
		SetText(urlStr).
		SetFieldBackgroundColor(tcell.ColorDarkSlateGray).
		SetFieldTextColor(tcell.ColorYellow)

	headersArea := tview.NewTextArea().
		SetLabel("Headers (Key: Value)").
		SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen)).
		SetPlaceholder("User-Agent: VaporTrace...")
	headersArea.SetText(headerStr, false)
	headersArea.SetBorder(true).SetTitleColor(tcell.ColorAqua)
	headersArea.SetBackgroundColor(tcell.ColorBlack)

	bodyArea := tview.NewTextArea().
		SetLabel("Body (Payload) [Ctrl+B to Brute]").
		SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite))
	bodyArea.SetText(bodyStr, false)
	bodyArea.SetBorder(true).SetTitleColor(tcell.ColorAqua)
	bodyArea.SetBackgroundColor(tcell.ColorBlack)

	// --- LOGIC CLOSURES ---

	closeUI := func() {
		pages.RemovePage("interceptor")
		if cmdInput != nil {
			app.SetFocus(cmdInput)
		}
	}

	forwardFunc := func() {
		newMethod := methodField.GetText()
		newURL := urlField.GetText()
		newBody := bodyArea.GetText()

		// Attempt to build new request
		newReq, err := http.NewRequest(newMethod, newURL, strings.NewReader(newBody))

		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]Interceptor Rebuild Error: %v. Forwarding original.", err))
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			payload.ResponseChan <- req
		} else {
			lines := strings.Split(headersArea.GetText(), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					newReq.Header.Add(key, val)
				}
			}
			payload.ResponseChan <- newReq
		}
		closeUI()
	}

	dropFunc := func() {
		utils.TacticalLog("[red]INTERCEPTOR:[-] Packet Dropped by Operator")
		payload.ResponseChan <- nil
		closeUI()
	}

	// Task 7 & Task 2: Sync to Vault Logic
	syncToVault := func() {
		utils.TacticalLog("[magenta]VAULT:[-] Snapshot synced to Loot Database.")
		// Add simple entry to in-memory Vault for immediate UI display
		logic.Vault = append(logic.Vault, logic.Finding{
			Type:   "INTERCEPT_SNAP",
			Value:  shorten(urlField.GetText(), 40),
			Source: "Operator-Sync",
		})
		utils.LogLoot("SNAP", shorten(urlField.GetText(), 30), "Interceptor")

		// Persist
		utils.RecordFinding(db.Finding{
			Phase:   "PHASE X: MANUAL",
			Command: "intercept",
			Target:  urlField.GetText(),
			Details: "Manual Snapshot synced via Interceptor (Ctrl+S)",
			Status:  "MANUAL_ENTRY",
		})
	}

	// Task 3: Neuro Brute Logic
	neuroBrute := func() {
		currentBody := bodyArea.GetText()
		if currentBody == "" {
			utils.TacticalLog("[yellow]NEURO:[-] Body empty. Cannot bruteforce empty context.")
			return
		}
		utils.TacticalLog("[blue]NEURO:[-] Fuzzing body content with High-Entropy Payloads...")
		go logic.GlobalNeuro.PerformNeuroBrute(currentBody)
	}

	// Task 3: Neuro Inverter Logic
	neuroInvert := func() {
		logic.NeuroInverterActive = !logic.NeuroInverterActive
		status := "OFF"
		if logic.NeuroInverterActive {
			status = "ACTIVE"
		}
		utils.TacticalLog(fmt.Sprintf("[magenta]NEURO-INV:[-] Logic Inversion Mode %s", status))
	}

	// --- BUTTONS (Fixed: Unchained calls) ---

	forwardBtn := tview.NewButton("FORWARD (Ctrl+F)").SetSelectedFunc(forwardFunc)
	forwardBtn.SetBackgroundColor(tcell.ColorGreen)
	forwardBtn.SetLabelColor(tcell.ColorBlack)

	dropBtn := tview.NewButton("DROP (Ctrl+D)").SetSelectedFunc(dropFunc)
	dropBtn.SetBackgroundColor(tcell.ColorDarkRed)
	dropBtn.SetLabelColor(tcell.ColorWhite)

	bruteBtn := tview.NewButton("NEURO BRUTE (Ctrl+B)").SetSelectedFunc(neuroBrute)
	bruteBtn.SetBackgroundColor(tcell.ColorDarkBlue)
	bruteBtn.SetLabelColor(tcell.ColorWhite)

	invertBtn := tview.NewButton("NEURO INV (Ctrl+N)").SetSelectedFunc(neuroInvert)
	invertBtn.SetBackgroundColor(tcell.ColorDarkMagenta)
	invertBtn.SetLabelColor(tcell.ColorWhite)

	syncBtn := tview.NewButton("SYNC VAULT (Ctrl+S)").SetSelectedFunc(syncToVault)
	syncBtn.SetBackgroundColor(tcell.ColorOlive)
	syncBtn.SetLabelColor(tcell.ColorWhite)

	// Layout for Buttons
	btnRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(forwardBtn, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(dropBtn, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(bruteBtn, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(invertBtn, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(syncBtn, 0, 1, false)

	// --- MAIN LAYOUT ---
	topRow := tview.NewFlex().
		AddItem(methodField, 12, 0, false).
		AddItem(urlField, 0, 1, false)

	formFrame.AddItem(topRow, 3, 1, true).
		AddItem(headersArea, 0, 3, false).
		AddItem(bodyArea, 0, 4, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(btnRow, 1, 0, false)

	// --- KEY BINDINGS ---
	formFrame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlF:
			forwardFunc()
			return nil
		case tcell.KeyCtrlD:
			dropFunc()
			return nil
		case tcell.KeyCtrlB:
			neuroBrute()
			return nil
		case tcell.KeyCtrlN:
			neuroInvert()
			return nil
		case tcell.KeyCtrlS:
			syncToVault()
			return nil
		}
		return event
	})

	pages.AddPage("interceptor", formFrame, true, true)
	app.SetFocus(formFrame)
}

func shorten(s string, limit int) string {
	if len(s) > limit {
		return s[:limit] + "..."
	}
	return s
}
