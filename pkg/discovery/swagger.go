package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// SwaggerDoc now captures BasePath to handle /v2 prefixing
type SwaggerDoc struct {
	BasePath string                 `json:"basePath"`
	Paths    map[string]interface{} `json:"paths"`
}

func ParseSwagger(url string, proxy string) ([]string, error) {
	client, err := utils.GetClient(proxy)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	var doc SwaggerDoc
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode Swagger JSON: %v", err)
	}

	var endpoints []string
	for path := range doc.Paths {
		// Prepend BasePath (e.g., /v2) to the endpoint (e.g., /pet)
		// This ensures WalkVersions can detect the versioning pattern
		fullPath := doc.BasePath + path
		endpoints = append(endpoints, fullPath)
	}

	return endpoints, nil
}

func WalkVersions(endpoints []string) []string {
	// Updated regex to find versioning anywhere in the path
	versionRegex := regexp.MustCompile(`v[0-9]+(\.[0-9]+)?|api|beta|dev`)
	substitutes := []string{"v1", "v2", "v3", "api", "dev"}
	
	candidates := make(map[string]bool)

	for _, path := range endpoints {
		if versionRegex.MatchString(path) {
			for _, sub := range substitutes {
				newPath := versionRegex.ReplaceAllString(path, sub)
				if newPath != path {
					candidates[newPath] = true
				}
			}
		}
	}

	var results []string
	for c := range candidates {
		results = append(results, c)
	}
	return results
}

// ProbeEndpoint performs a fast HEAD request through the proxy to verify if a route exists
func ProbeEndpoint(baseURL string, path string, proxy string) (int, error) {
	client, err := utils.GetClient(proxy)
	if err != nil {
		return 0, err
	}

	// Build the full target URL
	fullURL := baseURL + path
	
	// HEAD is faster for recon as it returns no body
	req, err := http.NewRequest(http.MethodHead, fullURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}