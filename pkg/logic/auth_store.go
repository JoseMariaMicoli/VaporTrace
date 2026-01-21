package logic

// SessionStore holds the runtime configuration for the attacker
type SessionStore struct {
	TargetURL     string
	AttackerToken string
	VictimToken   string
	Threads       int
}

// CurrentSession holds the active configuration
var CurrentSession = SessionStore{
	TargetURL: "http://localhost",
	Threads:   10,
}