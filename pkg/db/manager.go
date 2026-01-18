package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pterm/pterm"
)

var DB *sql.DB
var LogQueue = make(chan Finding, 100)

// Finding represents a tactical discovery to be persisted
type Finding struct {
	Phase   string
	Target  string
	Details string
	Status  string
}

// InitDB initializes the SQLite persistence layer and mission schema
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./vaportrace.db")
	if err != nil {
		pterm.Error.Printf("Database connection error: %v\n", err)
		return
	}

	// Schema aligned with report levels: P1 (Persistence) to P4 (Lateral Movement)
	schema := `
	CREATE TABLE IF NOT EXISTS findings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		phase TEXT,
		target TEXT,
		details TEXT,
		status TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS mission_state (
		key TEXT PRIMARY KEY,
		value TEXT
	);`

	_, err = DB.Exec(schema)
	if err != nil {
		pterm.Error.Printf("Failed to initialize schema: %v\n", err)
	} else {
		// Log mission start following Apex report style
		pterm.Success.Println("Persistence Active [DATABASE ID: 1]")
		DB.Exec("INSERT OR REPLACE INTO mission_state (key, value) VALUES ('start_time', ?)", time.Now().Format("2006-01-02 15:04:05"))
	}
}

// StartAsyncWorker processes background writes (Go routine based Async I/O)
func StartAsyncWorker() {
	for f := range LogQueue {
		_, err := DB.Exec("INSERT INTO findings (phase, target, details, status) VALUES (?, ?, ?, ?)",
			f.Phase, f.Target, f.Details, f.Status)
		if err != nil {
			pterm.Debug.Printf("Async commit failed: %v\n", err)
		}
	}
}

// ResetDB purges all mission evidence for a fresh operation
func ResetDB() {
	_, err := DB.Exec("DROP TABLE IF EXISTS findings; DROP TABLE IF EXISTS mission_state;")
	if err != nil {
		pterm.Error.Printf("Reset failed: %v\n", err)
	}
	InitDB()
}

var isClosed bool // Guard variable

// CloseDB gracefully terminates the SQLite connection and logging channel
func CloseDB() {
	if isClosed {
		return // Prevent double-closing panic
	}
	
	if DB != nil {
		// 1. Close the channel first to stop the worker
		close(LogQueue) 
		
		// 2. Close the database connection
		err := DB.Close()
		if err != nil {
			pterm.Error.Printf("Failed to close database: %v\n", err)
		}
		
		isClosed = true
		pterm.Debug.Println("Persistence layer decommissioned safely.")
	}
}