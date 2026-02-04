package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type SwaggerDoc struct {
	BasePath string                 `json:"basePath"`
	Paths    map[string]interface{} `json:"paths"`
	// Components/Definitions support for future graph mapping (Sprint 11)
	Definitions map[string]interface{} `json:"definitions"`
	Components  map[string]interface{} `json:"components"`
}

func ParseSwagger(url string, proxy string) ([]string, error) {
	client := logic.GlobalClient

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	var doc SwaggerDoc
	// Attempt decoding
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode Swagger JSON: %v", err)
	}

	utils.RecordFinding(db.Finding{
		Phase:   "PHASE II: DISCOVERY",
		Command: "map",
		Target:  url,
		Details: "Swagger/OpenAPI Documentation Found (Deep Parse)",
		Status:  "INFO",
	})

	var endpoints []string
	for path := range doc.Paths {
		// Normalize base path usage
		fullPath := path
		if doc.BasePath != "" && doc.BasePath != "/" {
			fullPath = doc.BasePath + path
		}
		endpoints = append(endpoints, fullPath)
		logic.GlobalDiscovery.AddEndpoint(fullPath)
	}

	return endpoints, nil
}

func WalkVersions(endpoints []string) []string {
	// Task 6: Enhanced Recon
	// Logic to predict shadow versions /v1 -> /v2, /beta
	versionRegex := regexp.MustCompile(`v[0-9]+(\.[0-9]+)?|api|beta|dev|prod`)
	substitutes := []string{"v1", "v2", "v3", "api", "dev", "beta", "staging", "internal"}
	candidates := make(map[string]bool)

	for _, path := range endpoints {
		if versionRegex.MatchString(path) {
			for _, sub := range substitutes {
				newPath := versionRegex.ReplaceAllString(path, sub)
				if newPath != path {
					candidates[newPath] = true
				}
			}
		} else {
			// Prepend assumption
			for _, sub := range substitutes {
				candidates["/"+sub+path] = true
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
	client := logic.GlobalClient

	fullURL := baseURL + path
	// Use GET instead of HEAD for some aggressive CDNs
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 404 {
		utils.RecordFinding(db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Command: "map",
			Target:  fullURL,
			Details: fmt.Sprintf("Shadow Endpoint Active (%d)", resp.StatusCode),
			Status:  "SUCCESS",
		})
	}

	return resp.StatusCode, nil
}
