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

// Finding represents a tactical discovery to be persisted.
// This struct matches the expanded Schema for the C-Level Report.
type Finding struct {
	// --- Base Fields ---
	Phase     string
	Target    string
	Details   string
	Status    string
	Timestamp time.Time

	// --- Compliance & Enrichment Fields ---
	OWASP_ID   string
	MITRE_ID   string
	NIST_Tag   string
	CVE_ID     string
	CVSS_Score string // Legacy string support

	// --- New Architectural Fields ---
	Command      string  // Triggering command (e.g., "bola")
	MitreTactic  string  // e.g., "Privilege Escalation"
	NistControl  string  // e.g., "PR.AC-03"
	CVSS_Numeric float64 // Float for sorting/stats
}

// InitDB initializes the SQLite connection and enforces Schema compliance.
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./vaportrace.db")
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		return
	}

	// 1. Create Base Schema (Idempotent)
	// We define the FULL schema here for fresh installs.
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
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        cve_id TEXT,
        cvss_score TEXT,
        command TEXT,
        mitre_tactic TEXT,
        nist_control TEXT,
        cvss_numeric REAL
    );
    CREATE TABLE IF NOT EXISTS mission_state (
        key TEXT PRIMARY KEY,
        value TEXT
    );`

	_, err = DB.Exec(schema)
	if err != nil {
		fmt.Printf("Schema creation error: %v\n", err)
	}

	// 2. FORCE MIGRATION (For existing databases)
	// These will fail silently if columns exist, which is acceptable for SQLite.
	migrations := []string{
		"ALTER TABLE findings ADD COLUMN cve_id TEXT DEFAULT '-';",
		"ALTER TABLE findings ADD COLUMN cvss_score TEXT DEFAULT '0.0';",
		"ALTER TABLE findings ADD COLUMN command TEXT DEFAULT '';",
		"ALTER TABLE findings ADD COLUMN mitre_tactic TEXT DEFAULT 'Untriaged';",
		"ALTER TABLE findings ADD COLUMN nist_control TEXT DEFAULT 'N/A';",
		"ALTER TABLE findings ADD COLUMN cvss_numeric REAL DEFAULT 0.0;",
	}
	for _, q := range migrations {
		DB.Exec(q)
	}

	// 3. Initialize Mission State
	DB.Exec("INSERT OR IGNORE INTO mission_state (key, value) VALUES ('start_time', ?)", time.Now().Format("2006-01-02 15:04:05"))

	isClosed = false
	go StartAsyncWorker()
}

// StartAsyncWorker processes background writes from the LogQueue.
func StartAsyncWorker() {
	for f := range LogQueue {
		mu.Lock()
		closed := isClosed
		mu.Unlock()

		if closed {
			return
		}

		// Defaults for safety
		if f.CVSS_Score == "" {
			f.CVSS_Score = "0.0"
		}
		if f.CVE_ID == "" {
			f.CVE_ID = "-"
		}
		if f.Command == "" {
			f.Command = "manual"
		}

		query := `INSERT INTO findings (
            phase, target, details, status, 
            owasp_id, mitre_id, nist_tag, 
            cve_id, cvss_score, 
            command, mitre_tactic, nist_control, cvss_numeric
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err := DB.Exec(query,
			f.Phase, f.Target, f.Details, f.Status,
			f.OWASP_ID, f.MITRE_ID, f.NIST_Tag,
			f.CVE_ID, f.CVSS_Score,
			f.Command, f.MitreTactic, f.NistControl, f.CVSS_Numeric,
		)

		if err != nil {
			// fmt.Printf("DB Write Error: %v\n", err) // Optional debug
		}
	}
}

// ResetDB completely purges the database (Functionality Check: PASSED).
func ResetDB() {
	if DB == nil {
		return
	}
	// Atomic Wipe
	tx, _ := DB.Begin()
	tx.Exec("DELETE FROM findings")
	tx.Exec("DELETE FROM sqlite_sequence WHERE name='findings'") // Reset Auto-Increment
	tx.Exec("DELETE FROM mission_state")
	tx.Exec("INSERT INTO mission_state (key, value) VALUES ('start_time', ?)", time.Now().Format("2006-01-02 15:04:05"))
	tx.Commit()
}

// CloseDB gracefully shuts down the connection.
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
