package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	// Ensure this matches your module name in go.mod
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Reverse engineer API endpoints",
	Long:  `VaporTrace map will scan for Swagger specs, JS files, and hidden routes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VaporTrace Mapping initialized...")

		// 1. Get the proxy address from the global flag
		// We use the persistent flag 'Proxy' defined in root.go
		proxyAddr, _ := cmd.Flags().GetString("proxy")

		// 2. Initialize our "Burp-aware" client from pkg/utils
		client, err := utils.GetClient(proxyAddr)
		if err != nil {
			fmt.Println("Error creating HTTP client:", err)
			return
		}

		// 3. Make a test request
		url := "http://example.com"
		fmt.Printf("Probing: %s via Proxy: %s\n", url, proxyAddr)

		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Request failed:", err)
			return
		}
		defer resp.Body.Close()

		fmt.Printf("Success! Status Code: %d\n", resp.StatusCode)
	},
}

func init() {
	rootCmd.AddCommand(mapCmd)
}