package cmd

import (
	"fmt"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/pterm/pterm" // Added for tactical UI
	"github.com/spf13/cobra"
)

var (
	targetURL string
	jsURL     string
	mineFlag  bool
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Reverse engineer API endpoints with real-time feedback",
	Run: func(cmd *cobra.Command, args []string) {
		if targetURL == "" && jsURL == "" {
			pterm.Error.Println("Target specification missing. Use -u (Swagger) or -j (JS Bundle)")
			return
		}

		pterm.DefaultHeader.WithFullWidth().Println("VaporTrace Mapper: Intelligence Gathering")

		// 1. Networking Setup
		proxyAddr, _ := cmd.Flags().GetString("proxy")
		utils.UpdateGlobalClient(proxyAddr)
		
		var allEndpoints []string

		// 2. Swagger Probing (Section 1)
		if targetURL != "" {
			spinner, _ := pterm.DefaultSpinner.Start("Parsing Swagger/OpenAPI Spec...")
			endpoints, err := discovery.ParseSwagger(targetURL, proxyAddr)
			if err != nil {
				spinner.Fail(fmt.Sprintf("Swagger Analysis Failed: %v", err))
			} else {
				spinner.Success(fmt.Sprintf("Identified %d endpoints from documentation", len(endpoints)))
				allEndpoints = append(allEndpoints, endpoints...)
			}
		}

		// 3. JS Scraper (Section 2 - The Phase 9.1 Core)
		if jsURL != "" {
			spinner, _ := pterm.DefaultSpinner.Start("Deep-Scraping JS Bundle for Hidden Routes...")
			jsEndpoints, err := discovery.ExtractJSPaths(jsURL, proxyAddr)
			if err != nil {
				spinner.Fail(fmt.Sprintf("JS Scraping Failed: %v", err))
			} else if len(jsEndpoints) == 0 {
				spinner.Warning("No API patterns matched in the JavaScript bundle")
			} else {
				spinner.Success(fmt.Sprintf("Harvested %d routes from client-side code", len(jsEndpoints)))
				
				// Render a tactical table of findings
				tableData := pterm.TableData{{"EXTRACTED PATH", "SOURCE"}}
				for _, je := range jsEndpoints {
					tableData = append(tableData, []string{je, "JS Static Analysis"})
					allEndpoints = append(allEndpoints, je)
				}
				pterm.DefaultTable.WithHasHeader().WithData(tableData).WithBoxed().Render()
			}
		}

		// 4. Parameter Mining (Section 3)
		if mineFlag && len(allEndpoints) > 0 {
			pterm.DefaultSection.Println("Phase 2.4: Tactical Parameter Mining")
			
			parts := strings.Split(targetURL, "/")
			if len(parts) >= 3 {
				baseURL := parts[0] + "//" + parts[2]
				
				// Progress bar for mining
				pb, _ := pterm.DefaultProgressbar.WithTotal(len(allEndpoints)).WithTitle("Mining Params").Start()
				for _, endpoint := range allEndpoints {
					pb.UpdateTitle("Mining: " + endpoint)
					discovery.MineParameters(baseURL, endpoint, proxyAddr)
					pb.Increment()
				}
				pb.Stop()
			}
		}

		pterm.Success.Println("Mapping Complete. Intelligence persisted to vaportrace.db")
	},
}

func init() {
	rootCmd.AddCommand(mapCmd)
	mapCmd.Flags().StringVarP(&targetURL, "url", "u", "", "URL of the Swagger/OpenAPI spec")
	mapCmd.Flags().StringVarP(&jsURL, "js", "j", "", "URL of a JS bundle to scrape")
	mapCmd.Flags().BoolVarP(&mineFlag, "mine", "m", false, "Enable parameter mining on discovered routes")
}