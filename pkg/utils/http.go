/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

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

// JoinURL properly joins a base URL with a path, avoiding double slashes.
// It handles cases where base ends with "/" and path starts with "/".
func JoinURL(baseURL, path string) string {
	if baseURL == "" {
		return path
	}
	if path == "" {
		return baseURL
	}

	// Ensure we don't have double slashes
	if baseURL[len(baseURL)-1] == '/' && path[0] == '/' {
		// Both have slashes, remove one
		return baseURL + path[1:]
	} else if baseURL[len(baseURL)-1] != '/' && path[0] != '/' {
		// Neither has a slash, add one
		return baseURL + "/" + path
	}
	// One has a slash, concatenate as is
	return baseURL + path
}
