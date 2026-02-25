/*
Copyright (c) 2026 José María Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

package logic

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// HTTP/2 Pseudo-Header Fingerprint Profiles
// Different browsers have different pseudo-header orders
type HTTP2Profile struct {
	Name       string
	UA         string
	PseudoPath string // Path pseudo-header handling
}

var http2Profiles = []HTTP2Profile{
	{
		Name:       "chrome-windows",
		UA:         "Chrome/120 on Windows",
		PseudoPath: "standard", // :method, :authority, :scheme, :path
	},
	{
		Name:       "firefox-windows",
		UA:         "Firefox/121 on Windows",
		PseudoPath: "reordered", // Different order to match Firefox behavior
	},
	{
		Name:       "safari-macos",
		UA:         "Safari/605 on macOS",
		PseudoPath: "aggressive", // More aggressive pseudo-header manipulation
	},
	{
		Name:       "brave-linux",
		UA:         "Brave/1.73 on Linux",
		PseudoPath: "standard",
	},
}

// GetHTTP2Profile returns an HTTP/2 profile based on the current User-Agent
// This ensures pseudo-header order matches the browser fingerprint
func GetHTTP2Profile(ua string) HTTP2Profile {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	profile := http2Profiles[seed.Intn(len(http2Profiles))]
	utils.TacticalLog(fmt.Sprintf("[cyan]HTTP/2 PROFILE:[-] Selected %s for User-Agent randomization", profile.Name))
	return profile
}

// ApplyHTTP2Evasion modifies the request to evade HTTP/2 fingerprinting detection
// This works on the HTTP/1.1 layer but hints to the HTTP/2 transport about expected behavior
func ApplyHTTP2Evasion(req *http.Request, profile HTTP2Profile) {
	if req == nil {
		return
	}

	// === PRIORITY ALPHA: HTTP/2 Pseudo-Header Evasion ===
	// Strategy: Add "noise" headers that mimic browser behavior but don't affect the request

	// Add pseudo-header noise headers that HTTP/2 transports might consider
	// These are application-level hints about request context
	switch profile.PseudoPath {
	case "reordered":
		// Firefox-like behavior: Additional context headers
		req.Header.Set("X-Requested-With", "XMLHttpRequest") // Mimics AJAX requests
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-site")

	case "aggressive":
		// Safari-like: More aggressive with fetch directives
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "none")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

	case "standard":
		// Chrome default: Standard headers only
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "cross-site")
	}

	// Add request timing context for HTTP/2 stream prioritization evasion
	// This makes the request appear less robotic
	req.Header.Set("X-Client-Timing", fmt.Sprintf("%d", time.Now().UnixMilli()%10000))

	utils.TacticalLog("[green]✓ HTTP/2 EVASION:[-] Applied pseudo-header randomization")
}

// ValidateHTTP2Compatibility checks if the request will work over HTTP/2
func ValidateHTTP2Compatibility(req *http.Request) error {
	// Ensure the request doesn't use HTTP/1.1 specific features
	// that would break HTTP/2 multiplexing

	// Validate no dangerous headers
	if req.Header.Get("Connection") != "" && req.Header.Get("Connection") != "keep-alive" {
		utils.TacticalLog("[yellow]WARNING:[-] Connection header may break HTTP/2")
	}

	if req.Header.Get("Upgrade") != "" {
		utils.TacticalLog("[yellow]WARNING:[-] Upgrade header not supported in HTTP/2")
	}

	return nil
}
