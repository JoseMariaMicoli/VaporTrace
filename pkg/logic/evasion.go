package logic

import (
	"math/rand"
	"time"
)

// UserAgents list to rotate fingerprints
var UserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/120.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
}

// GetRandomUA returns a random User-Agent from the pool
func GetRandomUA() string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	return UserAgents[seed.Intn(len(UserAgents))]
}

// ApplyJitter introduces a small, random delay to bypass simple rate-limiters
func ApplyJitter(baseDelayMs int) {
	if baseDelayMs == 0 {
		return
	}
	// Add up to 50% variance
	variation := rand.Intn(baseDelayMs / 2)
	time.Sleep(time.Duration(baseDelayMs+variation) * time.Millisecond)
}