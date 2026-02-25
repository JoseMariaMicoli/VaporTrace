![VaporTrace Logo](../../assets/images/VaporTrace_Logo.png)

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

# Strategic Planning & HITL Orchestration

**Document:** manuals/04_STRATEGIC_PLANNING.md  
**Version:** 3.0+  
**For:** Intermediate users  
**Read Time:** 20 minutes

---

## Overview

**Human-in-the-Loop (HITL) Orchestration** is VaporTrace's core feature that enables intelligent, semi-automated attack planning and execution. Instead of running independent modules, you create a **tactical plan**, review it, edit it, and then execute it with full control.

```
┌──────────────┐
│ DISCOVER     │ (map/swagger/scrape/mine)
│ Endpoints    │
└──────┬───────┘
       │
       ▼
┌──────────────────────┐
│ ANALYZE              │ Generate tactical plan
│ Generate Plan        │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ REVIEW & EDIT        │ list-plan, edit, drop
│ Tactical Actions     │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ COMMIT               │ Execute actions
│ Execute Attack       │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│ ANALYZE RESULTS      │ F3 LOOT, F4 ANALYSIS
│ Captured Loot        │
└──────────────────────┘
```

---

## Step 1: Discovery Phase

### Purpose
Enumerate all API endpoints to build attack surface.

### Commands

#### `target <url>`
Set your target URL.

```bash
> target https://api.example.com
[08:30:00] TARGET: https://api.example.com

> target https://api.example.com:8443
[08:30:00] TARGET: https://api.example.com:8443
```

#### `map [url]` (Recommended)
Automatic reconnaissance - runs all discovery modules.

```bash
> map

Runs:
1. Spider: Crawl all accessible endpoints
2. Swagger: Check for /swagger, /api-docs, /openapi.json
3. Scrape: Extract endpoints from JavaScript
4. Mine: Brute-force hidden parameters

Results appear in F2 (MAP) tab
```

**Alternative: Individual modules**
```bash
> swagger https://api.example.com/openapi.json
> scrape https://api.example.com/js/main.js
> mine https://api.example.com/api/users
```

### Expected Output (F1 LOGS)
```
[08:30:05] DISCOVER: Starting mapping...
[08:30:10] SPIDER: Spidering https://api.example.com
[08:30:15] SPIDER: Found 23 endpoints
[08:30:20] SWAGGER: Checking for Swagger specs...
[08:30:22] SWAGGER: Found /swagger.json
[08:30:25] SWAGGER: Parsed 15 endpoints
[08:30:30] SCRAPE: Extracting from JavaScript...
[08:30:35] SCRAPE: Found 5 more endpoints
[08:30:40] MINE: Testing parameters...
[08:30:50] MINE: Found debug=true parameter
[08:30:55] DISCOVER: Complete! Total 33 endpoints
```

### View Results
Press **F2** to see MAP tab:
```
Timestamp │ Source  │ Endpoint           │ Method │ Params
──────────────────────────────────────────────────
08:30:15  │ Spider  │ /api/users         │ GET    │ id
08:30:15  │ Spider  │ /api/users         │ POST   │ name, email
08:30:22  │ Swagger │ /api/admin         │ GET    │ -
08:30:35  │ Scrape  │ /api/billing       │ GET    │ invoice_id
08:30:50  │ Mine    │ /api/users?debug=1 │ GET    │ debug
```

### Tips
- Use `map` for complete reconnaissance (5-10 minutes)
- Check F2 to see all discovered endpoints before proceeding
- Use `swagger` directly if you know the spec location
- `mine` is slower but finds debug/test parameters

---

## Step 2: Analyze Phase

### Purpose
Generate tactical attack plan from discovered endpoints.

### Command

#### `analyze`
Artificial intelligence analyzes discovered endpoints and generates attack tactics.

```bash
> analyze

[08:31:00] ANALYZE: Generating tactical plan...
[08:31:05] NEURO: Evaluating endpoint /api/users...
[08:31:06] NEURO: Evaluating endpoint /api/admin...
[08:31:07] NEURO: Evaluating endpoint /api/billing...
...
[08:31:15] ANALYZE: Generated 8 tactical actions
[08:31:15] ANALYZE: Confidence scores range: LOW (40%) - HIGH (95%)
```

### What AI Does

**1. Endpoint Analysis**
- Examines each endpoint
- Identifies potential vulnerabilities (based on path, method, parameters)
- Assigns confidence scores

**2. Attack Selection**
- BOLA: Endpoints with [id], /{id}, ?id=, /resources/{resource_id}
- BFLA: Endpoints with /admin, /management, /private
- BOPLA: POST/PUT endpoints accepting objects
- SSRF: Endpoints with fetch, proxy, url parameters
- Flow: Multi-step chains combining modules

**3. Payload Generation**
- Selects appropriate test values
- For BOLA: Sequential IDs (1,2,3...) and random IDs
- For BFLA: Admin tokens, elevated role tokens
- For BOPLA: Mass assignment properties (admin=true, role=admin)
- For SSRF: Internal IPs (127.0.0.1, 169.254.169.254)

**4. Prioritization**
- High confidence attacks first
- CRITICAL severity endpoints prioritized
- Sequential execution order

### Expected Output (F5 PLAN Tab)
```
ID │ Status  │ Type  │ Target            │ Payload  │ Conf. │ Notes
───────────────────────────────────────────────────────────────
1  │ PENDING │ BOLA  │ /api/users/[id]   │ 1..100   │ HIGH  │ Sequential
2  │ PENDING │ BFLA  │ /api/admin        │ Escalate │ HIGH  │ Needs token
3  │ PENDING │ BOPLA │ /api/users        │ Role     │ MED.  │ Mass assign
4  │ PENDING │ SSRF  │ /api/fetch        │ Meta     │ HIGH  │ Metadata
5  │ PENDING │ EXHAUST │ /api/data       │ Rate     │ LOW   │ Rate limit
```

### Neural Engine Integration

If `neuro on`:
- AI generates payloads using LLM
- Confidence scores based on model predictions
- Alternative mutations suggested (view in F6 NEURO tab)

If `neuro off`:
- Uses built-in rule-based payload generation
- Faster but less sophisticated
- Confidence based on statistical patterns

---

## Step 3: Review & Edit Phase

### Purpose
Human review and modification before execution.

### View Plan

#### `list-plan`
Display all pending actions.

```bash
> list-plan

[08:31:20] ACTION 1: BOLA /api/users/[id] [HIGH] [PENDING]
[08:31:20]   Payload: Sequential IDs 1-100
[08:31:20]   Description: Test object-level authorization
[08:31:20]
[08:31:20] ACTION 2: BFLA /api/admin [HIGH] [PENDING]
[08:31:20]   Payload: Test with user token (escalation)
[08:31:20]   Description: Test function-level authorization
[08:31:20]
[08:31:20] ACTION 3: BOPLA /api/users [MEDIUM] [PENDING]
[08:31:20]   Payload: Inject admin=true, role=admin properties
[08:31:20]   Description: Test object property authorization
```

### Modify Actions

#### `edit <id> <new_payload>`
Override AI-suggested payload.

```bash
> list-plan
[Shows ACTION 1: BOLA /api/users/[id] with payload "1-100"]

> edit 1 /api/users/999999999
[08:31:25] ACTION 1 payload updated to: /api/users/999999999

> edit 2 /api/admin?role=superadmin
[08:31:26] ACTION 2 payload updated to: /api/admin?role=superadmin
```

**Why edit?**
- You know the target better than AI
- Want to test specific values
- Override noisy false-positives
- Use custom payloads (SQLi, XXE, etc)

#### `drop <id>`
Mark action to skip (won't execute).

```bash
> list-plan
[Shows ACTION 5: EXHAUST /api/data [LOW] - you think noisy]

> drop 5
[08:31:30] ACTION 5 marked as DROPPED

> list-plan
[Shows 4 PENDING actions, ACTION 5 now hidden]
```

**When to drop:**
- Noisy endpoints (generates false positives)
- Rate limiting endpoints (might block test)
- Already tested modules
- Out of scope testing

### View in F5 PLAN Tab

Press **F5** to see tactical actions:
```
ID │ Status   │ Type   │ Target          │ Payload      │ Confidence
─────────────────────────────────────────────────────────────────
1  │ PENDING  │ BOLA   │ /api/users/[id] │ 1..100       │ HIGH (95%)
2  │ PENDING  │ BFLA   │ /api/admin      │ Escalate     │ HIGH (85%)
3  │ PENDING  │ BOPLA  │ /api/users      │ Role inject  │ MED (70%)
4  │ DROPPED  │ EXHAUST│ /api/data       │ Rate limit   │ LOW (40%)
```

**Tab Navigation:**
- **Tab:** Move between actions
- **E:** Edit selected action
- **D:** Drop selected action
- **Arrow keys:** Navigate
- **Enter:** Preview action details

### Interception Mode (Pre-Execution)

Set up **Interceptor** before commit to manually modify requests:

```bash
> neuro-gen 5
[Shows 5 payload variants in F6 tab]

> proxy localhost:8080
[Route through Burp Suite for inspection]
```

---

## Step 4: Execution Phase

### Purpose
Execute the tactical plan and capture results.

### Command

#### `commit`
Execute all PENDING actions.

```bash
> commit

[08:31:40] EXEC: Starting tactical action execution...
[08:31:40] EXEC: Spawning 4 workers (1 per action, sequential)
[08:31:45] EXEC: ACTION_1 BOLA /api/users - RUNNING
[08:31:50] EXEC: ACTION_1 BOLA /api/users - SUCCESS (15 objects)
[08:31:50] LOOT: Captured user@example.com, user2@domain.com...
[08:31:55] EXEC: ACTION_2 BFLA /api/admin - RUNNING
[08:32:00] EXEC: ACTION_2 BFLA /api/admin - SUCCESS (admin panel)
[08:32:00] LOOT: Captured admin_token=xyz123...
[08:32:05] EXEC: ACTION_3 BOPLA /api/users - RUNNING
[08:32:10] EXEC: ACTION_3 BOPLA /api/users - FAILED (403 Forbidden)
[08:32:15] EXEC: All actions complete
```

### Real-Time Monitoring

#### View in F5 PLAN Tab
```
ID │ Status    │ Result
───────────────────────────────
1  │ SUCCESS   │ 15 objects
2  │ SUCCESS   │ Admin panel
3  │ FAILED    │ 403 Forbidden
4  │ DROPPED   │ Skipped
```

#### View in F1 LOGS Tab
```
[08:31:45] ACTION_1: BOLA /api/users - RUNNING
[08:31:50] REQUEST: GET /api/users/1
[08:31:50] RESPONSE: 200 OK (145 bytes)
[08:31:51] REQUEST: GET /api/users/2
[08:31:51] RESPONSE: 200 OK (142 bytes)
...
[08:31:55] RESULT: 15 accessible users found
[08:31:55] LOOT: 15 email addresses captured
```

#### View in F3 LOOT Tab
```
Type   │ Value                │ Source
───────────────────────────────────
EMAIL  │ user@example.com     │ /api/users
EMAIL  │ user2@domain.com     │ /api/users
JWT    │ eyJ0exAi...          │ /api/admin
AWS_KEY│ AKIA...              │ SSRF
```

### Interceptor During Execution

**Ctrl+I** opens interceptor modal:
- Shows requests before they're sent
- Manually modify payload
- Forward or drop request
- See response

```bash
[INTERCEPTOR MODAL]
REQUEST: GET /api/users/1
Headers: Authorization: Bearer token123

[Edit payload to: /api/users/999]
[Press Enter to forward]

[Response shows user data for ID 999]
```

### Concurrent Execution

Actions execute **sequentially** (one after another):
- More reliable (avoids rate limiting)
- Better for forensics (clear logs)
- Can pause with Ctrl+C

Alternative - Use `flow race` for timing attacks:
```bash
> flow race
Executes with timing randomization
Tests race condition vulnerabilities
```

---

## Step 5: Results & Analysis Phase

### Purpose
Review findings and captured loot.

### View Loot

Press **F3** to see captured secrets:
```
Type      │ Value
──────────────────────────────
EMAIL     │ admin@internal.com
PASSWORD  │ SecurePass123
JWT       │ eyJ0exAi...
API_KEY   │ sk-xxxxxxxxxx
```

Command line alternative:
```bash
> loot
> loot list
> loot export
```

### View Analysis

Press **F4** to see vulnerability findings:
```
Type     │ Severity │ Endpoint      │ Description
──────────────────────────────────────
BOLA     │ HIGH     │ /api/users    │ ID enumeration works
BFLA     │ HIGH     │ /api/admin    │ Privilege escalation
BOPLA    │ MEDIUM   │ /api/users    │ Mass assignment
```

### Generate Report

```bash
> report
[08:33:00] REPORT: Generating findings report...
[08:33:05] REPORT: Saved to ./reports/VAPORTRACE_PEN_TEST_20260208_0833.md
```

---

## Complete Example Walkthrough

### Scenario: 15-minute basic test

```
> target https://api.example.com
[08:30:00] TARGET set

> map
[08:30:30] Complete! 33 endpoints discovered
Press F2 to review

> analyze
[08:31:15] Generated 8 tactical actions
Press F5 to review

> list-plan
[Shows 8 actions with confidence]

> drop 5
[Drop low-confidence action]

> commit
[08:31:40] Executing 7 actions...
Press F1 to watch logs
Press F3 to see loot in real-time

[08:32:15] All actions complete

> loot list
[15 secrets captured]

> report
[Report saved]
```

### Scenario: Advanced test with AI mutations

```
> neuro on
[Neural engine activated]

> target https://api.example.com
> map
> analyze
> list-plan

[AI suggests payloads with confidence scores]

> press F6
[View NEURO tab with 5 payload variants]

> edit 1 [use AI variant from F6]

> commit
[Watch F1 for execution, F6 for mutations being tested]

[Results in F3, F4]

> report
```

---

## Advanced: Action Chains

### Combining Multiple Steps

```bash
> flow list
[Shows saved attack chains]

> flow run
[Executes: map → analyze → commit → report]

> flow race
[Timing-randomized execution for race conditions]
```

---

## Tips & Tricks

### Strategic Planning
1. Always `map` first - understand the attack surface
2. Review generated plan before `commit`
3. Edit/drop obvious false positives
4. Start with HIGH confidence actions
5. Lower confidence actions might need manual payload adjustment

### Efficiency
- Use `analyze` once, then re-run `commit` multiple times with edits
- `drop` noisy modules instead of re-analyzing
- `edit` specific IDs for targeted testing

### When to Use Interception
- Testing WAF bypasses (enable `neuro on`)
- Custom payloads beyond AI suggestions
- Request modification (headers, tokens)
- Response analysis

### Review Before Commit
- Never blindly run `commit`
- Review all 8+ actions with `list-plan`
- Understand what each action does
- Check confidence scores
- Drop/edit as needed

---

See also: [02_FIRST_RUN.md](02_FIRST_RUN.md) for step-by-step walkthrough, [18_COMMAND_REFERENCE.md](18_COMMAND_REFERENCE.md) for all commands, [05_RECONNAISSANCE.md](05_RECONNAISSANCE.md) for discovery details.

