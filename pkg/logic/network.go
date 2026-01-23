package logic

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
	"io"
    "bytes"

	"github.com/pterm/pterm"
)

// detectedProxy stores the Burp/ZAP address found during auto-detection.
var detectedProxy *url.URL

// InitializeRotaryClient sets up the GlobalClient with a dynamic proxy selector.
// This supports Phase 9.4 (Sensing) and Phase 6.2 (IP Rotation) simultaneously.
func InitializeRotaryClient() {
	// Ensure GlobalClient is initialized to avoid nil pointer if this is called early
	if GlobalClient == nil {
		GlobalClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Dynamic Proxy Selector: Evaluated for every single request.
		// This allows us to toggle between Burp and IP Rotation mid-session.
		Proxy: func(req *http.Request) (*url.URL, error) {
			// Priority 1: Phase 6.2 IP Rotation (Proxy Pool)
			// If the user loaded a proxy list, we prioritize rotation for stealth.
			poolProxy := GetRandomProxy()
			if poolProxy != "" {
				return url.Parse(poolProxy)
			}

			// Priority 2: Phase 9.4 Auto-detected Proxy (Burp/ZAP)
			// If no pool exists, we fall back to the intercepted research proxy.
			if detectedProxy != nil {
				return detectedProxy, nil
			}

			// Fallback: Direct Connection
			return nil, nil
		},
		// Performance tuning for industrialized scanning
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
			pterm.Success.Printfln("Phase 9.4: Auto-detected Burp/ZAP Proxy at %s", p)
			detectedProxy = proxyURL
			InitializeRotaryClient() // Re-initialize with the detected proxy
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

func SafeDo(req *http.Request, isHit bool, module string) (*http.Response, error) {
    ApplyEvasion(req)
    req.Header.Set("X-VaporTrace-Module", module)

    // --- PHASE 8.2: CLOUD PIVOT TRIGGER ---
    // Added trigger to detect metadata services before request execution
    TriggerCloudPivot(req.URL.String())
    // --------------------------------------
    
    if isHit {
        pterm.Info.Printfln("Mirroring tactical HIT to proxy history [%s]", module)
    }

    if GlobalClient.Transport == nil {
        InitializeRotaryClient()
    }

    resp, err := GlobalClient.Do(req)
    if err != nil {
        return nil, err
    }

    // --- PHASE 8.1: LOOT EXTRACTION HOOK ---
    // Read the body to scan it
    bodyBytes, _ := io.ReadAll(resp.Body)
    resp.Body.Close()

    // Pass to the Scanner (Implemented in logic/loot.go)
    if len(bodyBytes) > 0 {
        go ScanForLoot(string(bodyBytes), req.URL.String())
    }

    // Restore the body so the rest of the app can use it
    resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
    // ---------------------------------------

    return resp, nil
}