package discovery

import (
	"io"
	"regexp"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" 
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic" // Standardized to Logic for Phase 9.5
)

// PHASE 9.1 IMPROVEMENT: Pre-compile regex at the package level.
var (
	pathRegex = regexp.MustCompile(`"/(api|v[0-9]|rest)/[a-zA-Z0-9\-\_/]+ "`)
	urlRegex = regexp.MustCompile(`https?://[a-zA-Z0-9\.\-]+\.[a-z]{2,}/[a-zA-Z0-9\-\_/]+`)
)

func ExtractJSPaths(url string, proxy string) ([]string, error) {
	// Use GlobalClient from logic to support Phase 9.4 Proxy Sensing
	client := logic.GlobalClient

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
		path := m[1 : len(m)-1]
		cleaned = append(cleaned, path)

		// PHASE 9.5: Pipe to Global Store
		logic.GlobalDiscovery.AddEndpoint(path)

		db.LogQueue <- db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Target:  url,
			Details: "JS Endpoint Discovery (Relative): " + path,
			Status:  "INFO",
		}
	}

	// Process absolute URLs
	for _, u := range urlMatches {
		cleaned = append(cleaned, u)
		
		// PHASE 9.5: Pipe to Global Store
		logic.GlobalDiscovery.AddEndpoint(u)

		db.LogQueue <- db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Target:  url,
			Details: "JS Endpoint Discovery (Absolute): " + u,
			Status:  "INFO",
		}
	}

	return cleaned, nil
}