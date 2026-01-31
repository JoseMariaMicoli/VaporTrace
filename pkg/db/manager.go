package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pterm/pterm"
)

var (
	DB       *sql.DB
	LogQueue = make(chan Finding, 100)
	isClosed bool
	mu       sync.Mutex // Ensures thread-safe access to isClosed
)

// Finding represents a tactical discovery to be persisted
// PATCHED Phase 9.13: Added Framework Tagging Columns
type Finding struct {
	Phase    string
	Target   string
	Details  string
	Status   string
	OWASP_ID string // e.g., "API1:2023"
	MITRE_ID string // e.g., "T1548"
	NIST_Tag string // e.g., "DE.AE"
}

// InitDB initializes the SQLite persistence layer and mission schema
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./vaportrace.db")
	if err != nil {
		pterm.Error.Printf("Database connection error: %v\n", err)
		return
	}

	// Schema updated for Framework Compliance Columns
	schema := `
    CREATE TABLE IF NOT EXISTS findings (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        phase TEXT,
        target TEXT,
        details TEXT,
        status TEXT,
        owasp_id TEXT,
        mitre_id TEXT,
        nist_tag TEXT,
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
		// Log to console only if needed, otherwise handled by logger
		DB.Exec("INSERT OR REPLACE INTO mission_state (key, value) VALUES ('start_time', ?)", time.Now().Format("2006-01-02 15:04:05"))
	}

	// Reset state and start worker
	isClosed = false
	go StartAsyncWorker()
}

// StartAsyncWorker processes background writes
func StartAsyncWorker() {
	for f := range LogQueue {
		mu.Lock()
		closed := isClosed
		mu.Unlock()

		if closed {
			return
		}

		query := `INSERT INTO findings (phase, target, details, status, owasp_id, mitre_id, nist_tag) 
                  VALUES (?, ?, ?, ?, ?, ?, ?)`
		
		_, err := DB.Exec(query, f.Phase, f.Target, f.Details, f.Status, f.OWASP_ID, f.MITRE_ID, f.NIST_Tag)
		if err != nil {
			pterm.Debug.Printf("Async commit failed: %v\n", err)
		}
	}
}

// ResetDB purges all mission evidence for a fresh operation
func ResetDB() {
	if DB == nil {
		pterm.Error.Println("Database connection is not initialized. Run 'init_db' first.")
		return
	}

	spinner, _ := pterm.DefaultSpinner.Start("Wiping all records from findings table...")

	_, err := DB.Exec("DELETE FROM findings")
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to purge table: %v", err))
		return
	}

	_, err = DB.Exec("DELETE FROM mission_state")
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to purge mission state: %v", err))
		return
	}

	spinner.Success("Database purged successfully.")
}

// CloseDB gracefully terminates the SQLite connection and logging channel
func CloseDB() {
	mu.Lock()
	if isClosed {
		mu.Unlock()
		return
	}
	isClosed = true
	mu.Unlock()

	if DB != nil {
		// Close channel to drain worker
		close(LogQueue)

		err := DB.Close()
		if err != nil {
			pterm.Error.Printf("Failed to close database: %v\n", err)
		}
		pterm.Debug.Println("Persistence layer decommissioned safely.")
	}
}