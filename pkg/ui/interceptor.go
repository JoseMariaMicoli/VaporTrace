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
func ShowInterceptorModal(app *tview.Application, pages *tview.Pages, payload *logic.InterceptorPayload) {
	req := payload.Request

	// Read and reset body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	bodyStr := string(bodyBytes)

	// Format Headers
	var headerBuilder strings.Builder
	for k, v := range req.Header {
		headerBuilder.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, ",")))
	}
	headerStr := headerBuilder.String()

	// UI Components
	form := tview.NewForm()
	form.SetBorder(true).
		SetTitle(" [red]TACTICAL INTERCEPTOR (ACTIVE) [white]").
		SetTitleAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorBlack)

	// Editable Fields
	methodField := tview.NewInputField().SetLabel("Method").SetText(req.Method).SetFieldWidth(10)
	urlField := tview.NewInputField().SetLabel("URL").SetText(req.URL.String()).SetFieldWidth(60)

	headersArea := tview.NewTextArea().SetLabel("Headers (Key: Value)").SetText(headerStr, false)
	headersArea.SetBorder(true)

	bodyArea := tview.NewTextArea().SetLabel("Body").SetText(bodyStr, false)
	bodyArea.SetBorder(true)

	// Logic Closures
	forwardFunc := func() {
		// Reconstruct Request
		newReq, err := http.NewRequest(methodField.GetText(), urlField.GetText(), strings.NewReader(bodyArea.GetText()))
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]Interceptor Rebuild Error: %v", err))
			payload.ResponseChan <- req // Fallback to original
		} else {
			// Parse Headers from text area
			lines := strings.Split(headersArea.GetText(), "\n")
			for _, line := range lines {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					newReq.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
			payload.ResponseChan <- newReq
		}
		pages.RemovePage("interceptor")
		app.SetFocus(cmdInput) // defined in dashboard.go
	}

	dropFunc := func() {
		utils.TacticalLog("[red]INTERCEPTOR:[-] Packet Dropped by Operator")
		payload.ResponseChan <- nil // Signal drop
		pages.RemovePage("interceptor")
		app.SetFocus(cmdInput)
	}

	// Button: Inject Session
	injectBtn := tview.NewButton("Inject Active Session").SetSelectedFunc(func() {
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

	forwardBtn := tview.NewButton("FORWARD (Ctrl+F)").SetSelectedFunc(forwardFunc)
	dropBtn := tview.NewButton("DROP (Ctrl+D)").SetSelectedFunc(dropFunc)

	// Layout Construction
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(methodField, 14, 1, false).
			AddItem(urlField, 0, 4, false), 3, 1, true).
		AddItem(headersArea, 0, 3, false).
		AddItem(bodyArea, 0, 4, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(injectBtn, 24, 1, false).
			AddItem(nil, 0, 1, false), 3, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(forwardBtn, 0, 1, false).
			AddItem(dropBtn, 0, 1, false), 3, 1, false)

	// Set Input Capture for Hotkeys
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

	// Display
	pages.AddPage("interceptor", flex, true, true)
	app.SetFocus(flex)
}
