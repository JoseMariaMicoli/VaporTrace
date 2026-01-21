package logic

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pterm/pterm"
)

// GlobalClient is now managed by the central store to avoid redeclaration errors.
// This file focuses on the networking transport and the SafeDo gatekeeper.

// DetectAndSetProxy checks for common intercepting proxies via HTTP GET.
func DetectAndSetProxy() {
	proxies := []string{"http://127.0.0.1:8080", "http://127.0.0.1:8081"}
	
	for _, p := range proxies {
		proxyURL, _ := url.Parse(p)
		
		// Create a temporary transport to test the proxy functional capacity
		transport := &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: transport, Timeout: 2 * time.Second}
		
		// Try to reach a neutral endpoint through the proxy
		_, err := client.Get("http://httpbin.org/get")
		if err == nil {
			pterm.Success.Printfln("Phase 9.4: Auto-detected Burp/ZAP Proxy at %s", p)
			
			// Update the GlobalClient to use this functional proxy
			GlobalClient.Transport = transport
			return
		}
	}

	// Fallback to raw TCP sensing if the HTTP check fails
	for _, p := range proxies {
		u, _ := url.Parse(p)
		conn, err := net.DialTimeout("tcp", u.Host, 300*time.Millisecond)
		if err != nil {
			continue
		}
		conn.Close()

		GlobalClient.Transport = &http.Transport{
			Proxy:           http.ProxyURL(u),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		pterm.Success.Printfln("Phase 9.4: Linked to Proxy (TCP Sense) at %s", p)
		return
	}
	pterm.Info.Println("No Proxy detected. Running in Direct Mode.")
}

// SafeDo is the tactical gatekeeper for Phase 9.6.
// It executes the request and mirrors confirmed hits to the proxy for researcher visibility.
func SafeDo(req *http.Request, isHit bool, module string) (*http.Response, error) {
	// Add tracking headers for proxy history and IR triage
	req.Header.Set("X-VaporTrace-Module", module)
	
	if isHit {
		pterm.Info.Printfln("Mirroring tactical HIT to proxy history [%s]", module)
	}

	return GlobalClient.Do(req)
}