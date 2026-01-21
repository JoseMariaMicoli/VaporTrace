package discovery

import (
	"fmt"
	"net/http"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic" // Standardized to Logic
)

func MineParameters(baseURL string, endpoint string, proxy string) {
	params := []string{"debug", "admin", "test", "dev", "internal", "config", "role"}
	
	// Use the GlobalClient managed by the logic package
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
			fmt.Printf("    [!] Potential Hidden Param: %s (Status: %d)\n", p, resp.StatusCode)

			db.LogQueue <- db.Finding{
				Phase:   "PHASE II: DISCOVERY",
				Target:  fullURL,
				Details: fmt.Sprintf("Potential Hidden Parameter: %s", p),
				Status:  "VULNERABLE",
			}
		}
	}
}