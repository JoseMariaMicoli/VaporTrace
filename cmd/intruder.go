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
	"github.com/JoseMariaMicoli/VaporTrace/pkg/attack"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	intruderThreads int
)

// intruderCmd represents the intruder command
var intruderCmd = &cobra.Command{
	Use:   "intruder [mode] [url] [param] [wordlist]",
	Short: "Automated fuzzing engine (Sniper Mode)",
	Long: `VaporTrace Intruder: Tier 3 Attack Engine.
	
Sniper Mode iterates through a wordlist, replacing the value of a specific query parameter.
It automatically compares responses against a baseline to detect anomalies (status code changes, length variations).

Example:
  intruder sniper https://api.target.com/user?id=1 id /usr/share/wordlists/seclists/Fuzzing/big-list-of-naughty-strings.txt`,
	Args: cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		mode := args[0]
		target := args[1]
		param := args[2]
		wordlist := args[3]

		if mode != "sniper" {
			pterm.Error.Println("Currently only 'sniper' mode is supported in Day 1 implementation.")
			return
		}

		// Use global threads if flag not set specific
		if intruderThreads == 0 {
			intruderThreads = logic.CurrentSession.Threads
		}

		pterm.DefaultHeader.WithFullWidth().Println("VaporTrace Intruder: Tier 3 Engagement")
		pterm.Info.Printfln("Target: %s | Param: %s", target, param)
		pterm.Info.Printfln("Wordlist: %s | Threads: %d", wordlist, intruderThreads)

		config := attack.IntruderConfig{
			TargetURL:    target,
			Param:        param,
			WordlistPath: wordlist,
			Concurrency:  intruderThreads,
			Mode:         attack.Sniper,
		}

		// Execute
		attack.RunSniper(config)
	},
}

func init() {
	rootCmd.AddCommand(intruderCmd)
	intruderCmd.Flags().IntVarP(&intruderThreads, "threads", "t", 10, "Concurrency level")
}
