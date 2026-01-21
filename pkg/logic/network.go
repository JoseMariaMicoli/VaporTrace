package logic

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pterm/pterm"
)

// GlobalClient is the central HTTP engine for VaporTrace.
// It is declared here to be accessible by bola.go and other logic modules.
/*var GlobalClient = &http.Client{
	Timeout: 10 * time.Second,
}*/

// DetectAndSetProxy checks for common intercepting proxies via HTTP GET.
// This is a "heavy" check to ensure the proxy is actually processing requests.
func DetectAndSetProxy() {
	proxies := []string{"http://127.0.0.1:8080", "http://127.0.0.1:8081"}
	
	for _, p := range proxies {
		proxyURL, _ := url.Parse(p)
		
		// Create a temporary transport to test the proxy
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
	pterm.Info.Println("Phase 9.4: No active proxy detected. Running direct.")
}

// SenseEnvironment (Phase 9.4) is the high-speed version of proxy detection.
// It performs a fast TCP handshake to check for local ports (8080/8081).
func SenseEnvironment() {
	spinner, _ := pterm.DefaultSpinner.Start("Phase 9.4: Sensing Environment...")
	
	proxies := []string{
		"http://127.0.0.1:8080", // Burp Suite Default
		"http://127.0.0.1:8081", // ZAP Default
	}

	for _, p := range proxies {
		u, _ := url.Parse(p)
		
		// Perform a raw TCP dial (300ms timeout)
		conn, err := net.DialTimeout("tcp", u.Host, 300*time.Millisecond)
		if err != nil {
			continue
		}
		conn.Close()

		// If the port is open, link the GlobalClient
		proxyURL, _ := url.Parse(p)
		GlobalClient.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		spinner.Success(fmt.Sprintf("Auto-Proxy: Linked to %s", p))
		return
	}

	spinner.Info("No Proxy detected. Running in Direct Mode.")
}