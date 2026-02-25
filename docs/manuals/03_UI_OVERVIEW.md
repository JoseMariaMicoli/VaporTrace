# UI Overview & Navigation

**Document:** manuals/03_UI_OVERVIEW.md  
**Version:** 3.0+  
**For:** New & intermediate users  
**Read Time:** 15 minutes

---

## Dashboard Overview

The VaporTrace dashboard is organized into **7 interactive tabs** + **1 command prompt** at the bottom. Each tab displays real-time data from active attack modules.

```
┌─────────────────────────────────────────────────────────────┐
│ VaporTrace Dashboard                                        │
├─────────────────────────────────────────────────────────────┤
│ F1:LOGS │ F2:MAP │ F3:LOOT │ F4:ANALYSIS │ F5:PLAN │ F6:NEURO │ F7:SETTINGS
├─────────────────────────────────────────────────────────────┤
│                     [TAB CONTENT AREA]                       │
│                                                               │
│                    Live data display (800x400)              │
│                                                               │
├─────────────────────────────────────────────────────────────┤
│ TARGET: https://api.example.com / STATUS: Ready            │
├─────────────────────────────────────────────────────────────┤
│ > _                                                          │
│ (Command input - type commands here)                        │
└─────────────────────────────────────────────────────────────┘
```

---

## Tab Guide (F1-F7)

### F1 - LOGS Tab
**Hotkey:** F1  
**Shows:** Real-time tactical feed  
**Updates:** Every 200ms (batched)

**What you'll see:**
```
[08:30:00] TARGET: https://api.example.com
[08:30:05] DISCOVER: Starting mapping...
[08:30:15] SPIDER: Found 23 endpoints
[08:30:25] SWAGGER: Found spec at /api/v1/swagger.json
[08:30:35] SCRAPE: Found 5 more endpoints in JS
[08:30:40] DISCOVER: Mapping complete (33 endpoints)
[08:30:45] ANALYZE: Generating tactical plan...
[08:30:50] PLAN: Generated 8 tactical actions
[08:31:00] EXEC: Starting action execution...
[08:31:10] ACTION_1: BOLA /api/users - SUCCESS
[08:31:20] LOOT: Captured 15 user objects
```

**Key Features:**
- Timestamps for each event
- Color-coded by severity (✓ success, ⚠️  warning, ✗ error)
- Scrollable history (100+ events retained)
- Search with Ctrl+F
- Export with Ctrl+X

**When to use:** Monitor real-time attack progress, debug issues

---

### F2 - MAP Tab
**Hotkey:** F2  
**Shows:** Discovered API endpoints table  
**Updates:** After map/swagger/scrape/mine commands

**Table Structure:**
```
Timestamp   │ Source   │ Endpoint              │ Method │ Params
──────────────────────────────────────────────────────────────
08:30:15    │ Spider   │ /api/users            │ GET    │ id, role
08:30:20    │ Spider   │ /api/users            │ POST   │ name, email
08:30:25    │ Swagger  │ /api/admin/settings   │ GET    │ -
08:30:30    │ Swagger  │ /api/admin/reset      │ PUT    │ action
08:30:35    │ Scrape   │ /api/billing/invoice  │ GET    │ id, user_id
08:30:40    │ Mine     │ /api/users?debug=1    │ GET    │ debug
08:30:45    │ Swagger  │ /api/v2/users        │ GET    │ page, limit
```

**Key Information:**
- **Source:** How endpoint was discovered (Spider, Swagger, Scrape, Mine)
- **Endpoint:** Full API path
- **Method:** HTTP verb (GET, POST, PUT, DELETE, PATCH)
- **Params:** Known query/path parameters

**Actions in F2:**
- **Ctrl+F:** Search for endpoints (find "/admin", "POST", etc)
- **Ctrl+A:** Select all endpoints
- **Ctrl+B:** Batch operations on selected
- **Arrow Keys:** Navigate rows
- **Enter:** Test selected endpoint

**When to use:** Review discovered attack surface before exploitation

---

### F3 - LOOT Tab
**Hotkey:** F3  
**Shows:** Captured secrets & credentials  
**Updates:** Real-time during exploitation

**What you'll see:**
```
Type       │ Value                      │ Source              │ Time
──────────────────────────────────────────────────────────────
EMAIL      │ user@example.com           │ /api/users          │ 08:31:20
JWT        │ eyJ0exAiOiJIUzI1NiIsInR... │ /api/profile        │ 08:31:25
PASSWORD   │ SecurePass123              │ /api/admin/reset    │ 08:31:30
API_KEY    │ sk-xxxxxxxxxxxxxxxx        │ Metadata (SSRF)     │ 08:31:35
PHONE      │ +1-555-1234567            │ /api/users/1        │ 08:31:40
CREDIT_CC  │ 4111-1111-1111-1111       │ /api/billing        │ 08:31:45
SESSION    │ PHPSESSID=abc123...        │ /api/auth           │ 08:31:50
```

**Key Information:**
- **Type:** Secret classification (EMAIL, JWT, PASSWORD, API_KEY, PHONE, CREDIT_CC, SESSION, etc)
- **Value:** Actual secret (obfuscated for sensitive data)
- **Source:** Endpoint/module where captured
- **Time:** When captured

**Actions in F3:**
- **Ctrl+F:** Search for specific secrets
- **Ctrl+C:** Copy secret to clipboard
- **Ctrl+A:** Select all loot
- **Ctrl+X:** Export loot to JSON
- **Arrow Keys:** Navigate items

**Security Note:** Loot stored locally in SQLite. Sensitive data marked with ⚠️

**When to use:** Review captured credentials during or after exploitation

---

### F4 - ANALYSIS Tab
**Hotkey:** F4  
**Shows:** Vulnerabilities and security findings  
**Updates:** After exploitation modules run

**What you'll see:**
```
ID  │ Type    │ Endpoint          │ Severity │ Confidence │ Description
────────────────────────────────────────────────────────────────────
1   │ BOLA    │ /api/users/[id]   │ HIGH     │ 95%        │ ID enumeration works
2   │ BFLA    │ /api/admin/*      │ HIGH     │ 80%        │ Admin endpoint accessible
3   │ BOPLA   │ /api/users        │ MEDIUM   │ 70%        │ Mass assignment on role
4   │ SSRF    │ /api/fetch        │ CRITICAL │ 100%       │ Can access 169.254.169.254
5   │ CONFIG  │ /.well-known/*    │ MEDIUM   │ 60%        │ CORS: Allow-All
6   │ AUTH    │ /api/auth         │ LOW      │ 40%        │ Weak password policy
```

**Key Information:**
- **Type:** Vulnerability type (BOLA, BFLA, BOPLA, SSRF, CONFIG, AUTH, etc)
- **Endpoint:** Where vulnerability found
- **Severity:** CRITICAL, HIGH, MEDIUM, LOW (CVSS scale)
- **Confidence:** Likelihood finding is valid (false positive check)
- **Description:** What was found

**Color Coding:**
- 🔴 CRITICAL - Immediate exploitation risk
- 🟠 HIGH - Significant security issue
- 🟡 MEDIUM - Should be reviewed
- 🟢 LOW - Informational/low-risk

**Actions in F4:**
- **Ctrl+F:** Search findings
- **Enter:** View detailed analysis
- **Ctrl+X:** Export findings to report

**When to use:** Review vulnerability scan results

---

### F5 - PLAN Tab
**Hotkey:** F5  
**Shows:** Tactical action queue and execution status  
**Updates:** During execution (real-time)

**What you'll see:**
```
ID  │ Status     │ Type  │ Target            │ Payload      │ Confidence │ Result
─────────────────────────────────────────────────────────────────────────
1   │ PENDING    │ BOLA  │ /api/users/[id]   │ Sequential   │ HIGH       │ -
2   │ SUCCESS    │ BFLA  │ /api/admin/reset  │ Escalate     │ HIGH       │ 3 objects
3   │ FAILED     │ BOPLA │ /api/users        │ Mass assign  │ MEDIUM     │ 403 Forbidden
4   │ PENDING    │ SSRF  │ /api/fetch        │ Metadata     │ HIGH       │ -
5   │ DROPPED    │ Exhaust │ /api/data       │ Rate limit   │ LOW        │ -
```

**Action Status:**
- **PENDING:** Queued, waiting to execute
- **RUNNING:** Currently executing
- **SUCCESS:** Completed with findings
- **FAILED:** Completed but no findings
- **DROPPED:** Marked to skip (not execute)

**Key Information:**
- **Type:** Attack module (BOLA, BFLA, BOPLA, SSRF, etc)
- **Target:** Endpoint being tested
- **Payload:** Attack strategy
- **Confidence:** AI confidence in this action
- **Result:** Findings captured (object count, error, etc)

**Actions in F5:**
- **Tab:** Navigate between actions
- **E:** Edit payload of selected action
- **D:** Drop action (mark to skip)
- **Arrow Keys:** Navigate list
- **Ctrl+F:** Search actions

**When to use:** Monitor tactical action progress during `commit`

---

### F6 - NEURO Tab
**Hotkey:** F6  
**Shows:** AI-generated payload mutations  
**Updates:** When neural engine active

**What you'll see:**
```
Mutation │ Payload                              │ Confidence │ Status
──────────────────────────────────────────────────────────────────
Base     │ /api/users/1                         │ 100%       │ Base
Variant1 │ /api/users/999999999                 │ 95%        │ Tested
Variant2 │ /api/users/1/../../admin             │ 85%        │ Tested
Variant3 │ /api/users/1?role=admin              │ 70%        │ Ready
Variant4 │ /api%2fusers%2f1                     │ 60%        │ Ready
Variant5 │ /api/users/1;admin=true              │ 50%        │ Ready
```

**Payload Variants:**
- **Base:** Original payload
- **Variant (N):** AI mutation attempt
- **Obfuscated:** URL encoding, path traversal, parameter pollution
- **Status:** Ready to test, Testing, Tested, Failed

**Neural Engine Info:**
- Learns from response patterns
- Generates alternative bypasses
- Tests WAF evasion techniques
- Requires: NEURO_API_KEY configured

**Actions in F6:**
- **Ctrl+C:** Copy payload variant
- **Ctrl+X:** Export all variants
- **Enter:** Test selected variant
- **Arrow Keys:** Navigate mutations

**Status Indicator at top:**
```
NEURO: Connected ✓ | Latency: 245ms | Model: llama-3.1-8b | Tokens: 1,234/10,000
```

**When to use:** See alternative payloads, WAF bypass techniques

---

### F7 - SETTINGS Tab
**Hotkey:** F7  
**Shows:** Runtime configuration and options  
**Updates:** When you change settings

**Categories:**

#### Global Settings
```
API Target:          https://api.example.com
Auth Token:          eyJ0exAiOiJIUzI1NiIs...
Proxy:               localhost:8080 (active)
Debug Mode:          ON
Verbose Logging:     ON
```

#### Module Settings
```
Enabled Modules:     BOLA, BFLA, BOPLA, SSRF, Exhaust, Audit, Probe
Timeout (seconds):   30
Retries:            3
Rate Limit Delay:   1000ms (1 sec)
Follow Redirects:    ON
SSL Verify:          OFF (for testing)
```

#### Neural Engine
```
Provider:            Groq
API Key:             [••••••••••••••••••] (set)
Model:               llama-3.1-8b
Enabled:             ON
Temperature:         0.7 (creativity)
Max Tokens:          2000
```

#### Database
```
Location:            ~/.vaportrace/vaportrace.db
Size:               2.3 MB
Records:            4,532 endpoints, 1,234 findings
Auto-backup:        Every 6 hours
```

#### Export Defaults
```
Report Format:       Markdown (can be PDF)
Include Loot:        Yes
Include Screenshots: No
Compress Archive:    Yes
```

**Actions in F7:**
- **Tab:** Navigate between sections
- **Enter:** Edit selected setting
- **Ctrl+S:** Save settings
- **Ctrl+R:** Reset to defaults
- **F1:** Help for current setting

**When to use:** Configure API keys, proxies, modules before testing

---

## Command Prompt Area

Located at bottom of dashboard:

```
TARGET: https://api.example.com | STATUS: Ready | PID: 2847
> _
```

**Features:**
- Type commands and press Enter
- Command history with Up/Down arrow keys
- Tab auto-completion
- Supports all 40+ CLI commands

**Common commands:**
```
> target https://api.example.com    (set target)
> map                               (discover endpoints)
> analyze                           (generate attack plan)
> list-plan                         (show pending actions)
> commit                            (execute attacks)
> help                              (show all commands)
> exit                              (quit)
```

---

## Tab Interactions

### Typical Workflow Tabs
```
1. F2 (MAP)     ← Start here - review what you'll attack
2. F1 (LOGS)    ← Monitor real-time progress
3. F5 (PLAN)    ← Review tactical actions
4. F3 (LOOT)    ← Check captured secrets
5. F4 (ANALYSIS)← Review vulnerabilities found
```

### During Exploitation
```
Start command:   > commit
Watch: F5 (PLAN) - see action status
Watch: F1 (LOGS) - see execution logs
Pause: Ctrl+I to intercept requests
Check: F6 (NEURO) for AI alternatives
Review: F3 (LOOT) for real-time captures
```

### Configuration
```
Before attack:   Open F7 (SETTINGS)
Set: API keys, proxy, modules
Save: Ctrl+S
Verify: Click "Test Connection"
Ready: Close F7 and run attack
```

---

## Dashboard Features

### Real-Time Updates

**200ms Update Cycle:**
- Logs refresh every 200ms (batched)
- Endpoints update when discovered
- Loot appends in real-time
- Plan status updates live
- Performance optimized (no cascade)

### Status Bar

**Bottom left shows:**
```
TARGET: https://api.example.com
```

**Bottom right shows:**
```
STATUS: Ready | PROC: 3 | MEM: 256MB | UPTIME: 05:23
```

- **STATUS:** Ready/Running/Paused
- **PROC:** Active processes/goroutines
- **MEM:** Memory usage
- **UPTIME:** How long dashboard running

### Search & Filter

**Press Ctrl+F in any tab:**
- Logs: Search log messages
- Map: Find endpoints by path/method
- Loot: Search credential type/value
- Analysis: Find vulnerability type
- Plan: Search action type/target

---

## Keyboard Navigation Reference

| Key | Action | Context |
|-----|--------|---------|
| F1-F7 | Switch tabs | Dashboard |
| Tab | Next tab | Navigation |
| Shift+Tab | Previous tab | Navigation |
| Ctrl+H | Help modal | Anywhere |
| Ctrl+I | Interceptor | Anywhere |
| Ctrl+F | Search | Current tab |
| Ctrl+D | Debug mode | Anywhere |
| Ctrl+S | Save settings | Settings tab |
| Ctrl+X | Export | Current tab |
| Page Up | Scroll up | Lists/Modals |
| Page Down | Scroll down | Lists/Modals |
| Esc | Close modal | Modals |

---

## Performance Optimization

### Dashboard was optimized to prevent cascading updates:
- **Fixed:** Timer no longer calls updateTabs (caused cascading collapse)
- **Batch:** 200ms batching reduces UI redraws
- **Channel:** Data flows through dedicated channels (MapDataChan, LootDataChan, etc)
- **Buffer:** 100-item buffer per channel
- **Result:** Smooth 60 FPS, no lag even with 1000+ endpoints

**See:** [dev-logs/06_TUI_RENDERING.md](../dev-logs/06_TUI_RENDERING.md) for technical details

---

## Tips & Tricks

### Rapid Workflow
1. Open F2 first (see what you'll attack)
2. Use F1-F7 for instant tab switching (no menu)
3. Ctrl+H to remember hotkeys
4. Ctrl+F to search instead of scrolling

### During Active Attacks
1. F1 stays visible (monitor logs)
2. F5 shows action status
3. Ctrl+I intercepts if needed
4. F3 shows loot captured in real-time

### Review Findings
1. F4 shows vulnerabilities (click for details)
2. F3 shows loot captured
3. Ctrl+X to export all
4. `report` command generates final report

### Debug Issues
1. F7 to check configuration
2. Ctrl+D to enable debug mode
3. F1 will show verbose logs
4. Ctrl+F to search for errors

---

See also: [02_FIRST_RUN.md](02_FIRST_RUN.md) for step-by-step walkthrough, [17_KEYBOARD_SHORTCUTS.md](17_KEYBOARD_SHORTCUTS.md) for all hotkeys.

