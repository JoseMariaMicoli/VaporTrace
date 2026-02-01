```markdown
    __  __                         _____                    
    \ \ / /___  _ __  ___  _ __   |_   _| __ __ _  ___ ___ 
     \ V // _ `| '_ \/ _ \| '__|    | || '__/ _` |/ __/ _ \
      \  / (_| | |_)  (_) | |       | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|       |_||_|  \__,_|\___\___|
               |_|      [ Surgical API Exploitation Suite ]


```

**VaporTrace** is a high-performance Red Team framework engineered in Go for surgical reconnaissance, logic-first exploitation, and automated compliance mapping of API architectures. It specializes in uncovering "Shadow APIs," analyzing complex authorization logic (BOLA/BFLA), and providing a unified command-and-control environment through the **Hydra TUI**. With the integration of the **9.13 Reporting Engine** and the **AI Heuristic Brain**, VaporTrace transforms raw technical findings into executive-level risk intelligence.

> **Project Phase:** HYDRA (Sprint 10) - Unified Tactical TUI & AI Integration.
> **Current Version:** v3.1-Hydra (Stable)
> **Core Principle:** Middleware-First Interception & Logic-First Exploitation.

---

## ⚠️ FULL LEGAL DISCLAIMER & RULES OF ENGAGEMENT

**THIS TOOL IS FOR AUTHORIZED PENETRATION TESTING AND EDUCATIONAL PURPOSES ONLY.**

1. **Authorization Required:** Never use VaporTrace against targets you do not have explicit, written permission to test.
2. **No Liability:** The author and contributors assume no liability and are not responsible for any misuse, data loss, service degradation, or legal consequences caused by this program.
3. **Local Laws:** It is the user's responsibility to comply with all applicable local, state, and international laws.
4. **Logic Risk:** Be aware that automated BOLA/BFLA probing can modify server-side data. Always perform tests in a controlled environment.

**By compiling or running this software, you agree to these terms.**

---

## 🛡️ Framework Alignment (MITRE / OWASP / NIST)

VaporTrace implements a "Zero-Touch" tagging engine that aligns every finding with industry-standard frameworks for immediate NIST compliance auditing.

### **Tactical Mapping Matrix**

| Command | OWASP API 2023 | MITRE ID | MITRE Tactic | CVSS |
| --- | --- | --- | --- | --- |
| **bola** | API1: BOLA | T1594 | Exfiltration | 9.1 |
| **weaver** | API2: Auth | T1606.001 | Credential Access | 8.5 |
| **bopla** | API3: Property | T1548 | Privilege Escalation | 9.8 |
| **exhaust** | API4: Resource | T1499.004 | Impact (DoS) | 7.5 |
| **bfla** | API5: Function | T1548.003 | Privilege Escalation | 8.2 |
| **ssrf** | API7: SSRF | T1071.001 | Discovery | 9.2 |
| **audit** | API8: Misconfig | T1562.001 | Defense Evasion | 5.4 |
| **map** | API9: Inventory | T1595.002 | Reconnaissance | N/A |
| **probe** | API10: Consump. | T1190 | Initial Access | 8.1 |

### **NIST CSF v2.0 Mapping**

| Control ID | Category | Technical Implementation |
| --- | --- | --- |
| **ID.AM-07** | Asset Management | Automated inventory of all REST/Microservice endpoints. |
| **PR.AC-01** | Identity Mgmt | Vault storage and injection of credentials via Aggregator. |
| **PR.AC-03** | Access Control | Verification of Object-Level and Function-Level auth logic. |
| **PR.AC-05** | Least Privilege | Testing for functional escalation between roles. |
| **PR.DS-01** | Data Security | Assessment of SSRF and sensitive data exposure in headers/JS. |
| **PR.PS-01** | Protective Tech | Validation of WAF, SSL/TLS, and secure header configurations. |
| **DE.AE-02** | Detection | Simulation of adversarial patterns to test SOC/SIEM response. |
| **ID.RA-01** | Risk Assessment | Automated CVSS scoring and severity distribution reporting. |

---

## 🖥️ Deep Dive: The Hydra Tactical TUI

The legacy interactive shell has been deprecated in favor of the **Hydra TUI**. Built on `rivo/tview`, it provides a high-speed environment for parallel operations and deep packet manipulation.

### **Interface Navigation (Global Hotkeys)**

| Key | Context | Function |
| --- | --- | --- |
| **F1** | **LOGS** | Real-time system events, attack module feedback, and error streams. |
| **F2** | **TARGETS** | Active target management, scope definition, and global config. |
| **F3** | **TASKS** | Background process monitor for long-running scans and worker pools. |
| **F4** | **TRAFFIC** | **Deep Packet Inspection.** Split-view (Request/Response) of all middleware traffic. |
| **F5** | **CONTEXT** | **The Aggregator.** Displays auto-injected credentials and AI correlations. |
| **F6** | **INTERCEPTOR** | Toggles Global Interception Mode (Red footer indicates **INT-ON**). |

### **The Tactical Interceptor (Modal Commands)**

When Interception is enabled (**F6**), the `TacticalTransport` middleware pauses traffic and opens the Red Modal.

| Key | Action | Description |
| --- | --- | --- |
| **CTRL + F** | **FORWARD** | Injects the modified packet back into the pipeline and sends it to the target. |
| **CTRL + D** | **DROP** | Drops the packet immediately. The request never leaves the local machine. |
| **TAB** | **NAVIGATE** | Switch focus between Path, Headers, and Body input fields. |

---

## 🛠️ Tactical Command Reference

The Hydra TUI centralizes all commands through a unified command bar. Below is the full technical catalog of available modules.

| Command | Action | Technical Context | Framework Focus |
| --- | --- | --- | --- |
| `target <url>` | **Scope Definition** | Sets the global context for all modules. | General |
| `map -u <url>` | **Inventory** | Spidering, OpenAPI/JS mining, and route extraction. | OWASP API9 |
| `swagger <url>` | **Spec Parsing** | Ingests Swagger/OpenAPI definitions into the DB. | OWASP API9 |
| `scrape <url>` | **JS Mining** | Extracts hidden API paths from JavaScript bundles. | OWASP API9 |
| `mine <url>` | **Param Fuzz** | Brute-forces hidden parameters (debug, admin, test). | OWASP API9 |
| `bola <url>` | **ID Swap** | Broken Object Level Authorization testing. | OWASP API1 |
| `weaver <url>` | **Auth Forge** | JWT tampering, KID injection, and algorithm confusion. | OWASP API2 |
| `bopla <url>` | **Mass Assign** | Broken Object Property Level Authorization (Property injection). | OWASP API3 |
| `exhaust <url>` | **DoS Probe** | Testing resource limits (Payload size, pagination limits). | OWASP API4 |
| `bfla <url>` | **PrivEsc** | Broken Function Level Authorization (Method tampering). | OWASP API5 |
| `ssrf <url>` | **Infra Pivot** | Server-Side Request Forgery against Cloud Metadata (169.254). | OWASP API7 |
| `audit <url>` | **Config Check** | Header analysis, SSL/TLS checks, and CORS auditing. | OWASP API8 |
| `probe <url>` | **Integration** | Tests for unsafe consumption in webhooks/3rd party APIs. | OWASP API10 |
| `report` | **Generate** | Triggers the 9.13 Reporting Engine (Markdown/PDF). | Compliance |
| `init_db` | **Persistence** | Initializes the SQLite3 Framework-Tagged backend. | Infrastructure |
| `seed_db` | **Intelligence** | Populates the Aggregator with test/known credentials. | Infrastructure |
| `reset_db` | **Wipe** | Purges all mission data from the local database. | Infrastructure |

---

## 📊 Sprint 9.13: Automated Reporting Engine

The **9.13 Reporting Engine** automates the transition from exploitation to documentation, generating high-fidelity Markdown and PDF reports directly from the SQLite `ContextStore`.

### **Report Architecture:**

* **Executive Summary:** High-level risk overview with CVSS v3.1 distribution charts.
* **Vulnerability Distribution:** Graphical breakdown of Critical, High, Medium, and Low findings.
* **Remediation Priority Tracker:** A prioritized list for engineering teams focusing on critical path vulnerabilities.
* **Adversarial Methodology:** Technical logs including the exact VaporTrace command, timestamp, and target URL for every finding.
* **Framework Verification:** Direct mapping of each finding to MITRE ATT&CK, NIST, and OWASP identifiers.

---

## 🎭 Evasion Techniques & AI Heuristic Brain

### **Sprint 11/12: Deep Evasion Suite**

To bypass modern Next-Gen Firewalls (WAF) and behavioral analytics, VaporTrace implements:

* **Ghost Masquerade:** Process renaming to `kworker_system_auth` to hide from local host-based monitoring.
* **Dynamic Jitter:** Adds variable delay to requests to prevent signature-based rate-limiting detection.
* **Header Randomization:** Rotates User-Agents, Fingerprints, and non-essential headers for every request.
* **Proxy Pools:** Native support for rotation through SOCKS5/HTTP proxy lists to mask origin IPs.

### **Sprint 10.5: AI Heuristic Brain [ACTIVE]**

VaporTrace is currently integrating an AI-driven analysis layer to automate logic-flaw discovery:

* **Pattern Correlation:** The AI analyzes response structures across different roles to predict potential BOLA/BFLA endpoints without manual fuzzing.
* **Autonomous Payload Generation:** Uses heuristic models to generate payloads that specifically target identified technology stacks (e.g., SpringBoot, Express, Django).
* **Anomaly Scoring:** Flags unusual API behavior that might indicate hidden "Shadow" functionality or debug modes.

---

## 🚀 Strategic Roadmap

### **Part I: The Hardened Core & Intelligence [STABLE]**

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 1: Foundation** | 1.1 | Cobra CLI Engine: Subcommand-based architecture (map, scan, auth). | ✅ DONE |
|  | 1.2 | Interactive Shell UI: Advanced REPL with readline auto-completion. | ✅ DONE |
|  | 1.3 | The Burp Bridge: Industrial-strength HTTP client with native proxy support. | ✅ DONE |
|  | 1.4 | SSL/TLS Hardening: Automatic bypass of self-signed certs for proxies. | ✅ DONE |
|  | 1.5 | Global Config: Persistent flag management for headers and authentication. | ✅ DONE |
| **Sprint 2: Recon** | 2.1 | Spec Ingestion: Automated parsing of Swagger (v2) and OpenAPI (v3). | ✅ DONE |
|  | 2.2 | JS Route Scraper: Regex-based endpoint extraction from JS bundles. | ✅ DONE |
|  | 2.3 | Version Walker: Identification of deprecated versions (/v1/ vs /v2/). | ✅ DONE |
|  | 2.4 | Parameter Miner: Automatic identification of hidden query params/headers. | ✅ DONE |
| **Sprint 3: Auth Logic** | 3.1 | BOLA Prober (API1): Tactical ID-swapping engine with session stores. | ✅ DONE |
|  | 3.2 | BOPLA/Mass Assignment (API3): Fuzzing bodies for hidden properties. | ✅ DONE |
|  | 3.3 | BFLA Module (API5): Hierarchical access testing via method manipulation. | ✅ DONE |
| **Sprint 4: Injection** | 4.1 | Resource Exhaustion (API4): Probing pagination and payload limits. | ✅ DONE |
|  | 4.2 | SSRF Tracker (API7): Detecting OOB callbacks via URL-parameter injection. | ✅ DONE |
|  | 4.3 | Security Misconfig (API8): Automated CORS and Security Header audit. | ✅ DONE |
|  | 4.4 | Integration Probe (API10): Unsafe consumption in webhooks/3rd party. | ✅ DONE |
| **Sprint 5: Intel** | 5.1 | SQLite Persistence: Local-first mission database for session continuity. | ✅ DONE |
|  | 5.2 | Async Log Worker: Non-blocking background commitments of findings. | ✅ DONE |
|  | 5.3 | Classified Reporting: NIST-aligned Markdown/PDF debrief generator. | ✅ DONE |
|  | 5.4 | Database Management: Built-in init_db and reset_db control. | ✅ DONE |
| **Sprint 6: Evasion** | 6.1 | Header Randomization: Rotating User-Agents and JA3 fingerprints. | ✅ DONE |
|  | 6.2 | IP Rotation: Integration with proxy-chains and Tor. | ✅ DONE |
|  | 6.3 | Timing Attacks: Implementing jitter and "Sleepy Probes" for NHPP. | ✅ DONE |
| **Sprint 7: Flow & Logic** | 7.1 | Flow Engine Implementation: Command suite, recording, and replay. | ✅ DONE |
|  | 7.2 | State-Machine Mapping: Logical order enforcement & out-of-order testing. | ✅ DONE |
|  | 7.3 | Race Condition Engine: Multi-threaded "Turbo Intruder" probes. | ✅ DONE |
| **Sprint 8: Post-Exfil** | 8.1 | Discovery Vault: Real-time regex scanning of all responses for secrets. | ✅ DONE |
|  | 8.2 | Cloud Pivot Engine: Interception of IMDS (169.254.169.254) requests. | ✅ DONE |
|  | 8.3 | Ghost-Weaver Agent: OIDC interception and encrypted exfiltration. | ✅ DONE |
|  | 8.4 | NHPP Evasion: Masking data as "Deprecated Dependency" system logs. | ✅ DONE |
|  | 8.5 | OOB Validation: Automated validation for leaked tokens/infrastructure. | ✅ DONE |
| **Sprint 9: Hardening** | 9.1 | Report Engine: Refactored NIST generator with Vault integration. | ✅ DONE |
|  | 9.1.1 | Tactical UI: Integrated spinners and real-time feedback tables. | ✅ DONE |
|  | 9.2 | Surgical BOLA: Response Diffing engine to eliminate False Positives. | ✅ DONE |
|  | 9.3 | Concurrency Engine: High-speed channel-based worker pools. | ✅ DONE |
|  | 9.4 | Environment Sensing: Burp/ZAP detection with X-Header signaling. | ✅ DONE |
|  | 9.5 | Discovery-to-Engine: Automating map-to-scan handover pipeline. | ✅ DONE |
|  | 9.6 | Universal Proxy: Refactored SafeDo with multi-module mirroring. | ✅ DONE |
|  | 9.7 | BOLA Concurrency: Multi-threaded mass scanner upgrade. | ✅ DONE |
|  | 9.8 | Industrialized BOPLA: Concurrent JSON property fuzzing. | ✅ DONE |
|  | 9.9 | Industrialized BFLA: Method Matrix worker pool (Verb-Tampering). | ✅ DONE |
|  | 9.10 | Universal Concurrency: GenericExecutor standardization. | ✅ DONE |
|  | 9.11 | Ghost Masquerade: Process renaming to kworker_system_auth. | ✅ DONE |
|  | 9.13 | Refactor: Framework-Tagged DB (OWASP/MITRE/NIST) Integration | ✅ DONE |

### **Part II: The Hydra TUI & Autonomous Systems [ACTIVE]**

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 10: Hydra** | 10.1 | Universal Target Function (Global Context) | ✅ DONE |
|  | 10.2 | Project Mosaic: The Hydra-TUI Dashboard | ✅ DONE |
|  | 10.2.1 | Terminal Multi-Pane (Quadrants + F-Tabs Switcher) | ✅ DONE |
|  | 10.2.2 | Legacy Shell Fallback (CLI Flag Logic) | ✅ DONE |
|  | 10.3 | Contextual Aggregator & Information Gathering | ✅ DONE |
|  | 10.4 | Tactical Interceptor (F2 Modal Manipulation) | ✅ DONE |
|  | 10.5 | AI Base Integration (Heuristic Brain) | ❌ ACTIVE |
|  | 10.6 | AI Payload Generation & Autonomous Fuzzing | ❌ [PLANNED] |

### **Part III: The Future Evolution [NEW]**

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 11: Autonomy** | 11.1 | Dynamic Dependency Injection (DDI) | ❌ [NEW] |
|  | 11.2 | State-Machine driven payload selection | ❌ [NEW] |
|  | 11.3 | Autonomous lateral movement within API subnets | ❌ [NEW] |
| **Sprint 12: Evasion V2** | 12.1 | Deep Traffic Shaping: Mimicking legitimate API traffic | ❌ [NEW] |
|  | 12.2 | Encrypted OOB: Secure exfiltration via custom protocols | ❌ [NEW] |
|  | 12.3 | Behavioral Jitter: Randomized inter-packet timing | ❌ [NEW] |
| **Sprint 13: The Hive** | 13.1 | Hybrid C2 Architecture: gRPC Control Plane | ❌ [NEW] |
|  | 13.2 | RESTful Management API for the Hive Master | ❌ [NEW] |
|  | 13.3 | VaporTrace Console: Web-based Mission Dashboard | ❌ [NEW] |
| **Sprint 14: Pivot** | 14.1 | Cross-Tenant Leakage: Exploiting shared infrastructure | ❌ [NEW] |
|  | 14.2 | K8s Escape: API-to-Cluster orchestration pivoting | ❌ [NEW] |
|  | 14.3 | Serverless Poisoning: Attacking Lambda/Cloud-Function logic | ❌ [NEW] |
| **Sprint 15: Mastery** | 15.1 | Post-Quantum Cryptography for NHPP | ❌ [NEW] |
|  | 15.2 | Multi-Agent Swarm Logic (Coordinated BOLA) | ❌ [NEW] |

---

## 🛠️ Installation & Usage

### 1. Build from Source

```bash
go mod tidy
go build -o VaporTrace

```

### 2. Tactical Workflow Example (BOPLA / API3)

Identify a sensitive property and attempt to escalate privilege using Mass Assignment:

```bash
# 1. Launch the Hydra TUI
./VaporTrace

# 2. Initialize Persistence
:init_db
:seed_db

# 3. Set Scope
:target https://api.target.corp

# 4. Execute Mass Assignment Probe
:bopla https://api.target.corp/v1/user/me '{"name":"vapor"}'

```

### 3. Generate Tactical Report

Once the mission is complete, generate the 9.13 Framework-Aligned report:

```bash
:report

```

---

## 📑 Tactical Incident Response (IR) Template (NIST SP 800-61 R3 Alignment)

Use this unified template to document findings across the VaporTrace tactical phases.

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
> :{{executed_command}}
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
* **UI Engine:** `rivo/tview` / `gdamore/tcell` (Hydra TUI).
* **Database:** SQLite3 with Framework-Tagging (MITRE/OWASP/NIST).
* **Networking:** Middleware-Driven `http.RoundTripper` with native Proxy support.
* **Reporting:** NIST-aligned Markdown/PDF generator (Sprint 9.13).

**VaporTrace - Reveal the Invisible.**

```
---