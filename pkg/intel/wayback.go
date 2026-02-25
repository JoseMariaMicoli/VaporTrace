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

package intel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// WaybackResult represents a raw CDX API response line
type WaybackResult []string

// FetchWaybackHistory retrieves historical URLs for a domain
func FetchWaybackHistory(targetDomain string) {
	utils.TacticalLog(fmt.Sprintf("[magenta]INTEL:[-] Querying Wayback Machine for %s...", targetDomain))

	// Clean domain
	u, err := url.Parse(targetDomain)
	if err == nil && u.Hostname() != "" {
		targetDomain = u.Hostname()
	}

	// CDX API URL
	// fl=original: Only get the URL column
	// collapse=urlkey: Deduplicate results
	// filter=!statuscode:404: Ignore broken links (optional, but we want everything usually)
	apiURL := fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=json&fl=original&collapse=urlkey", targetDomain)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]INTEL ERROR:[-] Wayback query failed: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		utils.TacticalLog(fmt.Sprintf("[red]INTEL ERROR:[-] Wayback API returned status %d", resp.StatusCode))
		return
	}

	// Parse Response (Array of Arrays)
	var results [][]string
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &results); err != nil {
		utils.TacticalLog("[red]INTEL ERROR:[-] Failed to parse Wayback JSON")
		return
	}

	count := 0
	ignored := 0

	utils.TacticalLog(fmt.Sprintf("[blue]INTEL:[-] Processing %d raw historical records...", len(results)))

	for _, row := range results {
		if len(row) < 1 {
			continue
		}
		rawURL := row[0]

		// Filter static assets (Noise reduction)
		if isStaticAsset(rawURL) {
			ignored++
			continue
		}

		// Valid Candidate
		count++

		// 1. Add to Global Discovery (F2 Map)
		logic.GlobalDiscovery.AddEndpoint(rawURL)

		// 2. Log to Database
		utils.RecordFinding(db.Finding{
			Phase:        "TIER 4: OSINT",
			Command:      "intel-wayback",
			Target:       rawURL,
			Details:      "Historical endpoint discovered via Wayback Machine",
			Status:       "INFO",
			MitreTactic:  "Reconnaissance",
			CVSS_Numeric: 0.0,
		})
	}

	utils.TacticalLog(fmt.Sprintf("[green]✓ INTEL COMPLETE:[-] %d ghost endpoints added to map (%d ignored as static).", count, ignored))
}

func isStaticAsset(urlStr string) bool {
	lower := strings.ToLower(urlStr)
	extensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico",
		".css", ".woff", ".woff2", ".ttf", ".eot",
		".mp4", ".mp3", ".avi", ".mov",
		".pdf", ".doc", ".docx",
	}

	for _, ext := range extensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
		// Also check query params edge case "image.jpg?v=1"
		if strings.Contains(lower, ext+"?") {
			return true
		}
	}
	return false
}
