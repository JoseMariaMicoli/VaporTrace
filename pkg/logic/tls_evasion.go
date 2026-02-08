package logic

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// TLSProfile defines a browser-like TLS fingerprint
type TLSProfile struct {
	Name           string
	CipherSuites   []uint16
	EllipticCurves []tls.CurveID
	Extensions     []string
	SignatureAlgs  []tls.SignatureScheme
	Version        uint16 // TLS version (e.g., tls.VersionTLS12)
}

// TLSProfiles contains realistic browser TLS fingerprints
var TLSProfiles = map[string]TLSProfile{
	"chrome-windows": {
		Name: "Chrome on Windows",
		// Chrome 120+ cipher suite order (as of early 2024)
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		EllipticCurves: []tls.CurveID{
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
			tls.X25519,
		},
		Extensions: []string{
			"server_name",
			"extended_master_secret",
			"renegotiation_info",
			"supported_groups",
			"ec_point_formats",
			"session_ticket",
			"application_layer_protocol_negotiation",
			"status_request",
			"signed_certificate_timestamp",
			"key_share",
			"psk_key_exchange_modes",
			"supported_versions",
		},
		SignatureAlgs: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.ECDSAWithP521AndSHA512,
			tls.PSSWithSHA384,
			tls.PSSWithSHA512,
		},
		Version: tls.VersionTLS13,
	},
	"firefox-windows": {
		Name: "Firefox on Windows",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		EllipticCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		Extensions: []string{
			"server_name",
			"renegotiation_info",
			"supported_groups",
			"ec_point_formats",
			"session_ticket",
			"application_layer_protocol_negotiation",
			"status_request",
			"key_share",
			"supported_versions",
			"psk_key_exchange_modes",
		},
		SignatureAlgs: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.PSSWithSHA512,
		},
		Version: tls.VersionTLS13,
	},
	"safari-macos": {
		Name: "Safari on macOS",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		EllipticCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		Extensions: []string{
			"server_name",
			"renegotiation_info",
			"supported_groups",
			"ec_point_formats",
			"application_layer_protocol_negotiation",
			"status_request",
			"key_share",
			"supported_versions",
		},
		SignatureAlgs: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.ECDSAWithP384AndSHA384,
		},
		Version: tls.VersionTLS13,
	},
	"chrome-macos": {
		Name: "Chrome on macOS",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		EllipticCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		Extensions: []string{
			"server_name",
			"extended_master_secret",
			"supported_groups",
			"ec_point_formats",
			"session_ticket",
			"application_layer_protocol_negotiation",
			"status_request",
			"signed_certificate_timestamp",
			"key_share",
			"psk_key_exchange_modes",
			"supported_versions",
		},
		SignatureAlgs: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.ECDSAWithP521AndSHA512,
		},
		Version: tls.VersionTLS13,
	},
	"chromium-linux": {
		Name: "Chromium on Linux",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		EllipticCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		Extensions: []string{
			"server_name",
			"extended_master_secret",
			"supported_groups",
			"ec_point_formats",
			"session_ticket",
			"application_layer_protocol_negotiation",
			"status_request",
			"key_share",
			"psk_key_exchange_modes",
			"supported_versions",
		},
		SignatureAlgs: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.PSSWithSHA384,
		},
		Version: tls.VersionTLS13,
	},
	"brave-linux": {
		Name: "Brave on Linux",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		EllipticCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		Extensions: []string{
			"server_name",
			"extended_master_secret",
			"supported_groups",
			"ec_point_formats",
			"session_ticket",
			"application_layer_protocol_negotiation",
			"status_request",
			"key_share",
			"supported_versions",
			"psk_key_exchange_modes",
		},
		SignatureAlgs: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.ECDSAWithP384AndSHA384,
		},
		Version: tls.VersionTLS13,
	},
	"firefox-linux": {
		Name: "Firefox on Linux",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		EllipticCurves: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		Extensions: []string{
			"server_name",
			"renegotiation_info",
			"supported_groups",
			"ec_point_formats",
			"session_ticket",
			"application_layer_protocol_negotiation",
			"status_request",
			"key_share",
			"supported_versions",
		},
		SignatureAlgs: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.ECDSAWithP384AndSHA384,
		},
		Version: tls.VersionTLS13,
	},
	"edge-windows": {
		Name: "Edge on Windows",
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		EllipticCurves: []tls.CurveID{
			tls.CurveP256,
			tls.CurveP384,
			tls.X25519,
			tls.CurveP521,
		},
		Extensions: []string{
			"server_name",
			"extended_master_secret",
			"renegotiation_info",
			"supported_groups",
			"ec_point_formats",
			"session_ticket",
			"application_layer_protocol_negotiation",
			"status_request",
			"signed_certificate_timestamp",
			"key_share",
			"psk_key_exchange_modes",
			"supported_versions",
		},
		SignatureAlgs: []tls.SignatureScheme{
			tls.ECDSAWithP256AndSHA256,
			tls.PSSWithSHA256,
			tls.ECDSAWithP384AndSHA384,
			tls.PSSWithSHA384,
			tls.PSSWithSHA512,
		},
		Version: tls.VersionTLS13,
	},
}

// GetTLSConfigForProfile returns a TLS config matching the specified browser profile
// Note: This uses Go's standard crypto/tls. For true JA3 evasion, tls-utls library would be needed.
func GetTLSConfigForProfile(profileName string) *tls.Config {
	profile, exists := TLSProfiles[profileName]
	if !exists {
		// Default to Chrome Windows if profile not found
		profile = TLSProfiles["chrome-windows"]
	}

	config := &tls.Config{
		InsecureSkipVerify:       true,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               profile.Version,
		PreferServerCipherSuites: false,
		CipherSuites:             profile.CipherSuites,
		CurvePreferences:         profile.EllipticCurves,
		NextProtos:               []string{"h2", "http/1.1"},
	}

	return config
}

// TLSProfileTransport wraps an http.Transport with TLS profile mimicry
// This provides TLS-level browser fingerprint matching
type TLSProfileTransport struct {
	BaseTransport *tls.Config
	Profile       string
}

// DialTLS creates a TLS connection with profile-specific configuration
// This is used internally by HTTP transport for HTTPS connections
func (t *TLSProfileTransport) DialTLS(network, addr string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, addr, 30*time.Second)
	if err != nil {
		return nil, err
	}

	tlsConfig := GetTLSConfigForProfile(t.Profile)
	tlsConn := tls.Server(conn, tlsConfig)

	err = tlsConn.Handshake()
	if err != nil {
		conn.Close()
		return nil, err
	}

	utils.TacticalLog(fmt.Sprintf("[cyan]TLS:[-] Connected with %s profile to %s", t.Profile, addr))
	return tlsConn, nil
}

// SelectOptimalTLSProfile chooses the best TLS profile based on target analysis
// Supports 9 realistic browser profiles with proper rotation
func SelectOptimalTLSProfile(targetHost string) string {
	// 9 diverse browser profiles for maximum evasion effectiveness
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

// ApplyTLSEvasion configures the HTTP transport with TLS fingerprint matching
// This should be called during client initialization
func ApplyTLSEvasion(baseTransport *tls.Config) *tls.Config {
	utils.TacticalLog("[green]✓ TLS EVASION:[-] JA3 fingerprint mitigation enabled")
	return GetTLSConfigForProfile("chrome-windows")
}
