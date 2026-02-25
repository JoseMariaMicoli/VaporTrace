# First Run Guide

**Document:** 02_FIRST_RUN.md  
**Version:** 3.0+  
**Time:** 15 minutes  
**Prerequisites:** [Installation complete](01_INSTALLATION_SETUP.md)

---

## Your First VaporTrace Session

Let's walk through your first penetration test scenario step-by-step.

### Step 1: Launch VaporTrace (2 minutes)

```bash
cd /path/to/VaporTrace
./VaporTrace
```

**You should see:**
```
██╗   ██╗ █████╗ ██████╗  ██████╗ ██████╗ ████████╗██████╗  █████╗  ██████╗███████╗
██║   ██║██╔══██╗██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝██╔══██╗██╔══██╗██╔════╝██╔════╝
╚██╗ ██╔╝███████║██████╔╝██║   ██║██████╔╝   ██║   ██████╔╝███████║██║     █████╗  
 ╚████╔╝ ██╔══██║██╔═══╝ ██║   ██║██╔══██╗   ██║   ██╔══██╗██╔══██║██║     ██╔══╝  
  ╚██╔╝  ██║  ██║██║     ╚██████╔╝██║  ██║   ██║   ██║  ██║██║  ██║╚██████╗███████╗
   ╚═╝   ╚═╝  ╚═╝╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚══════╝

VAPOR/INT> _
```

The prompt `VAPOR/INT>` is ready for your first command!

### Step 2: Initialize Database (1 minute)

```bash
> init_db
```

**Expected output:**
```
[08:30:42] SYSTEM: Database initialized successfully.
[08:30:42] SYSTEM: Tables created (Endpoints, Findings, Loot, Sessions)
```

### Step 3: Set Your Target (1 minute)

Choose a test target. We'll use httpbin.org (a safe, public API for testing):

```bash
> target https://httpbin.org
```

**Expected output:**
```
[08:31:15] SYSTEM: Target set to https://httpbin.org
[08:31:15] PIPELINE & STATUS - TARGET: https://httpbin.org
```

**Look at the dashboard:**
- **Left side panel** shows "TARGET: https://httpbin.org" ✅

### Step 4: Discover Endpoints (3 minutes)

Run the automatic discovery (combines spider + Swagger parsing + JS scraping):

```bash
> map
```

**Expected output:**
```
[08:31:30] DISCOVER: Starting attack surface mapping...
[08:31:30] SPIDER: Crawling https://httpbin.org
[08:31:45] SPIDER: Found 15 endpoints
[08:32:00] SWAGGER: No Swagger/OpenAPI spec found (expected)
[08:32:05] SCRAPE: No JavaScript bundles (expected for simple API)
[08:32:05] DISCOVER: Mapping complete - 15 endpoints discovered
```

**Check Tab F2 (MAP):**
- Press **F2** to switch to the MAP tab
- You should see a table with discovered endpoints:
  ```
  TIMESTAMP | SOURCE | ENDPOINT              | META
  08:32:00  | Recon  | /robots.txt          | Static file
  08:32:01  | Recon  | /get                 | [GET]
  08:32:02  | Recon  | /post                | [POST]
  ...
  ```

### Step 5: Analyze for Vulnerabilities (3 minutes)

Use the tactical planner to generate attack scenarios:

```bash
> analyze
```

**Expected output:**
```
[08:32:15] TACTICAL_PLAN: Analyzing 15 endpoints...
[08:32:25] TACTICAL_PLAN: Generated 8 tactical actions
[08:32:25] ACTION_1: BOLA on /get [Confidence: MEDIUM]
[08:32:25] ACTION_2: BFLA on /post [Confidence: LOW]
[08:32:25] ACTION_3: EXHAUST on /status [Confidence: MEDIUM]
...
```

**Check Tab F5 (PLAN):**
- Press **F5** to view the planner
- You should see a table with tactical actions:
  ```
  ID | TYPE    | TARGET    | PAYLOAD | CONFIDENCE | STATUS
  1  | BOLA    | /get      | [id]    | MEDIUM     | PENDING
  2  | BFLA    | /post     | [role]  | LOW        | PENDING
  ...
  ```

### Step 6: Review Actions (2 minutes)

See all pending actions:

```bash
> list-plan
```

**Expected output:**
```
[08:32:35] TACTICAL_PLAN: Pending Actions (8 total)
ID: 1    TYPE: BOLA      TARGET: /get       STATUS: PENDING
ID: 2    TYPE: BFLA      TARGET: /post      STATUS: PENDING
ID: 3    TYPE: EXHAUST   TARGET: /status    STATUS: PENDING
...
```

### Step 7: Edit an Action (Optional - 2 minutes)

Customize a payload before execution:

```bash
> edit 1 /api/users/999
```

**Expected output:**
```
[08:32:50] ACTION: Action #1 payload updated to '/api/users/999'
```

### Step 8: Execute Attack (2 minutes)

Run all pending tactical actions:

```bash
> commit
```

**Expected output:**
```
[08:33:00] EXEC: Starting tactical action execution...
[08:33:01] EXEC: ACTION_1 BOLA /get - Attempting ID enumeration...
[08:33:05] EXEC: ACTION_1 BOLA /get - Found 3 objects [EXPLOITED]
[08:33:06] EXEC: ACTION_2 BFLA /post - Attempting privilege escalation...
[08:33:10] EXEC: ACTION_2 BFLA /post - No privilege vectors found [FAILED]
[08:33:11] EXEC: ACTION_3 EXHAUST /status - Testing resource limits...
[08:33:15] EXEC: ACTION_3 EXHAUST /status - Rate limited at 100 req/s [SUCCESS]
...
[08:33:30] EXEC: All actions completed.
```

**Watch Tab F5 & F6:**
- **F5:** Planner updates with "SUCCESS" or "FAILED" status
- **F6:** AI analysis if neuro engine is enabled

### Step 9: Review Findings (2 minutes)

Check what was discovered:

```bash
> list-plan
```

**Expected output:**
```
ID: 1    TYPE: BOLA    TARGET: /get    STATUS: SUCCESS    OBJECTS: 3
ID: 2    TYPE: BFLA    TARGET: /post   STATUS: FAILED     REASON: No admin role
ID: 3    TYPE: EXHAUST TARGET: /status STATUS: SUCCESS    LIMIT: 100/s
```

**Check Tab F3 (LOOT):**
- Press **F3** to view captured secrets
- You should see captured data:
  ```
  TYPE    | VALUE              | SOURCE
  EMAIL   | test@example.com   | /get response
  TOKEN   | eyJ0... (JWT)      | /post response
  ```

### Step 10: Generate Report (2 minutes)

Export your findings:

```bash
> report
```

**Expected output:**
```
[08:33:45] REPORT: Generating findings report...
[08:33:50] REPORT: Report generated as 'VAPORTRACE_PEN_TEST_20260208_0835.md'
[08:33:50] REPORT: File saved to ./reports/
```

**Check Tab F7 (REPORT):**
- Press **F7** to view the report in the editor
- You should see Markdown formatted findings

### Step 11: Exit Gracefully (1 minute)

```bash
> exit
```

**Or press:**
```
Esc
```

**You should see:**
```
Secure Shutdown Protocol?
[Yes] [No]
```

Select **Yes** to exit cleanly.

---

## Summary: What You Just Did

✅ **Initialized** a penetration test framework  
✅ **Set** a target API  
✅ **Discovered** 15 endpoints  
✅ **Generated** 8 tactical attack scenarios  
✅ **Executed** real exploitation attempts  
✅ **Captured** findings (loot, emails, tokens)  
✅ **Generated** a professional report  

**Total time:** 15 minutes from zero to report!

---

## Understanding the Dashboard

### The 7 Tabs

| Tab | Key | Purpose |
|-----|-----|---------|
| **F1 LOGS** | Tactical Feed | System messages & discovery logs |
| **F2 MAP** | Attack Surface | Discovered endpoints & metadata |
| **F3 LOOT** | Vault | Captured secrets, tokens, credentials |
| **F4 TRAFFIC** | Sniffer | HTTP requests & responses |
| **F5 PLAN** | Planner | Tactical actions & execution status |
| **F6 NEURO** | AI Engine | AI mutations & analysis (if enabled) |
| **F7 REPORT** | Debrief | Generated findings export |

### The Left Panel

Shows real-time pipeline status:
- **TARGET** - Your current scope
- **AUTH (ATK)** - Attacker token (for privilege escalation)
- **CONTEXTS** - Number of active attack contexts
- **PROXY** - Proxy status & pool size
- **INTERCEPTOR** - Interception mode status
- **NEURO-BRAIN** - AI engine status

---

## Quick Hotkeys

| Key | Function |
|-----|----------|
| **Ctrl+H** | Show all hotkeys in modal |
| **F1-F7** | Switch tabs |
| **Ctrl+I** | Toggle interceptor |
| **Page Up/Dn** | Scroll logs |
| **Esc** | Exit |

See [17_KEYBOARD_SHORTCUTS.md](17_KEYBOARD_SHORTCUTS.md) for complete reference.

---

## Next Steps

Now that you've completed your first session, explore:

1. **Understand modules better:** [06_EXPLOITATION.md](06_EXPLOITATION.md)
2. **Try request interception:** [08_INTERCEPTOR_MITM.md](08_INTERCEPTOR_MITM.md)
3. **Use AI for payloads:** [07_AI_NEURO_ENGINE.md](07_AI_NEURO_ENGINE.md)
4. **Build attack chains:** [09_ATTACK_CHAINS.md](09_ATTACK_CHAINS.md)
5. **Deep reconnaissance:** [05_RECONNAISSANCE.md](05_RECONNAISSANCE.md)

---

## Common First-Run Questions

**Q: Why did some actions FAIL?**  
A: httpbin.org is designed to be safe. It doesn't have privilege escalation vectors. FAILURES are expected and normal.

**Q: How do I test on a real target?**  
A: Just change the target: `> target https://your-api-here.com`

**Q: Can I run against localhost?**  
A: Yes: `> target http://localhost:8000`

**Q: How do I enable AI payloads?**  
A: See [07_AI_NEURO_ENGINE.md](07_AI_NEURO_ENGINE.md) for setup instructions.

**Q: Where are findings saved?**  
A: Reports go to `./reports/` - press F7 to view in editor.

---

**Next:** [03_UI_OVERVIEW.md](03_UI_OVERVIEW.md) →

