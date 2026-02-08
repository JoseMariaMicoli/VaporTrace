```
     __  __                         _____                    
    \ \ / /___  _ __  ___  _ __   |_   _| __ __ _  ___ ___ 
     \ V // _ `| '_ \/ _ \| '__|    | || '__/ _` |/ __/ _ \
      \  / (_| | |_)  (_) | |       | || | | (_| | (_|  __/
       \/ \__,_| .__/\___/|_|       |_||_|  \__,_|\___\___|
               |_|      [ Surgical API Exploitation Suite ]
```

# VaporTrace v3.1-Hydra

**Enterprise-Grade API Security Testing Platform**

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go) ![License](https://img.shields.io/badge/License-Custom-red?style=flat-square) ![Status](https://img.shields.io/badge/Status-Production-brightgreen?style=flat-square) ![OWASP](https://img.shields.io/badge/OWASP-API%20Top%2010-blue?style=flat-square) ![MITRE](https://img.shields.io/badge/MITRE-ATT%26CK-orange?style=flat-square) ![NIST](https://img.shields.io/badge/NIST-CSF%20v2.0-purple?style=flat-square)

---

## 🎯 Executive Summary

**VaporTrace** is a surgical API exploitation and risk intelligence platform engineered for **penetration testers**, **red teams**, and **security researchers**. It automates the discovery, analysis, and exploitation of OWASP API Top 10 vulnerabilities with human-in-the-loop orchestration and neural engine-driven payload generation.

### Key Capabilities

- **Automated API discovery** - Endpoint mapping, shadow API detection, parameter fuzzing
- **Authorization testing** - BOLA, BFLA, BOPLA, SSRF, exhaustion, misconfig, webhook injection
- **AI-powered payloads** - Neural engine with cloud (Groq) + local (Mistral) hybrid architecture
- **Tactical TUI dashboard** - 7-tab real-time monitoring with 19 keyboard shortcuts
- **Compliance reporting** - NIST CSF v2.0, MITRE ATT&CK, OWASP API Top 10, CVSS scoring
- **Enterprise features** - Multi-operator coordination, mission vault (SQLite), evasion techniques

---

## 🔧 Technology Stack

| Component | Technology | Details |
|-----------|-----------|---------|
| **Core Language** | ![Go](https://img.shields.io/badge/Go_1.21+-00ADD8?style=flat-square) | High-performance, concurrent, statically compiled |
| **TUI Dashboard** | ![rivo/tview](https://img.shields.io/badge/rivo/tview-Terminal_UI-blue?style=flat-square) | Real-time terminal interface with F-key navigation |
| **Database** | ![SQLite3](https://img.shields.io/badge/SQLite3-Mission_Vault-003B57?style=flat-square) | Persistent mission storage, audit logs, findings |
| **AI Engine** | ![Hybrid](https://img.shields.io/badge/Hybrid-Groq/Ollama-purple?style=flat-square) | Primary: Groq API (cloud LLM), Fallback: Mistral (local) |
| **HTTP Client** | ![Native](https://img.shields.io/badge/Native-http.RoundTripper-orange?style=flat-square) | Middleware-driven, proxy routing, TLS hardening |
| **Networking** | ![Multi-Protocol](https://img.shields.io/badge/HTTP/HTTPS/SOCKS5-green?style=flat-square) | Proxy rotation, IP masking, evasion support |
| **Reporting** | ![Compliance](https://img.shields.io/badge/Markdown/PDF-NIST_MITRE_OWASP-red?style=flat-square) | Framework-aligned findings export |

---

## ⚡ Core Modules

### **Reconnaissance & Discovery** (OWASP API9)
\`\`\`
target          Set global scope URL
map             Automated endpoint discovery (spider + swagger + scrape + mine)
swagger         Parse OpenAPI/Swagger specifications
scrape          Extract endpoints from JavaScript bundles
mine            Brute-force hidden query parameters
sessions        Manage authentication tokens and credentials
\`\`\`

### **Authorization Testing** (OWASP API1-5)
\`\`\`
bola            Broken Object Level Authorization (ID enumeration)
bfla            Broken Function Level Authorization (privilege escalation)
bopla           Broken Object Property Authorization (mass assignment)
flow            Orchestrate multi-step attack chains
\`\`\`

### **Infrastructure & Data Testing** (OWASP API4, API7, API8, API10)
\`\`\`
exhaust         Resource exhaustion & DoS testing
ssrf            Server-side request forgery (cloud metadata)
audit           Security configuration auditing
probe           Webhook & third-party integration testing
\`\`\`

### **Tactical Planning** (HITL Orchestration)
\`\`\`
analyze         AI-driven tactical plan generation
list-plan       Review pending tactical actions
edit <id>       Override AI payloads
drop <id>       Mark action to skip execution
commit          Execute all pending actions
\`\`\`

### **AI & Neural Engine**
\`\`\`
neuro on        Activate neural engine (LLM payloads)
neuro off       Disable AI mutations
neuro-gen <n>   Generate n alternative payloads
test-neuro      Test AI provider connectivity
ask <prompt>    Direct LLM query
\`\`\`

### **Infrastructure & Reporting**
\`\`\`
proxy           Set upstream proxy
proxies load    Load rotating proxy list
report          Generate findings report (Markdown/PDF)
loot            Manage captured secrets and credentials
init_db         Initialize mission database
reset_db        Clear all mission data
\`\`\`

**For complete command documentation, see:** [Command Reference](docs/manuals/18_COMMAND_REFERENCE.md)

---

## 🎮 Keyboard Shortcuts

| Key | Function | Scope |
|-----|----------|-------|
| **F1-F7** | Toggle tabs (LOGS, MAP, LOOT, ANALYSIS, PLAN, NEURO, SETTINGS) | Global |
| **Ctrl+H** | Help modal (all hotkeys) | Global |
| **Ctrl+I** | Interceptor MITM modal | Global |
| **Ctrl+F** | Search current tab | Global |
| **Ctrl+D** | Debug mode toggle | Global |
| **Ctrl+S** | Save session | Global |
| **Ctrl+A** | Select all items | Tab context |
| **Ctrl+B** | Batch operations | Tab context |
| **Ctrl+X** | Export current view | Tab context |
| **Esc** | Close modal | Modal context |
| **Page Up/Down** | Scroll lists | List context |

**For complete hotkey reference, see:** [Keyboard Shortcuts](docs/manuals/17_KEYBOARD_SHORTCUTS.md)

---

## 🛡️ Compliance Frameworks

### OWASP API Security Top 10 (2023)
- ✅ API1: Broken Object Level Authorization (BOLA)
- ✅ API2: Broken Authentication (BFLA)
- ✅ API3: Broken Object Property Level Authorization (BOPLA)
- ✅ API4: Unrestricted Resource Consumption (Exhaustion)
- ✅ API5: Broken Function Level Authorization
- ✅ API7: Server-Side Request Forgery (SSRF)
- ✅ API8: Security Misconfiguration (Audit)
- ✅ API9: Improper Inventory Management (Map/Discovery)
- ✅ API10: Unsafe Consumption of APIs (Probe)

### Framework Mapping
| Framework | VaporTrace Implementation | Coverage |
|-----------|--------------------------|----------|
| **MITRE ATT&CK®** | T-code tagging in reports | 100% |
| **NIST CSF v2.0** | Govern (GV), Detect (DE), Protect (PR) | Full |
| **CVSS v3.1/4.0** | Automated severity scoring | Complete |

---

## 📚 Documentation

**Complete documentation is located in the \`/docs\` folder:**

### User Guides (docs/manuals/)
- **[Installation & Setup](docs/manuals/01_INSTALLATION_SETUP.md)** - Installation methods, dependencies, configuration
- **[First Run Tutorial](docs/manuals/02_FIRST_RUN.md)** - Step-by-step 11-step walkthrough
- **[UI Dashboard Guide](docs/manuals/03_UI_OVERVIEW.md)** - Tabs, layout, navigation, performance
- **[Strategic Planning](docs/manuals/04_STRATEGIC_PLANNING.md)** - HITL workflow and orchestration
- **[Keyboard Shortcuts](docs/manuals/17_KEYBOARD_SHORTCUTS.md)** - All 19 hotkeys with examples
- **[Command Reference](docs/manuals/18_COMMAND_REFERENCE.md)** - 40+ commands with parameters and examples
- **[Master Index](docs/manuals/INDEX.md)** - Navigation hub for all 20 user guides

### Technical Documentation (docs/dev-logs/)
- **[Architecture Overview](docs/dev-logs/INDEX.md)** - System design, layers, data flows
- **[Documentation Status](docs/DOCUMENTATION_STATUS.md)** - Completion tracker

### Quick Start
1. **New users:** Start with [First Run Tutorial](docs/manuals/02_FIRST_RUN.md)
2. **Dashboard help:** See [UI Dashboard Guide](docs/manuals/03_UI_OVERVIEW.md)
3. **Attack planning:** See [Strategic Planning](docs/manuals/04_STRATEGIC_PLANNING.md)
4. **Commands:** See [Command Reference](docs/manuals/18_COMMAND_REFERENCE.md)
5. **Hotkeys:** See [Keyboard Shortcuts](docs/manuals/17_KEYBOARD_SHORTCUTS.md)

**This README covers architecture and overview. For complete operational details, refer to the documentation in \`/docs\`.**

---

## 🚀 Installation

### Quick Start
\`\`\`bash
# Clone repository
git clone https://github.com/JoseMariaMicoli/VaporTrace.git
cd VaporTrace

# Install dependencies
go mod tidy

# Build binary
go build -o VaporTrace main.go

# Run
./VaporTrace
\`\`\`

**For detailed installation guide, see:** [Installation & Setup](docs/manuals/01_INSTALLATION_SETUP.md)

---

## 🎯 Use Cases

### **Security Teams**
- Automated API security testing in SDLC
- Risk assessment of internal microservices
- Compliance validation (NIST, OWASP, MITRE)

### **Penetration Testers**
- Comprehensive API attack surface discovery
- Authorization logic validation
- WAF bypass testing via neural engine
- Tactical engagement orchestration with reporting

### **Red Teams**
- Coordinated API exploitation with HITL control
- Real-time interception and payload modification
- Cloud metadata exploitation (SSRF/IMDS)
- Evasion techniques (IP rotation, timing jitter, process masking)

### **Security Researchers**
- API vulnerability research and documentation
- Custom exploitation module development
- Attack technique validation and documentation
- Framework-aligned finding publication

---

## 🌐 Evasion & Detection Bypass

VaporTrace includes built-in evasion for modern defensive environments:

| Technique | Purpose | Implementation |
|-----------|---------|-----------------|
| **Header Randomization** | Bypass signature detection | Rotating User-Agents, JA3 fingerprints |
| **IP Rotation** | Mask origin | SOCKS5/HTTP proxy rotation |
| **Timing Jitter** | Evade rate-limiting | Dynamic inter-packet delays |
| **Process Masquerade** | Hide from EDR | Rename to \`kworker_system_auth\` |
| **Protocol Obfuscation** | Bypass IPS/WAF | Custom header injection |

---

## 🏗️ Architecture

### 6-Layer System Design
\`\`\`
┌─────────────────────────────────────────┐
│ 1. Tactical UI (rivo/tview)            │ Real-time dashboard
├─────────────────────────────────────────┤
│ 2. Command Engine                      │ 40+ commands, orchestration
├─────────────────────────────────────────┤
│ 3. Execution Engine                    │ Module coordination, concurrency
├─────────────────────────────────────────┤
│ 4. Exploitation Modules                │ BOLA, BFLA, BOPLA, SSRF, etc.
├─────────────────────────────────────────┤
│ 5. Networking Layer                    │ HTTP, proxy, TLS hardening
├─────────────────────────────────────────┤
│ 6. Persistence & AI                    │ SQLite, Groq/Ollama, reporting
└─────────────────────────────────────────┘
\`\`\`

### Data Flow
1. **Discovery** → Endpoints discovered → Stored in database
2. **Analysis** → AI examines endpoints → Generates tactical plan
3. **Planning** → Human reviews plan → Edits/approves actions
4. **Execution** → Actions executed → Findings captured
5. **Reporting** → Results → Framework-mapped report

### AI & Neural Engine
- **Primary Provider:** Groq (fast cloud LLM)
- **Fallback Provider:** Mistral via Ollama (local, private)
- **Failover:** Automatic cloud→local on quota/connection errors
- **Capabilities:** Contextual payloads, pattern analysis, WAF bypass strategies

---

## ⚠️ Legal & Authorization

**THIS TOOL IS FOR AUTHORIZED PENETRATION TESTING & EDUCATIONAL PURPOSES ONLY**

1. **Authorization Required** - Only test targets with explicit written permission
2. **Compliance** - You are responsible for all applicable laws and regulations
3. **No Liability** - Author assumes no responsibility for misuse or consequences
4. **Testing Safety** - Some modules modify server state; test in controlled environments

**By using this software, you agree to these terms.**

---

## 🔗 Contributing

- **Issues & Feedback:** GitHub Issues
- **Development:** Technical documentation in \`/docs/dev-logs/\`
- **Questions:** See FAQ in documentation

---

## 📄 License

Custom License - See LICENSE file

---

**VaporTrace** — Surgical API Exploitation Platform  
**Version:** 3.1.2-Hydra | **Status:** Production Ready

For complete documentation, see [docs/manuals/INDEX.md](docs/manuals/INDEX.md)
