package logic

import (
	"net/http"
)

// GlobalClient is the shared tactical client used by all logic modules
// It defaults to a standard client but gets overwritten by the Shell's 'proxy' command
//var GlobalClient *http.Client

func init() {
	// Initialize with a default client to prevent nil pointer panics
	GlobalClient = &http.Client{}
}

// SetGlobalClient updates the client for all logic modules (called by UI)
func SetGlobalClient(client *http.Client) {
	GlobalClient = client
}