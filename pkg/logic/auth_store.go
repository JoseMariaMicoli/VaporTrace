package logic

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/pterm/pterm"
)

// SessionStore holds the runtime configuration for the attacker.
// Patched to include legacy fields and thread-safety for Sprint 10.1.
type SessionStore struct {
	mu            sync.RWMutex
	TargetURL     string
	AttackerToken string // Restored for BOLA/BFLA modules
	VictimToken   string // Restored for multi-session logic
	Threads       int    // Restored for engine concurrency control
}

// CurrentSession holds the active configuration with tactical defaults.
var CurrentSession = &SessionStore{
	TargetURL: "http://localhost",
	Threads:   10,
}

// SetGlobalTarget validates the URL and synchronizes the framework's focus.
// This is the core function for Task Force 10.1.
func (s *SessionStore) SetGlobalTarget(rawURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("malformed URL: %v", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("URL missing protocol or host")
	}

	s.TargetURL = rawURL
	pterm.Success.Printfln("Global Target Locked: %s", s.TargetURL)
	return nil
}

// GetTarget provides thread-safe access to the mission target.
func (s *SessionStore) GetTarget() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.TargetURL
}