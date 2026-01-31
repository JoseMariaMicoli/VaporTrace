package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB       *sql.DB
	LogQueue = make(chan Finding, 500)
	isClosed bool
	mu       sync.Mutex
)

// Finding represents a tactical discovery to be persisted
type Finding struct {
	Phase      string
	Target     string
	Details    string
	Status     string
	OWASP_ID   string
	MITRE_ID   string
	NIST_Tag   string
	CVE_ID     string
	CVSS_Score string
}

// InitDB initializes the connection and ensures schema integrity.
// IT DOES NOT SEED DATA.
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./vaportrace.db")
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		return
	}

	// 1. Create Base Schema (Idempotent)
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
	DB.Exec(schema)

	// 2. FORCE SCHEMA MIGRATION (For Task 12 Compliance)
	// We blindly try to add columns. If they exist, SQLite returns an error which we ignore.
	// This ensures old DB files are instantly compatible with the new Report Generator.
	migrationQueries := []string{
		"ALTER TABLE findings ADD COLUMN cve_id TEXT DEFAULT '-';",
		"ALTER TABLE findings ADD COLUMN cvss_score TEXT DEFAULT '0.0';",
	}
	for _, q := range migrationQueries {
		DB.Exec(q)
	}

	// 3. Set Session Start Time (Only if not exists)
	DB.Exec("INSERT OR IGNORE INTO mission_state (key, value) VALUES ('start_time', ?)", time.Now().Format("2006-01-02 15:04:05"))

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

		// Sanitize inputs to prevent NULL constraints
		cve := f.CVE_ID
		if cve == "" {
			cve = "-"
		}
		cvss := f.CVSS_Score
		if cvss == "" {
			cvss = "0.0"
		}
		nist := f.NIST_Tag
		if nist == "" {
			nist = "N/A"
		}

		query := `INSERT INTO findings (phase, target, details, status, owasp_id, mitre_id, nist_tag, cve_id, cvss_score) 
                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err := DB.Exec(query, f.Phase, f.Target, f.Details, f.Status, f.OWASP_ID, f.MITRE_ID, nist, cve, cvss)
		if err != nil {
			// Fail silently to preserve UI stability
		}
	}
}

// ResetDB completely wipes the database for a fresh mission.
func ResetDB() {
	if DB == nil {
		return
	}

	// Atomic Wipe
	tx, _ := DB.Begin()
	tx.Exec("DELETE FROM findings")
	tx.Exec("DELETE FROM sqlite_sequence WHERE name='findings'") // Reset ID to 1
	tx.Exec("DELETE FROM mission_state")
	tx.Exec("INSERT INTO mission_state (key, value) VALUES ('start_time', ?)", time.Now().Format("2006-01-02 15:04:05"))
	tx.Commit()
}

func CloseDB() {
	mu.Lock()
	if isClosed {
		mu.Unlock()
		return
	}
	isClosed = true
	mu.Unlock()

	if DB != nil {
		close(LogQueue)
		DB.Close()
	}
}
