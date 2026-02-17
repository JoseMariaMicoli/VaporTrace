/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

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

var (
	visitedURLs sync.Map
	// Regex to find href="..." and src="..."
	hrefRegex = regexp.MustCompile(`href=["'](.*?)["']`)
	srcRegex  = regexp.MustCompile(`src=["'](.*?)["']`)
)

// FetchAndExtract grabs a URL, logs it, and returns discovered links
func FetchAndExtract(target string) ([]string, error) {
	logic.EnsureTransport()

	req, _ := http.NewRequest("GET", target, nil)
	logic.ApplyEvasion(req) // Apply Stealth/Evasion settings

	resp, err := logic.GlobalClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Logic: Record everything, even 403s/500s (useful intelligence)
	status := fmt.Sprintf("%d", resp.StatusCode)
	logic.GlobalDiscovery.AddEndpoint(target)
	utils.LogMap(target, "Spider", status)

	// Persist finding to DB so it appears in the report
	utils.RecordFinding(db.Finding{
		Phase:   "PHASE II: DISCOVERY",
		Command: "spider",
		Target:  target,
		Details: fmt.Sprintf("Spider crawled endpoint (%s)", status),
		Status:  "INFO",
	})

	if resp.StatusCode == 404 {
		return nil, nil
	}

	// Only parse HTML/JS/JSON for new links
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/html") && !strings.Contains(ct, "application/javascript") && !strings.Contains(ct, "application/json") {
		return nil, nil
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	var links []string

	// Extract hrefs
	matches := hrefRegex.FindAllStringSubmatch(body, -1)
	for _, m := range matches {
		if len(m) > 1 {
			links = append(links, m[1])
		}
	}
	// Extract srcs
	srcMatches := srcRegex.FindAllStringSubmatch(body, -1)
	for _, m := range srcMatches {
		if len(m) > 1 {
			links = append(links, m[1])
		}
	}

	return links, nil
}

// StartSpider initiates the recursive crawl
func StartSpider(startURL string, maxDepth int) {
	u, err := url.Parse(startURL)
	if err != nil {
		utils.TacticalLog("[red]SPIDER ERROR:[-] Invalid start URL.")
		return
	}

	domain := u.Hostname()
	scheme := u.Scheme
	if scheme == "" {
		scheme = "https"
	}

	utils.TacticalLog(fmt.Sprintf("[green]SPIDER:[-] Initializing crawler on %s (Depth: %d, Scope: %s)...", startURL, maxDepth, domain))

	var wg sync.WaitGroup
	// Semaphore to limit concurrency (prevent FD exhaustion)
	semaphore := make(chan struct{}, 10)

	// Recursive crawler function
	var crawl func(string, int)
	crawl = func(currentURL string, depth int) {
		defer wg.Done()

		if depth > maxDepth {
			return
		}

		// Normalize URL
		currentURL = strings.TrimSpace(currentURL)
		if strings.HasPrefix(currentURL, "//") {
			currentURL = scheme + ":" + currentURL
		} else if strings.HasPrefix(currentURL, "/") {
			currentURL = scheme + "://" + domain + currentURL
		} else if !strings.HasPrefix(currentURL, "http") {
			return // Skip javascript:, mailto:, etc
		}

		// Scope Check (Stay on target domain)
		cu, err := url.Parse(currentURL)
		if err != nil || cu.Hostname() != domain {
			return
		}

		// Deduplication (Atomic check-and-set)
		if _, loaded := visitedURLs.LoadOrStore(currentURL, true); loaded {
			return
		}

		// Acquire Token
		semaphore <- struct{}{}

		// Execute Request
		links, err := FetchAndExtract(currentURL)

		// Release Token
		<-semaphore

		if err != nil {
			return
		}

		// Recurse for found links
		for _, link := range links {
			wg.Add(1)
			go crawl(link, depth+1)
		}
	}

	// Seed the recursive loop
	wg.Add(1)
	go crawl(startURL, 0)

	// Wait in background to avoid blocking main thread (Dashboard needs to update)
	go func() {
		wg.Wait()
		utils.TacticalLog("[green]✓ SPIDER:[-] Crawl complete. Check F2 Map.")
	}()
}
