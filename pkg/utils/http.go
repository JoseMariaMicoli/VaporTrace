package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

func GetClient(proxyAddr string) (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Essential for Burp
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
		Timeout:   time.Second * 10,
	}, nil
}