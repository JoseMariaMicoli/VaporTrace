```markdown
    __   __                       _____                    
    \ \ / /___  _ __  ___  _ __  |_   _| __ __ _  ___ ___ 
     \ V // _ `| '_ \/ _ \| '__|   | || '__/ _` |/ __/ _ \
      \  / (_| | |_)  (_) | |      | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|      |_||_|  \__,_|\___\___|
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

### **Phase 4: Consumption & Injection (API4, API7, API8, API10) [ACTIVE]**

* [x] **Resource Exhaustion (API4):** Probing pagination limits and payload size constraints.
* [x] **SSRF Tracker (API7):** Detecting out-of-band callbacks via URL-parameter injection.
* [x] **Security Misconfig (API8):** Automated CORS, Security Header, and Verbose Error audit.
* [x] **Integration Probe (API10):** Identifying unsafe consumption in webhooks and 3rd party triggers.

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
| `auth` | Set identity tokens in the session store | `auth attacker <token>` |
| `sessions` | View currently loaded tokens | `sessions` |
| `bola` | Execute a live BOLA ID-swap probe | `bola <url> <id>` |
| `test-bola` | Run logic verification against httpbin | `test-bola` |
| `bopla` | Execute Mass Assignment fuzzing | `bopla <url> '{"id":1}'` |
| `test-bopla` | Verify BOPLA injection logic | `test-bopla` |
| `bfla` | Execute Method Shuffling / Verb Tampering | `bfla <url>` |
| `test-bfla` | Verify BFLA logic against httpbin | `test-bfla` |
| `map` | Execute full Phase 2 Recon | `map -u <url>` |
| `triage` | Scan local logs for leaked credentials | `triage` |
| `clear` | Reset the terminal view | `clear` |
| `exit` | Gracefully shutdown the suite | `exit` |

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

## 📑 Tactical Incident Response (IR) Template

Use this unified template to document findings across the VaporTrace tactical phases:

> **[VAPOR-TRACE-SECURITY-ADVISORY]**
> **FINDING ID:** VT-{{YEAR}}-{{ID}}
> **STRATEGIC PHASE:** {{Phase_1_to_4}}
> **TARGET ENDPOINT:** `{{target_url}}`
> **OWASP API TOP 10:** {{OWASP_Category}} (e.g., API4:2023 Resource Exhaustion)
> **TECHNICAL ANALYSIS:**
> * **Reconnaissance (P2):** Discovered via version walking / shadow API mining.
> * **Authorization Context (P3):** Identity swap performed between Attacker and Victim tokens.
> * **Injection/Consumption (P4):** Logic used to trigger SSRF or Resource Exhaustion.
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
> **IMPACT:** {{Data_Exfiltration / Service_Instability / Privilege_Escalation}}
> **REMEDIATION:** {{Engineering_Action_Plan}}

---

## 📡 The Technology Behind the Tracer

* **Language:** Golang (Concurrency-focused, statically linked).
* **UI Stack:** `pterm` for tactical dashboarding and `readline` for shell interactivity.
* **Network Stack:** Custom `net/http` wrapper with `crypto/tls` overrides and robust `net/url` path handling.

**VaporTrace - Reveal the Invisible.**

---