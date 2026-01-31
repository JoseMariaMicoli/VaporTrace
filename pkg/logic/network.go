package logic

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils" // Use unified utils
)

var detectedProxy *url.URL

func InitializeRotaryClient() {
	if GlobalClient == nil {
		GlobalClient = &http.Client{Timeout: 30 * time.Second}
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: func(req *http.Request) (*url.URL, error) {
			poolProxy := GetRandomProxy()
			if poolProxy != "" {
				return url.Parse(poolProxy)
			}
			if detectedProxy != nil {
				return detectedProxy, nil
			}
			return nil, nil
		},
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}
	GlobalClient.Transport = transport
}

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
			// FIXED: Use TacticalLog instead of pterm
			utils.TacticalLog(fmt.Sprintf("[green]✔[-] Phase 9.4: Auto-detected Proxy at %s", p))
			detectedProxy = proxyURL
			InitializeRotaryClient()
			return
		}
	}

	for _, p := range proxies {
		u, _ := url.Parse(p)
		conn, err := net.DialTimeout("tcp", u.Host, 300*time.Millisecond)
		if err != nil {
			continue
		}
		conn.Close()

		utils.TacticalLog(fmt.Sprintf("[green]✔[-] Phase 9.4: Linked to Proxy (TCP Sense) at %s", p))
		detectedProxy = u
		InitializeRotaryClient()
		return
	}

	utils.TacticalLog("[blue]i[-] No Proxy detected. Running in Direct Mode.")
	InitializeRotaryClient()
}

// SafeDo executes the request with evasion and triggers the loot scanner.
func SafeDo(req *http.Request, isHit bool, module string) (*http.Response, error) {
	ApplyEvasion(req)
	req.Header.Set("X-VaporTrace-Module", module)
	TriggerCloudPivot(req.URL.String())

	if isHit {
		utils.TacticalLog(fmt.Sprintf("[magenta]MIRROR[-] Confirmed hit via %s mirrored to proxy.", module))
	}

	if GlobalClient.Transport == nil {
		InitializeRotaryClient()
	}

	resp, err := GlobalClient.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return resp, err
	}
	resp.Body.Close()

	if len(bodyBytes) > 0 {
		go ScanForLoot(string(bodyBytes), req.URL.String())
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return resp, nil
}
