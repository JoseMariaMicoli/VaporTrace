/*
Copyright (c) 2026 Jose Maria Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
Study
Modify
Use for internal security testing

You may NOT:
Offer as a commercial service
Sell derived competing products
*/

package logic

// ResetRuntimeState purges volatile mission state kept in memory.
// It should be called together with db.ResetDB() when a full session reset is intended.
func ResetRuntimeState() {
	// Discovery inventory.
	GlobalDiscovery.mu.Lock()
	GlobalDiscovery.Inventory = make(map[string]*EndpointEntry)
	GlobalDiscovery.mu.Unlock()

	// Loot vault.
	vaultMux.Lock()
	Vault = []Finding{}
	vaultMux.Unlock()

	// Flow state.
	ActiveFlow = []FlowStep{}
	FlowContext = make(map[string]string)

	// Data silo used by autonomous chains.
	GlobalDataSilo.mu.Lock()
	GlobalDataSilo.data = make(map[string]interface{})
	GlobalDataSilo.mu.Unlock()

	// Traffic history used by planner heuristics.
	trafficMu.Lock()
	TrafficHistory = make(map[string]int)
	trafficMu.Unlock()

	// Traffic vault ring buffer.
	GlobalVault.mu.Lock()
	GlobalVault.Buffer = make([]*TrafficEntry, GlobalVault.Capacity)
	GlobalVault.Lookup = make(map[string]*TrafficEntry)
	GlobalVault.Head = 0
	GlobalVault.Count = 0
	GlobalVault.mu.Unlock()
}
