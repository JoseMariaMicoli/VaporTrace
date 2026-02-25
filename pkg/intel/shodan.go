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
	"net"
	"net/http"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// ShodanHost represents the API response structure
type ShodanHost struct {
	IP        string   `json:"ip_str"`
	Ports     []int    `json:"ports"`
	Hostnames []string `json:"hostnames"`
	Data      []struct {
		Port    int    `json:"port"`
		Product string `json:"product"`
		Data    string `json:"data"` // Banner
	} `json:"data"`
}

// QueryShodan fetches host info
func QueryShodan(target string) {
	apiKey := GetShodanKey()
	if apiKey == "" {
		utils.TacticalLog("[yellow]INTEL:[-] Shodan API key not configured. Run 'intel config shodan <key>'")
		return
	}

	// Resolve domain to IP if necessary
	ip := target
	if net.ParseIP(target) == nil {
		ips, err := net.LookupHost(target)
		if err == nil && len(ips) > 0 {
			ip = ips[0]
			utils.TacticalLog(fmt.Sprintf("[blue]INTEL:[-] Resolved %s to %s", target, ip))
		}
	}

	utils.TacticalLog(fmt.Sprintf("[magenta]INTEL:[-] Querying Shodan for IP %s...", ip))

	url := fmt.Sprintf("https://api.shodan.io/shodan/host/%s?key=%s", ip, apiKey)
	client := &http.Client{Timeout: 15 * time.Second}

	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]INTEL ERROR:[-] Connection failed: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		utils.TacticalLog("[yellow]INTEL:[-] No Shodan data found for this host.")
		return
	}
	if resp.StatusCode != 200 {
		utils.TacticalLog(fmt.Sprintf("[red]INTEL ERROR:[-] Shodan API status %d", resp.StatusCode))
		return
	}

	var host ShodanHost
	if err := json.NewDecoder(resp.Body).Decode(&host); err != nil {
		utils.TacticalLog("[red]INTEL ERROR:[-] Failed to decode response")
		return
	}

	// Process Results
	utils.TacticalLog(fmt.Sprintf("[green]✓ SHODAN HIT:[-] %s (%v)", host.IP, host.Hostnames))

	for _, item := range host.Data {
		msg := fmt.Sprintf("Port %d (%s)", item.Port, item.Product)
		utils.TacticalLog("  " + msg)

		// Record to DB
		utils.RecordFinding(db.Finding{
			Phase:       "TIER 4: OSINT",
			Command:     "intel-shodan",
			Target:      fmt.Sprintf("%s:%d", host.IP, item.Port),
			Details:     fmt.Sprintf("Open Port: %s | Banner: %s", item.Product, truncate(item.Data, 50)),
			Status:      "INFO",
			MitreTactic: "Reconnaissance",
		})
	}
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
