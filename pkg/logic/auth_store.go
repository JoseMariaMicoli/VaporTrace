package logic

import (
	"fmt"
	"net/url"
	"sync"
)

// SessionStore holds the runtime configuration for the attacker.
type SessionStore struct {
	mu            sync.RWMutex
	TargetURL     string
	AttackerToken string
	VictimToken   string
	Threads       int
}

// CurrentSession holds the active configuration with tactical defaults.
var CurrentSession = &SessionStore{
	TargetURL: "http://localhost",
	Threads:   10,
}

// SetGlobalTarget validates the URL and locks it in the session.
// CHANGE: Removed pterm printing. Engine handles feedback now.
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
	return nil
}

// GetTarget provides thread-safe access to the mission target.
func (s *SessionStore) GetTarget() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.TargetURL
}
