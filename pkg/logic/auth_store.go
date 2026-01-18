package logic

// SessionStore holds tokens configured via the 'auth' command
type SessionStore struct {
	VictimToken   string
	AttackerToken string
}

// CurrentSession is the global state for the active CLI session
var CurrentSession = &SessionStore{}