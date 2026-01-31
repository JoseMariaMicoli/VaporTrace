package logic

import (
	"bufio"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

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
	return ProxyPool[rand.Intn(len(ProxyPool))]
}

// Tactical Fingerprints: Browser-accurate User-Agents
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0",
}

// ApplyEvasion handles the tactical transformation of a request before it is sent
func ApplyEvasion(req *http.Request) {
	rand.Seed(time.Now().UnixNano())

	// 1. Header Randomization (Phase 6.1)
	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	// 2. Timing Attacks: Sleepy Probes & Jitter (Phase 6.3)
	jitter := time.Duration(rand.Intn(130)+20) * time.Millisecond
	time.Sleep(jitter)
}
