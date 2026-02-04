package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// Task 1: Deduplication for swagger findings
var swaggerCache sync.Map

// Common paths to probe if the user provides a root URL
var commonSwaggerPaths = []string{
	"/swagger.json",
	"/openapi.json",
	"/spec.json", // httpbin.org
	"/api-docs",
	"/v2/api-docs",
	"/api/v2/api-docs",
	"/api/swagger.json",
	"/rest/admin/application-configuration", // Juice Shop (Variant)
}

func ParseSwagger(url string, proxy string) ([]string, error) {
	// 1. Initial Attempt
	endpoints, err := fetchAndParse(url)
	if err == nil && len(endpoints) > 0 {
		return endpoints, nil
	}

	// 2. Auto-Discovery Fallback (Heuristics)
	// If the user provided a root url (e.g. http://site.com), try common suffixes
	baseURL := strings.TrimRight(url, "/")

	utils.TacticalLog(fmt.Sprintf("[yellow]DISCOVER:[-] Direct parse failed. Probing %d common spec locations...", len(commonSwaggerPaths)))

	for _, suffix := range commonSwaggerPaths {
		probeURL := baseURL + suffix
		eps, err := fetchAndParse(probeURL)
		if err == nil && len(eps) > 0 {
			utils.TacticalLog(fmt.Sprintf("[green]SUCCESS:[-] Found spec at %s", probeURL))
			return eps, nil
		}
	}

	return nil, fmt.Errorf("failed to locate valid Swagger/OpenAPI spec at %s (or common paths)", url)
}

// fetchAndParse handles the network and version-agnostic parsing
func fetchAndParse(url string) ([]string, error) {
	client := logic.GlobalClient

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "VaporTrace-Scanner/3.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	// Use generic map to support Swagger 2.0 and OpenAPI 3.0+
	var doc map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, err
	}

	// Determine Base Path
	basePath := ""

	// Swagger 2.0
	if bp, ok := doc["basePath"].(string); ok {
		basePath = bp
	}

	// OpenAPI 3.0 (check 'servers' array)
	if servers, ok := doc["servers"].([]interface{}); ok && len(servers) > 0 {
		if srv, ok := servers[0].(map[string]interface{}); ok {
			if u, ok := srv["url"].(string); ok {
				// Often URLs are just "/" or relative
				if u != "/" {
					basePath = u
				}
			}
		}
	}

	// Parse Paths
	paths, ok := doc["paths"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no 'paths' key found in JSON")
	}

	var endpoints []string

	utils.RecordFinding(db.Finding{
		Phase:   "PHASE II: DISCOVERY",
		Command: "map",
		Target:  url,
		Details: "Swagger/OpenAPI Documentation Found",
		Status:  "INFO",
	})

	for pathKey := range paths {
		fullPath := pathKey
		// Join basePath if it exists and isn't already included
		if basePath != "" && basePath != "/" {
			if strings.HasSuffix(basePath, "/") && strings.HasPrefix(pathKey, "/") {
				fullPath = basePath + pathKey[1:]
			} else if !strings.HasSuffix(basePath, "/") && !strings.HasPrefix(pathKey, "/") {
				fullPath = basePath + "/" + pathKey
			} else {
				fullPath = basePath + pathKey
			}
		}

		if _, exists := swaggerCache.Load(fullPath); !exists {
			swaggerCache.Store(fullPath, true)
			endpoints = append(endpoints, fullPath)
			logic.GlobalDiscovery.AddEndpoint(fullPath)

			// Log for UI Table
			utils.LogMap(fullPath, "OpenAPI Spec", "200")
		}
	}

	return endpoints, nil
}

func WalkVersions(endpoints []string) []string {
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

// ProbeEndpoint logic restored for use by cmd/shell.go or other modules
func ProbeEndpoint(baseURL string, path string, proxy string) (int, error) {
	client := logic.GlobalClient

	fullURL := baseURL + path
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
