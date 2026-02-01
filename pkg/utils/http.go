package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GlobalClient is the shared state. It is updated by logic.SetProxy.
var GlobalClient *http.Client

func init() {
	// Initialize with a safe default
	GlobalClient, _ = GetClient("")
}

// GetClient returns a standard HTTP client.
// Used internally by logic package to construct the base transport or by legacy calls.
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
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		} else {
			return nil, err
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}, nil
}
