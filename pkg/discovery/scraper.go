package discovery

import (
	"io"
	"regexp"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// ExtractJSPaths uses regex to find potential API paths in JS content
func ExtractJSPaths(url string, proxy string) ([]string, error) {
	client, err := utils.GetClient(proxy)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	// Regex for finding paths starting with /api/ or common patterns
	// This matches strings like "/api/v1/user" or "/v2/config"
	pathRegex := regexp.MustCompile(`"/(api|v[0-9]|rest)/[a-zA-Z0-9\-\_/]+ "`)
	matches := pathRegex.FindAllString(string(body), -1)

	var cleaned []string
	for _, m := range matches {
		// Remove the quotes
		cleaned = append(cleaned, m[1:len(m)-1])
	}
	return cleaned, nil
}