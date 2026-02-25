package logic

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	utls "github.com/refraction-networking/utls"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// TLSProfileMap maps profile keys to uTLS ClientHelloID constants
// Each profile represents a realistic browser fingerprint
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

// GetTLSClientHelloID returns the uTLS ClientHelloID for a given profile
// Maps VaporTrace profile names to corresponding uTLS presets
func GetTLSClientHelloID(profileName string) utls.ClientHelloID {
	helloID, exists := TLSProfileMap[profileName]
	if !exists {
		// Default to Chrome Auto if profile not found
		helloID = utls.HelloChrome_Auto
	}
	return helloID
}

// StochasticJitter introduces randomized timing delays before TLS dial
// This evades behavioral analysis and rate-limiting detection
// Delay range: 50ms to 250ms with stochastic distribution
func StochasticJitter() {
	minDelay := 50.0
	maxDelay := 250.0
	// Use exponential distribution for more realistic delay patterns
	delay := minDelay + (maxDelay-minDelay)*rand.Float64()
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// TLSProfileTransport wraps net.Dialer with uTLS for client-side TLS fingerprinting
// This provides realistic browser-like TLS ClientHello matching with uTLS presets
type TLSProfileTransport struct {
	BaseDialer *net.Dialer
	Profile    string
}

// DialTLS creates a TLS connection using uTLS with the specified profile
// Deprecated: Use DialTLSContext instead for proper context handling
func (t *TLSProfileTransport) DialTLS(network, addr string) (net.Conn, error) {
	return t.DialTLSContext(context.Background(), network, addr)
}

// DialTLSContext creates a TLS connection using uTLS with proper SNI and ALPN
// - Extracts host from addr for SNI configuration
// - Forces ALPN negotiation for h2 and http/1.1 to prevent WAF detection
// - Applies stochastic jitter before dial to evade behavioral analysis
// - Applies uTLS preset and performs handshake before returning connection
func (t *TLSProfileTransport) DialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// Introduce stochastic jitter for evasion
	StochasticJitter()

	// Extract host from addr (format: "host:port")
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to split host:port %s: %w", addr, err)
	}

	// Create base TCP connection with timeout
	dialer := &net.Dialer{Timeout: 30 * time.Second}
	if t.BaseDialer != nil {
		dialer = t.BaseDialer
	}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", addr, err)
	}

	// Get uTLS ClientHelloID for the profile
	helloID := GetTLSClientHelloID(t.Profile)

	// Create uTLS connection with SNI and ALPN
	uconn := utls.UClient(conn, &utls.Config{
		ServerName:         host,                       // Proper SNI with target hostname
		InsecureSkipVerify: true,                       // Skip cert verification for penetration testing
		NextProtos:         []string{"h2", "http/1.1"}, // Force ALPN negotiation
	}, helloID)

	// Perform TLS handshake
	err = uconn.Handshake()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("TLS handshake failed: %w", err)
	}

	utils.TacticalLog(fmt.Sprintf("[cyan]TLS:[-] Connected with %s profile to %s:%s", t.Profile, host, port))
	return uconn, nil
}

// SelectOptimalTLSProfile chooses the best TLS profile based on target analysis
// Supports 8 realistic browser profiles with proper rotation and deterministic selection
// Returns consistent profile for same target host to avoid connection inconsistencies
func SelectOptimalTLSProfile(targetHost string) string {
	// 8 diverse browser profiles for maximum evasion effectiveness
	profiles := []string{
		"chrome-windows",  // Chrome on Windows (most common)
		"firefox-windows", // Firefox on Windows
		"safari-macos",    // Safari on macOS
		"chrome-macos",    // Chrome on macOS
		"chromium-linux",  // Chromium on Linux
		"brave-linux",     // Brave on Linux
		"firefox-linux",   // Firefox on Linux
		"edge-windows",    // Edge on Windows
	}

	// Use hash of target to deterministically select profile
	// Same target always gets same profile (for consistency)
	sum := 0
	for _, char := range targetHost {
		sum += int(char)
	}

	selectedProfile := profiles[sum%len(profiles)]
	utils.TacticalLog(fmt.Sprintf("[blue]TLS PROFILE:[-] Selected [%s] for %s", selectedProfile, targetHost))
	return selectedProfile
}

// ApplyTLSEvasion returns a TLSProfileTransport configured for evasion
// This should be called during HTTP client initialization
func ApplyTLSEvasion(profileName string) *TLSProfileTransport {
	utils.TacticalLog("[green]✓ TLS EVASION:[-] uTLS fingerprint spoofing enabled with stochastic jitter")

	if profileName == "" {
		profileName = "chrome-windows"
	}

	return &TLSProfileTransport{
		BaseDialer: &net.Dialer{Timeout: 30 * time.Second},
		Profile:    profileName,
	}
}
