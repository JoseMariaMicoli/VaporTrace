package discovery

import (
	"io"
	"regexp"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" 
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// PHASE 9.1 IMPROVEMENT: Pre-compile regex at the package level.
// This prevents high CPU usage and memory fragmentation during large-scale scraping.
var (
	// PathRegex targets common API patterns: /api/v1, /rest/, etc.
	pathRegex = regexp.MustCompile(`"/(api|v[0-9]|rest)/[a-zA-Z0-9\-\_/]+ "`)
	
	// PHASE 10 PREVIEW: Added URL regex to catch full cross-domain API calls
	urlRegex = regexp.MustCompile(`https?://[a-zA-Z0-9\.\-]+\.[a-z]{2,}/[a-zA-Z0-9\-\_/]+`)
)

func ExtractJSPaths(url string, proxy string) ([]string, error) {
	// Use GlobalClient (Phase 1/9 Infrastructure)
	client := utils.GlobalClient

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	bodyStr := string(body)
	matches := pathRegex.FindAllString(bodyStr, -1)
	urlMatches := urlRegex.FindAllString(bodyStr, -1)

	var cleaned []string
	
	// Process relative API paths
	for _, m := range matches {
		// Strip the surrounding quotes
		path := m[1 : len(m)-1]
		cleaned = append(cleaned, path)

		// PERSISTENCE: Log finding to SQLite via Async Worker
		db.LogQueue <- db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Target:  url,
			Details: "JS Endpoint Discovery (Relative): " + path,
			Status:  "INFO",
		}
	}

	// Process absolute URLs (New in 9.1)
	for _, u := range urlMatches {
		cleaned = append(cleaned, u)
		db.LogQueue <- db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Target:  url,
			Details: "JS Endpoint Discovery (Absolute): " + u,
			Status:  "INFO",
		}
	}

	return cleaned, nil
}