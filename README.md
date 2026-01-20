```markdown
    __  __                         _____                    
    \ \ / /___  _ __  ___  _ __   |_   _| __ __ _  ___ ___ 
     \ V // _ `| '_ \/ _ \| '__|    | || '__/ _` |/ __/ _ \
      \  / (_| | |_)  (_) | |       | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|       |_||_|  \__,_|\___\___|
               |_|      [ Surgical API Exploitation Suite]

```

**VaporTrace** is a high-performance Red Team framework engineered in Go for surgical reconnaissance and exploitation of API architectures. It specializes in uncovering "Shadow APIs," analyzing authorization logic (BOLA/BFLA), and mapping the entire attack surface of modern REST/Microservice environments.

---

## ⚠️ FULL LEGAL DISCLAIMER & RULES OF ENGAGEMENT

**THIS TOOL IS FOR AUTHORIZED PENETRATION TESTING AND EDUCATIONAL PURPOSES ONLY.**

1. **Authorization Required:** Never use VaporTrace against targets you do not have explicit, written permission to test.
2. **No Liability:** The author and contributors assume no liability and are not responsible for any misuse, data loss, service degradation, or legal consequences caused by this program.
3. **Local Laws:** It is the user's responsibility to comply with all applicable local, state, and international laws.
4. **Logic Risk:** Be aware that automated BOLA/BFLA probing can modify server-side data. Always perform tests in a controlled staging environment when possible.

**By compiling or running this software, you agree to these terms.**

---

## 🛡️ Strategic Mapping: MITRE ATT&CK®

VaporTrace operations are mapped across the full attack lifecycle to provide stakeholders with clear visibility into adversary emulation:

| PHASE | TACTIC | TECHNIQUE | VAPORTRACE MODULE |
| --- | --- | --- | --- |
| **P1: Foundation** | Command and Control | T1105: Ingress Tool Transfer | `Burp Bridge / Proxy Config` |
| **P2: Discovery** | Reconnaissance | T1595.002: Active Scanning (API) | `map`, `mine`, `version-walker` |
| **P3: Auth Logic** | Privilege Escalation | T1548: Abuse Elevation Control | `bopla`, `bfla`, `bola` |
| **P4: Injection** | Impact | T1499: Endpoint DoS | `resource-exhaustion (API4)` |
| **P4: Injection** | Discovery | T1046: Network Service Discovery | `ssrf-tracker (API7)` |
| **P5: Reporting** | Reporting | T1592: Gather Victim Host Info | `persistence (SQLite) / report` |

---

## 🖥️ The Tactical Shell: Persistence & Context

The **VaporTrace Shell** is the core differentiator of this framework. Unlike standard one-shot CLI tools, the shell provides a **Persistent Security Context** required for complex logic testing.

### Strategic Use Case: The "Auth Pivot"

In modern API pentesting, most vulnerabilities aren't found in a single request, but in the **logical relationship** between two accounts.

* **Identity Management:** The shell maintains a global state for `Attacker` and `Victim` tokens. You configure them once, and the engine automatically handles the "Identity Swap" during probes.
* **Speed:** No need to re-type complex JWTs or headers for every command.
* **Real-time Triage:** Integrated `pterm` tables provide immediate feedback on whether a request was blocked (403), missing (404), or successfully leaked (200 OK).

To enter the interactive tactical mode, execute:

```bash
./VaporTrace shell

```

---

## 🚀 Strategic Roadmap

### **Phase 1: The Foundation [STABLE]**

* [x] **Cobra CLI Engine:** Subcommand-based architecture (`map`, `scan`, `auth`).
* [x] **Interactive Shell UI:** Advanced REPL with `readline` auto-completion and `pterm` styling.
* [x] **The Burp Bridge:** Industrial-strength HTTP client with native proxy support.
* [x] **SSL/TLS Hardening:** Automatic bypass of self-signed certs for intercepting proxies.
* [x] **Global Config:** Persistent flag management for headers and authentication.

### **Phase 2: Discovery & Inventory (API9) [STABLE]**

* [x] **Spec Ingestion:** Automated parsing of Swagger (v2) and OpenAPI (v3) definitions.
* [x] **JS Route Scraper:** Regex-based endpoint extraction from client-side JavaScript bundles.
* [x] **Version Walker:** Identification of deprecated versions (e.g., `/v1/` vs `/v2/`) to find unpatched logic.
* [x] **Parameter Miner:** Automatic identification of hidden query parameters and headers.

### **Phase 3: Authorization & Logic (API1, API3, API5) [STABLE]**

* [x] **BOLA Prober (API1):** Tactical ID-swapping engine with persistent session stores for Attacker/Victim contexts.
* [x] **BOPLA/Mass Assignment (API3):** Fuzzing JSON bodies for administrative or hidden properties.
* [x] **BFLA Module (API5):** Testing hierarchical access via HTTP method manipulation (GET vs DELETE).

### **Phase 4: Consumption & Injection (API4, API7, API8, API10) [STABLE]**

* [x] **Resource Exhaustion (API4):** Probing pagination limits and payload size constraints.
* [x] **SSRF Tracker (API7):** Detecting out-of-band callbacks via URL-parameter injection.
* [x] **Security Misconfig (API8):** Automated CORS, Security Header, and Verbose Error audit.
* [x] **Integration Probe (API10):** Identifying unsafe consumption in webhooks and 3rd party triggers.

### **Phase 5: Intelligence & Persistence [STABLE]**

* [x] **SQLite Persistence:** Local-first mission database to prevent data loss on session termination.
* [x] **Async Log Worker:** Non-blocking background commitments of tactical findings.
* [x] **Classified Reporting:** Automated generation of professional "Mission Debrief" reports in Markdown/PDF.
* [x] **Database Management:** Built-in `init_db` and `reset_db` commands for mission lifecycle control.

### **Phase 6: Advanced Evasion & Rate-Limit Bypassing [UPCOMING]**

* [ ] **Header Randomization:** Rotating User-Agents and JA3 fingerprints to bypass WAFs.
* [ ] **IP Rotation:** Integration with proxy-chains and Tor for distributed probing.
* [ ] **Timing Attacks:** Implementing jitter and "Sleepy Probes" to stay under SOC thresholds.

### **Phase 7: Business Logic & Workflow Fuzzing [UPCOMING]**

* [ ] **State-Machine Mapping:** Identifying logical order (e.g., Pay -> Download) and testing out-of-order execution.
* [ ] **Race Condition Engine:** Multi-threaded "Turbo Intruder" style probes for currency/credit exploits.

### **Phase 8: Post-Exploitation & Data Exfiltration [UPCOMING]**

* [ ] **Automated PII Scanner:** Scanning response bodies for sensitive data (Credit Cards, SSN, JWTs).
* [ ] **Secret Leaks:** Automatic detection of Cloud Keys (AWS/Azure) in verbose error messages.

### **Phase 9: Engineering & Hardening [ACTIVE]**

* [x] **9.1: Scraper Refinement:** Pre-compiled global regex for high-performance scraping.
* [x] **9.1.1: Tactical UI:** Integrated spinners and real-time tables for immediate feedback.
* [ ] **9.2: Surgical BOLA:** Response Diffing engine (Levenshtein) to eliminate False Positives.
* [ ] **9.3: Concurrency Engine:** High-speed worker pools for massive enumeration.

### **Phase 10: The Vanguard (Future)**

* [ ] **AI-Driven Fuzzing:** Context-aware payload generation using local LLM integration.
* [ ] **Auto-Exploit PoC:** Standalone script generation for verified vulnerabilities.

---

## 🛠️ Installation & Usage

### 1. Build from Source

```bash
go mod tidy
go build -o VaporTrace

```

### 2. Interactive Shell Usage

Launch the shell with `./VaporTrace shell` and use the following tactics:

| COMMAND | DESCRIPTION | EXAMPLE |
| --- | --- | --- |
| **Identity & Sessions** |  |  |
| `auth` | Set identity tokens (JWT/Cookies) in the session store | `auth attacker <token>` |
| `sessions` | View currently loaded tokens for Victim/Attacker | `sessions` |
| **Discovery & Recon** |  |  |
| `map` | Execute full Phase 2 Recon (Endpoint mapping) | `map -u <url>` |
| `swagger` | Parse OpenAPI/Swagger JSON to map attack surface | `swagger <url>` |
| `scrape` | Extract hidden API paths from JavaScript files | `scrape <url>` |
| `mine` | Fuzz for hidden parameters (debug, admin, etc.) | `mine <url> /users` |
| `proxy` | Route all tactical traffic through Burp Suite | `proxy http://127.0.0.1:8080` |
| `proxy off` | Disable the interceptor and go direct | `proxy off` |
| **Logic Exploitation** |  |  |
| `bola` | Execute a live BOLA ID-swap probe (API1) | `bola <url> <id>` |
| `bopla` | Execute Mass Assignment / BOPLA fuzzing (API3) | `bopla <url> '{"id":1}'` |
| `bfla` | Execute Method Shuffling / Verb Tampering (API5) | `bfla <url>` |
| `exhaust` | Execute Phase 4.1 Resource Exhaustion (API4) | `exhaust <url> <param>` |
| `ssrf` | Execute Phase 4.2 SSRF Tracking (API7) | `ssrf <url> <param> <cb>` |
| `audit` | Execute Phase 4.3 Security Misconfig Audit (API8) | `audit <url>` |
| `probe` | Execute Phase 4.4 Integration Probe (API10) | `probe <url> stripe` |
| **Logic Verification** |  |  |
| `test-bola` | Run BOLA logic verification against httpbin | `test-bola` |
| `test-bopla` | Verify BOPLA/Mass-Assignment injection engine | `test-bopla` |
| `test-bfla` | Verify BFLA/Verb-tampering logic | `test-bfla` |
| `test-exhaust` | Verify pagination fuzzing and latency detection | `test-exhaust` |
| `test-ssrf` | Verify SSRF redirect/tracking logic | `test-ssrf` |
| `test-audit` | Verify the Misconfig/CORS scanner | `test-audit` |
| `test-probe` | Verify Webhook/Integration spoofing logic | `test-probe` |
| **System & Debrief** |  |  |
| `init_db` | Initialize Phase 5 SQLite Persistence & Logging | `init_db` |
| `reset_db` | **Wipe all** local mission data (Purge) | `reset_db` |
| `report` | Generate Classified Markdown Mission Report | `report` |
| `clear` | Reset the terminal view/banner | `clear` |
| `exit` | Gracefully shutdown the tactical suite | `exit` |

---

### Final Patch Checklist

To ensure the commands above function as expected:

1. **Global Client:** Ensure `pkg/logic/context.go` and `pkg/discovery/discovery.go` both export `GlobalClient`.
2. **Removal of Local Clients:** In your uploaded files (e.g., `bola.go`, `exhaustion.go`), ensure you have removed the `client := &http.Client{...}` blocks so they default to the global proxied client.
3. **UI Bridge:** Ensure your `shell.go` switch-case triggers these functions.

Would you like me to generate the **Phase 5 Intelligence** logic to handle the `report` generation from the database?

### 3. Tactical Workflow Example (BOPLA / API3)

Identify a sensitive property and attempt to escalate:

```bash
# 1. Enter the shell
./VaporTrace shell

# 2. Set the Attacker Context
vapor@trace:~$ auth attacker eyJhbGciOiJIUzI1...

# 3. Target a user-settings endpoint with a base JSON object
# The engine will attempt to inject 'is_admin', 'role', etc.
vapor@trace:~$ bopla https://api.target.com/v1/user/me '{"name":"vapor"}'

```

---

## 📑 Tactical Incident Response (IR) Template

Use this unified template to document findings across the VaporTrace tactical phases:

> **[VAPOR-TRACE-SECURITY-ADVISORY]**
> **FINDING ID:** VT-{{YEAR}}-{{ID}}
> **STRATEGIC PHASE:** {{Phase_1_to_5}}
> **DATABASE ID:** {{DB_Session_ID}}
> **TARGET ENDPOINT:** `{{target_url}}`
> **OWASP API TOP 10:** {{OWASP_Category}} (e.g., API4:2023 Resource Exhaustion)
> **TECHNICAL ANALYSIS:**
> * **Reconnaissance (P2):** Discovered via version walking / shadow API mining.
> * **Authorization Context (P3):** Identity swap performed between Attacker and Victim tokens.
> * **Injection/Consumption (P4):** Logic used to trigger SSRF or Resource Exhaustion.
> * **Persistence (P5):** All tactical logs committed to SQLite for debrief.
> 
> 
> **REPRODUCTION LOG:**
> ```bash
> vapor@trace:~$ {{executed_command}}
> [RESULT] {{server_response_code}} | {{latency_ms}}ms
> 
> ```
> 
> 

> ```
> 
> **IMPACT:** {{Data\_Exfiltration / Service\_Instability / Privilege\_Escalation}}
> **REMEDIATION:** {{Engineering\_Action\_Plan}}
> 
> ```
> 
> 

---

## 📡 The Technology Behind the Tracer

* **Language:** Golang (Concurrency-focused, statically linked).
* **Database:** SQLite3 with async I/O worker pool for persistent mission tracking.
* **UI Stack:** `pterm` for tactical dashboarding and `readline` for shell interactivity.
* **Network Stack:** Custom `net/http` wrapper with `crypto/tls` overrides and robust `net/url` path handling.

**VaporTrace - Reveal the Invisible.**

---