package logic

import (
	"fmt"
	"net/http"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type MisconfigContext struct {
	TargetURL string
}

func (m *MisconfigContext) Audit() {
	utils.TacticalLog(fmt.Sprintf("[blue]API8 Audit: Scanning %s for Security Misconfigurations...[-]", m.TargetURL))

	req, _ := http.NewRequest("GET", m.TargetURL, nil)
	req.Header.Set("Origin", "https://evil-attacker.com")

	resp, err := GlobalClient.Do(req)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]Audit Error:[-] %v", err))
		return
	}
	defer resp.Body.Close()

	cors := resp.Header.Get("Access-Control-Allow-Origin")
	if cors == "*" || cors == "https://evil-attacker.com" {
		utils.RecordFinding(db.Finding{
			Phase:    "PHASE II: DISCOVERY",
			Target:   m.TargetURL,
			Details:  fmt.Sprintf("Weak CORS Policy: %s", cors),
			Status:   "VULNERABLE",
			OWASP_ID: "API8:2023",
			MITRE_ID: "T1562",
			NIST_Tag: "PR.IP",
		})
	}

	headers := []string{"Strict-Transport-Security", "Content-Security-Policy", "X-Content-Type-Options"}
	for _, h := range headers {
		if resp.Header.Get(h) == "" {
			utils.RecordFinding(db.Finding{
				Phase:    "PHASE II: DISCOVERY",
				Target:   m.TargetURL,
				Details:  fmt.Sprintf("Missing Header: %s", h),
				Status:   "WEAK CONFIG",
				OWASP_ID: "API8:2023",
				MITRE_ID: "T1592",
				NIST_Tag: "PR.IP",
			})
		}
	}

	m.TriggerVerboseError()
}

func (m *MisconfigContext) TriggerVerboseError() {
	req, _ := http.NewRequest("TRACE", m.TargetURL, nil)
	resp, err := GlobalClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		utils.RecordFinding(db.Finding{
			Phase:    "PHASE II: DISCOVERY",
			Target:   m.TargetURL,
			Details:  "Verbose Error / Stack Trace",
			Status:   "INFO",
			OWASP_ID: "API8:2023",
			MITRE_ID: "T1592",
			NIST_Tag: "DE.CM",
		})
	}
}
