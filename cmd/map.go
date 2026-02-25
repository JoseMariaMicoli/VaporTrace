package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils" // Import Central Utils
)

var (
	targetURL string
	jsURL     string
	mineFlag  bool
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Reverse engineer API endpoints",
	Long:  `VaporTrace map will scan for Swagger specs, JS files, and hidden routes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if targetURL == "" && jsURL == "" {
			fmt.Println("[!] Error: Please provide a target spec URL (-u) or a JS URL (-j)")
			return
		}

		fmt.Println("VaporTrace Mapping initialized...")

		// Get the proxy address from the global flag
		proxyAddr, _ := cmd.Flags().GetString("proxy")

		// PATCH: Activate the Global Proxy immediately.
		// All subsequent discovery calls will now route through this automatically.
		utils.UpdateGlobalClient(proxyAddr)

		// --- SECTION 1: SWAGGER & SHADOW PROBING ---
		var allEndpoints []string

		if targetURL != "" {
			fmt.Printf("[*] Probing spec: %s via Proxy: %s\n", targetURL, proxyAddr)
			
			// Note: We pass proxyAddr just to satisfy the signature, but the logic uses GlobalClient internally now.
			endpoints, err := discovery.ParseSwagger(targetURL, proxyAddr)
			if err != nil {
				fmt.Printf("[!] Mapping failed: %v\n", err)
			} else {
				allEndpoints = endpoints
				fmt.Printf("\n[+] Success! Discovered %d unique endpoints\n", len(endpoints))
				for _, e := range endpoints {
					fmt.Printf("    -> %s\n", e)
				}

				fmt.Println("\n[*] Phase 2.1: Analyzing paths for Shadow API versions (API9)...")
				shadowRoutes := discovery.WalkVersions(endpoints)

				if len(shadowRoutes) > 0 {
					fmt.Printf("[!] Found %d potential shadow routes to investigate.\n", len(shadowRoutes))
					
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
							}
						}
					}
				}
			}
		}

		// --- SECTION 2: JS ROUTE SCRAPER ---
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

		// --- SECTION 3: PARAMETER MINER ---
		if mineFlag && len(allEndpoints) > 0 {
			fmt.Println("\n[*] Phase 2.4: Mining hidden parameters on discovered endpoints...")
			
			parts := strings.Split(targetURL, "/")
			if len(parts) >= 3 {
				baseURL := parts[0] + "//" + parts[2]
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

		fmt.Println("\n[*] Recon Complete. Check Burp Suite 'HTTP History'.")
	},
}

func init() {
	rootCmd.AddCommand(mapCmd)
	mapCmd.Flags().StringVarP(&targetURL, "url", "u", "", "URL of the Swagger/OpenAPI spec")
	mapCmd.Flags().StringVarP(&jsURL, "js", "j", "", "URL of a JavaScript bundle to scrape")
	mapCmd.Flags().BoolVarP(&mineFlag, "mine", "m", false, "Mine endpoints for hidden parameters")
}