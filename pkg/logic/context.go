package logic

import (
	"net/http"
	"sync"
)

// GlobalClient is declared here ONCE for the whole 'logic' package
var GlobalClient *http.Client

type MissionDiscovery struct {
	Endpoints []string
	mu        sync.Mutex
}

var GlobalDiscovery = &MissionDiscovery{
	Endpoints: []string{},
}

func (md *MissionDiscovery) AddEndpoint(path string) {
	md.mu.Lock()
	defer md.mu.Unlock()
	for _, e := range md.Endpoints {
		if e == path { return }
	}
	md.Endpoints = append(md.Endpoints, path)
}

func init() {
	// Initialize it here
	GlobalClient = &http.Client{}
}

func SetGlobalClient(client *http.Client) {
	GlobalClient = client
}