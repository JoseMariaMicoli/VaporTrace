package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pterm/pterm"
)

// GlobalClient is the shared state across ALL modules
var GlobalClient *http.Client

func init() {
	GlobalClient, _ = GetClient("")
}

// GetClient constructs a tactical HTTP client
func GetClient(proxyAddr string) (*http.Client, error) {
	if proxyAddr == "" {
		proxyAddr = os.Getenv("HTTP_PROXY")
	}

	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
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

// UpdateGlobalClient handles network reconfiguration safely for TUI/CLI
func UpdateGlobalClient(proxyAddr string) {
	msg := fmt.Sprintf("Reconfiguring network stack for Interceptor (%s)...", proxyAddr)

	// CLI Mode: Use Spinner
	if UIMode == "CLI" {
		spinner, _ := pterm.DefaultSpinner.Start(msg)
		newClient, err := GetClient(proxyAddr)
		if err != nil {
			spinner.Fail("Network reconfiguration failed: " + err.Error())
			return
		}
		GlobalClient = newClient
		time.Sleep(500 * time.Millisecond)

		if proxyAddr == "" {
			spinner.Success("Traffic is now DIRECT (Proxy Disabled)")
		} else {
			spinner.Success("Traffic routed to INTERCEPTOR: " + proxyAddr)
		}
	} else {
		// TUI Mode: Use Logger Channel (No stdout/spinners!)
		TacticalLog("[blue]⠋[-] " + msg)
		newClient, err := GetClient(proxyAddr)
		if err != nil {
			TacticalLog("[red]✖[-] Network config failed: " + err.Error())
			return
		}
		GlobalClient = newClient
		// Simulate slight delay for user feedback
		time.Sleep(200 * time.Millisecond)

		if proxyAddr == "" {
			TacticalLog("[green]✔[-] Traffic is now DIRECT (Proxy Disabled)")
		} else {
			TacticalLog(fmt.Sprintf("[green]✔[-] Traffic routed to INTERCEPTOR: %s", proxyAddr))
		}
	}
}
