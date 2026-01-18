package cmd

import (
	"fmt"
	"strings" // Added for URL parsing

	"github.com/spf13/cobra"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
)

var targetURL string

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Reverse engineer API endpoints",
	Long:  `VaporTrace map will scan for Swagger specs, JS files, and hidden routes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if targetURL == "" {
			fmt.Println("[!] Error: Please provide a target spec URL with -u")
			return
		}

		fmt.Println("VaporTrace Mapping initialized...")
		proxyAddr, _ := cmd.Flags().GetString("proxy")

		// 1. Parse the Spec
		endpoints, err := discovery.ParseSwagger(targetURL, proxyAddr)
		if err != nil {
			fmt.Printf("[!] Mapping failed: %v\n", err)
			return
		}

		fmt.Printf("\n[+] Success! Discovered %d unique endpoints (prefixed with basePath)\n", len(endpoints))

		// 2. Generate Shadow Routes
		shadowRoutes := discovery.WalkVersions(endpoints)
		
		// 3. Live Validation
		if len(shadowRoutes) > 0 {
			fmt.Printf("[*] Phase 2.2: Probing %d shadow candidates via Proxy...\n", len(shadowRoutes))
			
			// Identify the Base URL (e.g., https://api.example.com)
			// We split by "://" to keep the protocol, then take the domain
			parts := strings.Split(targetURL, "/")
			baseURL := parts[0] + "//" + parts[2]

			for _, s := range shadowRoutes {
				status, err := discovery.ProbeEndpoint(baseURL, s, proxyAddr)
				if err != nil {
					continue 
				}

				// Focus on 200, 401, or 403 (anything but 404 is interesting)
				if status != 404 {
					fmt.Printf("    [!!!] INTERESTING: %s (Status: %d)\n", s, status)
				} else {
					fmt.Printf("    [-] Missing: %s (Status: 404)\n", s)
				}
			}
		}

		fmt.Println("\n[*] Recon Complete. Check Burp Suite 'HTTP History' for all validated paths.")
	},
}

func init() {
	rootCmd.AddCommand(mapCmd)
	mapCmd.Flags().StringVarP(&targetURL, "url", "u", "", "URL of the Swagger/OpenAPI spec")
}