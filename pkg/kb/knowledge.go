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

package kb

import (
	"fmt"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// AttackPattern represents a learned exploit vector
type AttackPattern struct {
	ID           int
	VulnType     string
	Payload      string
	SuccessCount int
	LastUsed     time.Time
}

// LearnPattern saves or updates a successful payload in the Knowledge Base.
// This is the "Short-Term Memory" converting to "Long-Term Memory".
func LearnPattern(vulnType, payload string) {
	if db.DB == nil {
		return
	}

	// Normalize
	vulnType = strings.ToUpper(strings.TrimSpace(vulnType))
	payload = strings.TrimSpace(payload)

	if payload == "" {
		return
	}

	// Upsert Logic: Insert new, or increment success_count if exists
	query := `
		INSERT INTO attack_patterns (vuln_type, payload, success_count, last_used)
		VALUES (?, ?, 1, CURRENT_TIMESTAMP)
		ON CONFLICT(vuln_type, payload) 
		DO UPDATE SET 
			success_count = success_count + 1,
			last_used = CURRENT_TIMESTAMP
	`

	_, err := db.DB.Exec(query, vulnType, payload)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]KB ERROR:[-] Failed to learn pattern: %v", err))
	} else {
		utils.TacticalLog(fmt.Sprintf("[green]KB:[-] Learned new %s pattern (Saved to Memory).", vulnType))
	}
}

// GetTopPayloads retrieves the most successful payloads for a specific vulnerability type.
// Used by the NeuroEngine to seed the context.
func GetTopPayloads(vulnType string, limit int) []string {
	if db.DB == nil {
		return []string{}
	}

	// If type is vague ("injection"), search broadly
	query := `SELECT payload FROM attack_patterns WHERE vuln_type LIKE ? ORDER BY success_count DESC LIMIT ?`
	searchType := "%" + vulnType + "%"
	if vulnType == "" {
		searchType = "%"
	}

	rows, err := db.DB.Query(query, searchType, limit)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	var payloads []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		payloads = append(payloads, p)
	}
	return payloads
}

// ListPatterns returns all stored patterns for display
func ListPatterns() ([]AttackPattern, error) {
	if db.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := db.DB.Query("SELECT id, vuln_type, payload, success_count, last_used FROM attack_patterns ORDER BY success_count DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AttackPattern
	for rows.Next() {
		var ap AttackPattern
		var ts string
		rows.Scan(&ap.ID, &ap.VulnType, &ap.Payload, &ap.SuccessCount, &ts)
		// Parse sqlite time string if needed, or just keep string.
		// For simplicity in display we ignore precise parsing here.
		results = append(results, ap)
	}
	return results, nil
}

// Stats returns KB statistics
func Stats() string {
	if db.DB == nil {
		return "DB Offline"
	}
	var count int
	db.DB.QueryRow("SELECT COUNT(*) FROM attack_patterns").Scan(&count)

	var topType string
	db.DB.QueryRow("SELECT vuln_type FROM attack_patterns ORDER BY success_count DESC LIMIT 1").Scan(&topType)
	if topType == "" {
		topType = "None"
	}

	return fmt.Sprintf("Total Patterns: %d | Top Vector: %s", count, topType)
}
