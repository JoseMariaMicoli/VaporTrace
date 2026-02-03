package ui

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowInterceptorModal builds and displays the F2 Interceptor UI
// Optimized for strict blocking/resuming of the Logic Thread with visible styling.
func ShowInterceptorModal(app *tview.Application, pages *tview.Pages, payload *logic.InterceptorPayload) {
	req := payload.Request

	// 1. Safe Body Reading
	// We read the stream, cache it, and immediately restore it to the request so logic flows aren't broken.
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
	// Go http.Request.Method can be empty for GET. We must be explicit for the UI.
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
		SetTitle(" [red::b]TACTICAL INTERCEPTOR (ACTIVE) [white]").
		SetTitleAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorBlack)

	// --- Input Fields (Styled) ---

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
		SetLabel("Body (Payload)").
		SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite))
	bodyArea.SetText(bodyStr, false)
	bodyArea.SetBorder(true).SetTitleColor(tcell.ColorAqua)
	bodyArea.SetBackgroundColor(tcell.ColorBlack)

	// --- LOGIC CLOSURES ---

	closeUI := func() {
		pages.RemovePage("interceptor")
		// Safely attempt to refocus the command input if accessible
		// Note: cmdInput is global in dashboard.go; ensuring focus returns there is UX critical.
		if cmdInput != nil {
			app.SetFocus(cmdInput)
		}
	}

	forwardFunc := func() {
		// Reconstruct Request
		newMethod := methodField.GetText()
		newURL := urlField.GetText()
		newBody := bodyArea.GetText()

		// Attempt to build new request
		newReq, err := http.NewRequest(newMethod, newURL, strings.NewReader(newBody))

		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]Interceptor Rebuild Error: %v. Forwarding original.", err))
			// Recover: Send original req back to unblock thread
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Ensure body is reset again
			payload.ResponseChan <- req
		} else {
			// Parse Headers from Text Area
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
			// Unblock Logic Thread
			payload.ResponseChan <- newReq
		}
		closeUI()
	}

	dropFunc := func() {
		utils.TacticalLog("[red]INTERCEPTOR:[-] Packet Dropped by Operator")
		payload.ResponseChan <- nil // Signal drop to logic layer
		closeUI()
	}

	// --- BUTTONS (Fixed: Unchained calls to avoid interface panic) ---

	injectBtn := tview.NewButton("Inject Active Session (Auth)").SetSelectedFunc(func() {
		token := logic.CurrentSession.AttackerToken
		if token != "" {
			current := headersArea.GetText()
			if current != "" && !strings.HasSuffix(current, "\n") {
				current += "\n"
			}
			newHeaders := current + fmt.Sprintf("Authorization: Bearer %s", token)
			headersArea.SetText(newHeaders, false)
		}
	})
	injectBtn.SetBackgroundColor(tcell.ColorDarkBlue)
	injectBtn.SetLabelColor(tcell.ColorWhite)

	forwardBtn := tview.NewButton("FORWARD (Ctrl+F)").SetSelectedFunc(forwardFunc)
	forwardBtn.SetBackgroundColor(tcell.ColorGreen)
	forwardBtn.SetLabelColor(tcell.ColorBlack)

	dropBtn := tview.NewButton("DROP (Ctrl+D)").SetSelectedFunc(dropFunc)
	dropBtn.SetBackgroundColor(tcell.ColorDarkRed)
	dropBtn.SetLabelColor(tcell.ColorWhite)

	// Layout for Buttons
	btnRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(forwardBtn, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false). // Spacer
		AddItem(dropBtn, 0, 1, false)

	// --- MAIN LAYOUT ---

	// Row 1: Method (Fixed 12) + URL (Flex)
	topRow := tview.NewFlex().
		AddItem(methodField, 12, 0, false).
		AddItem(urlField, 0, 1, false)

	formFrame.AddItem(topRow, 3, 1, true).
		AddItem(headersArea, 0, 3, false).
		AddItem(bodyArea, 0, 4, false).
		AddItem(injectBtn, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false). // Vertical Spacer
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
		}
		return event
	})

	// Mount and Force Focus
	pages.AddPage("interceptor", formFrame, true, true)
	app.SetFocus(formFrame)
}
