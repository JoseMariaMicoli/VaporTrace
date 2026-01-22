package logic

import (
	"regexp"
	"sync"
	"github.com/pterm/pterm"
)

type Finding struct {
	Type    string
	Value   string
	Source  string
}

var (
	Vault    []Finding
	vaultMux sync.Mutex

	// Automated PII & Secret Patterns
	Patterns = map[string]*regexp.Regexp{
		"AWS_KEY":       regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		"JWT_TOKEN":     regexp.MustCompile(`eyJ[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*`),
		"CREDIT_CARD":   regexp.MustCompile(`(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12})`),
		"EMAIL":         regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		"SLACK_TOKEN":   regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`),
		// Phase 8.2 Additions - Wrapped in MustCompile to fix type mismatch
		"DB_CONNECTION": regexp.MustCompile(`(postgres|mysql|mongodb):\/\/[a-zA-Z0-9]+:[a-zA-Z0-9]+@[a-z0-9.-]+:\d+\/\w+`),
		"STACK_TRACE":   regexp.MustCompile(`(?i)(stacktrace|exception|at\s+[\w\.]+\([\w\.]+\.java:\d+\)|File\s+".+\.py",\s+line\s+\d+)`),
		"METADATA_IP":   regexp.MustCompile(`169\.254\.169\.254`),
	}
)

func ScanForLoot(body string, url string) {
	vaultMux.Lock()
	defer vaultMux.Unlock()

	for label, re := range Patterns {
		matches := re.FindAllString(body, -1)
		for _, m := range matches {
			// Check for duplicates
			exists := false
			for _, v := range Vault {
				if v.Value == m { exists = true; break }
			}
			
			if !exists {
				Vault = append(Vault, Finding{Type: label, Value: m, Source: url})
				pterm.DefaultBasicText.WithStyle(pterm.NewStyle(pterm.FgBlack, pterm.BgLightYellow)).
					Printfln(" LOOT DISCOVERED [%s] at %s ", label, url)
			}
		}
	}
}