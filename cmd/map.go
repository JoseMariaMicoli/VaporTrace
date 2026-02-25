package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
)

var (
	targetURL string
	jsURL     string
	mineFlag  bool
)

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Reverse engineer API endpoints",
	Long:  `VaporTrace map will scan for Swagger specs, JS files, and hidden routes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Logic check: ensure at least one source is provided
		if targetURL == "" && jsURL == "" {
			fmt.Println("[!] Error: Please provide a target spec URL (-u) or a JS URL (-j)")
			return
		}

		fmt.Println("VaporTrace Mapping initialized...")

		// Get the proxy address from the global flag
		proxyAddr, _ := cmd.Flags().GetString("proxy")

		// --- SECTION 1: SWAGGER & SHADOW PROBING (Phase 2.1 & 2.2) ---
		var allEndpoints []string

		if targetURL != "" {
			fmt.Printf("[*] Probing spec: %s via Proxy: %s\n", targetURL, proxyAddr)
			endpoints, err := discovery.ParseSwagger(targetURL, proxyAddr)
			if err != nil {
				fmt.Printf("[!] Mapping failed: %v\n", err)
			} else {
				allEndpoints = endpoints
				fmt.Printf("\n[+] Success! Discovered %d unique endpoints (prefixed with basePath)\n", len(endpoints))
				for _, e := range endpoints {
					fmt.Printf("    -> %s\n", e)
				}

				// Generate Shadow Routes (Version Walker)
				fmt.Println("\n[*] Phase 2.1: Analyzing paths for Shadow API versions (API9)...")
				shadowRoutes := discovery.WalkVersions(endpoints)

				if len(shadowRoutes) > 0 {
					fmt.Printf("[!] Found %d potential shadow routes to investigate.\n", len(shadowRoutes))
					fmt.Printf("[*] Phase 2.2: Probing shadow candidates via Proxy...\n")

					// Extract Base URL (e.g., https://petstore.swagger.io)
					parts := strings.Split(targetURL, "/")
					if len(parts) >= 3 {
						baseURL := parts[0] + "//" + parts[2]

						for _, s := range shadowRoutes {
							status, err := discovery.ProbeEndpoint(baseURL, s, proxyAddr)
							if err != nil {
								continue
							}

							if status != 404 {
								fmt.Printf("    [!!!] INTERESTING: %s (Status: %d)\n", s, status)
							} else {
								// Visible 404s for debugging and progress tracking
								fmt.Printf("    [-] Missing: %s (Status: 404)\n", s)
							}
						}
					}
				} else {
					fmt.Println("[-] No versioning patterns detected for automated walking.")
				}
			}
		}

		// --- SECTION 2: JS ROUTE SCRAPER (Phase 2.3) ---
		if jsURL != "" {
			fmt.Printf("\n[*] Phase 2.3: Scraping JS Bundle: %s\n", jsURL)
			jsEndpoints, err := discovery.ExtractJSPaths(jsURL, proxyAddr)
			if err != nil {
				fmt.Printf("[!] JS Scraping failed: %v\n", err)
			} else {
				fmt.Printf("[+] Found %d potential endpoints in JS:\n", len(jsEndpoints))
				for _, je := range jsEndpoints {
					fmt.Printf("    -> %s (extracted from JS)\n", je)
					allEndpoints = append(allEndpoints, je)
				}
			}
		}

		// --- SECTION 3: PARAMETER MINER (Phase 2.4) ---
		if mineFlag && len(allEndpoints) > 0 {
			fmt.Println("\n[*] Phase 2.4: Mining hidden parameters on discovered endpoints...")
			
			// We need a baseURL for mining. We use the one from targetURL or a default.
			parts := strings.Split(targetURL, "/")
			if len(parts) >= 3 {
				baseURL := parts[0] + "//" + parts[2]

				// Limit mining to the first 5 endpoints to avoid rate limits during testing
				limit := 5
				for i, endpoint := range allEndpoints {
					if i >= limit {
						break
					}
					fmt.Printf("    [*] Mining: %s\n", endpoint)
					discovery.MineParameters(baseURL, endpoint, proxyAddr)
				}
			}
		}

		fmt.Println("\n[*] Recon Complete. Check Burp Suite 'HTTP History' for all validated paths.")
	},
}

func init() {
	rootCmd.AddCommand(mapCmd)

	// Phase 2 Flags
	mapCmd.Flags().StringVarP(&targetURL, "url", "u", "", "URL of the Swagger/OpenAPI spec")
	mapCmd.Flags().StringVarP(&jsURL, "js", "j", "", "URL of a JavaScript bundle to scrape")
	mapCmd.Flags().BoolVarP(&mineFlag, "mine", "m", false, "Mine endpoints for hidden parameters")
}