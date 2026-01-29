package report

import (
	"fmt"
	"os"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/pterm/pterm"
)

// GenerateMissionDebrief compiles findings into a professional tactical report.
// Aligned with Hydra-Worm Documentation Standards.
func GenerateMissionDebrief() {
	pterm.DefaultHeader.WithFullWidth(false).Println("Compiling Strategic Mission Report")

	var startTime string
	db.DB.QueryRow("SELECT value FROM mission_state WHERE key = 'start_time'").Scan(&startTime)

	reportDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("VAPOR_DEBRIEF_%s.md", reportDate)
	
	f, err := os.Create(fileName)
	if err != nil {
		pterm.Error.Printf("FileSystem Error: %v\n", err)
		return
	}
	defer f.Close()

	// I. HEADER - RED TEAM OPERATIONAL FORMAT
	f.WriteString("# VAPORTRACE TACTICAL DEBRIEF\n")
	f.WriteString(fmt.Sprintf("> **OPERATIONAL STATUS:** COMPLETED\n"))
	f.WriteString(fmt.Sprintf("> **DATABASE ID:** 1 | **GEN TIME:** %s\n", time.Now().Format("15:04:05")))
	f.WriteString(fmt.Sprintf("> **START TIME:** %s\n\n", startTime))
	f.WriteString("---\n\n")

	// II. HARVESTED ARTIFACTS (DISCOVERY VAULT)
	// Integration with Phase 8.1 Loot Scanner
	f.WriteString("## I. HARVESTED ARTIFACTS (DISCOVERY VAULT)\n\n")
	f.WriteString("| TYPE | SOURCE | VALUE (REDACTED) | TIMESTAMP |\n")
	f.WriteString("| :--- | :--- | :--- | :--- |\n")
	
	logic.GetVaultLock().Lock()
	vaultData := logic.GetVault()
	if len(vaultData) == 0 {
		f.WriteString("| - | - | *NO LOOT DISCOVERED* | - |\n")
	} else {
		for _, item := range vaultData {
			f.WriteString(fmt.Sprintf("| %s | %s | `%s` | %s |\n", 
				item.Type, item.Source, redactValue(item.Value), time.Now().Format("15:04:05")))
		}
	}
	logic.GetVaultLock().Unlock()
	f.WriteString("\n---\n\n")

	// III. MISSION PHASES SUMMARY
	f.WriteString("## II. MISSION PHASES SUMMARY\n\n")
	
	phases := []string{
		"PHASE II: DISCOVERY",
		"PHASE III: AUTH LOGIC",
		"PHASE IV: INJECTION",
		"PHASE VIII: EXFILTRATION",
	}

	for _, phase := range phases {
		f.WriteString(fmt.Sprintf("### %s\n", phase))
		f.WriteString("| ATTACK VECTOR | TARGET ENDPOINT | RESULT | TIMESTAMP |\n")
		f.WriteString("| :--- | :--- | :--- | :--- |\n")

		rows, _ := db.DB.Query("SELECT details, target, status, timestamp FROM findings WHERE phase = ?", phase)
		
		hasData := false
		for rows.Next() {
			hasData = true
			var details, target, status, ts string
			rows.Scan(&details, &target, &status, &ts)
			f.WriteString(fmt.Sprintf("| %s | `%s` | **%s** | %s |\n", details, target, status, ts))
		}
		rows.Close()

		if !hasData {
			f.WriteString("| - | - | *NO LOGS FOUND* | - |\n")
		}
		f.WriteString("\n")
	}

	// IV. MITRE ATT&CK OPERATIONAL MAPPING
	f.WriteString("## III. ADVERSARY EMULATION MAPPING\n\n")
	f.WriteString("| TACTIC | TECHNIQUE | RESULT |\n")
	f.WriteString("| :--- | :--- | :--- |\n")
	f.WriteString("| Reconnaissance | T1595.002 | Active API Scanning Successful |\n")
	f.WriteString("| Privilege Escalation | T1548 | BOLA/BFLA Logic Probing Completed |\n")
	f.WriteString("| Credential Access | T1552.005 | Cloud Instance Metadata Harvested |\n")
	f.WriteString("| Exfiltration | T1041 | Ghost-Weaver Encrypted Outbound |\n\n")

	// V. DFIR RESPONSE TEMPLATE (NIST SP 800-61)
	f.WriteString("---\n## IV. DFIR RESPONSE GUIDANCE\n\n")
	f.WriteString("### 1. Detection & Analysis\n")
	f.WriteString("* **Network Artifacts:** Monitor for `X-VaporTrace-Signal` headers in proxy logs.\n")
	f.WriteString("* **Endpoint Artifacts:** Audit for processes masquerading as `kworker_system_auth`.\n")
	f.WriteString("\n### 2. Containment & Eradication\n")
	f.WriteString("* **Logic Patching:** Implement strict UUID validation and object-level permission checks.\n")
	f.WriteString("* **Exfiltration Block:** Sinkhole OOB callback domains used during SSRF probes.\n\n")

	// Footer
	f.WriteString("---\n")
	f.WriteString("**CONFIDENTIAL // HYDRA-WORM INTEGRITY PROTOCOL // RED TEAM USE ONLY**\n")
}

// redactValue masks sensitive strings to maintain report integrity without leaking secrets
func redactValue(val string) string {
	if len(val) <= 8 {
		return "****"
	}
	return val[:4] + "...." + val[len(val)-4:]
}