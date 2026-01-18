package discovery

import (
	"fmt"
	"net/http"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// MineParameters checks an endpoint for hidden query parameters
func MineParameters(baseURL string, endpoint string, proxy string) {
	// High-value "spicy" parameters
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

		// If the server doesn't 404/400, it's worth a manual look in Burp
		if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusBadRequest {
			fmt.Printf("    [!] Potential Hidden Param: %s (Status: %d)\n", p, resp.StatusCode)
		}
	}
}