package discovery

import (
	"io"
	"regexp"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

var (
	pathRegex = regexp.MustCompile(`"/(api|v[0-9]|rest)/[a-zA-Z0-9\-\_/]+ "`)
	urlRegex  = regexp.MustCompile(`https?://[a-zA-Z0-9\.\-]+\.[a-z]{2,}/[a-zA-Z0-9\-\_/]+`)
)

func ExtractJSPaths(url string, proxy string) ([]string, error) {
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

	for _, m := range matches {
		path := m[1 : len(m)-1]
		cleaned = append(cleaned, path)

		logic.GlobalDiscovery.AddEndpoint(path)

		// ROUTING: Send to F2 Map Tab
		utils.LogMap(path, "JS-Scrape", "200")

		utils.RecordFinding(db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Command: "map", // Zero-Touch Trigger
			Target:  url,
			Details: "JS Endpoint Discovery (Relative): " + path,
			Status:  "INFO",
		})
	}

	for _, u := range urlMatches {
		cleaned = append(cleaned, u)
		logic.GlobalDiscovery.AddEndpoint(u)

		// ROUTING: Send to F2 Map Tab
		utils.LogMap(u, "JS-Scrape", "200")

		utils.RecordFinding(db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Command: "map", // Zero-Touch Trigger
			Target:  url,
			Details: "JS Endpoint Discovery (Absolute): " + u,
			Status:  "INFO",
		})
	}

	return cleaned, nil
}
