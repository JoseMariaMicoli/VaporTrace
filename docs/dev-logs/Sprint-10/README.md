/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

# Sprint 10: Hydra - TUI Dashboard & Autonomous Systems

**Status:** ✅ COMPLETE | **Version:** v2.0-Hydra | **Released:** January 2026

---

## 🎯 Sprint Overview

Sprint 10 marks the transition from CLI-only tool to enterprise-grade TUI dashboard. This sprint delivers the Hydra multi-pane terminal interface with real-time monitoring, tactical planning visualization, and autonomous attack execution. The TUI provides 7 tabs with integrated exploitation engines and live feedback.

**Slogan:** "From Command Line to Mission Control"

---

## 📋 Deliverables

### 10.1: Universal Target Function (Global Context) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/logic/environment.go`

**Features Delivered:**
- **Global Context Management** - Single target for all modules
- **Persistent State** - Target saved across sessions
- **Endpoint Mapping** - Discovered endpoints tied to target
- **Credential Association** - Auth tokens per target
- **Scope Management** - Include/exclude patterns
- **Multi-Target Support** - Switch between targets

**Target Context:**
```go
type GlobalContext struct {
    TargetURL      string              // Primary target (https://api.example.com)
    TargetHost     string              // Extracted host
    TargetPort     int                 // Port (80, 443, etc.)
    Scheme         string              // http or https
    
    // Discovered endpoints
    Endpoints      []APIEndpoint       // All found endpoints
    
    // Authentication
    Sessions       []AuthSession       // API keys, tokens, cookies
    DefaultAuth    string              // Default session to use
    
    // Scope
    IncludePattern []string            // Endpoints to test
    ExcludePattern []string            // Endpoints to skip
    
    // State
    StartedAt      time.Time           // Mission start
    CurrentPhase   string              // reconnaissance, exploitation, etc.
}
```

**Usage:**
```bash
> target https://api.example.com
[green]TARGET SET:[-] api.example.com

> map
[cyan]DISCOVERY:[-] Reconnaissance phase...
(endpoints mapped to this target)

> analyze
[cyan]PLANNING:[-] Tactical analysis using this target...
(attack plan generated)

> bola
[cyan]EXPLOITATION:[-] BOLA testing against this target...
```

**Status:** ✅ Production-ready with full context

---

### 10.2: Project Mosaic - Hydra TUI Dashboard ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/ui/dashboard.go`

**Features Delivered:**
- **Multi-Pane Layout** - 4 quadrants + footer
- **Tab Navigation** - F1-F7 for 7 tabs
- **Real-Time Updates** - Live feedback during attacks
- **Color Coding** - Status indicators (red/yellow/green)
- **Responsive Design** - Adapts to terminal size
- **Keyboard Navigation** - Vim and arrow keys

**Dashboard Layout:**
```
╔═══════════════════════════════════════════════════════════════════════════╗
║                         🕷️  VAPORTRACE 3.1-HYDRA                          ║
╠═════════════════════════════════════════════╦═════════════════════════════╣
║                                             ║                             ║
║  F1: TACTICAL FEED                          ║  F2: INTERCEPTOR/MAP       ║
║  [cyan]SYSTEM SYNC[□][━━━━━━━━━][-]           ║  ┌─────────────────────────┐ ║
║  [08:120:00] Target: api.example.com        ║  │ Discovered Endpoints:   │ ║
║  [08:120:05] [cyan]RECON:[-] 47 endpoints found  ║  │ GET  /api/users         │ ║
║  [08:120:10] [green]EXPLOIT:[-] BOLA started     ║  │ POST /api/auth/login    │ ║
║  [08:120:15] [yellow]INFO:[-] 12 findings so far  ║  │ PUT  /api/admin/config  │ ║
║  [08:120:20] [red]ERROR:[-] Rate limited (429)   ║  │ DELETE /api/users/{id}  │ ║
║                                             ║  └─────────────────────────┘ ║
║                                             ║                             ║
╠═════════════════════════════════════════════╬═════════════════════════════╣
║                                             ║                             ║
║  F3: LOOT VAULT                             ║  F4: FINDINGS              ║
║  ┌──────────────────────────────────────┐   ║  ┌──────────────────────────┐║
║  │ Type          │ Count │ Confidence  │   ║  │ Severity │ Count │ Status││
║  ├──────────────────────────────────────┤   ║  ├──────────────────────────┤║
║  │ API Keys      │  8    │ 95%         │   ║  │ CRITICAL │  12   │ ✓     ││
║  │ JWT Tokens    │  12   │ 98%         │   ║  │ HIGH     │  18   │ ✓     ││
║  │ Database      │  4    │ 87%         │   ║  │ MEDIUM   │  7    │ ✓     ││
║  │ AWS Creds     │  2    │ 100%        │   ║  │ LOW      │  2    │ ✓     ││
║  │ Infrastructure│  21   │ 92%         │   ║  └──────────────────────────┘║
║  └──────────────────────────────────────┘   ║                             ║
║                                             ║                             ║
╠═════════════════════════════════════════════════════════════════════════════╣
║ [SYSTEM SYNC ░░░░░] [08:120:25] [⚡INTERCEPTOR:ON] [green]🧠NEURO:ON[-]    ║
║ [cyan]AGGRESSIVE[-] [yellow]1.8x[-] [green]J[-][green]T[-][green]B[-][green]O[-][green]E[-]                    ║
╚═════════════════════════════════════════════════════════════════════════════╝

F1: TACTICAL FEED    F2: MAP/INTERCEPTOR   F3: LOOT VAULT    F4: FINDINGS
F5: PLAN             F6: SESSIONS          F7: REPORT        Ctrl+H: Help
```

**Status:** ✅ Production-ready with full TUI

---

### 10.2.1: Terminal Multi-Pane (Quadrants + Tab Switching) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/ui/dashboard.go`

**Features Delivered:**
- **4 Quadrants** - Tactical Feed, Map, Loot, Findings
- **7 Full Tabs** - F1-F7 navigation
- **Real-Time Updates** - 250ms refresh rate
- **Responsive Layout** - Adjusts to terminal size
- **Fullscreen Support** - Expand any pane
- **Split Pane Navigation** - Arrow keys to move between panes

**Quadrant Functions:**

1. **F1: Tactical Feed** - Live command execution log
   ```
   - Real-time output from running commands
   - Color-coded: GREEN (success), YELLOW (info), RED (error)
   - Scrollable history (1000+ lines)
   - Search functionality
   ```

2. **F2: Interceptor/Map** - Discovered endpoints and interception
   ```
   - List of discovered endpoints (from map)
   - HTTP method, path, parameters
   - Click to inspect or test
   - Real-time updates during reconnaissance
   ```

3. **F3: Loot Vault** - Captured secrets and credentials
   ```
   - API keys, tokens, database passwords
   - Confidence scores for each item
   - Count by type
   - Export functionality
   ```

4. **F4: Findings** - Vulnerabilities discovered
   ```
   - Organized by severity (CRITICAL, HIGH, MEDIUM, LOW)
   - Attack vector and proof
   - Remediation suggestions
   - Status tracking
   ```

5. **F5: Tactical Plan** - Attack orchestration
   ```
   - Queued and executed actions
   - Status of each step
   - Undo/rollback options
   - Dependency visualization
   ```

6. **F6: Sessions** - Authentication management
   ```
   - Active API keys and tokens
   - User/role information
   - Session timing and expiration
   - Add/remove credentials
   ```

7. **F7: Report** - Real-time reporting
   ```
   - Live report generation
   - NIST/OWASP framework mapping
   - Export as PDF/Markdown/JSON
   - Executive summary
   ```

**Status:** ✅ Production-ready with all 7 tabs

---

### 10.2.2: Legacy Shell Fallback (CLI Flag Logic) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `main.go`

**Features Delivered:**
- **--shell Flag** - Force CLI mode instead of TUI
- **--batch Mode** - Non-interactive scripting
- **--headless** - No UI for automation
- **Automatic Detection** - TUI when interactive, CLI when piped
- **Compatibility** - All commands work in both modes
- **Graceful Degradation** - TUI features disabled in CLI

**Command Examples:**
```bash
# TUI mode (default)
$ vaportrace
# Opens interactive Hydra dashboard

# CLI mode (explicit)
$ vaportrace --shell
> target https://api.example.com
> map
> analyze
> exit

# Headless/batch mode
$ vaportrace --headless --batch << EOF
target https://api.example.com
map
analyze
bola
report --format pdf
EOF

# Piped input (auto-detects CLI)
$ cat attack_plan.txt | vaportrace
```

**Status:** ✅ Production-ready with full fallback

---

### 10.3: Contextual Aggregator & Information Gathering ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/engine/core.go`

**Features Delivered:**
- **End-to-End Reconnaissance** - map command orchestrates all discovery
- **Endpoint Aggregation** - Combines Swagger, scraper, parameter mining
- **Information Correlation** - Links endpoints to vulnerabilities
- **Attack Surface Map** - Visual representation of endpoints
- **Automated Pipeline** - Discovery-to-exploitation flow

**Pipeline Flow:**
```
1. Target set (https://api.example.com)
   ↓
2. Reconnaissance (map command)
   ├─ Try to fetch Swagger/OpenAPI spec
   ├─ Scrape JavaScript for endpoints
   ├─ Mine for hidden parameters
   ├─ Detect API versions
   └─ Populate F2 Map table
   ↓
3. Information gathering complete
   → 47 endpoints discovered
   → Ready for exploitation phase
```

**Status:** ✅ Production-ready with full aggregation

---

### 10.4: Tactical Interceptor (F2 Modal Manipulation) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/ui/dashboard.go`

**Features Delivered:**
- **Request Interception Modal** - Modify requests mid-flight
- **Response Inspection** - View full request/response
- **Manual Testing** - Edit and resend modified requests
- **Parameter Fuzzing** - In-modal payload editing
- **Response Comparison** - Side-by-side diff
- **Keyboard Controls** - Vim-style navigation

**Interceptor Modal:**
```
┌───────────────────────────────────────────────────────────────────┐
│                    REQUEST INTERCEPTOR                             │
├───────────────────────────────────────────────────────────────────┤
│ Method: [POST]  URL: /api/users                                   │
│                                                                   │
│ Headers:                                                          │
│ Authorization: Bearer eyJhbGc...                                 │
│ Content-Type: application/json                                   │
│                                                                   │
│ Body:                                                             │
│ {                                                                │
│   "name": "John",                                               │
│   "email": "john@example.com",                                  │
│   "admin": false                        ← Edit: Change to true   │
│ }                                                                │
│                                                                   │
│ Actions: [Send] [Drop] [Edit] [Clone] [Cancel]                 │
└───────────────────────────────────────────────────────────────────┘
```

**Status:** ✅ Production-ready with full editing

---

### 10.5: AI Base Integration (Heuristic Brain) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/ai/client.go`

**Features Delivered:**
- **Groq AI Integration** - Fast LLM for payload generation
- **Payload Mutation** - AI-driven exploitation attempts
- **Fallback Logic** - Local Ollama if external LLM fails
- **Prompt Engineering** - Security-focused prompts
- **Token Budgeting** - Cost-effective API usage
- **Streaming Response** - Real-time token generation

**AI Integration:**
```bash
> neuro on
[cyan]NEURO ENGINE:[-] Connecting to Groq AI...
[green]CONNECTED:[-] Using Groq Mixtral-8x7b

> neuro-gen 10
[cyan]AI:[-] Generating 10 payloads for /api/users/{id}
[08:125:00] Payload 1: {"id": "1' OR '1'='1"}
[08:125:01] Payload 2: {"id": "1; DROP TABLE users;--"}
[08:125:02] Payload 3: {"id": "../../../etc/passwd"}
...
[green]GENERATED:[-] 10 high-entropy payloads ready for testing
```

**Status:** ✅ Production-ready with AI-driven attacks

---

### 10.6: AI Payload Generation & Autonomous Fuzzing ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/ai/client.go`

**Features Delivered:**
- **Autonomous Payload Generation** - AI creates attack payloads
- **Vulnerability-Aware** - Payloads target specific vulns
- **Adaptive Fuzzing** - Learns from responses
- **Mutation Engine** - Combines successful variations
- **Confidence Scoring** - Rates payload effectiveness
- **Recursive Testing** - Chains exploits together

**Autonomous Exploitation Workflow:**
```bash
> analyze
[cyan]PLANNING:[-] Generating attack plan...
[green]Found 3 BOLA endpoints: /api/users/{id}, /api/orders/{order_id}, etc.

> commit
[cyan]EXECUTING:[-] Running tactical plan...
[08:130:00] Action 1: BOLA /api/users/
  → AI generates 20 payloads
  → Tests IDs: 1, 2, 100, admin, ..., 999
  → Finds: Users 50, 123, 456 accessible
  → Confidence: 95%

[08:130:15] Action 2: BFLA /api/admin/settings
  → AI generates 15 payloads
  → Tests methods: GET, POST, PUT, DELETE, PATCH
  → Finds: DELETE accessible via PUT
  → Confidence: 87%

[08:130:30] Action 3: BOPLA /api/profile
  → AI generates 25 payloads
  → Tests properties: admin, role, verified, email_verified
  → Finds: Can set admin=true, verified=true
  → Confidence: 92%

[green]PLAN COMPLETE:[-] 12 findings with 94% average confidence
```

**Status:** ✅ Production-ready with full autonomy

---

## 🔄 Current Status by Sub-Phase

| Sub-Phase | Deliverable | Status | Completion |
|-----------|-------------|--------|------------|
| **10.1** | Universal Target | ✅ DONE | 100% |
| **10.2** | Hydra TUI Dashboard | ✅ DONE | 100% |
| **10.2.1** | Multi-Pane Quadrants | ✅ DONE | 100% |
| **10.2.2** | Shell Fallback | ✅ DONE | 100% |
| **10.3** | Information Gathering | ✅ DONE | 100% |
| **10.4** | Tactical Interceptor | ✅ DONE | 100% |
| **10.5** | AI Integration | ✅ DONE | 100% |
| **10.6** | Autonomous Fuzzing | ✅ DONE | 100% |

---

## 📊 Code Metrics

| Metric | Value |
|--------|-------|
| **TUI Tabs** | 7 (F1-F7) |
| **Quadrants** | 4 (Tactical, Map, Loot, Findings) |
| **Dashboard Lines** | ~2000 LOC |
| **Refresh Rate** | 250ms |
| **AI Models** | Groq + Ollama fallback |
| **Keyboard Shortcuts** | 19+ keybindings |

---

## 🎓 Architecture Decisions

### Hydra TUI Design
- Multi-pane layout maximizes information density
- Real-time updates for live monitoring
- F1-F7 tabs for comprehensive feature access
- Color coding for quick status assessment

### CLI Fallback Strategy
- Auto-detection (interactive vs. piped)
- Explicit flags for different modes
- Full feature parity between UI and CLI
- Batch/headless for automation

### AI Integration
- Groq for speed and cost efficiency
- Ollama fallback for offline capability
- Security-focused prompt engineering
- Token budgeting for cost control

---

## 🚀 Next Steps

Sprint 11 focuses on autonomy hardening:
- Dynamic Dependency Injection (DDI)
- State-machine driven selection
- Lateral movement within API subnets
- Race condition fixes and synchronization

---

## 📚 References

- **TUI Framework:** https://github.com/rivo/tview
- **Groq API:** https://console.groq.com/docs/quickstart
- **Ollama:** https://ollama.ai/
- **Terminal Multiplexing:** https://en.wikipedia.org/wiki/Terminal_multiplexer
