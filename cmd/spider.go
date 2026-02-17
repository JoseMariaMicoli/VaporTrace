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

package cmd

import (
	"github.com/JoseMariaMicoli/VaporTrace/pkg/discovery"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	spiderDepth int
)

var spiderCmd = &cobra.Command{
	Use:   "spider [url]",
	Short: "Recursively crawl a target to build the attack map",
	Long:  `Active reconnaissance spider that extracts links, API endpoints, and JS files from HTML. All findings are automatically fed into the Strategic Buffer.`,
	Run: func(cmd *cobra.Command, args []string) {
		target := logic.CurrentSession.GetTarget()
		if len(args) > 0 {
			target = args[0]
		}

		if target == "" || target == "http://localhost" {
			pterm.Error.Println("No target specified. Usage: spider <url> or set global target.")
			return
		}

		pterm.DefaultHeader.WithFullWidth().Println("VaporTrace Spider: Active Reconnaissance")
		pterm.Info.Printfln("Target: %s | Max Depth: %d", target, spiderDepth)
		pterm.Info.Println("Starting background crawler... Check F2 (Map) for findings.")

		// Launch the spider logic (async handled internally or here)
		go discovery.StartSpider(target, spiderDepth)
	},
}

func init() {
	rootCmd.AddCommand(spiderCmd)
	spiderCmd.Flags().IntVarP(&spiderDepth, "depth", "d", 2, "Maximum crawl depth")
}
