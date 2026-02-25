package discovery

import (
	"net/http"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// SetGlobalClient updates the shared utility client.
// This ensures that if the UI calls this, it updates the same client used by miner/scraper.
func SetGlobalClient(client *http.Client) {
	utils.GlobalClient = client
}
