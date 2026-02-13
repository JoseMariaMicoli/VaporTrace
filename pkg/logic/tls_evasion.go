package logic

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	utls "github.com/refraction-networking/utls"
)

// TLSProfileMap maps profile keys to uTLS ClientHelloID constants
var TLSProfileMap = map[string]utls.ClientHelloID{
	"chrome-windows":  utls.HelloChrome_Auto,
	"firefox-windows": utls.HelloFirefox_Auto,
	"safari-macos":    utls.HelloSafari_Auto,
	"chrome-macos":    utls.HelloChrome_Auto,
	"chromium-linux":  utls.HelloChrome_Auto,
	"brave-linux":     utls.HelloChrome_Auto,
	"firefox-linux":   utls.HelloFirefox_Auto,
	"edge-windows":    utls.HelloChrome_Auto,
}

// JA4Generator manages browser profiles for TLS handshakes
type JA4Generator struct {
	UserAgent string
}

func NewJA4Generator(ua string) *JA4Generator {
	return &JA4Generator{UserAgent: ua}
}

// GetClientHelloID matches the UA string to a specific uTLS ClientHello ID
// This ensures the TCP/TLS fingerprint matches the HTTP Application layer headers.
func (j *JA4Generator) GetClientHelloID() utls.ClientHelloID {
	ua := strings.ToLower(j.UserAgent)

	if strings.Contains(ua, "chrome") {
		// Return Chrome Auto (Robust for v1.6.7+)
		return utls.HelloChrome_Auto
	} else if strings.Contains(ua, "firefox") {
		// Return Firefox Auto
		return utls.HelloFirefox_Auto
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ios") {
		// Return iOS Auto (Fixes undefined HelloIOS_16)
		return utls.HelloIOS_Auto
	} else if strings.Contains(ua, "android") {
		// Return Android 11
		return utls.HelloAndroid_11_OkHttp
	}

	// Default randomized to avoid Go-default signature
	return utls.HelloRandomized
}

// Dial creates a camouflaged connection
func (j *JA4Generator) Dial(network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: 30 * time.Second}
	rawConn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)

	// Create uTLS Client
	uConn := utls.UClient(rawConn, &utls.Config{
		ServerName:         host,
		InsecureSkipVerify: true, // Standard for Proxy/Offsec tools
	}, j.GetClientHelloID())

	if err := uConn.Handshake(); err != nil {
		rawConn.Close()
		return nil, fmt.Errorf("uTLS handshake failed: %v", err)
	}

	return uConn, nil
}

// GetTLSClientHelloID is kept for backward compatibility
func GetTLSClientHelloID(profileName string) utls.ClientHelloID {
	helloID, exists := TLSProfileMap[profileName]
	if !exists {
		helloID = utls.HelloChrome_Auto
	}
	return helloID
}

// StochasticJitter introduces randomized timing delays before TLS dial
func StochasticJitter() {
	minDelay := 50.0
	maxDelay := 250.0
	delay := minDelay + (maxDelay-minDelay)*rand.Float64()
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// ApplyTLSEvasion returns a TLSProfileTransport (legacy wrapper)
// Note: This logic is partially superseded by InitializeRotaryClient using JA4Generator directly
func ApplyTLSEvasion(profileName string) *JA4Generator {
	// Map profile name to a UA string approximation for the generator
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	if strings.Contains(profileName, "firefox") {
		ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0"
	}
	return NewJA4Generator(ua)
}
