package logic

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

var detectedProxy *url.URL

// TrafficHistory stores the last known status code for an endpoint
// Used by the Strategic Engine to detect 403s/500s
var (
	TrafficHistory = make(map[string]int)
	trafficMu      sync.RWMutex
)

// GetConfiguredProxy returns the current string representation of the single upstream proxy
// Used by the UI Pipeline Quadrant to display status.
func GetConfiguredProxy() string {
	if detectedProxy != nil {
		return detectedProxy.String()
	}
	return ""
}

// GetTrafficHistory exports a snapshot of traffic for the Engine
// This fixes: undefined: logic.GetTrafficHistory in engine/core.go
func GetTrafficHistory() map[string]int {
	trafficMu.RLock()
	defer trafficMu.RUnlock()
	snapshot := make(map[string]int)
	for k, v := range TrafficHistory {
		snapshot[k] = v
	}
	return snapshot
}

// --- INTERCEPTOR STATE (HYDRA PROTOCOL) ---
// InterceptorActive is the global toggle for the Blocking Interceptor logic.
var InterceptorActive bool = false

// InterceptorChan acts as the synchronous bridge between the Logic Thread and the UI Thread.
var InterceptorChan = make(chan *InterceptorPayload)

// InterceptorPayload contains the request to be edited and the channel to return the decision.
// === SPRINT 11 FIX: Now carries RequestBodyBytes to ensure body is available to UI ===
type InterceptorPayload struct {
	Request          *http.Request
	RequestBodyBytes []byte // EXPLICITLY CARRY BODY TO UI
	ResponseChan     chan *http.Request
}

// TacticalTransport is the middleware that forces all traffic through the suite's logic
type TacticalTransport struct {
	Base http.RoundTripper
}

// RoundTrip executes the interceptor pipeline for EVERY request (Map, Scan, Exploit)
func (t *TacticalTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// === SPRINT 11: CRITICAL FIX ===
	// 0. CAPTURE BODY IMMEDIATELY (Before any sensor consumes it)
	// This ensures the body is available to all downstream processors
	var requestBodyBytes []byte
	if req.Body != nil {
		var err error
		requestBodyBytes, err = io.ReadAll(req.Body)
		if err == nil {
			// RE-STUFF the body immediately so it can be read by sensors
			req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		}
	}

	// 1. Content Aggregator: Contextual Enrichment (Phase 10.3)
	EnrichCommandRequest(req)
	TriggerCloudPivot(req.URL.String())

	// 2. Tactical Interceptor Hook (Phase 10.4 - Hydra Audit Fix)
	if InterceptorActive {
		respChan := make(chan *http.Request)

		utils.TacticalLog(fmt.Sprintf("[red]INTERCEPT:[-] Pausing %s request to %s for Editor...", req.Method, req.URL.Path))

		// === SPRINT 11 FIX: Pass captured body bytes to UI ===
		InterceptorChan <- &InterceptorPayload{
			Request:          req,
			RequestBodyBytes: requestBodyBytes, // PASS BODY EXPLICITLY
			ResponseChan:     respChan,
		}

		modifiedReq := <-respChan

		if modifiedReq == nil {
			utils.TacticalLog("[red]DROP:[-] Request dropped by operator.")
			return nil, fmt.Errorf("request dropped by operator")
		}

		req = modifiedReq
		// === RE-RESTORE BODY AFTER INTERCEPTOR MODAL ===
		// Ensure body is preserved through the modal interaction
		if req.Body != nil {
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		utils.TacticalLog("[green]RESUME:[-] Request modified and forwarded to wire.")
	}

	// 3. Capture Request Dump (Body is now guaranteed to be available)
	reqDump, _ := httputil.DumpRequestOut(req, true)

	// 4. RE-RESTORE BODY ONE FINAL TIME before wire transmission
	if len(requestBodyBytes) > 0 {
		req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	}

	// 5. Execute via Base Transport (The Wire)
	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// --- TRAFFIC SENSOR CAPTURE ---
	trafficMu.Lock()
	TrafficHistory[req.URL.Path] = resp.StatusCode
	trafficMu.Unlock()
	// ------------------------------

	// 5. Capture Response Body & Dump
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	resDump, _ := httputil.DumpResponse(resp, true)

	// 6. Send to Traffic Logger
	reqStr := string(reqDump)
	resStr := string(resDump)

	reqParts := splitDump(reqStr)
	resParts := splitDump(resStr)

	utils.LogTraffic(reqParts[0], reqParts[1], resParts[0], resParts[1])

	// 7. Loot Scanning
	if len(bodyBytes) > 0 {
		go ScanForLoot(string(bodyBytes), req.URL.String())
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return resp, nil
}

// SafeDo executes the request with Context Enrichment, Evasion, Interception, and Traffic Logging.
func SafeDo(req *http.Request, isHit bool, module string) (*http.Response, error) {
	ApplyEvasion(req)
	req.Header.Set("X-VaporTrace-Module", module)

	if GlobalClient.Transport == nil {
		InitializeRotaryClient()
	}

	return GlobalClient.Do(req)
}

// InitializeRotaryClient sets up the HTTP client with the Tactical Middleware
func InitializeRotaryClient() {
	if GlobalClient == nil {
		GlobalClient = &http.Client{Timeout: 30 * time.Second}
	}

	baseTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: func(req *http.Request) (*url.URL, error) {
			poolProxy := GetRandomProxy()
			if poolProxy != "" {
				u, _ := url.Parse(poolProxy)
				return u, nil
			}
			if detectedProxy != nil {
				return detectedProxy, nil
			}
			return nil, nil
		},
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}

	tacticalTransport := &TacticalTransport{Base: baseTransport}
	GlobalClient.Transport = tacticalTransport
	utils.GlobalClient = GlobalClient
}

// DetectAndSetProxy attempts to auto-discover local proxies
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
