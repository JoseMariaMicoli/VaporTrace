```markdown

    __   __                     _____                      
    \ \ / /___ _ __  ___  _ __ |_   _| __ __ _  ___ ___ 
     \ V // _ `| '_ \/ _ \| '__|  | || '__/ _` |/ __/ _ \
      \  / (_| | |_)  (_) | |     | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|     |_||_|  \__,_|\___\___|
               |_|      [ API INFRASTRUCTURE TRACER ]

```

**VaporTrace** is a high-performance Red Team framework engineered in Go for surgical reconnaissance and exploitation of API architectures. It specializes in uncovering "Shadow APIs," analyzing authorization logic (BOLA/BFLA), and mapping the entire attack surface of modern REST/Microservice environments.

---

## ⚠️ FULL LEGAL DISCLAIMER & RULES OF ENGAGEMENT

**THIS TOOL IS FOR AUTHORIZED PENETRATION TESTING AND EDUCATIONAL PURPOSES ONLY.**

1. **Authorization Required:** Never use VaporTrace against targets you do not have explicit, written permission to test.
2. **No Liability:** The author (JoseMariaMicoli) and contributors assume no liability and are not responsible for any misuse, data loss, service degradation, or legal consequences caused by this program.
3. **Local Laws:** It is the user's responsibility to comply with all applicable local, state, and international laws.
4. **Logic Risk:** Be aware that automated BOLA/BFLA probing can modify server-side data. Always perform tests in a controlled staging environment when possible.

**By compiling or running this software, you agree to these terms.**

---

## 🚀 Strategic Roadmap

### **Phase 1: The Foundation [STABLE]**

* [x] **Cobra CLI Engine:** Subcommand-based architecture (`map`, `scan`, `auth`).
* [x] **The Burp Bridge:** Industrial-strength HTTP client with native proxy support.
* [x] **SSL/TLS Hardening:** Automatic bypass of self-signed certs for intercepting proxies.
* [x] **Global Config:** Persistent flag management for headers and authentication.

### **Phase 2: Discovery & Inventory (API9)**

* [x] **Spec Ingestion:** Automated parsing of Swagger (v2) and OpenAPI (v3) definitions.
* [ ] **JS Route Scraper:** Regex-based endpoint extraction from client-side JavaScript bundles.
* [x] **Version Walker:** Identification of deprecated versions (e.g., `/v1/` vs `/v2/`) to find unpatched logic.
* [ ] **Parameter Miner:** Automatic identification of hidden query parameters and headers.

### **Phase 3: Authorization & Logic (API1, API3, API5)**

* [ ] **BOLA Prober (API1):** High-speed UUID/ID swapping using multi-token authentication contexts.
* [ ] **BOPLA/Mass Assignment (API3):** Fuzzing JSON bodies for administrative or hidden properties.
* [ ] **BFLA Module (API5):** Testing hierarchical access via HTTP method manipulation (GET vs DELETE).

### **Phase 4: Consumption & Injection (API4, API7, API8, API10)**

* [ ] **Resource Exhaustion (API4):** Probing pagination limits and payload size constraints.
* [ ] **SSRF Tracker (API7):** Detecting out-of-band callbacks via URL-parameter injection.
* [ ] **Security Misconfig (API8):** Automated CORS, Security Header, and Verbose Error audit.
* [ ] **Integration Probe (API10):** Identifying unsafe consumption in webhooks and 3rd party triggers.

---

## 🛠️ Installation & Usage

### 1. Build from Source

```bash
go mod tidy
go build -o VaporTrace

```

### 2. Verify Infrastructure (Proxy Test)

Ensure your intercepting proxy (e.g., Burp Suite) is listening on `127.0.0.1:8080`.

```bash
./VaporTrace map --proxy [http://127.0.0.1:8080](http://127.0.0.1:8080)

```

---

## 📡 The Technology Behind the Tracer

* **Language:** Golang (Concurrency-focused, statically linked).
* **CLI Framework:** Cobra (POSIX-compliant flags and subcommands).
* **Network Stack:** Custom `net/http` wrapper with `crypto/tls` overrides for Red Team flexibility.

**VaporTrace - Reveal the Invisible.**

---