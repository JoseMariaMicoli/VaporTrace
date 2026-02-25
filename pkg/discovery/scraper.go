package discovery

import (
	"io"
	"regexp"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" 
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils" // Import Central Utils
)

func ExtractJSPaths(url string, proxy string) ([]string, error) {
	// PATCH: Use GlobalClient
	client := utils.GlobalClient

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	pathRegex := regexp.MustCompile(`"/(api|v[0-9]|rest)/[a-zA-Z0-9\-\_/]+ "`)
	matches := pathRegex.FindAllString(string(body), -1)

	var cleaned []string
	for _, m := range matches {
		path := m[1 : len(m)-1]
		cleaned = append(cleaned, path)

		// PERSISTENCE HOOK
		db.LogQueue <- db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Target:  url,
			Details: "Extracted API Path from JS: " + path,
			Status:  "INFO",
		}
	}
	return cleaned, nil
}