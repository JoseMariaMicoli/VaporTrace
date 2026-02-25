package enrichment

import (
	"fmt"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
)

// EnrichmentProfile defines the strict security metadata structure
type EnrichmentProfile struct {
	OWASP       string
	MitreID     string
	MitreTactic string
	NistControl string // NIST CSF v2.0
	CVE         string
	CVSS        float64
}

// sourceOfTruth holds the immutable mapping table
var sourceOfTruth = map[string]EnrichmentProfile{
	"bola": {
		OWASP:       "API1:2023 Broken Object Level Auth",
		MitreID:     "T1594",
		MitreTactic: "Collection",
		NistControl: "PR.AC-03",
		CVE:         "CVE-2024-BOLA-GENERIC",
		CVSS:        9.1,
	},
	"weaver": {
		OWASP:       "API2:2023 Broken Authentication",
		MitreID:     "T1552.004",
		MitreTactic: "Credential Access",
		NistControl: "PR.AC-01",
		CVE:         "CVE-2024-AUTH-BYPASS",
		CVSS:        8.5,
	},
	"bopla": {
		OWASP:       "API3:2023 Broken Object Property Level Auth",
		MitreID:     "T1592.001",
		MitreTactic: "Privilege Escalation",
		NistControl: "PR.DS-01",
		CVE:         "CVE-2022-23131",
		CVSS:        9.8,
	},
	"exhaust": {
		OWASP:       "API4:2023 Unrestricted Resource Consumption",
		MitreID:     "T1499.004",
		MitreTactic: "Impact (DoS)",
		NistControl: "DE.AE-02",
		CVE:         "CVE-2023-44487",
		CVSS:        7.5,
	},
	"bfla": {
		OWASP:       "API5:2023 Broken Function Level Auth",
		MitreID:     "T1548.003",
		MitreTactic: "Privilege Escalation",
		NistControl: "PR.AC-05",
		CVE:         "CVE-2023-30533",
		CVSS:        8.2,
	},
	"pipeline": {
		OWASP:       "API6:2023 Unrestricted Access to Business Flows",
		MitreID:     "T1068",
		MitreTactic: "Persistence",
		NistControl: "PR.DS-02",
		CVE:         "CVE-2024-LOGIC-FLAW",
		CVSS:        7.8,
	},
	"ssrf": {
		OWASP:       "API7:2023 Server Side Request Forgery",
		MitreID:     "T1071.001",
		MitreTactic: "Command & Control",
		NistControl: "PR.DS-01",
		CVE:         "CVE-2021-26855",
		CVSS:        9.2,
	},
	"audit": {
		OWASP:       "API8:2023 Security Misconfiguration",
		MitreID:     "T1562.001",
		MitreTactic: "Defense Evasion",
		NistControl: "PR.PS-01",
		CVE:         "CVE-2024-AUDIT",
		CVSS:        5.4,
	},
	"map": {
		OWASP:       "API9:2023 Improper Inventory Management",
		MitreID:     "T1595.002",
		MitreTactic: "Reconnaissance",
		NistControl: "ID.AM-07",
		CVE:         "N/A",
		CVSS:        0.0,
	},
	"probe": {
		OWASP:       "API10:2023 Unsafe Consumption of APIs",
		MitreID:     "T1190",
		MitreTactic: "Initial Access",
		NistControl: "PR.DS-02",
		CVE:         "CVE-2024-PROBE",
		CVSS:        8.1,
	},
}

// EnrichFinding applies the Zero-Touch automated tagging logic.
func EnrichFinding(f *db.Finding) {
	// Normalize the command key (e.g., "BOLA" -> "bola")
	key := strings.ToLower(strings.TrimSpace(f.Command))

	if data, exists := sourceOfTruth[key]; exists {
		// Auto-populate strict framework fields
		f.OWASP_ID = data.OWASP
		f.MITRE_ID = data.MitreID
		f.MitreTactic = data.MitreTactic
		f.NistControl = data.NistControl
		f.CVE_ID = data.CVE
		f.CVSS_Numeric = data.CVSS

		// FIX: Correct float-to-string conversion
		// Backward compatibility for legacy string field
		if f.CVSS_Score == "" || f.CVSS_Score == "0.0" {
			f.CVSS_Score = fmt.Sprintf("%.1f", data.CVSS)
		}
	} else {
		// Fallback for custom or unknown commands
		if f.OWASP_ID == "" {
			f.OWASP_ID = "Unknown"
		}
		if f.MitreTactic == "" {
			f.MitreTactic = "Untriaged"
		}
		if f.NistControl == "" {
			f.NistControl = "N/A"
		}
	}
}
