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

	// Schema aligned with report levels
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
		pterm.Success.Println("Persistence Active [DATABASE ID: 1]")
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

		_, err := DB.Exec("INSERT INTO findings (phase, target, details, status) VALUES (?, ?, ?, ?)",
			f.Phase, f.Target, f.Details, f.Status)
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