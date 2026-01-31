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
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type BOLAContext struct {
	BaseURL       string
	VictimID      string
	AttackerID    string
	AttackerToken string
}

func ExecuteMassBOLA(concurrency int) {
	utils.TacticalLog("[cyan::b]PHASE 9.7: INDUSTRIALIZED BOLA ENGINE STARTED[-:-:-]")

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
		utils.TacticalLog("[yellow]No BOLA-vulnerable patterns detected in discovery phase.[-]")
		return
	}

	testIDs := []string{"1", "2", "10", "100", "101", "1000", "1337"}

	for _, t := range targets {
		ctx := &BOLAContext{
			BaseURL: CurrentSession.TargetURL + t,
		}
		ctx.MassProbe(testIDs, concurrency)
	}
	utils.TacticalLog("[green::b]BOLA Engine Execution Completed.[-:-:-]")
}

func (b *BOLAContext) getResource(resourceID string, token string) (int, string, error) {
	u, err := url.Parse(b.BaseURL)
	if err != nil {
		return 0, "", err
	}

	target := u.String()

	re := regexp.MustCompile(`\{.*?\}`)
	if re.MatchString(target) {
		target = re.ReplaceAllString(target, resourceID)
	} else {
		u.Path = path.Join(u.Path, resourceID)
		target = u.String()
	}

	req, _ := http.NewRequest("GET", target, nil)

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	} else if CurrentSession.AttackerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", CurrentSession.AttackerToken))
	}

	req.Header.Set("User-Agent", "VaporTrace/2.1.0 (Phase 9.10 Industrialized)")

	resp, err := SafeDo(req, false, "BOLA-ENGINE")
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}

func (b *BOLAContext) MassProbe(idList []string, concurrency int) {
	idChan := make(chan string, concurrency)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idChan {
				instance := *b
				instance.VictimID = id
				instance.ProbeSilent()
			}
		}()
	}

	for _, id := range idList {
		idChan <- id
	}
	close(idChan)
	wg.Wait()
}

func (b *BOLAContext) ProbeSilent() {
	activeToken := b.AttackerToken
	if activeToken == "" {
		activeToken = CurrentSession.AttackerToken
	}

	code, body, err := b.getResource(b.VictimID, activeToken)
	if err != nil || code != 200 {
		return
	}

	lowerBody := strings.ToLower(body)
	if strings.Contains(lowerBody, "not found") || strings.Contains(lowerBody, "error") || len(body) < 2 {
		return
	}

	utils.RecordFinding(db.Finding{
		Phase:      "PHASE III: AUTH LOGIC",
		Target:     b.BaseURL,
		Details:    fmt.Sprintf("BOLA ID Swap Success: ID %s returned 200 OK", b.VictimID),
		Status:     "VULNERABLE",
		OWASP_ID:   "API1:2023",
		MITRE_ID:   "T1548",
		NIST_Tag:   "DE.AE",
		CVE_ID:     "CVE-202X-BOLA-AUTH",
		CVSS_Score: "7.5", // CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N
	})
}
