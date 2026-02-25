package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
	"os"

	"github.com/pterm/pterm"
)

// GlobalClient is the shared state across ALL modules (Discovery, Logic, UI)
var GlobalClient *http.Client

func init() {
	// Initialize with a direct connection by default
	GlobalClient, _ = GetClient("")
}

// GetClient constructs a tactical HTTP client with TLS-bypass for interception
func GetClient(proxyAddr string) (*http.Client, error) {
	// If no proxy passed, check system environment
	if proxyAddr == "" {
		proxyAddr = os.Getenv("HTTP_PROXY")
	}

	transport := &http.Transport{
		// InsecureSkipVerify is CRITICAL for Burp's self-signed CA certificates
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, 
		// Disabling KeepAlives helps Burp intercept individual requests cleanly
		DisableKeepAlives: true,
	}

	if proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}

// UpdateGlobalClient is the key function called by the UI/CLI to redirect traffic
func UpdateGlobalClient(proxyAddr string) {
	pterm.DefaultSpinner.Start("Reconfiguring network stack for Interceptor...")
	
	newClient, err := GetClient(proxyAddr)
	if err != nil {
		pterm.DefaultSpinner.Fail("Network reconfiguration failed: " + err.Error())
		return
	}

	// ATOMIC UPDATE: This changes the pointer for everyone importing 'utils'
	GlobalClient = newClient
	
	// Tactical delay for UI feedback
	time.Sleep(500 * time.Millisecond)
	
	if proxyAddr == "" {
		pterm.DefaultSpinner.Success("Traffic is now DIRECT (Proxy Disabled)")
	} else {
		pterm.DefaultSpinner.Success("Traffic routed to INTERCEPTOR: " + proxyAddr)
	}
}