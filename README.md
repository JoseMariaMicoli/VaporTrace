```markdown
    __   __                    _____                   
    \ \ / /___  _ __  ___  _ __  |_   _| __ __ _  ___ ___ 
     \ V // _ `| '_ \/ _ \| '__|   | || '__/ _` |/ __/ _ \
      \  / (_| | |_)  (_) | |      | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|      |_||_|  \__,_|\___\___|
               |_|      [ API INFRASTRUCTURE TRACER ]
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

### **Phase 3: Authorization & Logic (API1, API3, API5) [ACTIVE]**

* [x] **BOLA Prober (API1):** Tactical ID-swapping engine with persistent session stores for Attacker/Victim contexts.
* [ ] **BOPLA/Mass Assignment (API3):** Fuzzing JSON bodies for administrative or hidden properties.
* [ ] **BFLA Module (API5):** Testing hierarchical access via HTTP method manipulation (GET vs DELETE).

### **Phase 4: Consumption & Injection (API4, API7, API8, API10) [BACKLOG]**

* [ ] **Resource Exhaustion (API4):** Probing pagination limits and payload size constraints.
* [ ] **SSRF Tracker (API7):** Detecting out-of-band callbacks via URL-parameter injection.
* [ ] **Security Misconfig (API8):** Automated CORS, Security Header, and Verbose Error audit.
* [ ] **Integration Probe (API10):** Identifying unsafe consumption in webhooks and 3rd party triggers.

---

## 🖥️ The Tactical Shell: Why Use It?

The **VaporTrace Shell** is designed for the "Pivot & Exploit" phase of an engagement. Unlike one-shot CLI tools, the shell maintains a **Persistent Security Context**.

### Use Case: The "Auth Pivot"
During an API audit, you often find a resource (e.g., `/api/v1/docs/777`) that belongs to **User A**. To test for BOLA (API1), you need to request that same resource using the session of **User B**. 

1. **Context Persistence:** The shell stores your `Attacker` and `Victim` tokens globally. You set them once with `auth`, and every subsequent probe uses them automatically.
2. **Speed:** No need to re-type complex JWTs or headers for every command.
3. **Real-time Triage:** Immediately see formatted tables and VULN alerts as you swap IDs and methods.

---

## 🛠️ Installation & Usage

### 1. Build from Source
```bash
go mod tidy
go build -o VaporTrace

```

### 2. Interactive Shell Commands

| COMMAND | DESCRIPTION | EXAMPLE |
| --- | --- | --- |
| `auth` | Set identity tokens in the session store | `auth attacker <token>` |
| `sessions` | View currently loaded tokens | `sessions` |
| `bola` | Execute a live BOLA ID-swap probe | `bola <url> <id>` |
| `test-bola` | Run logic verification against httpbin | `test-bola` |
| `map` | Execute full Phase 2 Recon | `map -u <url>` |
| `triage` | Scan local logs for leaked credentials | `triage` |
| `clear` | Reset the terminal view | `clear` |
| `exit` | Gracefully shutdown the suite | `exit` |

### 3. Real-World BOLA Workflow

```bash
# 1. Configure the Attacker Identity (User B)
vapor@trace:~$ auth attacker eyJhbGciOiJIUzI1...

# 2. Configure the Victim Identity (User A)
vapor@trace:~$ auth victim eyJhbGciOiJIUzI1...

# 3. Target a sensitive endpoint with the Victim's Resource ID
vapor@trace:~$ bola [https://api.target.com/v1/user/profile](https://api.target.com/v1/user/profile) 501

```

---

## 📡 The Technology Behind the Tracer

* **Language:** Golang (Concurrency-focused, statically linked).
* **UI Stack:** `pterm` for tactical dashboarding and `readline` for shell interactivity.
* **Network Stack:** Custom `net/http` wrapper with `crypto/tls` overrides and robust `net/url` path handling.

**VaporTrace - Reveal the Invisible.**

---