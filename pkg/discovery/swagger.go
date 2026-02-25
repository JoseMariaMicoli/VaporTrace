package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" // Added Persistence
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

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

	// PERSISTENCE HOOK: Log successful discovery of documentation
	db.LogQueue <- db.Finding{
		Phase:   "PHASE II: DISCOVERY",
		Target:  url,
		Details: "Swagger/OpenAPI Documentation Found",
		Status:  "INFO",
	}

	var endpoints []string
	for path := range doc.Paths {
		fullPath := doc.BasePath + path
		endpoints = append(endpoints, fullPath)
	}

	return endpoints, nil
}

func WalkVersions(endpoints []string) []string {
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

func ProbeEndpoint(baseURL string, path string, proxy string) (int, error) {
	client, err := utils.GetClient(proxy)
	if err != nil {
		return 0, err
	}

	fullURL := baseURL + path
	req, err := http.NewRequest(http.MethodHead, fullURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// PERSISTENCE HOOK: Log discovered live routes
	if resp.StatusCode == http.StatusOK {
		db.LogQueue <- db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Target:  fullURL,
			Details: "Live API Route Discovered",
			Status:  "SUCCESS",
		}
	}

	return resp.StatusCode, nil
}