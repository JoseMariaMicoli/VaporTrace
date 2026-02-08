package logic

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// === PRIORITY GAMMA: CONTEXTUAL THINKING TIME ===

// RequestContext categorizes the type of request for behavioral simulation
type RequestContext int

const (
	Discovery      RequestContext = iota // GET, OPTION requests (quick)
	Reconnaissance                       // Headers, timing probes (medium)
	Exploitation                         // POST, PUT, DELETE (slow/thinking)
	PixelRequest                         // Tracking pixels, stylesheets (instant)
)

// ContextualThinkingTime returns a delay based on request type and method
// Simulates human behavior:
// - GET (discovery): Quick (10-50ms)
// - POST (form submission): Slow (500-2000ms - thinking time)
// - PUT/DELETE (risky): Very slow (1-5s)
func ContextualThinkingTime(method string, path string) time.Duration {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	context := categorizeRequest(method, path)

	var minDelay, maxDelay int

	switch context {
	case Discovery:
		// GET requests: Quick scan (10-50ms)
		minDelay = 10
		maxDelay = 50
		utils.TacticalLog("[blue]THINKING:[-] Discovery request: Light jitter (10-50ms)")

	case Reconnaissance:
		// HEAD, OPTIONS, TRACE: Medium delay (50-300ms)
		minDelay = 50
		maxDelay = 300
		utils.TacticalLog("[blue]THINKING:[-] Reconnaissance request: Medium jitter (50-300ms)")

	case Exploitation:
		// POST, PUT, DELETE: Heavy delay (800-3000ms)
		// Simulates "filling out forms" or "processing data"
		minDelay = 800
		maxDelay = 3000
		utils.TacticalLog("[blue]THINKING:[-] Exploitation request: Heavy jitter (800-3000ms)")

	case PixelRequest:
		// Static resources, tracking pixels: Instant (0-5ms)
		minDelay = 0
		maxDelay = 5
		utils.TacticalLog("[blue]THINKING:[-] Pixel/static request: No jitter")
	}

	// Use Gaussian distribution for realistic delays
	delay := gaussianJitter(minDelay, maxDelay, seed)
	return time.Duration(delay) * time.Millisecond
}

// categorizeRequest determines the request context
func categorizeRequest(method, path string) RequestContext {
	switch method {
	case "GET":
		return Discovery
	case "HEAD", "OPTIONS", "TRACE":
		return Reconnaissance
	case "POST", "PUT", "DELETE", "PATCH":
		return Exploitation
	default:
		return Discovery
	}
}

// gaussianJitter applies Gaussian distribution with clamping
func gaussianJitter(min, max int, seed *rand.Rand) int {
	if min >= max {
		return min
	}

	mean := float64(min+max) / 2.0
	stddev := float64(max-min) / 4.0

	u1 := rand.Float64()
	u2 := rand.Float64()
	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	jitterValue := mean + (stddev * z)

	if jitterValue < float64(min) {
		jitterValue = float64(min)
	}
	if jitterValue > float64(max) {
		jitterValue = float64(max)
	}

	return int(jitterValue)
}

// ApplyContextualBehavior wraps the thinking time logic with SafeSleep integration
// Called before sending any request to simulate human behavior
func ApplyContextualBehavior(req *http.Request) {
	if req == nil {
		return
	}

	delay := ContextualThinkingTime(req.Method, req.URL.Path)

	if delay > 0 {
		// Use context from request (with a 30-second timeout for safety)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		utils.TacticalLog(fmt.Sprintf("[cyan]BEHAVIOR:[-] Contextual thinking time: %dms", delay.Milliseconds()))
		SafeSleep(ctx, delay, &globalStealthConfig.EnableThinkingTime)
	}
}

// SimulatePauses introduces random "human pause" between different attack phases
// WAF bot detection often looks for continuous scanning without breaks
// Now respects EnableThinkingTime toggle
func SimulatePauses(minSeconds, maxSeconds int) {
	if minSeconds >= maxSeconds {
		return
	}

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	seconds := seed.Intn(maxSeconds-minSeconds) + minSeconds

	duration := time.Duration(seconds) * time.Second
	ctx := context.Background()

	utils.TacticalLog(fmt.Sprintf("[yellow]PAUSE:[-] Human behavior simulation: sleeping for %d seconds", seconds))
	SafeSleep(ctx, duration, &globalStealthConfig.EnableThinkingTime)
}

// IsLikelyBotPattern detects if our current request pattern is too robotic
// Returns true if pattern suggests we should take a pause
func IsLikelyBotPattern(requestCount int, timeWindow int) bool {
	// If we've made > 20 requests in < 5 seconds, we're being too aggressive
	threshold := 20
	if requestCount > threshold {
		utils.TacticalLog(fmt.Sprintf("[red]BOT PATTERN DETECTED:[-] %d requests in %ds window", requestCount, timeWindow))
		return true
	}
	return false
}
