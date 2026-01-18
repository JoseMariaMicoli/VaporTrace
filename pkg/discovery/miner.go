package discovery

import (
	"fmt"
	"net/http"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" // Added Persistence
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

func MineParameters(baseURL string, endpoint string, proxy string) {
	params := []string{"debug", "admin", "test", "dev", "internal", "config", "role"}
	client, _ := utils.GetClient(proxy)

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

			// PERSISTENCE HOOK: Log hidden parameter success
			db.LogQueue <- db.Finding{
				Phase:   "PHASE II: DISCOVERY",
				Target:  fullURL,
				Details: fmt.Sprintf("Potential Hidden Parameter: %s", p),
				Status:  "VULNERABLE",
			}
		}
	}
}