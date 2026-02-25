package logic

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	utls "github.com/refraction-networking/utls"
)

var detectedProxy *url.URL

// TrafficHistory stores the last known status code for an endpoint
// Used by the Strategic Engine to detect 403s/500s
var (
	TrafficHistory = make(map[string]int)
	trafficMu      sync.RWMutex
)

// GetConfiguredProxy returns the current string representation of the single upstream proxy
// Used by the UI Pipeline Quadrant to display status.
func GetConfiguredProxy() string {
	if detectedProxy != nil {
		return detectedProxy.String()
	}
	return ""
}

// GetTrafficHistory exports a snapshot of traffic for the Engine
func GetTrafficHistory() map[string]int {
	trafficMu.RLock()
	defer trafficMu.RUnlock()
	snapshot := make(map[string]int)
	for k, v := range TrafficHistory {
		snapshot[k] = v
	}
	return snapshot
}

// === SPRINT 11.3 & 12: EVASION TECHNIQUES FOR AUTONOMY ===

// ApplyJitter adds Gaussian-distributed delay to request timing
// Prevents detection of automated tools through traffic analysis
// baseDelay: base milliseconds to add jitter to
func ApplyJitter(baseDelay int) time.Duration {
	// Use Gaussian (normal) distribution with mean=baseDelay, stddev=20%
	mean := float64(baseDelay)
	stddev := mean * 0.2 // 20% variation

	// Box-Muller transform for Gaussian distribution
	u1 := rand.Float64()
	u2 := rand.Float64()
	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	// Apply Gaussian jitter
	jitterValue := mean + (stddev * z)

	// Clamp to minimum of 10ms
	if jitterValue < 10 {
		jitterValue = 10
	}

	return time.Duration(int64(jitterValue)) * time.Millisecond
}

// MimicTraffic sets HTTP headers to mimic real browser/client traffic
// targetProfile: "iOS", "Chrome-MacOS", "Firefox-Windows", "EdgeBrowser", "Safari", "Bot"
func MimicTraffic(req *http.Request, targetProfile string) {
	if req == nil {
		return
	}

	profiles := map[string]map[string]string{
		"iOS": {
			"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate, br",
			"DNT":             "1",
		},
		"Chrome-MacOS": {
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate, br",
			"Sec-Fetch-Site":  "cross-site",
			"Sec-Fetch-Mode":  "cors",
		},
		"Firefox-Windows": {
			"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
			"Accept-Language":           "en-US,en;q=0.9",
			"Accept-Encoding":           "gzip, deflate, br",
			"Upgrade-Insecure-Requests": "1",
		},
		"EdgeBrowser": {
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate, br",
		},
		"Safari": {
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.9",
			"Accept-Encoding": "gzip, deflate, br",
		},
		"Bot": {
			"User-Agent":      "Mozilla/5.0 (compatible; VaporTrace/3.1; Security Testing)",
			"Accept":          "*/*",
			"Accept-Language": "en-US",
			"Accept-Encoding": "gzip, deflate",
		},
	}

	if headers, ok := profiles[targetProfile]; ok {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
}

// --- INTERCEPTOR STATE (HYDRA PROTOCOL) ---
// InterceptorActive is the global toggle for the Blocking Interceptor logic.
var InterceptorActive bool = false

// InterceptorChan acts as the synchronous bridge between the Logic Thread and the UI Thread.
var InterceptorChan = make(chan *InterceptorPayload)

// InterceptorPayload contains the request to be edited and the channel to return the decision.
type InterceptorPayload struct {
	Request          *http.Request
	RequestBodyBytes []byte // EXPLICITLY CARRY BODY TO UI
	ResponseChan     chan *http.Request
}

// TacticalTransport is the middleware that forces all traffic through the suite's logic
type TacticalTransport struct {
	Base http.RoundTripper
}

// RoundTrip executes the interceptor pipeline for EVERY request (Map, Scan, Exploit)
func (t *TacticalTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 0. CAPTURE BODY IMMEDIATELY
	var requestBodyBytes []byte
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Failed to read request body: %v", err))
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		} else {
			requestBodyBytes = bodyBytes
			req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		}
	}

	// EVASION LAYER (Before Enrichment)
	jitterDelay := ApplyJitter(100)
	time.Sleep(jitterDelay)

	// 1. Content Aggregator: Contextual Enrichment
	EnrichCommandRequest(req)
	TriggerCloudPivot(req.URL.String())

	// 2. Tactical Interceptor Hook
	if InterceptorActive {
		respChan := make(chan *http.Request)

		utils.TacticalLog(fmt.Sprintf("[red]INTERCEPT:[-] Pausing %s request to %s for Editor...", req.Method, req.URL.Path))

		InterceptorChan <- &InterceptorPayload{
			Request:          req,
			RequestBodyBytes: requestBodyBytes,
			ResponseChan:     respChan,
		}

		modifiedReq := <-respChan

		if modifiedReq == nil {
			utils.TacticalLog("[red]DROP:[-] Request dropped by operator.")
			return nil, fmt.Errorf("request dropped by operator")
		}

		req = modifiedReq
		if req.Body != nil {
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		utils.TacticalLog("[green]RESUME:[-] Request modified and forwarded to wire.")
	}

	// 3. Capture Request Dump
	reqDump, _ := httputil.DumpRequestOut(req, true)

	// 4. RE-RESTORE BODY ONE FINAL TIME
	if len(requestBodyBytes) > 0 {
		req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	}

	// 5. Execute via Base Transport (The Wire)
	// CAPTURE START TIME
	startTime := time.Now()
	resp, err := t.Base.RoundTrip(req)
	// CAPTURE DURATION
	duration := time.Since(startTime)

	if err != nil {
		return nil, err
	}

	// --- TRAFFIC SENSOR CAPTURE ---
	trafficMu.Lock()
	TrafficHistory[req.URL.Path] = resp.StatusCode
	trafficMu.Unlock()
	// ------------------------------

	// 5. Capture Response Body & Dump
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	resDump, _ := httputil.DumpResponse(resp, true)

	// === CRITICAL FIX: CENTRALIZED VAULT STORAGE ===
	// Store ALL traffic passing through the transport in the GlobalVault
	go func() {
		GlobalVault.Add(req, resp, duration)
	}()

	// 6. Send to Traffic Logger
	reqStr := string(reqDump)
	resStr := string(resDump)

	reqParts := splitDump(reqStr)
	resParts := splitDump(resStr)

	utils.LogTraffic(reqParts[0], reqParts[1], resParts[0], resParts[1])

	// 7. Loot Scanning
	if len(bodyBytes) > 0 {
		go ScanForLoot(string(bodyBytes), req.URL.String())
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return resp, nil
}

// SafeDo executes the request with Context Enrichment, Evasion, Interception, and Traffic Logging.
func SafeDo(req *http.Request, isHit bool, module string) (*http.Response, error) {
	// TASK 1: SCOPE HARD-ENFORCEMENT
	if !GlobalScope.IsInScope(req.URL.String()) {
		utils.TacticalLog(fmt.Sprintf("[yellow]SCOPE:[-] Blocked out-of-scope target: %s", req.URL.String()))
		return nil, fmt.Errorf("Target Out of Scope")
	}

	utils.TacticalLog(fmt.Sprintf("[blue]REQUEST:[-] %s %s (module: %s)", req.Method, req.URL, module))

	// === PRIORITY ALPHA: HTTP/2 PSEUDO-HEADER RANDOMIZATION ===
	if globalStealthConfig.EnablePathObfuscation {
		profile := GetHTTP2Profile(req.Header.Get("User-Agent"))
		ApplyHTTP2Evasion(req, profile)
		utils.TacticalLog(fmt.Sprintf("[cyan]EVASION:[-] Applied HTTP/2 profile: %s", profile.Name))
	}

	// === PRIORITY BETA: PATH OBFUSCATION ===
	if (req.Method == "GET" || req.Method == "POST") && globalStealthConfig.EnablePathObfuscation {
		obfuscationStrategy := SelectObfuscationStrategy()
		originalPath := req.URL.Path
		req.URL.Path = ObfuscatePath(originalPath, obfuscationStrategy)
		utils.TacticalLog(fmt.Sprintf("[cyan]EVASION:[-] Path obfuscation applied: %s → %s", originalPath, req.URL.Path))
	}

	// === PRIORITY EPSILON: PAYLOAD ENCODING ===
	if (req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH") && globalStealthConfig.EnablePayloadEncoding {
		if req.Body != nil {
			bodyBytes, _ := io.ReadAll(req.Body)
			encodingTechnique := SelectRandomEncoding()
			transformedBody, contentEncoding := TransformPayload(bodyBytes, encodingTechnique)

			if contentEncoding != "identity" {
				req.Header.Set("Content-Encoding", contentEncoding)
				utils.TacticalLog(fmt.Sprintf("[cyan]EVASION:[-] Payload encoding applied: %s", contentEncoding))
			}

			// FIX: Update Content-Length to match transformed body size
			newSize := int64(len(transformedBody))
			req.ContentLength = newSize
			req.Header.Set("Content-Length", fmt.Sprintf("%d", newSize))

			req.Body = io.NopCloser(bytes.NewReader(transformedBody))
			if len(transformedBody) != len(bodyBytes) {
				utils.TacticalLog(fmt.Sprintf("[cyan]EVASION:[-] Payload transformed (size: %d → %d bytes)", len(bodyBytes), len(transformedBody)))
			}
		}
	}

	// === PRIORITY GAMMA: CONTEXTUAL THINKING TIME ===
	if globalStealthConfig.EnableThinkingTime {
		delay := ContextualThinkingTime(req.Method, req.URL.Path)
		if delay > 0 {
			utils.TacticalLog(fmt.Sprintf("[cyan]BEHAVIOR:[-] Contextual thinking time: %dms", delay.Milliseconds()))
			time.Sleep(delay)
		}
	}

	ApplyEvasion(req)
	req.Header.Set("X-VaporTrace-Module", module)

	// FIX: Inject User-Agent into Context for JA4 generation
	ua := req.Header.Get("User-Agent")
	if ua != "" {
		ctx := context.WithValue(req.Context(), "ua", ua)
		req = req.WithContext(ctx)
	}

	if GlobalClient.Transport == nil {
		utils.TacticalLog("[yellow]WARN:[-] GlobalClient.Transport is nil, initializing...")
		InitializeRotaryClient()
	}

	resp, err := GlobalClient.Do(req)

	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Request failed: %v", err))
		return resp, err
	}

	utils.TacticalLog(fmt.Sprintf("[green]✓ RESPONSE:[-] %s %d", req.URL, resp.StatusCode))

	// === PRIORITY DELTA: RATE-LIMIT BACKOFF ===
	if resp.StatusCode == 429 || resp.StatusCode == 403 || resp.StatusCode == 503 {
		backoffDelay := HandleRateLimit(resp.StatusCode, resp.Header)
		if backoffDelay > 0 && globalStealthConfig.EnableBackoff {
			utils.TacticalLog(fmt.Sprintf("[red]BACKOFF:[-] Rate-limit triggered. Waiting %.0f seconds before retry...", backoffDelay.Seconds()))
			time.Sleep(backoffDelay)
			utils.TacticalLog("[green]✓ BACKOFF:[-] Cooldown expired. Resuming operations with rotated identity.")
		} else if !globalStealthConfig.EnableBackoff {
			utils.TacticalLog("[yellow]⚠ RATE-LIMIT:[-] 4xx-5xx detected but backoff disabled (stealth mode off)")
		}
	}

	return resp, nil
}

// EnsureTransport guarantees the GlobalClient has the TacticalTransport middleware.
func EnsureTransport() {
	if GlobalClient.Transport == nil {
		InitializeRotaryClient()
	}
}

// InitializeRotaryClient sets up the HTTP client with the Tactical Middleware
// CRITICAL FIX: Segregates TCP (DialContext) from TLS (DialTLSContext)
func InitializeRotaryClient() {
	if GlobalClient == nil {
		GlobalClient = &http.Client{Timeout: 30 * time.Second}
	}

	// Read stealth config safely
	config := GetStealthConfig()
	enableJA4 := config.EnableJA4Fingerprint

	// 1. Standard Network Dialer for TCP (handles HTTP, localhost, and TCP handshake)
	netDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	var baseTransport *http.Transport

	// 3. BRANCHING LOGIC: JA4 vs STANDARD
	if enableJA4 {
		// === JA4 EVASION MODE ===
		baseTransport = &http.Transport{
			DialContext: netDialer.DialContext,
			// DISABLE KEEPALIVES for JA4 to force fresh fingerprint per request
			DisableKeepAlives: true,
			// CRITICAL: Disable Go's built-in HTTP/2 upgrade logic
			ForceAttemptHTTP2: false,
			// CRITICAL: Ensure Go doesn't expect H2 connection
			TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),

			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// 1. Raw TCP
				rawConn, err := netDialer.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}

				host, _, _ := net.SplitHostPort(addr)

				// 2. Force HTTP/1.1 in uTLS to prevent binary frame crash
				uConfig := &utls.Config{
					ServerName:         host,
					InsecureSkipVerify: true,
					NextProtos:         []string{"http/1.1"}, // STRICT HTTP/1.1
				}

				// 3. Get UA from Context
				ua, ok := ctx.Value("ua").(string)
				if !ok || ua == "" {
					ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
				}

				// 4. Fingerprint
				generator := NewJA4Generator(ua)
				utils.TacticalLog(fmt.Sprintf("[magenta]JA4:[-] Negotiating TLS (uTLS/H1) for %s", addr))

				uConn := utls.UClient(rawConn, uConfig, generator.GetClientHelloID())

				// 5. Handshake
				if err := uConn.Handshake(); err != nil {
					rawConn.Close()
					return nil, fmt.Errorf("uTLS handshake failed: %v", err)
				}

				return uConn, nil
			},
			Proxy: func(req *http.Request) (*url.URL, error) {
				poolProxy := GetRandomProxy()
				if poolProxy != "" {
					u, _ := url.Parse(poolProxy)
					return u, nil
				}
				if detectedProxy != nil {
					return detectedProxy, nil
				}
				return nil, nil
			},
		}
		utils.TacticalLog("[green]✓ HTTP CLIENT:[-] Initialized. Mode: [magenta]JA4 Evasion (uTLS/H1)[-]")

	} else {
		// === STANDARD GO TLS MODE (Stealth Off) ===
		baseTransport = &http.Transport{
			DialContext: netDialer.DialContext,
			// Standard TLS (HTTP/2 capable via internal upgrade)
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			// Keep-Alives allowed in standard mode for speed
			DisableKeepAlives: false,
			Proxy: func(req *http.Request) (*url.URL, error) {
				poolProxy := GetRandomProxy()
				if poolProxy != "" {
					u, _ := url.Parse(poolProxy)
					return u, nil
				}
				if detectedProxy != nil {
					return detectedProxy, nil
				}
				return nil, nil
			},
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		}
		utils.TacticalLog("[green]✓ HTTP CLIENT:[-] Initialized. Mode: [cyan]Standard Go TLS[-]")
	}

	tacticalTransport := &TacticalTransport{Base: baseTransport}
	GlobalClient.Transport = tacticalTransport
	utils.GlobalClient = GlobalClient
}

// DetectAndSetProxy attempts to auto-discover local proxies
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
			utils.TacticalLog(fmt.Sprintf("[green]✔[-] Phase 9.4: Auto-detected Proxy at %s", p))
			detectedProxy = proxyURL
			InitializeRotaryClient()
			return
		}
	}
	utils.TacticalLog("[blue]i[-] No Proxy detected. Running in Direct Mode.")
	InitializeRotaryClient()
}

// SetProxy allows manual configuration from CLI commands
func SetProxy(proxyAddr string) {
	if proxyAddr != "" {
		u, err := url.Parse(proxyAddr)
		if err == nil {
			detectedProxy = u
			utils.TacticalLog(fmt.Sprintf("[green]NETWORK:[-] Proxy manually set to %s", proxyAddr))
		}
	} else {
		detectedProxy = nil
		utils.TacticalLog("[blue]NETWORK:[-] Proxy disabled (Direct Mode)")
	}
	InitializeRotaryClient()
}

func splitDump(dump string) []string {
	parts := bytes.SplitN([]byte(dump), []byte("\r\n\r\n"), 2)
	if len(parts) < 2 {
		return []string{string(dump), ""}
	}
	return []string{string(parts[0]), string(parts[1])}
}
