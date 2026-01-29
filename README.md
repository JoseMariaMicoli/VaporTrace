```markdown
    __  __                         _____                    
    \ \ / /___  _ __  ___  _ __   |_   _| __ __ _  ___ ___ 
     \ V // _ `| '_ \/ _ \| '__|    | || '__/ _` |/ __/ _ \
      \  / (_| | |_)  (_) | |       | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|       |_||_|  \__,_|\___\___|
               |_|      [ Surgical API Exploitation Suite]

```
**VaporTrace** is a high-performance Red Team framework engineered in Go for surgical reconnaissance and exploitation of API architectures. It specializes in uncovering "Shadow APIs," analyzing authorization logic (BOLA/BFLA), and mapping the entire attack surface of modern REST/Microservice environments.

> **Project Phase:** Engineering & Hardening: Refactoring Surgical Reporting & Vault Integration.
> **Research Status:** RED TEAM R&D / API SECURITY GAP ANALYSIS
> **Core Principle:** Logic-First Exploitation & Non-Destructive Ingestion

---

## ⚠️ FULL LEGAL DISCLAIMER & RULES OF ENGAGEMENT

**THIS TOOL IS FOR AUTHORIZED PENETRATION TESTING AND EDUCATIONAL PURPOSES ONLY.**

1. **Authorization Required:** Never use VaporTrace against targets you do not have explicit, written permission to test.
2. **No Liability:** The author and contributors assume no liability and are not responsible for any misuse, data loss, service degradation, or legal consequences caused by this program.
3. **Local Laws:** It is the user's responsibility to comply with all applicable local, state, and international laws.
4. **Logic Risk:** Be aware that automated BOLA/BFLA probing can modify server-side data. Always perform tests in a controlled staging environment when possible.

**By compiling or running this software, you agree to these terms.**

---

### **🛡️ MITRE ATT&CK Mapping (Full Suite)**

VaporTrace operations are mapped across the full attack lifecycle to provide stakeholders with clear visibility into adversary emulation.

| PHASE | TACTIC | TECHNIQUE | VAPORTRACE COMPONENT |
| --- | --- | --- | --- |
| **P1: Foundation** | Command and Control | T1105: Ingress Tool Transfer | `Burp Bridge / Proxy Config` |
| **P2: Discovery** | Reconnaissance | T1595.002: Active Scanning | `map`, `swagger`, `mine` |
| **P2: Discovery** | Reconnaissance | T1592: Victim Info | **`pipeline` (Endpoint Categorization)** |
| **P3: Auth Logic** | **Privilege Escalation** | **T1548: Abuse Elevation** | **`bola --pipeline` (Mass Engine)** |
| **P3: Auth Logic** | **Privilege Escalation** | **T1548.002: Mass Assignment** | **`bopla --pipeline` (Property Fuzzer)** |
| **P3: Auth Logic** | Privilege Escalation | T1548: Abuse Elevation | `scan-bfla` (Verb Tampering) |
| **P4: Injection** | Impact | T1499: Endpoint DoS | `resource-exhaustion (API4)` |
| **P4: Injection** | Discovery | T1046: Network Service Discovery | `ssrf-tracker (API7)` |
| **P5: Reporting** | Reporting | T1592: Victim Info | `persistence (SQLite) / report` |
| **P8: Exfiltration** | **Exfiltration** | **T1041: Exfiltration Over C2** | **`weaver` (Ghost-Weaver)** |
| **P8: Discovery** | **Credential Access** | **T1552.001: Files/Env** | **`loot` (Discovery Vault)** |
| **P8: Discovery** | **Credential Access** | **T1552.005: Cloud Provider** | **`TriggerCloudPivot`** |
| **Standardization** | **Exfiltration** | **T1071.001: Web Protocols** | **`SafeDo` (Universal Mirroring)** |
| **Standardization** | Credential Access | T1557: AiTM | **`X-VaporTrace-Signal`** |


---

## 🖥️ The Tactical Shell: Persistence & Context

The **VaporTrace Shell** is the core differentiator of this framework. Unlike standard one-shot CLI tools, the shell provides a **Persistent Security Context** required for complex logic testing.

### Strategic Use Case: The "Ghost-Weaver" Pivot

In modern API pentesting, most vulnerabilities aren't found in a single request, but in the **logical relationship** between two accounts and background persistence.

* **Identity Management:** The shell maintains a global state for `Attacker` and `Victim` tokens. You configure them once, and the engine automatically handles the "Identity Swap" during probes.
* **Background Sovereignty:** The `weaver` command spawns a background agent that monitors for OIDC tokens and exfiltrates discovered loot via AES-256-GCM encrypted channels.
* **NHPP Evasion:** Every exfiltrated packet is masked as a `[WARN] Deprecated dependency` log to bypass basic automated traffic analysis.
* **Real-time Triage:** Integrated `pterm` tables provide immediate feedback on whether a request was blocked (403), missing (404), or successfully leaked (200 OK).

To enter the interactive tactical mode, execute:

```bash
./VaporTrace shell

```

---

## 🚀 Strategic Roadmap

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 1: Foundation** |  |  |  |
| ✅ | 1.1 | **Cobra CLI Engine:** Subcommand-based architecture (`map`, `scan`, `auth`). | DONE |
| ✅ | 1.2 | **Interactive Shell UI:** Advanced REPL with `readline` auto-completion. | DONE |
| ✅ | 1.3 | **The Burp Bridge:** Industrial-strength HTTP client with native proxy support. | DONE |
| ✅ | 1.4 | **SSL/TLS Hardening:** Automatic bypass of self-signed certs for proxies. | DONE |
| ✅ | 1.5 | **Global Config:** Persistent flag management for headers and authentication. | DONE |
| **Sprint 2: Recon** |  |  |  |
| ✅ | 2.1 | **Spec Ingestion:** Automated parsing of Swagger (v2) and OpenAPI (v3). | DONE |
| ✅ | 2.2 | **JS Route Scraper:** Regex-based endpoint extraction from JS bundles. | DONE |
| ✅ | 2.3 | **Version Walker:** Identification of deprecated versions (`/v1/` vs `/v2/`). | DONE |
| ✅ | 2.4 | **Parameter Miner:** Automatic identification of hidden query params/headers. | DONE |
| **Sprint 3: Auth Logic** |  |  |  |
| ✅ | 3.1 | **BOLA Prober (API1):** Tactical ID-swapping engine with session stores. | DONE |
| ✅ | 3.2 | **BOPLA/Mass Assignment (API3):** Fuzzing bodies for hidden properties. | DONE |
| ✅ | 3.3 | **BFLA Module (API5):** Hierarchical access testing via method manipulation. | DONE |
| **Sprint 4: Injection** |  |  |  |
| ✅ | 4.1 | **Resource Exhaustion (API4):** Probing pagination and payload limits. | DONE |
| ✅ | 4.2 | **SSRF Tracker (API7):** Detecting OOB callbacks via URL-parameter injection. | DONE |
| ✅ | 4.3 | **Security Misconfig (API8):** Automated CORS and Security Header audit. | DONE |
| ✅ | 4.4 | **Integration Probe (API10):** Unsafe consumption in webhooks/3rd party. | DONE |
| **Sprint 5: Intel** |  |  |  |
| ✅ | 5.1 | **SQLite Persistence:** Local-first mission database for session continuity. | DONE |
| ✅ | 5.2 | **Async Log Worker:** Non-blocking background commitments of findings. | DONE |
| ✅ | 5.3 | **Classified Reporting:** NIST-aligned Markdown/PDF debrief generator. | DONE |
| ✅ | 5.4 | **Database Management:** Built-in `init_db` and `reset_db` control. | DONE |
| **Sprint 6: Evasion** |  |  |  |
| ✅ | 6.1 | **Header Randomization:** Rotating User-Agents and JA3 fingerprints. | DONE |
| ✅ | 6.2 | **IP Rotation:** Integration with proxy-chains and Tor. | DONE |
| ✅ | 6.3 | **Timing Attacks:** Implementing jitter and "Sleepy Probes" for NHPP. | DONE |
| **Sprint 7: Flow & Logic** |  |  |  |
| ✅ | 7.1 | **Flow Engine Implementation:** Command suite, recording, and replay. | DONE |
| ✅ | 7.2 | **State-Machine Mapping:** Logical order enforcement & out-of-order testing. | DONE |
| ✅ | 7.3 | **Race Condition Engine:** Multi-threaded "Turbo Intruder" probes. | DONE |
| **Sprint 8: Post-Exfil** |  |  |  |
| ✅ | 8.1 | **Discovery Vault:** Real-time regex scanning of all responses for secrets. | DONE |
| ✅ | 8.2 | **Cloud Pivot Engine:** Interception of IMDS (169.254.169.254) requests. | DONE |
| ✅ | 8.3 | **Ghost-Weaver Agent:** OIDC interception and encrypted exfiltration. | DONE |
| ✅ | 8.4 | **NHPP Evasion:** Masking data as "Deprecated Dependency" system logs. | DONE |
| ✅ | 8.5 | **OOB Validation:** Automated validation for leaked tokens/infrastructure. | DONE |
| **Sprint 9: Hardening** |  |  |  |
| ✅ | 9.1 | **Report Engine:** Refactored NIST generator with Vault integration. | DONE |
| ✅ | 9.1.1 | **Tactical UI:** Integrated spinners and real-time feedback tables. | DONE |
| ✅ | 9.2 | **Surgical BOLA:** Response Diffing engine to eliminate False Positives. | DONE |
| ✅ | 9.3 | **Concurrency Engine:** High-speed channel-based worker pools. | DONE |
| ✅ | 9.4 | **Environment Sensing:** Burp/ZAP detection with X-Header signaling. | DONE |
| ✅ | 9.5 | **Discovery-to-Engine:** Automating map-to-scan handover pipeline. | DONE |
| ✅ | 9.6 | **Universal Proxy:** Refactored `SafeDo` with multi-module mirroring. | DONE |
| ✅ | 9.7 | **BOLA Concurrency:** Multi-threaded mass scanner upgrade. | DONE |
| ✅ | 9.8 | **Industrialized BOPLA:** Concurrent JSON property fuzzing. | DONE |
| ✅ | 9.9 | **Industrialized BFLA:** Method Matrix worker pool (Verb-Tampering). | DONE |
| ✅ | 9.10 | **Universal Concurrency:** `GenericExecutor` standardization. | DONE |
| ✅ | 9.11 | **Ghost Masquerade:** Process renaming to `kworker_system_auth`. | DONE |
| ❌ | **9.12** | **UI Sovereignty:** Migration to `tview` multi-pane dashboarding. | **ACTIVE** |
| **Sprint 10: Vanguard** |  |  |  |
| ❌ | 10.1 | **AI-Driven Fuzzing:** Context-aware payload generation (Local LLM). | PLANNED |
| ❌ | 10.2 | **Auto-Exploit PoC:** Standalone script generation for verified vulns. | PLANNED |

---

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
| `proxies load <f>` | Ingests proxy list for IP rotation | `proxies load list.txt` |
| `proxies reset` | Flushes pool (Returns to Direct/Burp mode) | `proxies reset` |
| `target <url>` | Locks base URL for automated pipeline | `target https://api.target.com` |
| `pipeline` | Categorize targets for BOLA/BFLA/BOPLA | `pipeline` |
| **Logic Exploitation** |  |  |
| `flow add` | Record business logic sequence (Interactive) | `flow add` |
| `flow run` | Replay sequence with variable injection | `flow run` |
| `flow step` | Tests prerequisite bypasses. | `flow step <id>` |
| `flow race` | High-concurrency synchronized TOCTOU attack. | `flow race <id> <threads>` |
| `flow clear` | Reset flow variables. | `flow clear` |
| `bola` | Execute a live BOLA ID-swap probe (API1) | `bola <url> <id>` |
| `bopla` | Execute Mass Assignment / BOPLA fuzzing (API3) | `bopla <url> '{"id":1}'` |
| `bfla` | Execute Method Shuffling / Verb Tampering (API5) | `bfla <url>` |
| `exhaust` | Execute Phase 4.1 Resource Exhaustion (API4) | `exhaust <url> <param>` |
| `ssrf` | Execute Phase 4.2 SSRF Tracking (API7) | `ssrf <url> <param> <cb>` |
| `audit` | Execute Phase 4.3 Security Misconfig Audit (API8) | `audit <url>` |
| `probe` | Execute Phase 4.4 Integration Probe (API10) | `probe <url> stripe` |
| **Data & Exfiltration** |  |  |
| `weaver <int>` | Deploy Ghost-Weaver background agent with exfil interval | `weaver 60` |
| `loot list` | View all discovered secrets (AWS Keys, JWTs, IPs) | `loot list` |
| `loot clear` | Purge the in-memory discovery vault | `loot clear` |
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

### 3. Tactical Workflow Example (BOPLA / API3)

Identify a sensitive property and attempt to escalate:

```bash
# 1. Enter the shell
./VaporTrace shell

# 2. Set the Attacker Context
vapor@trace:~$ auth attacker eyJhbGciOiJIUzI1...

# 3. Target a user-settings endpoint with a base JSON object
# The engine will attempt to inject 'is_admin', 'role', etc.
vapor@trace:~$ bopla [https://api.target.com/v1/user/me](https://api.target.com/v1/user/me) '{"name":"vapor"}'

```

---

### **📑 Tactical Incident Response (IR) Template (NIST SP 800-61 R3 Alignment)**

Use this unified template to document findings across the VaporTrace tactical phases. Note the new **Mirroring** section for P9.6.

> **[VAPOR-TRACE-SECURITY-ADVISORY]**
> **FINDING ID:** VT-{{YEAR}}-{{ID}}
> **STRATEGIC PHASE:** {{Phase_1_to_8}}
> **DATABASE ID:** {{DB_Session_ID}}
> **TARGET ENDPOINT:** `{{target_url}}`
> **OWASP API TOP 10:** {{OWASP_Category}} (e.g., API1:2023 BOLA)
> **TECHNICAL ANALYSIS:**
> * **Reconnaissance (P2):** Discovered via version walking / shadow API mining.
> * **Tactical Pipeline (P9.5):** Categorized as {{Engine_Type}} based on route heuristics.
> * **Authorization Context (P3):** Identity swap performed between Attacker and Victim tokens.
> * **Exfiltration (P8.3):** Verified via Ghost-Weaver masked encryption.
> * **Mirroring (P9.6):** Request captured in proxy history via `X-VaporTrace-Signal`.
> * **Injection/Consumption (P4):** Logic used to trigger SSRF or Resource Exhaustion.
> * **Persistence (P5):** All tactical logs committed to SQLite for debrief.
> 
> 
> **REPRODUCTION LOG:**
> ```bash
> vapor@trace:~$ {{executed_command}}
> [MIRROR] Confirmed hit via {{Module}} mirrored to proxy.
> [RESULT] {{server_response_code}} | {{latency_ms}}ms
> 
> ```
> 
> 
> **IMPACT:** {{Data_Exfiltration / Service_Instability / Privilege_Escalation}}
> **REMEDIATION:** {{Engineering_Action_Plan}}

---

## IV. DFIR RESPONSE GUIDANCE

#### **1. Detection & Analysis (ID.AN)**

* **Network Artifacts:** Monitor for anomalous traffic mirroring or the presence of the `X-VaporTrace-Signal` header. Watch for SSRF patterns targeting internal IP metadata ranges.
* **Endpoint Artifacts:** Audit for background processes renamed to `kworker_system_auth` or unauthorized access to `/proc/net/arp` and OIDC cache files.

#### **2. Containment & Eradication (PR.PT)**

* **Logic Hardening:** Implement Object-Level Authorization (OLA) at the middleware layer to mitigate BOLA.
* **Metadata Protection:** Enforce IMDSv2 with session-oriented headers to prevent unauthenticated credential harvesting.

---

## 📡 The Technology Behind the Tracer

* **Language:** Golang (Concurrency-focused, statically linked).
* **Database:** SQLite3 with async I/O worker pool for persistent mission tracking.
* **UI Stack:** `pterm` for tactical dashboarding and `readline` for shell interactivity.
* **Network Stack:** Custom `net/http` wrapper with `crypto/tls` overrides and robust `net/url` path handling.

**VaporTrace - Reveal the Invisible.**

---