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
	fuzzType string // "params" or "paths"
)

var fuzzCmd = &cobra.Command{
	Use:   "fuzz [url]",
	Short: "Brute-force discovery for parameters and hidden paths",
	Long:  `Uses embedded top-100 wordlists to find hidden attack surface. Automatically detects anomalies in response size and status code.`,
	Run: func(cmd *cobra.Command, args []string) {
		target := logic.CurrentSession.GetTarget()
		if len(args) > 0 {
			target = args[0]
		}

		if target == "" || target == "http://localhost" {
			pterm.Error.Println("No target specified. Usage: fuzz <url> --type <params|paths>")
			return
		}

		pterm.DefaultHeader.WithFullWidth().Println("VaporTrace Fuzzer: Anomaly Detection")

		if fuzzType == "paths" {
			pterm.Info.Printfln("Starting Path Enumeration on %s", target)
			go discovery.FuzzPaths(target, nil) // nil uses default wordlist
		} else if fuzzType == "params" {
			pterm.Info.Printfln("Starting Parameter Mining on %s", target)
			go discovery.FuzzParams(target, nil) // nil uses default wordlist
		} else {
			pterm.Error.Println("Invalid type. Use --type params OR --type paths")
		}
	},
}

func init() {
	rootCmd.AddCommand(fuzzCmd)
	fuzzCmd.Flags().StringVarP(&fuzzType, "type", "t", "params", "Type of fuzzing: 'params' or 'paths'")
}
