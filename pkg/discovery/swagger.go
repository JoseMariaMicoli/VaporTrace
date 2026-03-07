/*
Copyright (c) 2026 José María Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

package discovery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
		"/openapi.yaml",
		"/openapi.yml",
		"/swagger.yaml",
		"/swagger.yml",
		"/api-docs",
		"/v2/api-docs", // Spring Boot default
		"/v3/api-docs", // Spring Boot OpenApi 3
		"/spec.json",
		"/docs.json",
	}
)

func ParseSwagger(url string, proxy string) ([]string, error) {
	logic.EnsureTransport()

	// 1. Initial Attempt (User provided specific URL)
	endpoints, err := fetchAndParse(url)
	if err == nil && len(endpoints) > 0 {
		return endpoints, nil
	}

	// 2. Auto-Discovery Fallback (Combinatorial Heuristics on inferred base roots)
	candidates := generateCandidates(url)

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

func shouldTrimSwaggerLeaf(segment string) bool {
	s := strings.ToLower(strings.TrimSpace(segment))
	if s == "" {
		return false
	}

	switch s {
	case "index.html", "swagger-ui.html", "redoc.html", "docs", "doc", "swagger", "api-docs", "openapi":
		return true
	}

	return strings.HasPrefix(s, "openapi.")
}

func normalizeBaseRoots(raw string) []string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		fallback := strings.TrimRight(strings.TrimSpace(raw), "/")
		if fallback == "" {
			return []string{}
		}
		return []string{fallback}
	}

	origin := parsed.Scheme + "://" + parsed.Host
	roots := []string{origin}

	path := strings.Trim(parsed.Path, "/")
	if path == "" {
		return roots
	}

	segments := strings.Split(path, "/")

	if len(segments) > 0 && shouldTrimSwaggerLeaf(segments[len(segments)-1]) {
		segments = segments[:len(segments)-1]
	}

	if len(segments) == 0 {
		return roots
	}

	// Add progressively broader path roots, starting from the deepest.
	for i := len(segments); i >= 1; i-- {
		roots = append(roots, origin+"/"+strings.Join(segments[:i], "/"))
	}

	// If common docs markers are present, also include up-to-marker root.
	for i, segment := range segments {
		low := strings.ToLower(segment)
		if low == "swagger" || low == "docs" || low == "doc" || low == "api-docs" || low == "openapi" {
			if i == 0 {
				roots = append(roots, origin)
			} else {
				roots = append(roots, origin+"/"+strings.Join(segments[:i], "/"))
			}
		}
	}

	seen := make(map[string]bool)
	out := make([]string, 0, len(roots))
	for _, root := range roots {
		root = strings.TrimRight(root, "/")
		if root == "" || seen[root] {
			continue
		}
		seen[root] = true
		out = append(out, root)
	}

	return out
}

// generateCandidates builds the list of URLs to probe.
// It treats the input as a base URL and also infers parent roots for doc/UI paths.
func generateCandidates(raw string) []string {
	uniqueMap := make(map[string]bool)
	var candidates []string

	// Helper to add unique paths
	add := func(p string) {
		if !uniqueMap[p] {
			uniqueMap[p] = true
			candidates = append(candidates, p)
		}
	}

	roots := normalizeBaseRoots(raw)
	for _, base := range roots {
		// Priority 1: Highest-probability OpenAPI locations across common frameworks.
		priority := []string{
			"/openapi.json",
			"/swagger.json",
			"/openapi.yaml",
			"/openapi.yml",
			"/swagger.yaml",
			"/swagger.yml",
			"/v3/api-docs",
			"/v2/api-docs",
			"/api-docs",
			"/api/openapi.json",
			"/api/swagger.json",
			"/docs/openapi.json",
			"/docs/swagger.json",
			"/swagger/openapi.json",
			"/swagger/swagger.json",
			"/openapi/v1.json",
			"/q/openapi", // Quarkus
		}
		for _, suffix := range priority {
			add(base + suffix)
		}

		// Priority 2: Combinatorial generation for edge deployments.
		// Structure: Base + Prefix + Version + Filename.
		for _, prefix := range swaggerPrefixes {
			for _, version := range swaggerVersions {
				for _, file := range swaggerFilenames {
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
	}

	return candidates
}

// fetchAndParse handles the network and version-agnostic parsing
func fetchAndParse(url string) ([]string, error) {
	utils.TacticalLog(fmt.Sprintf("[cyan]SWAGGER:[-] Starting fetchAndParse for %s", url))
	logic.EnsureTransport()

	client := logic.GlobalClient
	utils.TacticalLog(fmt.Sprintf("[cyan]SWAGGER:[-] GlobalClient transport is nil? %v", client.Transport == nil))

	req, _ := http.NewRequest("GET", url, nil)
	logic.ApplyEvasion(req) // Apply rotating User-Agent
	req.Header.Set("Accept", "application/json, */*")

	utils.TacticalLog(fmt.Sprintf("[cyan]SWAGGER:[-] About to make request to %s", url))

	resp, err := client.Do(req)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]SWAGGER ERROR:[-] Request failed: %v", err))
		return nil, err
	}
	utils.TacticalLog(fmt.Sprintf("[green]✓ SWAGGER:[-] Got response: %d", resp.StatusCode))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 1) JSON Swagger/OpenAPI
	var doc map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &doc); err == nil {
		return extractEndpointsFromDocument(url, doc)
	}

	// 2) YAML OpenAPI (best-effort parser for paths block)
	yamlEndpoints := extractYAMLPaths(string(bodyBytes))
	if len(yamlEndpoints) > 0 {
		for _, fullPath := range yamlEndpoints {
			if _, exists := swaggerCache.Load(fullPath); !exists {
				swaggerCache.Store(fullPath, true)
				logic.GlobalDiscovery.AddEndpoint(fullPath)
				utils.LogMap(fullPath, "OpenAPI YAML", "200")
			}
		}
		utils.RecordFinding(db.Finding{
			Phase:   "PHASE II: DISCOVERY",
			Command: "map",
			Target:  url,
			Details: "OpenAPI YAML Documentation Found",
			Status:  "INFO",
		})
		return yamlEndpoints, nil
	}

	// 3) Swagger UI / ReDoc HTML pages that reference JSON spec URLs
	specRefs := extractSpecRefs(url, string(bodyBytes))
	for _, ref := range specRefs {
		eps, refErr := fetchAndParse(ref)
		if refErr == nil && len(eps) > 0 {
			return eps, nil
		}
	}

	return nil, fmt.Errorf("response is not parseable as JSON/YAML OpenAPI doc")
}

func extractEndpointsFromDocument(sourceURL string, doc map[string]interface{}) ([]string, error) {

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
		Target:  sourceURL,
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

func extractYAMLPaths(body string) []string {
	lines := strings.Split(body, "\n")
	pathsIndent := -1
	inPaths := false
	var endpoints []string

	for _, raw := range lines {
		line := strings.TrimRight(raw, " \t\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))

		if !inPaths {
			if trimmed == "paths:" {
				inPaths = true
				pathsIndent = indent
			}
			continue
		}

		if indent <= pathsIndent {
			break
		}

		if strings.HasPrefix(strings.TrimSpace(line), "/") && strings.HasSuffix(trimmed, ":") {
			pathKey := strings.TrimSuffix(strings.TrimSpace(line), ":")
			endpoints = append(endpoints, pathKey)
		}
	}

	return endpoints
}

func extractSpecRefs(baseURL, body string) []string {
	refRegexes := []*regexp.Regexp{
		regexp.MustCompile(`(?i)url\s*:\s*["']([^"']*(openapi|swagger|api-docs)[^"']*)["']`),
		regexp.MustCompile(`(?i)["']([^"']*(openapi|swagger|api-docs)[^"']*)["']`),
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var refs []string
	for _, re := range refRegexes {
		matches := re.FindAllStringSubmatch(body, -1)
		for _, m := range matches {
			if len(m) < 2 {
				continue
			}
			raw := strings.TrimSpace(m[1])
			if raw == "" || strings.HasPrefix(raw, "javascript:") {
				continue
			}
			u, err := url.Parse(raw)
			if err != nil {
				continue
			}
			resolved := base.ResolveReference(u).String()
			if seen[resolved] {
				continue
			}
			seen[resolved] = true
			refs = append(refs, resolved)
		}
	}

	return refs
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
