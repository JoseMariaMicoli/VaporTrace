```markdown
    __  __                         _____                    
    \ \ / /___  _ __  ___  _ __   |_   _| __ __ _  ___ ___ 
     \ V // _ `| '_ \/ _ \| '__|    | || '__/ _` |/ __/ _ \
      \  / (_| | |_)  (_) | |       | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|       |_||_|  \__,_|\___\___|
               |_|      [ Surgical API Exploitation Suite ]

```
### **VaporTrace v3.1-Hydra**

### [ Advanced API Risk Intelligence & Heuristic Analysis Suite ]

**VaporTrace** is a high-performance security ecosystem engineered in Go, designed to deliver **high-level risk intelligence** through surgical reconnaissance and logic-first exploitation. It bridges the gap between deep technical vulnerabilities and business impact by transforming raw API telemetry into actionable defensive insights and structured compliance evidence.

With the successful completion of **Sprint 10 (HYDRA)**, the suite provides a unified tactical interface for real-time monitoring of complex authorization logic and automated alignment with global cybersecurity standards.

> **Project Phase:** HYDRA (Sprint 10) - Stable Release.
> **Current Version:** v3.1-Hydra (Stable)
> **Core Principle:** Middleware-First Interception & Heuristic Logic Analysis.

---

### **⚡ Core Capabilities & Framework Alignment**

* **Full OWASP API Top 10 (2023) Validation:** VaporTrace is engineered to rigorously test against the entire **OWASP API Security Top 10 2023** spectrum, ensuring comprehensive validation of the modern API attack surface, from **BOLA (API1)** to **Unsafe Consumption (API10)**.
* **AI Heuristic Brain (Hybrid Stack):** Utilizing a hybrid architecture (Gemini Cloud + Local Mistral via Ollama), VaporTrace identifies "Logic Gaps" and non-linear vulnerabilities that traditional scanners overlook.
* **9.13 Reporting & Compliance Engine:** Automatically synthesizes technical findings into **executive-ready documentation**. Every discovery is natively mapped to:
* **MITRE ATT&CK®:** Identification of adversary tactics and techniques (e.g., T1594, T1020).
* **NIST CSF v2.0:** Direct correlation with Identify (ID), Protect (PR), and Detect (DE) functions.
* **CVE & CVSS v3.1/4.0 Scoring:** Automated assignment of **CVE** references and **CVSS** vectors.


* **Shadow API Discovery:** Advanced reconnaissance modules designed to uncover undocumented, legacy, or "zombie" endpoints, providing a comprehensive map of the hidden attack surface.
* **Hydra Tactical TUI:** A sophisticated, real-time operational dashboard built with `rivo/tview`. It centralizes all interceptor telemetry and mission data into a single, high-visibility interface.
* **Mission Vault Persistence:** A hardened **SQLite3** backend that maintains a persistent record of all engagement logs, ensuring that every finding is auditable and ready for **DFIR**.

---

### **📈 Strategic Roadmap: The Path Forward (Sprints 11–16)**

VaporTrace is committed to continuous evolution, moving toward autonomous operations and sector-specific regulatory excellence.

* **Sprint 11 | Graph-Based Attack Surface Mapping:** Implementation of Directed Acyclic Graphs (DAG) to visualize service dependencies and identify multi-step "Chained Logic Flaws."
* **Sprint 12 | Cloud-Native Metadata & Pivot Modules:** Automated exploitation of Cloud Metadata Services (IMDSv2) to harvest credentials via SSRF vectors in AWS, Azure, and GCP environments.
* **Sprint 13 | Tactical Stealth & Evasion:** Integration of dynamic JA3 TLS fingerprint randomization and header jittering to emulate sophisticated adversaries and bypass modern WAF/IPS detection.
* **Sprint 14 | GxP & Pharmaceutical Regulatory Deep-Dive:** Expansion of the reporting engine to include **GxP (Good Practice)** and **HIPAA** validation checks, specifically tailored for the pharmaceutical and healthcare sectors.
* **Sprint 15 | Distributed Hydra (Multi-Operator Mode):** Transition to a synchronized, encrypted backend allowing multiple operators to collaborate on a single "Mission Vault" in real-time.
* **Sprint 16 | Autonomous Heuristic Remediation:** Development of a "Blue-Team Mirror" that uses the AI Brain to suggest specific code-level middleware fixes for discovered BOLA/BFLA vulnerabilities.

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

### **Tactical Mapping Matrix (CLI Modules)**

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

### **Neural Engine & Interceptor Mappings (Real-Time)**

New capabilities introduced in **v3.1-Hydra** map to the following techniques:

| Feature | MITRE ID | Tactic | Technique Name |
| --- | --- | --- | --- |
| **Shadow Recon** | **T1595** | Reconnaissance | Active Scanning (Hash-router & Shadow API prediction) |
| **Neuro Brute** | **T1110** | Credential Access | Brute Force (Neural-driven Mutation) |
| **Interceptor Snap** | **T1041** | Exfiltration | Exfiltration Over C2 Channel (Traffic Snapshot) |

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
| **F6** | **NEURAL** | **Neural Engine View.** Displays AI mutation logic and fuzzy results. |
| **CTRL+I** | **INTERCEPTOR** | **Global Toggle.** Activates/Deactivates the Tactical Interceptor Modal. |

### **The Tactical Interceptor (Modal Commands)**

When Interception is enabled (`Ctrl+I`), the `TacticalTransport` middleware pauses outgoing traffic. Use these shortcuts within the Red Modal to manipulate or synchronize data.

| Key | Action | Description |
| --- | --- | --- |
| **CTRL + F** | **FORWARD** | Injects the modified packet back into the pipeline. |
| **CTRL + D** | **DROP** | Drops the packet immediately; the request never leaves the local machine. |
| **CTRL + A** | **SYNC TO NEURAL** | **Global Action:** Mirrors the current buffer (Tab 4) to **Tab 7 (NEURAL AI)**. |
| **CTRL + N** | **NEURO INV** | **Neural Inverter:** Toggles AI-assisted payload mutation for logic bypass. |
| **CTRL + B** | **NEURO BRUTE** | **Neural Brute:** Triggers a high-speed, entropy-aware fuzzer on the selected field. |
| **TAB** | **NAVIGATE** | Switches focus between Path, Headers, and Body input fields. |

---

## 🧠 System Architecture & Operational Logic

### **1. Traffic Flow Diagram**

The following represents the packet journey through the TNI system, from the TUI interface to the Hybrid AI Brain and finally to the target.

```ascii
[USER / TUI]  <--->  [ TACTICAL INTERCEPTOR (Ctrl+I) ]  <--->  [ MIDDLEWARE TRANSPORT ]
      ^                          |                                      |
      | (Ctrl+A/B)               | (Traffic Snapshot)                   | (HTTP Req)
      v                          v                                      v
[ NEURO ENGINE ] <--- [ HYBRID BRAIN ] ---> [ TARGET API ]
      |
      +--> Primary: Gemini (Cloud)
      +--> Fallback: Mistral (Local / Ollama)

```

### **2. Hybrid Brain Decision Logic**

The system prioritizes intelligence depth while maintaining zero-latency availability:

1. **Primary Execution:** TNI attempts to contact the Cloud Provider (**Gemini/OpenAI**).
2. **Health Check:** If the system encounters a `429 (Quota Exceeded)`, `Connection Refused`, or `Empty Response`, the `ExecuteQuery` function triggers an automatic failover.
3. **Local Fallback:** The query is immediately re-routed to the **Local Secondary Provider** (Mistral/Ollama). A red warning is logged in the Tactical Log to notify the operator of the switch.

---

## 🧪 Operational Verification (Testing Battery)

To verify the installation and new operational features, perform the following manual test battery:

### **1. Verify Interceptor Toggle**

* **Action:** Press `Ctrl+I` within the TUI.
* **Test:** Open a separate terminal and run `curl -x http://127.0.0.1:8080 http://example.com`.
* **Verification:** Confirm the TUI modal appears (Red Border), pausing the request.

### **2. Verify Neuro Brute**

* **Action:** Inside the Interceptor, use `TAB` to select the Request Body.
* **Trigger:** Press `Ctrl+B`.
* **Verification:** Switch to the Neural Tab (`F6`) and confirm that mutation payloads are populating.

### **3. Verify Hybrid Fallback**

* **Action:** Temporarily disable your internet connection or provide an invalid API key for Gemini in `config.yaml`.
* **Trigger:** Run an analysis via `Ctrl+A`.
* **Verification:** Check the Logs (`F1`). You should see a warning *"Switching to Fallback"* followed by results generated via Ollama.

---

### **🛠️ Tactical Command Reference**

The Hydra TUI centralizes all commands through a unified command bar. Below is the complete, untruncated technical catalog of available modules.

| Command | Action | Technical Context | Framework Focus |
| --- | --- | --- | --- |
| `target <url>` | **Scope Definition** | Sets the global context for all modules. | General |
| `map -u <url>` | **Inventory** | Spidering, OpenAPI mining, and route extraction. | OWASP API9 |
| `swagger <url>` | **Spec Parsing** | Ingests Swagger/OpenAPI definitions into the DB. | OWASP API9 |
| `scrape <url>` | **JS Mining** | Extracts hidden API paths from JavaScript bundles. | OWASP API9 |
| `mine <url>` | **Param Fuzz** | Brute-forces hidden parameters (debug, admin, test). | OWASP API9 |
| `bola <url>` | **ID Swap** | Broken Object Level Authorization testing. | OWASP API1 |
| **`weaver`** | **Auth Forge** | Intercepts OIDC tokens and masks data exfiltration. | OWASP API2 |
| `bopla <url>` | **Mass Assign** | Broken Object Property Level Authorization (Property injection). | OWASP API3 |
| `exhaust <url>` | **DoS Probe** | Testing resource limits (Payload size, pagination limits). | OWASP API4 |
| `bfla <url>` | **PrivEsc** | Broken Function Level Authorization (Method tampering). | OWASP API5 |
| `ssrf <url>` | **Infra Pivot** | SSRF against Cloud Metadata (169.254.169.254). | OWASP API7 |
| `audit <url>` | **Config Check** | Header analysis, SSL/TLS checks, and CORS auditing. | OWASP API8 |
| `probe <url>` | **Integration** | Tests for unsafe consumption in webhooks/3rd party APIs. | OWASP API10 |
| **`proxy`** | **Routing** | Enables/disables traffic routing (Default: Burp @ 127.0.0.1:8080). | Infrastructure |
| **`proxies load`** | **Rotation** | Loads a list of proxies for rotation to bypass rate limiting. | Evasion |
| **`proxies reset`** | **Clear** | Clears the current proxy rotation list. | Infrastructure |
| **`sessions`** | **Context** | Manages active authentication sessions and stored cookies. | Authentication |
| **`neuro on`** | **Enable Engine** | Activates the Neural Mutation layer for all traffic. | Logic Bypass |
| **`neuro off`** | **Disable Engine** | Reverts to standard manual or static payloads. | General |
| **`neuro config`** | **LLM Settings** | Opens the configuration modal for LLM provider endpoints. | Infrastructure |
| **`test-neuro`** | **Engine Diag** | Runs connectivity and latency tests to the AI provider. | Infrastructure |
| **`test-bola`** | **Logic Diag** | Tests BOLA using a dummy ID (e.g., "999") against httpbin. | Logic Testing |
| **`test-bopla`** | **Logic Diag** | Runs a mass assignment test against a patch endpoint. | Logic Testing |
| **`test-bfla`** | **Logic Diag** | Tests BFLA by attempting admin verbs with low-priv sessions. | Logic Testing |
| **`test-neuro`** | **Latency Test** | Tests response speed of the mistral/ollama local engine. | Neural Config |
| `report` | **Generate** | Triggers the 9.13 Reporting Engine (Markdown/PDF). | Compliance |
| `init_db` | **Persistence** | Initializes the SQLite3 Framework-Tagged backend. | Infrastructure |
| `seed_db` | **Intelligence** | Populates the Aggregator with test/known credentials. | Infrastructure |
| `reset_db` | **Wipe** | Purges all mission data from the local database. | Infrastructure |

### **Framework Compliance Context**

* **MITRE ATT&CK:** Every synchronization via `CTRL + A + S` automatically tags findings for **T1552** (Credential Access) or **T1562.001** (Defense Evasion) within the persistence layer.
* **OWASP API Security:** The `neuro` engine is specifically optimized for testing logic-heavy vulnerabilities like **API1:2023** (BOLA) and **API2:2023** (Broken Authentication) that require context-aware mutations beyond simple brute force.

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

**Report Examples:**
Inside the /reports directory, you will find two demonstration files:

```
Automated Seed Report: Generated using the seed_db command to showcase the engine's formatting and dummy data capabilities.

Hybrid Engagement Report: A specialized report combining synthetic data with real-world findings from a dedicated Pen-Test API, demonstrating how the tool handles live reconnaissance.

```

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
| **Sprint 1: Foundation** |
|  | 1.1 | Cobra CLI Engine: Subcommand-based architecture (map, scan, auth). | ✅ DONE |
|  | 1.2 | Interactive Shell UI: Advanced REPL with readline auto-completion. | ✅ DONE |
|  | 1.3 | The Burp Bridge: Industrial-strength HTTP client with native proxy support. | ✅ DONE |
|  | 1.4 | SSL/TLS Hardening: Automatic bypass of self-signed certs for proxies. | ✅ DONE |
|  | 1.5 | Global Config: Persistent flag management for headers and authentication. | ✅ DONE |
| **Sprint 2: Recon** |
|  | 2.1 | Spec Ingestion: Automated parsing of Swagger (v2) and OpenAPI (v3). | ✅ DONE |
|  | 2.2 | JS Route Scraper: Regex-based endpoint extraction from JS bundles. | ✅ DONE |
|  | 2.3 | Version Walker: Identification of deprecated versions (/v1/ vs /v2/). | ✅ DONE |
|  | 2.4 | Parameter Miner: Automatic identification of hidden query params/headers. | ✅ DONE |
| **Sprint 3: Auth Logic** |
|  | 3.1 | BOLA Prober (API1): Tactical ID-swapping engine with session stores. | ✅ DONE |
|  | 3.2 | BOPLA/Mass Assignment (API3): Fuzzing bodies for hidden properties. | ✅ DONE |
|  | 3.3 | BFLA Module (API5): Hierarchical access testing via method manipulation. | ✅ DONE |
| **Sprint 4: Injection** |
|  | 4.1 | Resource Exhaustion (API4): Probing pagination and payload limits. | ✅ DONE |
|  | 4.2 | SSRF Tracker (API7): Detecting OOB callbacks via URL-parameter injection. | ✅ DONE |
|  | 4.3 | Security Misconfig (API8): Automated CORS and Security Header audit. | ✅ DONE |
|  | 4.4 | Integration Probe (API10): Unsafe consumption in webhooks/3rd party. | ✅ DONE |
| **Sprint 5: Intel** |
|  | 5.1 | SQLite Persistence: Local-first mission database for session continuity. | ✅ DONE |
|  | 5.2 | Async Log Worker: Non-blocking background commitments of findings. | ✅ DONE |
|  | 5.3 | Classified Reporting: NIST-aligned Markdown/PDF debrief generator. | ✅ DONE |
|  | 5.4 | Database Management: Built-in init_db and reset_db control. | ✅ DONE |
| **Sprint 6: Evasion** |
|  | 6.1 | Header Randomization: Rotating User-Agents and JA3 fingerprints. | ✅ DONE |
|  | 6.2 | IP Rotation: Integration with proxy-chains and Tor. | ✅ DONE |
|  | 6.3 | Timing Attacks: Implementing jitter and "Sleepy Probes" for NHPP. | ✅ DONE |
| **Sprint 7: Flow & Logic** |
|  | 7.1 | Flow Engine Implementation: Command suite, recording, and replay. | ✅ DONE |
|  | 7.2 | State-Machine Mapping: Logical order enforcement & out-of-order testing. | ✅ DONE |
|  | 7.3 | Race Condition Engine: Multi-threaded "Turbo Intruder" probes. | ✅ DONE |
| **Sprint 8: Post-Exfil** |
|  | 8.1 | Discovery Vault: Real-time regex scanning of all responses for secrets. | ✅ DONE |
|  | 8.2 | Cloud Pivot Engine: Interception of IMDS (169.254.169.254) requests. | ✅ DONE |
|  | 8.3 | Ghost-Weaver Agent: OIDC interception and encrypted exfiltration. | ✅ DONE |
|  | 8.4 | NHPP Evasion: Masking data as "Deprecated Dependency" system logs. | ✅ DONE |
|  | 8.5 | OOB Validation: Automated validation for leaked tokens/infrastructure. | ✅ DONE |
| **Sprint 9: Hardening** |
|  | 9.1 | Report Engine: Refactored NIST generator with Vault integration. | ✅ DONE |
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

### **Part II: The Hydra TUI & Autonomous Systems [TESTING]**

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 10: Hydra** |
|  | 10.1 | Universal Target Function (Global Context) | ✅ DONE |
|  | 10.2 | Project Mosaic: The Hydra-TUI Dashboard | ✅ DONE |
|  | 10.2.1 | Terminal Multi-Pane (Quadrants + F-Tabs Switcher) | ✅ DONE |
|  | 10.2.2 | Legacy Shell Fallback (CLI Flag Logic) | ✅ DONE |
|  | 10.3 | Contextual Aggregator & Information Gathering | ✅ DONE |
|  | 10.4 | Tactical Interceptor (F2 Modal Manipulation) | ✅ DONE |
|  | 10.5 | AI Base Integration (Heuristic Brain) | ✅ DONE |
|  | 10.6 | AI Payload Generation & Autonomous Fuzzing | ✅ DONE |

### **Part III: The Future Evolution [NEW]**

| Phase | Sub-Phase | Focus / Technical Deliverable | Status |
| --- | --- | --- | --- |
| **Sprint 11: Autonomy** |  
|  | **11.1** | **Dynamic Dependency Injection (DDI)** |
|  | 11.2 | State-Machine driven payload selection | ❌ [NEW] |
|  | 11.3 | Autonomous lateral movement within API subnets | ❌ [NEW] |
| **Sprint 12: Evasion V2** |
|  |  | 12.1 | Deep Traffic Shaping: Mimicking legitimate API traffic |
|  | 12.2 | Encrypted OOB: Secure exfiltration via custom protocols | ❌ [NEW] |
|  | 12.3 | Behavioral Jitter: Randomized inter-packet timing | ❌ [NEW] |
| **Sprint 13: The Hive** |
|  |  | 13.1 | Hybrid C2 Architecture: gRPC Control Plane |
|  | 13.2 | RESTful Management API for the Hive Master | ❌ [NEW] |
|  | 13.3 | VaporTrace Console: Web-based Mission Dashboard | ❌ [NEW] |
| **Sprint 14: Pivot** |  
|  | 14.1 | Cross-Tenant Leakage: Exploiting shared infrastructure |
|  | 14.2 | K8s Escape: API-to-Cluster orchestration pivoting | ❌ [NEW] |
|  | 14.3 | Serverless Poisoning: Attacking Lambda/Cloud-Function logic | ❌ [NEW] |
| **Sprint 15: Mastery** | 
|  | 15.1 | Post-Quantum Cryptography for NHPP |
|  | 15.2 | Multi-Agent Swarm Logic (Coordinated BOLA) | ❌ [NEW] |

---

### 🔧 III. Installation & Setup

Follow these steps to configure the development environment, initialize the AI engine, and compile the **VaporTrace** binary.

#### **1. Clone the Repository**

```bash
git clone git@github.com:JoseMariaMicoli/VaporTrace.git
cd VaporTrace


```

#### **2. Install & Start Ollama (AI Engine)**

VaporTrace leverages **Ollama** running the **Mistral** model for localized, private heuristic analysis.

* **Installation:**
* **Linux/macOS:** `curl -fsSL https://ollama.com/install.sh | sh`
* **Windows:** Download the official installer at [ollama.com](https://ollama.com).
* **Start the Service:**
The Ollama server must be active before running VaporTrace:

```bash
ollama serve


```

*(Keep this terminal open or ensure the service is running as a background daemon).*

* **Pull the Model:**
In a separate terminal, download the Mistral model:

```bash
ollama pull mistral


```

#### **3. Build from Source**

Requires **Go 1.21+**. The build process automatically resolves dependencies for the **Hydra TUI** and the networking stack.

```bash
# Synchronize and verify dependencies
go mod tidy

# Compile the tactical binary
go build -o VaporTrace main.go


```

#### **4. Execution**

Initialize the tactical suite:

```bash
./VaporTrace


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

> ```
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

> ```
> 
> ```
> 
> 

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
* **Reporting:** NIST-aligned Markdown generator (Sprint 9.13).

**VaporTrace - Reveal the Invisible.**

---
