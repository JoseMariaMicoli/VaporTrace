/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

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
