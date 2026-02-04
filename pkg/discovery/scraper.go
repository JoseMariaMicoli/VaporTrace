package discovery

import (
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

var (
	// Expanded Regex for SPA routing and relative API calls
	pathRegex = regexp.MustCompile(`["'](/[a-zA-Z0-9_\-\.\/]+)["']`)
	urlRegex  = regexp.MustCompile(`https?://[a-zA-Z0-9\.\-]+\.[a-z]{2,}/[a-zA-Z0-9\-\_/]+`)
	// Hash router support (Task 6)
	spaRegex = regexp.MustCompile(`[#][\/]([a-zA-Z0-9\-\_\/]+)`)

	// Task 1: Deduplication for scraper findings within a session
	scrapedCache sync.Map
)

func ExtractJSPaths(url string, proxy string) ([]string, error) {
	client := logic.GlobalClient

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JS bundle: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JS body: %v", err)
	}

	bodyStr := string(body)
	matches := pathRegex.FindAllString(bodyStr, -1)
	urlMatches := urlRegex.FindAllString(bodyStr, -1)
	spaMatches := spaRegex.FindAllString(bodyStr, -1)

	var cleaned []string

	addPath := func(p string, typeStr string) {
		// Task 1: Deduplication logic using sync.Map
		if _, exists := scrapedCache.Load(p); !exists && len(p) > 2 {
			scrapedCache.Store(p, true)
			cleaned = append(cleaned, p)
			logic.GlobalDiscovery.AddEndpoint(p)

			// Task 2 & 1: Send structured finding info via LogMap
			// Dashboard will regex parse this format
			utils.LogMap(p, fmt.Sprintf("%s (JS-Scraper)", typeStr), "200")

			utils.RecordFinding(db.Finding{
				Phase:   "PHASE II: DISCOVERY",
				Command: "scrape",
				Target:  url,
				Details: fmt.Sprintf("JS Discovery (%s): %s", typeStr, p),
				Status:  "INFO",
			})
		}
	}

	for _, m := range matches {
		// Strip quotes
		clean := m[1 : len(m)-1]
		if len(clean) > 200 {
			continue
		} // Noise filter
		// Filter common false positives
		if clean == "/json" || clean == "/application" {
			continue
		}
		addPath(clean, "Relative")
	}

	for _, u := range urlMatches {
		addPath(u, "Absolute")
	}

	for _, s := range spaMatches {
		// Clean the hash '#/' prefix to get the logic path
		clean := s[2:]
		addPath(clean, "SPA-Route")
	}

	return cleaned, nil
}
