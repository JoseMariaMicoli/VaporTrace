package logic

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// === PRIORITY DELTA: RATE-LIMIT BACKOFF STRATEGY ===

// RateLimitState tracks 429 and WAF challenge responses
type RateLimitState struct {
	mu               sync.RWMutex
	LastRateLimit    time.Time
	RateLimitCount   int
	BackoffActive    bool
	NextRetryTime    time.Time
	CurrentProxy     string
	CurrentUserAgent string
}

var globalRateLimitState = &RateLimitState{
	RateLimitCount: 0,
	BackoffActive:  false,
}

// HandleRateLimit processes 429 responses and triggers backoff strategy
// Returns the delay before the next request can be made
// Now respects EnableBackoff toggle
func HandleRateLimit(statusCode int, responseHeaders map[string][]string) time.Duration {
	globalRateLimitState.mu.Lock()
	defer globalRateLimitState.mu.Unlock()

	// Detect rate limit or WAF challenge
	if statusCode == 429 || (statusCode >= 400 && statusCode < 430) {
		globalRateLimitState.RateLimitCount++
		globalRateLimitState.LastRateLimit = time.Now()

		// Calculate exponential backoff
		backoffDelay := calculateExponentialBackoff(globalRateLimitState.RateLimitCount)

		utils.TacticalLog(fmt.Sprintf("[red]RATE_LIMIT:[-] HTTP %d detected. Backoff #%d: %d seconds",
			statusCode, globalRateLimitState.RateLimitCount, backoffDelay))

		// Set backoff active
		globalRateLimitState.BackoffActive = true
		globalRateLimitState.NextRetryTime = time.Now().Add(time.Duration(backoffDelay) * time.Second)

		// Trigger evasion rotation
		rotateEvasionIdentity()

		return time.Duration(backoffDelay) * time.Second
	}

	// Reset on success
	if statusCode == 200 || statusCode == 201 {
		globalRateLimitState.RateLimitCount = 0
		globalRateLimitState.BackoffActive = false
	}

	return 0
}

// calculateExponentialBackoff computes backoff time with jitter
// Backoff: 2^attempt + random jitter (30-60s range)
func calculateExponentialBackoff(attempt int) int {
	if attempt <= 0 {
		attempt = 1
	}

	// Exponential backoff: 2^attempt
	baseDelay := math.Pow(2, float64(attempt-1))
	if baseDelay > 120 { // Max 2 minutes
		baseDelay = 120
	}

	// Add random jitter (±30%)
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	jitter := float64(seed.Intn(100)-50) / 100.0 // -50% to +50%
	finalDelay := baseDelay * (1 + jitter)

	if finalDelay < 30 {
		finalDelay = 30 // Minimum 30 seconds
	}

	return int(finalDelay)
}

// rotateEvasionIdentity changes proxy and User-Agent after rate limit
// This makes resumption less suspicious
func rotateEvasionIdentity() {
	utils.TacticalLog("[yellow]EVASION:[-] Rotating proxy and User-Agent identity...")

	// Rotate proxy if pool available
	if len(ProxyPool) > 0 {
		newProxy := GetRandomProxy()
		globalRateLimitState.CurrentProxy = newProxy
		utils.TacticalLog(fmt.Sprintf("[cyan]PROXY:[-] Switched to %s", newProxy))
		SetProxy(newProxy)
	}

	utils.TacticalLog("[cyan]USER-AGENT:[-] Next requests will use rotated fingerprint")
}

// ApplyBackoffWithSleep applies the rate-limit backoff using SafeSleep
// This respects the EnableBackoff toggle and context cancellation
func ApplyBackoffWithSleep(ctx context.Context, delay time.Duration) bool {
	if delay <= 0 {
		return true
	}

	utils.TacticalLog(fmt.Sprintf("[yellow]BACKOFF:[-] Entering exponential backoff for %d seconds", int(delay.Seconds())))

	// Use SafeSleep with EnableBackoff toggle
	completed := SafeSleep(ctx, delay, &globalStealthConfig.EnableBackoff)

	if !completed {
		utils.TacticalLog("[yellow]BACKOFF:[-] Backoff interrupted or skipped")
	} else {
		utils.TacticalLog("[green]BACKOFF:[-] Exponential backoff completed, resuming requests")
	}

	return completed
}

// IsBackoffActive returns whether we're currently in cooldown
func IsBackoffActive() bool {
	globalRateLimitState.mu.RLock()
	defer globalRateLimitState.mu.RUnlock()

	if !globalRateLimitState.BackoffActive {
		return false
	}

	// Check if backoff period has expired
	if time.Now().After(globalRateLimitState.NextRetryTime) {
		return false
	}

	return true
}

// GetBackoffWaitTime returns how long to wait before next request
func GetBackoffWaitTime() time.Duration {
	globalRateLimitState.mu.RLock()
	defer globalRateLimitState.mu.Unlock()

	if !globalRateLimitState.BackoffActive {
		return 0
	}

	waitTime := time.Until(globalRateLimitState.NextRetryTime)
	if waitTime < 0 {
		return 0
	}

	return waitTime
}

// ResetRateLimitState manually resets the backoff (for testing)
func ResetRateLimitState() {
	globalRateLimitState.mu.Lock()
	defer globalRateLimitState.mu.Unlock()

	globalRateLimitState.RateLimitCount = 0
	globalRateLimitState.BackoffActive = false
	globalRateLimitState.NextRetryTime = time.Time{}

	utils.TacticalLog("[green]✓ RATE_LIMIT:[-] Backoff state reset")
}

// GetRateLimitStatus returns current backoff statistics
func GetRateLimitStatus() map[string]interface{} {
	globalRateLimitState.mu.RLock()
	defer globalRateLimitState.mu.RUnlock()

	return map[string]interface{}{
		"active":            globalRateLimitState.BackoffActive,
		"count":             globalRateLimitState.RateLimitCount,
		"last_trigger":      globalRateLimitState.LastRateLimit,
		"next_retry":        globalRateLimitState.NextRetryTime,
		"wait_time_seconds": GetBackoffWaitTime().Seconds(),
	}
}
