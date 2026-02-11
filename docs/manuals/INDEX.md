# VaporTrace User Manual - Complete Index

**Version:** 3.0+ (Post Sprint 11)  
**Last Updated:** February 8, 2026  
**Status:** Complete & Comprehensive

---

## Table of Contents

### 🚀 Getting Started
1. [Installation & Setup](01_INSTALLATION_SETUP.md) - Initial configuration and dependencies
2. [First Run Guide](02_FIRST_RUN.md) - Your first VaporTrace session
3. [UI Overview](03_UI_OVERVIEW.md) - Dashboard tabs and layout

### 🎯 Core Features
4. [Strategic Planning (HITL)](04_STRATEGIC_PLANNING.md) - Human-in-the-loop attack orchestration
5. [Reconnaissance & Discovery](05_RECONNAISSANCE.md) - Building your attack surface map
6. [Exploitation Engine](06_EXPLOITATION.md) - OWASP API Top 10 attack modules
7. [AI & Neural Engine](07_AI_NEURO_ENGINE.md) - ML-driven payload generation
7a. [Neuro Quick Usage Guide](NEURO_QUICK_USAGE_GUIDE.md) - Complete workflows and examples for AI features
8. [Interceptor & MITM](08_INTERCEPTOR_MITM.md) - Request manipulation and analysis

### 🛠️ Advanced Features
9. [Attack Chains (Flow)](09_ATTACK_CHAINS.md) - Orchestrated attack sequences
10. [Ghost Weaver (Evasion)](10_GHOST_WEAVER.md) - Token forgery and data masking
11. [Loot Vault (Exfiltration)](11_LOOT_VAULT.md) - Secret capture and management
12. [Proxy & Network](12_PROXY_NETWORK.md) - Upstream proxies and traffic routing
12a. [QUICK_START_RACE.md](QUICK_START_RACE.md) - Race Condition Testing (Tier 3) ⭐ NEW

### ⭐ Tier 4: Advanced Orchestration (Sprint 20)
23. [Intelligence Layer (OSINT)](23_INTEL_OSINT.md) - Wayback Machine, Shodan, external intelligence
24. [Chain Reactor](24_CHAIN_REACTOR.md) - Multi-step stateful attack workflows
25. [Value Extractor](25_EXTRACTOR.md) - Data extraction from responses (JSON, Regex, Cookies)
26. [Knowledge Base (Institutional Memory)](26_KNOWLEDGE_BASE.md) - Record and learn from successful exploits

### 📊 Reporting & Analysis
13. [Report Generation](13_REPORTING.md) - Export findings and generate reports
14. [Dashboard Analytics](14_ANALYTICS.md) - Metrics and insights

### 🔧 Configuration & Troubleshooting
15. [Configuration](15_CONFIGURATION.md) - Settings, proxies, authentication
16. [Troubleshooting](16_TROUBLESHOOTING.md) - Common issues and solutions
17. [Keyboard Shortcuts](17_KEYBOARD_SHORTCUTS.md) - Complete hotkey reference

### 📚 Reference
18. [Command Reference](18_COMMAND_REFERENCE.md) - All CLI commands with examples
**22. [Advanced Discovery Guide](22_DISCOVERY_GUIDE.md) - Spider & Fuzz techniques (Tier 2) ⭐ NEW**
19. [API Module Documentation](19_API_MODULES.md) - Detailed module descriptions
20. [FAQ & Tips](20_FAQ_TIPS.md) - Frequently asked questions and pro tips
21. [WAF Evasion Techniques](21_WAF_EVASION_TECHNIQUES.md) - Advanced WAF bypass strategies

---

## Quick Navigation by Use Case

### 🔍 "I want to discover API endpoints"
**Read:** [05_RECONNAISSANCE.md](05_RECONNAISSANCE.md) and **[22_DISCOVERY_GUIDE.md](22_DISCOVERY_GUIDE.md)** ⭐ NEW
- Target management
- Automatic endpoint discovery (map, swagger, scrape)
- **Advanced: Domain crawling with spider (Tier 2) ⭐ NEW**
- **Advanced: Brute-force fuzzing with fuzz (Tier 2) ⭐ NEW**
- Parameter fuzzing
- Swagger/OpenAPI parsing

### 🕷️ "I want to crawl a domain for endpoints"
**Read:** **[22_DISCOVERY_GUIDE.md](22_DISCOVERY_GUIDE.md)** ⭐ NEW
- Spider command for recursive crawling
- Depth control and optimization
- WAF evasion for crawling
- Performance tuning
- Real-world examples

### 🔎 "I want to fuzz for hidden paths/parameters"
**Read:** **[22_DISCOVERY_GUIDE.md](22_DISCOVERY_GUIDE.md)** ⭐ NEW
- Fuzz command for path enumeration
- Parameter discovery with anomaly detection
- Embedded wordlists
- Speed vs stealth tradeoffs
- Troubleshooting guides

### ⚔️ "I want to test for vulnerabilities"
**Read:** [06_EXPLOITATION.md](06_EXPLOITATION.md)
- BOLA (Broken Object Level Auth)
- BFLA (Broken Function Level Auth)
- BOPLA (Broken Object Property Level Auth)
- SSRF, Resource Exhaustion, etc.

### 🧠 "I want to use AI to generate payloads"
**Read:** [07_AI_NEURO_ENGINE.md](07_AI_NEURO_ENGINE.md) and [NEURO_QUICK_USAGE_GUIDE.md](NEURO_QUICK_USAGE_GUIDE.md)
- Neural Engine setup and configuration
- AI payload generation with multiple providers
- Complete workflows and examples
- Troubleshooting and best practices
- 5 AI providers: Groq, OpenAI, Gemini, Ollama, Hybrid mode

### 🎯 "I want to orchestrate complex attacks"
**Read:** [04_STRATEGIC_PLANNING.md](04_STRATEGIC_PLANNING.md) & [09_ATTACK_CHAINS.md](09_ATTACK_CHAINS.md)
- Tactical planning
- Action editing
- Attack chains
- Commit workflows

### 🚨 "I want to intercept and modify requests"
**Read:** [08_INTERCEPTOR_MITM.md](08_INTERCEPTOR_MITM.md)
- Interceptor configuration
- Request/response editing
- Payload injection
- Modal controls

### 👻 "I want to evade detection"
**Read:** [10_GHOST_WEAVER.md](10_GHOST_WEAVER.md)
- OIDC token interception
- Data masking
- Evasion techniques
- Payload obfuscation

### 💎 "I want to manage captured secrets"
**Read:** [11_LOOT_VAULT.md](11_LOOT_VAULT.md)
- Loot discovery
- Secret management
- Credential usage
- Export options

### 🌐 "I want to find ghost endpoints with OSINT" (Tier 4 Day 1)
**Read:** [23_INTEL_OSINT.md](23_INTEL_OSINT.md) ⭐ NEW
- Wayback Machine historical URLs
- Shodan infrastructure discovery
- Legacy API version detection
- Ghost endpoint findings
- OSINT workflow integration

### 🔗 "I want to orchestrate multi-step attacks" (Tier 4 Day 2)
**Read:** [24_CHAIN_REACTOR.md](24_CHAIN_REACTOR.md) ⭐ NEW
- Stateful attack chains
- Variable extraction and injection
- Authentication flow automation
- CSRF bypass workflows
- Privilege escalation sequences

### 📊 "I want to extract data from responses" (Tier 4 Day 2)
**Read:** [25_EXTRACTOR.md](25_EXTRACTOR.md) ⭐ NEW
- JSON path extraction
- Regular expression patterns
- Cookie/Header extraction
- Variable storage and reuse
- Integration with chains

### 🏃 "I want to test for race conditions" (Tier 3)
**Read:** [QUICK_START_RACE.md](QUICK_START_RACE.md) ⭐ NEW
- Race condition basics
- Synchronization gate pattern
- TOCTOU vulnerability detection
- Common race scenarios
- Troubleshooting guide

### 🔫 "I want to fuzz parameters with custom wordlists" (Tier 3)
**Read:** [QUICK_START_RACE.md](QUICK_START_RACE.md) & [18_COMMAND_REFERENCE.md](18_COMMAND_REFERENCE.md) ⭐ NEW
- Intruder sniper mode
- Anomaly detection
- Custom wordlist loading
- Baseline comparison

---

## Step-by-Step Workflows

### Workflow 1: Basic API Reconnaissance (15 minutes)
1. [Installation & Setup](01_INSTALLATION_SETUP.md) - Get started
2. [First Run Guide](02_FIRST_RUN.md) - Launch VaporTrace
3. [05_RECONNAISSANCE.md](05_RECONNAISSANCE.md) - Target, Map, Swagger

**Commands:**
```bash
target https://api.example.com
map
swagger https://api.example.com/swagger.json
```

### Workflow 2: HITL Attack Planning & Execution (30 minutes)
1. [04_STRATEGIC_PLANNING.md](04_STRATEGIC_PLANNING.md) - Understand HITL
2. [05_RECONNAISSANCE.md](05_RECONNAISSANCE.md) - Discover endpoints
3. [04_STRATEGIC_PLANNING.md](04_STRATEGIC_PLANNING.md) - Analyze & plan
4. [06_EXPLOITATION.md](06_EXPLOITATION.md) - Execute attacks

**Commands:**
```bash
target https://api.example.com
map
analyze
list-plan
commit
```

### Workflow 3: Race Condition Testing (Tier 3) (15 minutes)
1. [QUICK_START_RACE.md](QUICK_START_RACE.md) - Learn race testing basics
2. [Identify target endpoint](#) - Find state-changing endpoints
3. Run race test - Execute concurrent requests
4. Analyze results - Check for timing-based vulnerabilities

**Commands:**
```bash
target https://api.example.com
race https://api.example.com/api/claim?code=WINNER 30
report
```

### Workflow 4: AI-Driven Exploitation (20 minutes)
1. [07_AI_NEURO_ENGINE.md](07_AI_NEURO_ENGINE.md) - Setup neural engine
2. [06_EXPLOITATION.md](06_EXPLOITATION.md) - Run attacks with AI
3. [13_REPORTING.md](13_REPORTING.md) - Generate report

**Commands:**
```bash
neuro on
target https://api.example.com
ssrf
report
```

### Workflow 4: AI-Driven Exploitation (20 minutes)
1. [07_AI_NEURO_ENGINE.md](07_AI_NEURO_ENGINE.md) - Setup neural engine
2. [06_EXPLOITATION.md](06_EXPLOITATION.md) - Run attacks with AI
3. [13_REPORTING.md](13_REPORTING.md) - Generate report

**Commands:**
```bash
neuro on
target https://api.example.com
ssrf
report
```

### Workflow 5: Learning & Scaling with Knowledge Base (Tier 4 Day 3) (30 minutes)
1. [26_KNOWLEDGE_BASE.md](26_KNOWLEDGE_BASE.md) - Learn from successful exploits
2. [06_EXPLOITATION.md](06_EXPLOITATION.md) - Test target A
3. [26_KNOWLEDGE_BASE.md](26_KNOWLEDGE_BASE.md) - Record successful vectors
4. [06_EXPLOITATION.md](06_EXPLOITATION.md) - Test target B with learned patterns
5. [07_AI_NEURO_ENGINE.md](07_AI_NEURO_ENGINE.md) - AI mutates based on KB

**Commands:**
```bash
target https://api.target-a.com
bola /users/999                    # Success
kb add BOLA /users/{id} GET id=999
neuro on
target https://api.target-b.com
neuro-gen BOLA 5                   # AI learns from KB
bola /users/999                    # Success again (faster!)
```

### Workflow 6: Request Interception & Manipulation (25 minutes)
1. [08_INTERCEPTOR_MITM.md](08_INTERCEPTOR_MITM.md) - Setup interceptor
2. [06_EXPLOITATION.md](06_EXPLOITATION.md) - Run attack with interception
3. [11_LOOT_VAULT.md](11_LOOT_VAULT.md) - Review captured loot

**Commands:**
```bash
Ctrl+I (toggle interceptor)
target https://api.example.com
bola /api/users/
[Intercept each request, modify, forward]
```

### Workflow 6: Learning From Exploits (Tier 4 Day 3) (15 minutes)
1. [26_KNOWLEDGE_BASE.md](26_KNOWLEDGE_BASE.md) - Setup Knowledge Base
2. [06_EXPLOITATION.md](06_EXPLOITATION.md) - Execute successful attacks
3. [26_KNOWLEDGE_BASE.md](26_KNOWLEDGE_BASE.md) - Record vectors and AI learning

**Commands:**
```bash
target https://api.example.com
bola /users/999                 # ✓ Success
kb add BOLA /users/{id} GET id=999

# On next target:
neuro on
neuro-gen BOLA 10              # AI learns from KB
bola /users/999                # Faster exploitation!
```

---

## Feature Matrix

| Feature | Category | Difficulty | Time | Read |
|---------|----------|------------|------|------|
| Endpoint Discovery | Recon | Easy | 5min | [05](05_RECONNAISSANCE.md) |
| BOLA Testing | Exploit | Easy | 10min | [06](06_EXPLOITATION.md) |
| Tactical Planning | HITL | Medium | 15min | [04](04_STRATEGIC_PLANNING.md) |
| AI Payload Gen | AI | Medium | 20min | [07](07_AI_NEURO_ENGINE.md) |
| Request Interception | Advanced | Medium | 20min | [08](08_INTERCEPTOR_MITM.md) |
| Attack Chains | Advanced | Hard | 30min | [09](09_ATTACK_CHAINS.md) |
| Token Forgery | Evasion | Hard | 25min | [10](10_GHOST_WEAVER.md) |
| Multi-Module Orchestration | Advanced | Hard | 40min | [04](04_STRATEGIC_PLANNING.md) + [09](09_ATTACK_CHAINS.md) |

---

## Module Quick Reference

### Reconnaissance Modules
- **map** - Full auto-discovery (spiderer + swagger + JS scraping)
- **swagger** - Parse OpenAPI specifications
- **scrape** - Extract endpoints from JavaScript
- **mine** - Brute-force parameters
- **sessions** - Manage auth tokens

### Exploitation Modules
- **bola** - ID enumeration
- **bfla** - Privilege escalation
- **bopla** - Mass assignment
- **ssrf** - Internal access
- **exhaust** - Resource DoS
- **audit** - Configuration audit
- **probe** - Webhook injection
- **flow** - Attack chain orchestration

### AI & Advanced
- **neuro** - Enable/disable neural engine
- **neuro-gen** - Generate n payloads
- **test-neuro** - Connectivity test
- **ask** - Direct LLM query
- **weaver** - Token forgery & evasion

---

## Keyboard Shortcuts Quick Ref

| Shortcut | Function |
|----------|----------|
| **Ctrl+H** | Show keybindings popup |
| **F1-F7** | Switch tabs |
| **Ctrl+I** | Toggle Interceptor |
| **Ctrl+A** | Analyze with AI |
| **Esc** | Exit VaporTrace |

**Full Reference:** [17_KEYBOARD_SHORTCUTS.md](17_KEYBOARD_SHORTCUTS.md)

---

## Common Questions

**Q: How do I start a penetration test?**  
A: See [Workflow 2](INDEX.md#workflow-2-hitl-attack-planning--execution-30-minutes) - HITL Attack Planning

**Q: How do I intercept requests?**  
A: See [08_INTERCEPTOR_MITM.md](08_INTERCEPTOR_MITM.md) - Toggle Ctrl+I

**Q: How do I use AI to generate payloads?**  
A: See [07_AI_NEURO_ENGINE.md](07_AI_NEURO_ENGINE.md) - Neural Engine Setup

**Q: How do I export findings?**  
A: See [13_REPORTING.md](13_REPORTING.md) - Report Generation

**For more Q&A:** [20_FAQ_TIPS.md](20_FAQ_TIPS.md)

---

## Getting Help

1. **In-App Help:**
   - Press `Ctrl+H` for keybindings
   - Run `help <command>` for command help
   - Run `usage` for all commands

2. **User Documentation:**
   - Browse this manual (files listed above)
   - Check [Troubleshooting](16_TROUBLESHOOTING.md)
   - See [FAQ & Tips](20_FAQ_TIPS.md)

3. **Technical Documentation:**
   - See [dev-logs/](../dev-logs/INDEX.md) for architecture
   - See [dev-logs/IMPLEMENTATION.md](../dev-logs/IMPLEMENTATION.md) for deep dives

---

## Document Organization

```
docs/
├── DOCUMENTATION_STATUS.md
├── TUI.png
├── diagram_mermaid.png
├── manuals/
│   ├── INDEX.md (you are here)
│   ├── 01_INSTALLATION_SETUP.md
│   ├── 02_FIRST_RUN.md
│   ├── 03_UI_OVERVIEW.md
│   ├── 04_STRATEGIC_PLANNING.md
│   ├── 05_RECONNAISSANCE.md
│   ├── 06_EXPLOITATION.md
│   ├── 07_AI_NEURO_ENGINE.md
│   ├── NEURO_QUICK_USAGE_GUIDE.md (✨ NEW - Complete neuro workflows)
│   ├── 08_INTERCEPTOR_MITM.md
│   ├── 09_ATTACK_CHAINS.md
│   ├── 10_GHOST_WEAVER.md
│   ├── 11_LOOT_VAULT.md
│   ├── 12_PROXY_NETWORK.md
│   ├── 12_SPRINT12_EVASION_V2.md
│   ├── 13_REPORTING.md
│   ├── 13_SPRINT13_EVASION_HARDENING.md
│   ├── 14_ANALYTICS.md
│   ├── 15_CONFIGURATION.md
│   ├── 16_TROUBLESHOOTING.md
│   ├── 17_KEYBOARD_SHORTCUTS.md
│   ├── 18_COMMAND_REFERENCE.md
│   ├── 19_API_MODULES.md
│   ├── 20_FAQ_TIPS.md
│   ├── 21_WAF_EVASION_TECHNIQUES.md
│   ├── KEYBINDINGS_QUICK_REFERENCE.md
│   └── SPRINT12_INTEGRATION_GUIDE.md
└── dev-logs/
    ├── INDEX.md
    ├── 00_START_HERE.md (✨ NEW - Navigation hub)
    ├── NEURO_QUICK_REFERENCE.md (✨ NEW - One-page reference)
    ├── NEURO_SOURCE_FIXES_NEEDED.md (✨ NEW - Implementation guide)
    ├── NEURO_COMPLETE_STATUS.md (✨ NEW - Status & workflows)
    ├── NEURO_AUDIT_REPORT.md
    ├── NEURO_AUDIT_EXECUTIVE_SUMMARY.md
    ├── NEURO_AUDIT_VISUAL_MAP.md
    ├── AUDIT_SUMMARY_FOR_TEAM.md
    ├── Dev-Roadmap.md
    ├── BUG-Fixes/
    ├── Sprint-01/
    ├── Sprint-02/
    ├── Sprint-03/
    ├── Sprint-04/
    ├── Sprint-05/
    ├── Sprint-06/
    ├── Sprint-07/
    ├── Sprint-08/
    ├── Sprint-09/
    ├── Sprint-10/
    ├── Sprint-11/
    ├── Sprint-12/
    ├── Sprint-13/
    ├── Sprint-14/
    ├── Sprint-15/
    └── Sprint-16/
```

---

## How to Use This Manual

1. **New User?** Start with [01_INSTALLATION_SETUP.md](01_INSTALLATION_SETUP.md) then [02_FIRST_RUN.md](02_FIRST_RUN.md)
2. **Want to do something specific?** Use "Quick Navigation by Use Case" above
3. **Need step-by-step?** Follow the "Workflows" section
4. **Looking for reference?** Jump to [18_COMMAND_REFERENCE.md](18_COMMAND_REFERENCE.md)
5. **Want to use AI/Neuro?** Check [NEURO_QUICK_USAGE_GUIDE.md](NEURO_QUICK_USAGE_GUIDE.md) for complete workflows
6. **Technical details?** See [dev-logs/INDEX.md](../dev-logs/INDEX.md) or [dev-logs/00_START_HERE.md](../dev-logs/00_START_HERE.md)
7. **Need quick AI reference?** Use [dev-logs/NEURO_QUICK_REFERENCE.md](../dev-logs/NEURO_QUICK_REFERENCE.md)

---

## 🎯 New: Comprehensive Neuro Engine Documentation (February 2026)

We've completely refreshed the Neuro Engine documentation and guides:

### For Users
- **[NEURO_QUICK_USAGE_GUIDE.md](NEURO_QUICK_USAGE_GUIDE.md)** - Start here! Complete workflows, provider setup, troubleshooting
  - 4 real-world attack scenarios
  - Setup instructions for 5 providers (Groq, OpenAI, Gemini, Ollama, Hybrid)
  - 2,500+ words of actionable guidance
  
### For Developers
- **[dev-logs/NEURO_SOURCE_FIXES_NEEDED.md](../dev-logs/NEURO_SOURCE_FIXES_NEEDED.md)** - Implementation guide

---

## 🔥 NEW: Tier 3 Offensive Capability Upgrade (February 11, 2026)

VaporTrace now includes **advanced fuzzing and race condition testing** for sophisticated logic flaw detection.

### What's New (Sprint 20)
- ✅ **Intruder Engine** - Automated fuzzing with anomaly detection
- ✅ **Race Condition Engine** - TOCTOU vulnerability testing with synchronization gate
- ✅ **AI Payload Generation** - Groq-driven fuzzing suggestions (Sprint 19)
- ✅ **Reporting Integration** - Tier 3 findings in F7 reports

### For Users
- **[QUICK_START_RACE.md](QUICK_START_RACE.md)** ⭐ **START HERE FOR TIER 3**
  - Race condition testing basics
  - Common TOCTOU vulnerabilities
  - Step-by-step usage guide
  - Troubleshooting tips

### For Developers
- **[../dev-logs/TIER_3_IMPLEMENTATION_SUMMARY.md](../dev-logs/TIER_3_IMPLEMENTATION_SUMMARY.md)** - Architecture overview
- **[../dev-logs/YOUR_ACTION_ITEMS.md](../dev-logs/YOUR_ACTION_ITEMS.md)** - Implementation checklist

### Quick Commands (Tier 3)
```bash
# Fuzzing with custom wordlists
intruder sniper https://api.example.com/user?id=1 id ./payloads.txt

# Race condition testing (20 threads by default)
race https://api.example.com/api/claim?code=WINNER

# High-intensity race test (50 threads)
race https://api.example.com/api/claim?code=WINNER 50

# View findings in report
report
```

### Severity & Remediation
- **Intruder Findings:** Medium/High severity
- **Race Condition Findings:** CRITICAL (CVSS 8.5+)
- **Remediation:** Flagged as `**ARCHITECTURAL FIX REQ**` (not simple patches)

---

## 📁 Tier 3 Documentation Structure

```
docs/
├── manuals/
│   ├── QUICK_START_RACE.md                    ⭐ User guide
│   └── INDEX.md                                (you are here)
└── dev-logs/
    ├── TIER_3_IMPLEMENTATION_SUMMARY.md       ⭐ Architecture
    ├── YOUR_ACTION_ITEMS.md                    Implementation checklist
    └── Dev-Roadmap.md                          Updated with Tier 3 & Tier 4 planning
```

---

## 🎯 Version History

- **v3.2-Chimera (Feb 11, 2026)** - Tier 3 Complete: Race Condition & Intruder Engines
- **v3.2-Chimera (Feb 8, 2026)** - Sprint 16-17: Blue-Team Mirror + WAF Hardening
- **v3.1-Hydra (Feb 1, 2026)** - Sprint 11 Complete: Full Autonomy & Neuro Integration
- **v3.0 (Jan 15, 2026)** - Core APIs + TUI Dashboard (Sprint 1-10)
  - All 6 critical code issues documented
  - Before/after code examples
  - 1.5-2 hour implementation roadmap
  
### Quick Reference
- **[dev-logs/NEURO_QUICK_REFERENCE.md](../dev-logs/NEURO_QUICK_REFERENCE.md)** - One-page reference
  - All commands at a glance
  - Feature status table
  
### Navigation Hub
- **[dev-logs/00_START_HERE.md](../dev-logs/00_START_HERE.md)** - Central navigation
  - Links organized by role (user/dev/manager)
  - Complete workflow diagrams

---

## Version History

- **3.1 (Feb 10, 2026)** - Neuro Engine refinement: Fixed race conditions, nil checks, improved prompts. Complete documentation refresh.
- **3.0+ (Feb 8, 2026)** - Post Sprint 11: TUI fixes, comprehensive documentation, AI integration stable
- **2.5 (Jan 2026)** - Added neural engine and AI payloads
- **2.0 (Dec 2025)** - HITL strategic planning added
- **1.0 (Oct 2025)** - Initial release with core exploitation modules

---

**Next Step:** [Installation & Setup](01_INSTALLATION_SETUP.md) →

