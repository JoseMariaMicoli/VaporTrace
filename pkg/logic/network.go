package logic

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

var detectedProxy *url.URL

// GetConfiguredProxy returns the current string representation of the single upstream proxy
// Used by the UI Pipeline Quadrant to display status.
func GetConfiguredProxy() string {
	if detectedProxy != nil {
		return detectedProxy.String()
	}
	return ""
}

// --- INTERCEPTOR STATE ---
var InterceptorActive bool = false
var InterceptorChan = make(chan *InterceptorPayload)

type InterceptorPayload struct {
	Request      *http.Request
	ResponseChan chan *http.Request
}

// TacticalTransport is the middleware that forces all traffic through the suite's logic
type TacticalTransport struct {
	Base http.RoundTripper
}

// RoundTrip executes the interceptor pipeline for EVERY request (Map, Scan, Exploit)
func (t *TacticalTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 1. Content Aggregator: Contextual Enrichment (Phase 10.3)
	// We do this at the transport level to catch Discovery traffic too
	EnrichCommandRequest(req)
	TriggerCloudPivot(req.URL.String())

	// 2. Tactical Interceptor Hook (Phase 10.4)
	if InterceptorActive {
		respChan := make(chan *http.Request)

		// Notify UI and Block
		utils.TacticalLog(fmt.Sprintf("[red]INTERCEPT:[-] Pausing request to %s for F2 Modal...", req.URL.Path))

		InterceptorChan <- &InterceptorPayload{
			Request:      req,
			ResponseChan: respChan,
		}

		// Wait for Operator Action
		modifiedReq := <-respChan
		if modifiedReq == nil {
			utils.TacticalLog("[red]DROP:[-] Request dropped by operator.")
			return nil, fmt.Errorf("request dropped by operator")
		}
		req = modifiedReq
		utils.TacticalLog("[green]RESUME:[-] Request modified and forwarded.")
	}

	// 3. Capture Request Dump (For F4 Upper View)
	reqDump, _ := httputil.DumpRequestOut(req, true)

	// 4. Execute via Base Transport
	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// 5. Capture Response Body & Dump (For F4 Lower View)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	resp.Body.Close() // Close original stream

	// Reconstruct body for downstream use
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	resDump, _ := httputil.DumpResponse(resp, true)

	// 6. Send to Traffic Logger (F4)
	reqStr := string(reqDump)
	resStr := string(resDump)

	reqParts := splitDump(reqStr)
	resParts := splitDump(resStr)

	utils.LogTraffic(reqParts[0], reqParts[1], resParts[0], resParts[1])

	// 7. Loot Scanning (Phase 8)
	if len(bodyBytes) > 0 {
		go ScanForLoot(string(bodyBytes), req.URL.String())
	}

	// Reset body for caller
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return resp, nil
}

// SafeDo executes the request with Context Enrichment, Evasion, Interception, and Traffic Logging.
// RESTORED: Required by bfla.go, bola.go, etc.
func SafeDo(req *http.Request, isHit bool, module string) (*http.Response, error) {
	// 1. Evasion & Headers (Specific to Attack Modules)
	ApplyEvasion(req)
	req.Header.Set("X-VaporTrace-Module", module)

	// 2. Ensure Client is Initialized with Interceptor Transport
	if GlobalClient.Transport == nil {
		InitializeRotaryClient()
	}

	// 3. Execute (Triggers TacticalTransport.RoundTrip)
	return GlobalClient.Do(req)
}

func InitializeRotaryClient() {
	if GlobalClient == nil {
		GlobalClient = &http.Client{Timeout: 30 * time.Second}
	}

	// Base Transport Configuration
	baseTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: func(req *http.Request) (*url.URL, error) {
			// Phase 6.2 Priority:
			// 1. Proxy Pool (Rotation)
			poolProxy := GetRandomProxy()
			if poolProxy != "" {
				u, _ := url.Parse(poolProxy)
				return u, nil
			}
			// 2. Static Proxy (Burp/ZAP)
			if detectedProxy != nil {
				return detectedProxy, nil
			}
			// 3. Direct
			return nil, nil
		},
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}

	// Wrap in Tactical Middleware
	tacticalTransport := &TacticalTransport{Base: baseTransport}
	GlobalClient.Transport = tacticalTransport

	// Sync Utils Client
	utils.GlobalClient = GlobalClient
}

func DetectAndSetProxy() {
	proxies := []string{"http://127.0.0.1:8080", "http://127.0.0.1:8081"}

	for _, p := range proxies {
		proxyURL, _ := url.Parse(p)
		transport := &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: transport, Timeout: 2 * time.Second}

		_, err := client.Get("http://httpbin.org/get")
		if err == nil {
			utils.TacticalLog(fmt.Sprintf("[green]✔[-] Phase 9.4: Auto-detected Proxy at %s", p))
			detectedProxy = proxyURL
			InitializeRotaryClient()
			return
		}
	}
	utils.TacticalLog("[blue]i[-] No Proxy detected. Running in Direct Mode.")
	InitializeRotaryClient()
}

// SetProxy allows manual configuration from CLI commands
func SetProxy(proxyAddr string) {
	if proxyAddr != "" {
		u, err := url.Parse(proxyAddr)
		if err == nil {
			detectedProxy = u
			utils.TacticalLog(fmt.Sprintf("[green]NETWORK:[-] Proxy manually set to %s", proxyAddr))
		}
	} else {
		detectedProxy = nil
		utils.TacticalLog("[blue]NETWORK:[-] Proxy disabled (Direct Mode)")
	}
	InitializeRotaryClient()
}

func splitDump(dump string) []string {
	parts := bytes.SplitN([]byte(dump), []byte("\r\n\r\n"), 2)
	if len(parts) < 2 {
		return []string{string(dump), ""}
	}
	return []string{string(parts[0]), string(parts[1])}
}
