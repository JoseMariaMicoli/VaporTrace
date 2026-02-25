package discovery

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// FuzzPaths performs forced browsing (dirbusting)
func FuzzPaths(baseURL string, customList []string) {
	wordlist := Top100Paths
	if len(customList) > 0 {
		wordlist = customList
	}

	utils.TacticalLog(fmt.Sprintf("[cyan]FUZZ:[-] Path enumeration starting on %s (%d payloads)", baseURL, len(wordlist)))

	// Worker Pool
	jobs := make(chan string, len(wordlist))
	var wg sync.WaitGroup

	// 5 Concurrent threads
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				target := utils.JoinURL(baseURL, path)

				req, _ := http.NewRequest("GET", target, nil)
				logic.ApplyEvasion(req) // Evasion logic

				resp, err := logic.GlobalClient.Do(req)
				if err != nil {
					continue
				}
				resp.Body.Close()

				// Detection Logic: Not 404
				if resp.StatusCode != 404 {
					utils.LogMap(target, "Path-Fuzz", fmt.Sprintf("%d", resp.StatusCode))
					utils.TacticalLog(fmt.Sprintf("[green]FOUND:[-] %s [%d]", path, resp.StatusCode))
					logic.GlobalDiscovery.AddEndpoint(target) // Feed into analysis buffer

					// Auto-record finding
					utils.RecordFinding(db.Finding{
						Phase:   "PHASE II: DISCOVERY",
						Command: "fuzz-paths",
						Target:  target,
						Details: fmt.Sprintf("Hidden path discovered (Status: %d)", resp.StatusCode),
						Status:  "INFO",
					})
				}
			}
		}()
	}

	for _, p := range wordlist {
		jobs <- p
	}
	close(jobs)
	wg.Wait()
	utils.TacticalLog("[green]✓ FUZZ:[-] Path enumeration complete.")
}

// FuzzParams performs parameter discovery
func FuzzParams(targetURL string, customList []string) {
	wordlist := Top100Params
	if len(customList) > 0 {
		wordlist = customList
	}

	utils.TacticalLog(fmt.Sprintf("[cyan]FUZZ:[-] Parameter mining on %s (%d payloads)", targetURL, len(wordlist)))

	// Baseline Request (to compare against)
	baseReq, _ := http.NewRequest("GET", targetURL, nil)
	logic.ApplyEvasion(baseReq)
	baseResp, err := logic.GlobalClient.Do(baseReq)
	var baseLen int64 = 0
	if err == nil && baseResp != nil {
		baseLen = baseResp.ContentLength
		baseResp.Body.Close()
	}

	jobs := make(chan string, len(wordlist))
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for param := range jobs {
				// Construct URL: url?param=fuzz
				u, _ := url.Parse(targetURL)
				q := u.Query()
				q.Set(param, "VaporTrace") // Payload value
				u.RawQuery = q.Encode()

				req, _ := http.NewRequest("GET", u.String(), nil)
				logic.ApplyEvasion(req)

				resp, err := logic.GlobalClient.Do(req)
				if err != nil {
					continue
				}
				resp.Body.Close()

				// Anomaly Detection
				// 1. Status Code Difference?
				// 2. Significant Length Difference?
				if baseResp != nil && resp.StatusCode != baseResp.StatusCode {
					recordParam(targetURL, param, "Status Code Anomaly")
				} else if baseLen > 0 && resp.ContentLength != -1 {
					diff := resp.ContentLength - baseLen
					if diff < 0 {
						diff = -diff
					}
					if diff > 100 { // Arbitrary threshold for "significant"
						recordParam(targetURL, param, fmt.Sprintf("Size Anomaly (%d bytes)", diff))
					}
				}
			}
		}()
	}

	for _, p := range wordlist {
		jobs <- p
	}
	close(jobs)
	wg.Wait()
	utils.TacticalLog("[green]✓ FUZZ:[-] Parameter mining complete.")
}

func recordParam(url, param, reason string) {
	utils.TacticalLog(fmt.Sprintf("[green]FOUND:[-] Parameter '%s' on %s (%s)", param, url, reason))
	utils.LogMap(url+"?"+param, "Param-Fuzz", reason)

	utils.RecordFinding(db.Finding{
		Phase:   "PHASE II: DISCOVERY",
		Command: "fuzz-params",
		Target:  url,
		Details: fmt.Sprintf("Hidden parameter '%s' discovered via %s", param, reason),
		Status:  "VULNERABLE",
	})
}
