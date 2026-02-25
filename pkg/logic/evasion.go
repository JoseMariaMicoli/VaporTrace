package logic

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// === STEALTH CONTROLLER: DYNAMIC EVASION CONFIGURATION ===

// StealthLevel defines the current evasion configuration
type StealthLevel struct {
	mu                      sync.RWMutex
	EnableJitter            bool
	EnableThinkingTime      bool
	EnableBackoff           bool
	EnablePathObfuscation   bool
	EnablePayloadEncoding   bool
	GlobalEvasionMultiplier float64
	Mode                    string // "Aggressive", "Fast", "Silent", "Custom"
}

// globalStealthConfig holds the global evasion configuration
var globalStealthConfig = &StealthLevel{
	EnableJitter:            true,
	EnableThinkingTime:      true,
	EnableBackoff:           true,
	EnablePathObfuscation:   true,
	EnablePayloadEncoding:   true,
	GlobalEvasionMultiplier: 1.0,
	Mode:                    "Aggressive",
}

// ProxyPool holds a list of proxy URLs (SOCKS5/Tor/HTTP) for IP rotation
var ProxyPool []string

// LoadProxiesFromFile reads a line-separated file of proxy URLs and populates the pool.
func LoadProxiesFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var newPool []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			// Ensure protocol is present; default to http if missing
			if !strings.HasPrefix(line, "http") && !strings.HasPrefix(line, "socks5") {
				line = "http://" + line
			}
			newPool = append(newPool, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	ProxyPool = newPool
	// Logging is now handled by the Engine calling this function
	return nil
}

// GetRandomProxy returns a proxy from the pool or an empty string if none are configured
func GetRandomProxy() string {
	if len(ProxyPool) == 0 {
		return ""
	}
	// Create new random source each time (thread-safe)
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ProxyPool[seed.Intn(len(ProxyPool))]
}

// Tactical Fingerprints: 20+ diverse browser-accurate User-Agents across Windows, macOS, Linux, iOS, Android, Brave, Chromium
var userAgents = []string{
	// WINDOWS (5)
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",

	// macOS (5)
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",

	// LINUX (3)
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",

	// iOS - SAFARI & MOBILE BROWSERS (5)
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.7 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/120.0.6099.129 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Braveios/1.64 Mobile/15E148 Safari/605.1",

	// ANDROID - CHROME & BRAVE MOBILE (5)
	"Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 14; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 14) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 Brave/1.73",

	// BRAVE BROWSER - DESKTOP (2)
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Brave/1.73",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Brave/1.73",
}

// ApplyEvasion handles the tactical transformation of a request before it is sent
func ApplyEvasion(req *http.Request) {
	// Create new random source each time (thread-safe, no global state mutation)
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 1. Header Randomization (Phase 6.1) - Only apply if not stealth off
	// Even in stealth off mode, we apply basic User-Agent rotation for minimal anonymity
	if len(userAgents) > 0 {
		ua := userAgents[seed.Intn(len(userAgents))]
		req.Header.Set("User-Agent", ua)
		// Only add evasion headers if stealth is NOT completely off
		if globalStealthConfig.EnableJitter || globalStealthConfig.EnableThinkingTime || globalStealthConfig.EnablePathObfuscation || globalStealthConfig.EnablePayloadEncoding {
			req.Header.Set("Accept", "application/json, text/plain, */*")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9")
			req.Header.Set("Cache-Control", "no-cache")
			req.Header.Set("Connection", "keep-alive")
		}
	}

	// 2. Timing Attacks: Gaussian Stochastic Jitter (Phase 6.3)
	// Use Gaussian distribution to evade rate-limit detection
	// Now respects EnableJitter toggle and context
	if globalStealthConfig.EnableJitter {
		jitterDelay := StochasticJitterMS(20, 150) // 20-150ms with Gaussian distribution
		if jitterDelay > 0 {
			ctx := context.Background()
			SafeSleep(ctx, time.Duration(jitterDelay)*time.Millisecond, &globalStealthConfig.EnableJitter)
			utils.TacticalLog("[blue]EVASION:[-] Applied stochastic jitter")
		}
	}
}

// StochasticJitterMS returns a Gaussian-distributed random millisecond delay
// min: minimum milliseconds
// max: maximum milliseconds
// Uses Box-Muller transform for proper Gaussian distribution
func StochasticJitterMS(min, max int) int {
	if min >= max {
		return min
	}

	// Calculate mean and standard deviation for Gaussian distribution
	mean := float64(min+max) / 2.0
	stddev := float64(max-min) / 4.0 // Approx 95% of values within [min, max]

	// Box-Muller transform for Gaussian distribution (thread-safe)
	u1 := rand.Float64()
	u2 := rand.Float64()
	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	// Apply Gaussian jitter
	jitterValue := mean + (stddev * z)

	// Clamp to [min, max] range
	if jitterValue < float64(min) {
		jitterValue = float64(min)
	}
	if jitterValue > float64(max) {
		jitterValue = float64(max)
	}

	return int(jitterValue)
}

// === INTERRUPTIBLE SLEEP: Context-Aware & Toggleable ===

// SafeSleep performs a context-aware, interruptible sleep that respects stealth toggles
// - Returns immediately if toggle is false
// - Respects context cancellation (for UI responsiveness)
// - Applies global evasion multiplier
// Returns true if sleep completed, false if cancelled or skipped
func SafeSleep(ctx context.Context, duration time.Duration, toggle *bool) bool {
	globalStealthConfig.mu.RLock()
	multiplier := globalStealthConfig.GlobalEvasionMultiplier
	globalStealthConfig.mu.RUnlock()

	// If toggle is false, skip sleep entirely
	if toggle != nil && !*toggle {
		utils.TacticalLog("[yellow]STEALTH:[-] Sleep skipped (toggle disabled)")
		return false
	}

	// Apply global evasion multiplier
	adjustedDuration := time.Duration(float64(duration) * multiplier)

	// Use select to check context cancellation
	select {
	case <-time.After(adjustedDuration):
		// Sleep completed successfully
		return true
	case <-ctx.Done():
		// Context cancelled mid-sleep
		utils.TacticalLog("[yellow]STEALTH:[-] Sleep interrupted by context cancellation")
		return false
	}
}

// === UI BRIDGE: STEALTH MODE PRESETS ===

// SetStealthMode sets preset evasion configurations for quick UI activation
// Modes: "Aggressive" (all on, 1x), "Fast" (reduced delays, 0.5x), "Silent" (max stealth, 2x), "Debug" (all off)
func SetStealthMode(mode string) {
	globalStealthConfig.mu.Lock()
	defer globalStealthConfig.mu.Unlock()

	switch strings.ToLower(mode) {
	case "aggressive":
		globalStealthConfig.EnableJitter = true
		globalStealthConfig.EnableThinkingTime = true
		globalStealthConfig.EnableBackoff = true
		globalStealthConfig.EnablePathObfuscation = true
		globalStealthConfig.EnablePayloadEncoding = true
		globalStealthConfig.GlobalEvasionMultiplier = 1.0
		globalStealthConfig.Mode = "Aggressive"
		utils.TacticalLog("[red::b]STEALTH:[-] Mode set to [red]AGGRESSIVE[-] (all evasion enabled, 1x speed)[-:-:-]")

	case "fast":
		globalStealthConfig.EnableJitter = true
		globalStealthConfig.EnableThinkingTime = true
		globalStealthConfig.EnableBackoff = true
		globalStealthConfig.EnablePathObfuscation = true
		globalStealthConfig.EnablePayloadEncoding = false
		globalStealthConfig.GlobalEvasionMultiplier = 0.5
		globalStealthConfig.Mode = "Fast"
		utils.TacticalLog("[blue::b]STEALTH:[-] Mode set to [blue]FAST[-] (reduced delays, 0.5x speed)[-:-:-]")

	case "silent":
		globalStealthConfig.EnableJitter = true
		globalStealthConfig.EnableThinkingTime = true
		globalStealthConfig.EnableBackoff = true
		globalStealthConfig.EnablePathObfuscation = true
		globalStealthConfig.EnablePayloadEncoding = true
		globalStealthConfig.GlobalEvasionMultiplier = 2.0
		globalStealthConfig.Mode = "Silent"
		utils.TacticalLog("[green::b]STEALTH:[-] Mode set to [green]SILENT[-] (maximum evasion, 2x speed)[-:-:-]")

	case "debug":
		globalStealthConfig.EnableJitter = false
		globalStealthConfig.EnableThinkingTime = false
		globalStealthConfig.EnableBackoff = true // FIXED: Keep backoff enabled even in debug to prevent pipeline failures
		globalStealthConfig.EnablePathObfuscation = false
		globalStealthConfig.EnablePayloadEncoding = false
		globalStealthConfig.GlobalEvasionMultiplier = 1.0
		globalStealthConfig.Mode = "Debug"
		utils.TacticalLog("[yellow::b]STEALTH:[-] Mode set to [yellow]DEBUG[-] (evasion disabled except backoff for pipeline stability)[-:-:-]")

	default:
		utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Unknown stealth mode: %s. Use: aggressive|fast|silent|debug", mode))
	}
}

// GetStealthConfig returns current evasion configuration (read-only snapshot)
func GetStealthConfig() *StealthLevel {
	globalStealthConfig.mu.RLock()
	defer globalStealthConfig.mu.RUnlock()

	return globalStealthConfig
}

// SetEvasionToggle updates a specific evasion toggle
func SetEvasionToggle(feature string, enabled bool) {
	globalStealthConfig.mu.Lock()
	defer globalStealthConfig.mu.Unlock()

	switch strings.ToLower(feature) {
	case "jitter":
		globalStealthConfig.EnableJitter = enabled
		utils.TacticalLog(fmt.Sprintf("[cyan]STEALTH:[-] Jitter %s", map[bool]string{true: "[green]ENABLED", false: "[red]DISABLED"}[enabled]))

	case "thinking":
		globalStealthConfig.EnableThinkingTime = enabled
		utils.TacticalLog(fmt.Sprintf("[cyan]STEALTH:[-] Thinking Time %s", map[bool]string{true: "[green]ENABLED", false: "[red]DISABLED"}[enabled]))

	case "backoff":
		globalStealthConfig.EnableBackoff = enabled
		utils.TacticalLog(fmt.Sprintf("[cyan]STEALTH:[-] Rate-Limit Backoff %s", map[bool]string{true: "[green]ENABLED", false: "[red]DISABLED"}[enabled]))

	case "obfuscation":
		globalStealthConfig.EnablePathObfuscation = enabled
		utils.TacticalLog(fmt.Sprintf("[cyan]STEALTH:[-] Path Obfuscation %s", map[bool]string{true: "[green]ENABLED", false: "[red]DISABLED"}[enabled]))

	case "encoding":
		globalStealthConfig.EnablePayloadEncoding = enabled
		utils.TacticalLog(fmt.Sprintf("[cyan]STEALTH:[-] Payload Encoding %s", map[bool]string{true: "[green]ENABLED", false: "[red]DISABLED"}[enabled]))

	default:
		utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Unknown feature: %s", feature))
	}
}

// SetGlobalMultiplier scales all evasion sleep durations (0.1 = 10x speed, 2.0 = 2x slowdown)
func SetGlobalMultiplier(multiplier float64) {
	if multiplier < 0.1 || multiplier > 5.0 {
		utils.TacticalLog(fmt.Sprintf("[red]ERROR:[-] Multiplier must be between 0.1 and 5.0, got: %f", multiplier))
		return
	}

	globalStealthConfig.mu.Lock()
	globalStealthConfig.GlobalEvasionMultiplier = multiplier
	globalStealthConfig.mu.Unlock()

	utils.TacticalLog(fmt.Sprintf("[cyan]STEALTH:[-] Global evasion multiplier set to %.1fx", multiplier))
}

// GetStealthStatus returns a human-readable status string
func GetStealthStatus() string {
	config := GetStealthConfig()
	return fmt.Sprintf(
		"[cyan]STEALTH STATUS[-]\n  Mode: %s | Multiplier: %.1fx\n  Jitter: %s | Thinking: %s | Backoff: %s | Obfuscation: %s | Encoding: %s",
		config.Mode,
		config.GlobalEvasionMultiplier,
		map[bool]string{true: "[green]ON", false: "[red]OFF"}[config.EnableJitter],
		map[bool]string{true: "[green]ON", false: "[red]OFF"}[config.EnableThinkingTime],
		map[bool]string{true: "[green]ON", false: "[red]OFF"}[config.EnableBackoff],
		map[bool]string{true: "[green]ON", false: "[red]OFF"}[config.EnablePathObfuscation],
		map[bool]string{true: "[green]ON", false: "[red]OFF"}[config.EnablePayloadEncoding],
	)
}
