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

// Task 1: Deduplication for swagger findings within a session
var swaggerCache sync.Map

// Define the building blocks for Heuristic Discovery
var (
	swaggerPrefixes = []string{
		"",         // root
		"/api",     // common
		"/rest",    // legacy
		"/swagger", // explicit
		"/doc",     // documentation
		"/docs",
		"/service",
	}

	swaggerVersions = []string{
		"",        // no version
		"/v1",     // version 1
		"/v2",     // version 2
		"/v3",     // version 3
		"/v1.0",   // explicit float
		"/api/v1", // nested common
		"/api/v2",
	}

	swaggerFilenames = []string{
		"/swagger.json",
		"/openapi.json",
		"/api-docs",
		"/v2/api-docs", // Spring Boot default
		"/v3/api-docs", // Spring Boot OpenApi 3
		"/spec.json",
		"/docs.json",
	}
)

func ParseSwagger(url string, proxy string) ([]string, error) {
	// 1. Initial Attempt (User provided specific URL)
	endpoints, err := fetchAndParse(url)
	if err == nil && len(endpoints) > 0 {
		return endpoints, nil
	}

	// 2. Auto-Discovery Fallback (Combinatorial Heuristics)
	baseURL := strings.TrimRight(url, "/")

	// Remove common file extensions from input URL if present to get true root
	if strings.HasSuffix(baseURL, ".json") || strings.HasSuffix(baseURL, ".yaml") {
		lastSlash := strings.LastIndex(baseURL, "/")
		if lastSlash != -1 {
			baseURL = baseURL[:lastSlash]
		}
	}

	// Generate Candidate List
	candidates := generateCandidates(baseURL)

	utils.TacticalLog(fmt.Sprintf("[yellow]DISCOVER:[-] Direct parse failed. Engaging Heuristic Engine to probe %d potential locations...", len(candidates)))

	// Iterate through candidates
	// Note: In a future sprint, this could be parallelized with a worker pool.
	for _, probeURL := range candidates {
		// Log debug only if verbose mode or necessary, otherwise keep UI clean
		// utils.TacticalLog(fmt.Sprintf("[gray]Probing: %s[-]", probeURL))

		eps, err := fetchAndParse(probeURL)
		if err == nil && len(eps) > 0 {
			utils.TacticalLog(fmt.Sprintf("[green]SUCCESS:[-] Found valid spec at %s", probeURL))
			return eps, nil
		}
	}

	return nil, fmt.Errorf("failed to locate valid Swagger/OpenAPI spec after probing %d paths", len(candidates))
}

// generateCandidates builds the list of URLs to probe
func generateCandidates(base string) []string {
	uniqueMap := make(map[string]bool)
	var candidates []string

	// Helper to add unique paths
	add := func(p string) {
		if !uniqueMap[p] {
			uniqueMap[p] = true
			candidates = append(candidates, p)
		}
	}

	// Logic 1: Standard Static Paths (High Probability)
	add(base + "/swagger.json")
	add(base + "/openapi.json")
	add(base + "/api/swagger.json")
	add(base + "/v2/api-docs")

	// Logic 2: Combinatorial Generation
	// Structure: Base + Prefix + Version + Filename
	// Example: http://site.com + /api + /v1 + /swagger.json

	for _, prefix := range swaggerPrefixes {
		for _, version := range swaggerVersions {
			for _, file := range swaggerFilenames {
				// Avoid double slashes logic if strings are empty
				path := base

				if prefix != "" {
					path += prefix
				}
				if version != "" {
					path += version
				}
				path += file

				add(path)
			}
		}
	}

	return candidates
}

// fetchAndParse handles the network and version-agnostic parsing
func fetchAndParse(url string) ([]string, error) {
	client := logic.GlobalClient

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "VaporTrace-Scanner/3.1")
	req.Header.Set("Accept", "application/json, */*")

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

	// Validation: Ensure it's actually a Swagger doc
	isSwagger := false
	if _, ok := doc["swagger"]; ok {
		isSwagger = true
	}
	if _, ok := doc["openapi"]; ok {
		isSwagger = true
	}

	if !isSwagger {
		return nil, fmt.Errorf("valid JSON but missing 'swagger' or 'openapi' keys")
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

	// Log findings to database
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
