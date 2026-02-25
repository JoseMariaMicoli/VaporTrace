package discovery

import (
	"fmt"
	"net/http"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

func MineParameters(baseURL string, endpoint string, proxy string) {
	params := []string{"debug", "admin", "test", "dev", "internal", "config", "role"}

	client := logic.GlobalClient

	for _, p := range params {
		fullURL := fmt.Sprintf("%s%s?%s=true", baseURL, endpoint, p)

		req, _ := http.NewRequest("GET", fullURL, nil)
		resp, err := client.Do(req)

		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusBadRequest {
			// Task 2: Log for UI
			utils.LogMap(fmt.Sprintf("%s?%s", endpoint, p), "Param Mining", fmt.Sprintf("%d", resp.StatusCode))

			utils.RecordFinding(db.Finding{
				Phase:    "PHASE II: DISCOVERY",
				Target:   fullURL,
				Details:  fmt.Sprintf("Potential Hidden Parameter: %s", p),
				Status:   "VULNERABLE",
				OWASP_ID: "API3:2023",
				MITRE_ID: "T1596", // Search Open Technical Databases
				NIST_Tag: "ID.RA",
			})
		}
	}
}
