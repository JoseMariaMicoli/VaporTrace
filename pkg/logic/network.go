package logic

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pterm/pterm"
)

// detectedProxy stores the Burp/ZAP address found during auto-detection.
var detectedProxy *url.URL

// InitializeRotaryClient sets up the GlobalClient with a dynamic proxy selector.
// REGLA DE ORO: Utilizes the GlobalClient declared in store.go.
func InitializeRotaryClient() {
	if GlobalClient == nil {
		GlobalClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Dynamic Proxy Selector for Phase 9.4 (Sensing) and Phase 6.2 (IP Rotation)
		Proxy: func(req *http.Request) (*url.URL, error) {
			// Priority 1: Phase 6.2 IP Rotation
			poolProxy := GetRandomProxy()
			if poolProxy != "" {
				return url.Parse(poolProxy)
			}

			// Priority 2: Phase 9.4 Auto-detected Proxy
			if detectedProxy != nil {
				return detectedProxy, nil
			}

			return nil, nil
		},
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	GlobalClient.Transport = transport
}

// DetectAndSetProxy checks for common intercepting proxies via HTTP GET.
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
			pterm.Success.Printfln("Phase 9.4: Auto-detected Proxy at %s", p)
			detectedProxy = proxyURL
			InitializeRotaryClient()
			return
		}
	}

	// Fallback to raw TCP sensing
	for _, p := range proxies {
		u, _ := url.Parse(p)
		conn, err := net.DialTimeout("tcp", u.Host, 300*time.Millisecond)
		if err != nil {
			continue
		}
		conn.Close()

		pterm.Success.Printfln("Phase 9.4: Linked to Proxy (TCP Sense) at %s", p)
		detectedProxy = u
		InitializeRotaryClient()
		return
	}
	
	pterm.Info.Println("No Proxy detected. Running in Direct Mode.")
	InitializeRotaryClient()
}

// SafeDo executes the request with evasion and triggers the loot scanner.
// This is the core ingestion point for the Discovery Vault.
func SafeDo(req *http.Request, isHit bool, module string) (*http.Response, error) {
	ApplyEvasion(req)
	req.Header.Set("X-VaporTrace-Module", module)

	// --- PHASE 8.2: CLOUD PIVOT TRIGGER ---
	// Pre-execution trigger to catch Metadata targets before the request is sent
	TriggerCloudPivot(req.URL.String())
	
	if isHit {
		pterm.Info.Printfln("Mirroring tactical HIT to proxy history [%s]", module)
	}

	// Ensure transport is ready
	if GlobalClient.Transport == nil {
		InitializeRotaryClient()
	}

	resp, err := GlobalClient.Do(req)
	if err != nil {
		return nil, err
	}

	// --- PHASE 8.1: LOOT EXTRACTION HOOK (REGLA DE ORO) ---
	// Non-destructive capture: Read the entire body into memory
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// Even if read fails, ensure we close the body and return the response
		resp.Body.Close()
		return resp, err
	}
	resp.Body.Close()

	// Trigger the asynchronous scanner if data exists
	if len(bodyBytes) > 0 {
		// pterm.Debug.Printfln("Ingesting %d bytes for scanning from %s", len(bodyBytes), req.URL.String())
		go ScanForLoot(string(bodyBytes), req.URL.String())
	}

	// Restore the body for downstream consumers (the probe command logic)
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return resp, nil
}