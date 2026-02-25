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

package engine

import (
	"fmt"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// HeuristicRule defines the contract for any tactical detection logic.
type HeuristicRule struct {
	Name        string
	Description string
	Evaluate    func(endpoint string, status int, loot logic.LootSummary) (bool, float64, string) // Returns: triggered, confidence, reason
}

// ActiveRules holds the registry of deterministic detection logic.
var ActiveRules = []HeuristicRule{
	{
		Name:        "BOLA-Detector",
		Description: "Detects numerical or UUID resource IDs in paths.",
		Evaluate: func(endpoint string, status int, loot logic.LootSummary) (bool, float64, string) {
			// Pattern Matching for BOLA indicators (Integers or UUIDs)
			if strings.Contains(endpoint, "{id}") ||
				strings.Contains(endpoint, "/user/") ||
				strings.Contains(endpoint, "/order/") ||
				strings.Contains(endpoint, "/invoice/") ||
				strings.Contains(endpoint, "/account/") {

				score := 50.0 // Base confidence (Medium)
				reason := "Endpoint pattern accepts variable Resource IDs."

				// Enrich confidence if we have credentials to exploit it (Phase 3 Requisite)
				if loot.HasJWT || loot.Credential != "" {
					score = 85.0 // High Confidence
					reason += " Valid Authorization Token available for privilege swapping."
				}

				// Additional boost if status is 200 (endpoint is reachable)
				if status == 200 {
					score = 75.0
					reason += " Status 200 confirms endpoint accessibility."
				}

				return true, score, reason
			}
			return false, 0.0, ""
		},
	},
	{
		Name:        "Admin-PrivEsc",
		Description: "Identifies protected administrative interfaces.",
		Evaluate: func(endpoint string, status int, loot logic.LootSummary) (bool, float64, string) {
			path := strings.ToLower(endpoint)

			// Detect High Value Targets
			if strings.Contains(path, "admin") ||
				strings.Contains(path, "config") ||
				strings.Contains(path, "dashboard") ||
				strings.Contains(path, "settings") ||
				strings.Contains(path, "superuser") {

				score := 70.0
				reason := "Administrative keyword detected in path."

				// If we saw a 403 Forbidden, we know it exists but is blocked.
				// This is the ideal candidate for BFLA/Verb Tampering.
				if status == 403 {
					score = 95.0 // Critical
					reason += " History confirms 403 Forbidden: Target exists and enforces ACLs (BFLA candidate)."
				}

				// If we have admin credentials, increase score
				if strings.Contains(loot.Credential, "admin") || strings.Contains(loot.Credential, "root") {
					score = 90.0
					reason += " Administrator credentials available for access attempt."
				}

				return true, score, reason
			}
			return false, 0.0, ""
		},
	},
	{
		Name:        "Cloud-SSRF",
		Description: "Detects parameters prone to Server-Side Request Forgery.",
		Evaluate: func(endpoint string, status int, loot logic.LootSummary) (bool, float64, string) {
			path := strings.ToLower(endpoint)

			// Check for Callback Parameters (OOB)
			if strings.Contains(path, "url=") ||
				strings.Contains(path, "hook") ||
				strings.Contains(path, "callback") ||
				strings.Contains(path, "redirect") ||
				strings.Contains(path, "webhook") {

				score := 60.0
				reason := "Callback parameter detected in URL structure."

				// If we already stole AWS keys, SSRF is CRITICAL for lateral movement
				if loot.HasAWS {
					score = 99.0
					reason += " AWS Keys already exfiltrated. High risk of Metadata access via SSRF."
				}
				return true, score, reason
			}
			return false, 0.0, ""
		},
	},
	{
		Name:        "Info-Exposure",
		Description: "Detects sensitive file extensions and backup artifacts.",
		Evaluate: func(endpoint string, status int, loot logic.LootSummary) (bool, float64, string) {
			path := strings.ToLower(endpoint)

			// Look for backup files, logs, or configs
			if strings.HasSuffix(path, ".log") ||
				strings.HasSuffix(path, ".bak") ||
				strings.HasSuffix(path, ".sql") ||
				strings.HasSuffix(path, ".env") ||
				strings.Contains(path, ".git/") {

				score := 90.0
				reason := "High-risk file extension detected."

				if status == 200 {
					score = 100.0
					reason += " Confirmed 200 OK Access."
				}
				return true, score, reason
			}
			return false, 0.0, ""
		},
	},
	{
		Name:        "BOPLA-Detector",
		Description: "Detects broken object property level authorization.",
		Evaluate: func(endpoint string, status int, loot logic.LootSummary) (bool, float64, string) {
			path := strings.ToLower(endpoint)

			// Look for endpoints that modify object properties
			if strings.Contains(path, "/update") ||
				strings.Contains(path, "/edit") ||
				strings.Contains(path, "/modify") ||
				strings.Contains(path, "/patch") ||
				strings.Contains(path, "/properties") {

				score := 55.0
				reason := "Endpoint modifies object properties. BOPLA (Broken Object Property Level Auth) possible."

				if loot.HasJWT {
					score = 70.0
					reason += " Auth token available for testing privilege escalation."
				}
				return true, score, reason
			}
			return false, 0.0, ""
		},
	},
	{
		Name:        "Injection-Vectors",
		Description: "Detects endpoints vulnerable to injection attacks.",
		Evaluate: func(endpoint string, status int, loot logic.LootSummary) (bool, float64, string) {
			path := strings.ToLower(endpoint)

			// Look for common injection parameters
			if strings.Contains(path, "search") ||
				strings.Contains(path, "query") ||
				strings.Contains(path, "filter") ||
				strings.Contains(path, "sql") ||
				strings.Contains(path, "cmd") ||
				strings.Contains(path, "input") {

				score := 40.0
				reason := "Endpoint accepts user input. Injection vulnerability possible."

				if status == 200 {
					score = 60.0
					reason += " Endpoint is accessible. Ready for injection testing."
				}
				return true, score, reason
			}
			return false, 0.0, ""
		},
	},
}

// RunHeuristics evaluates a single endpoint against all heuristic rules and returns
// the highest confidence finding as a TacticalAction. This is the primary heuristic
// scoring engine used by the strategic planner.
func RunHeuristics(endpoint string, lastStatus int, loot logic.LootSummary) *TacticalAction {
	var bestMatch *TacticalAction
	highestScore := 0.0

	for _, rule := range ActiveRules {
		triggered, score, reason := rule.Evaluate(endpoint, lastStatus, loot)

		// Priority Logic: Higher scores overwrite lower scores for the same endpoint
		if triggered && score > highestScore {
			highestScore = score

			// Translate Numeric Score to Confidence String
			confStr := "LOW"
			if score >= 80 {
				confStr = "CRITICAL"
			} else if score >= 60 {
				confStr = "HIGH"
			} else if score >= 40 {
				confStr = "MEDIUM"
			}

			// Map Heuristic to Tactical Type used by Core Engine
			actionType := "GENERIC"
			if rule.Name == "BOLA-Detector" {
				actionType = "BOLA"
			} else if rule.Name == "Admin-PrivEsc" {
				actionType = "BFLA"
			} else if rule.Name == "Cloud-SSRF" {
				actionType = "SSRF"
			} else if rule.Name == "Info-Exposure" {
				actionType = "AUDIT"
			} else if rule.Name == "BOPLA-Detector" {
				actionType = "BOPLA"
			} else if rule.Name == "Injection-Vectors" {
				actionType = "INJECTION"
			}

			// Generate the candidate action with default payload based on type
			payload := GeneratePayload(actionType)

			// Generate the candidate action
			bestMatch = &TacticalAction{
				Type:       actionType,
				Target:     endpoint, // Context-relative path (will be enriched in Neuro)
				Payload:    payload,
				Confidence: confStr,
				Reasoning:  reason,
				Status:     "PENDING",
			}
		}
	}
	return bestMatch
}

// GeneratePayload creates a reasonable default payload based on the attack type
func GeneratePayload(actionType string) string {
	switch actionType {
	case "BOLA":
		return "ID: 1337"
	case "BFLA":
		return "Method: DELETE"
	case "BOPLA":
		return "Property-Name: admin"
	case "SSRF":
		return "URL: http://169.254.169.254/latest/meta-data/"
	case "INJECTION":
		return "Input: ' OR '1'='1"
	case "AUDIT":
		return "Action: LIST"
	default:
		return "Auto-Generated Payload"
	}
}

// EvaluateConfidence converts a numeric score (0-100) to a confidence level string
func EvaluateConfidence(score float64) string {
	if score >= 80 {
		return "CRITICAL"
	} else if score >= 60 {
		return "HIGH"
	} else if score >= 40 {
		return "MEDIUM"
	}
	return "LOW"
}

// GetHeuristicByName retrieves a specific heuristic rule by name
func GetHeuristicByName(name string) *HeuristicRule {
	for i := range ActiveRules {
		if ActiveRules[i].Name == name {
			return &ActiveRules[i]
		}
	}
	return nil
}

// RunAllHeuristics evaluates all endpoints against all heuristic rules and returns
// all matching actions (not just the best match per endpoint)
func RunAllHeuristics(endpoints []string, trafficHistory map[string]int, loot logic.LootSummary) []TacticalAction {
	var allActions []TacticalAction

	for _, endpoint := range endpoints {
		status := trafficHistory[endpoint]
		action := RunHeuristics(endpoint, status, loot)
		if action != nil {
			action.ID = len(allActions) + 1
			allActions = append(allActions, *action)
		}
	}

	utils.TacticalLog(fmt.Sprintf("[blue]HEURISTICS:[-] Evaluated %d endpoints, generated %d actions.", len(endpoints), len(allActions)))
	return allActions
}

// ScoreEndpoint returns the heuristic confidence score for an endpoint without generating an action
func ScoreEndpoint(endpoint string, status int, loot logic.LootSummary) float64 {
	maxScore := 0.0
	for _, rule := range ActiveRules {
		triggered, score, _ := rule.Evaluate(endpoint, status, loot)
		if triggered && score > maxScore {
			maxScore = score
		}
	}
	return maxScore
}
