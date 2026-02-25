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
8. [Interceptor & MITM](08_INTERCEPTOR_MITM.md) - Request manipulation and analysis

### 🛠️ Advanced Features
9. [Attack Chains (Flow)](09_ATTACK_CHAINS.md) - Orchestrated attack sequences
10. [Ghost Weaver (Evasion)](10_GHOST_WEAVER.md) - Token forgery and data masking
11. [Loot Vault (Exfiltration)](11_LOOT_VAULT.md) - Secret capture and management
12. [Proxy & Network](12_PROXY_NETWORK.md) - Upstream proxies and traffic routing

### 📊 Reporting & Analysis
13. [Report Generation](13_REPORTING.md) - Export findings and generate reports
14. [Dashboard Analytics](14_ANALYTICS.md) - Metrics and insights

### 🔧 Configuration & Troubleshooting
15. [Configuration](15_CONFIGURATION.md) - Settings, proxies, authentication
16. [Troubleshooting](16_TROUBLESHOOTING.md) - Common issues and solutions
17. [Keyboard Shortcuts](17_KEYBOARD_SHORTCUTS.md) - Complete hotkey reference

### 📚 Reference
18. [Command Reference](18_COMMAND_REFERENCE.md) - All CLI commands with examples
19. [API Module Documentation](19_API_MODULES.md) - Detailed module descriptions
20. [FAQ & Tips](20_FAQ_TIPS.md) - Frequently asked questions and pro tips

---

## Quick Navigation by Use Case

### 🔍 "I want to discover API endpoints"
**Read:** [05_RECONNAISSANCE.md](05_RECONNAISSANCE.md)
- Target management
- Automatic endpoint discovery
- Parameter fuzzing
- Swagger/OpenAPI parsing

### ⚔️ "I want to test for vulnerabilities"
**Read:** [06_EXPLOITATION.md](06_EXPLOITATION.md)
- BOLA (Broken Object Level Auth)
- BFLA (Broken Function Level Auth)
- BOPLA (Broken Object Property Level Auth)
- SSRF, Resource Exhaustion, etc.

### 🧠 "I want to use AI to generate payloads"
**Read:** [07_AI_NEURO_ENGINE.md](07_AI_NEURO_ENGINE.md)
- Neural Engine setup
- AI payload generation
- Mutation strategies
- Configuration

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

### Workflow 3: AI-Driven Exploitation (20 minutes)
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

### Workflow 4: Request Interception & Manipulation (25 minutes)
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
├── manuals/
│   ├── INDEX.md (you are here)
│   ├── 01_INSTALLATION_SETUP.md
│   ├── 02_FIRST_RUN.md
│   ├── 03_UI_OVERVIEW.md
│   ├── 04_STRATEGIC_PLANNING.md
│   ├── 05_RECONNAISSANCE.md
│   ├── 06_EXPLOITATION.md
│   ├── 07_AI_NEURO_ENGINE.md
│   ├── 08_INTERCEPTOR_MITM.md
│   ├── 09_ATTACK_CHAINS.md
│   ├── 10_GHOST_WEAVER.md
│   ├── 11_LOOT_VAULT.md
│   ├── 12_PROXY_NETWORK.md
│   ├── 13_REPORTING.md
│   ├── 14_ANALYTICS.md
│   ├── 15_CONFIGURATION.md
│   ├── 16_TROUBLESHOOTING.md
│   ├── 17_KEYBOARD_SHORTCUTS.md
│   ├── 18_COMMAND_REFERENCE.md
│   ├── 19_API_MODULES.md
│   └── 20_FAQ_TIPS.md
└── dev-logs/
    ├── INDEX.md (architecture & technical overview)
    ├── 01_ARCHITECTURE.md
    ├── 02_IMPLEMENTATION.md
    ├── 03_MODULES_DETAILED.md
    ├── 04_DATA_FLOW.md
    ├── 05_AI_INTEGRATION.md
    └── DIAGRAMS/ (ASCII flow diagrams)
```

---

## How to Use This Manual

1. **New User?** Start with [01_INSTALLATION_SETUP.md](01_INSTALLATION_SETUP.md) then [02_FIRST_RUN.md](02_FIRST_RUN.md)
2. **Want to do something specific?** Use "Quick Navigation by Use Case" above
3. **Need step-by-step?** Follow the "Workflows" section
4. **Looking for reference?** Jump to [18_COMMAND_REFERENCE.md](18_COMMAND_REFERENCE.md)
5. **Technical details?** See [dev-logs/INDEX.md](../dev-logs/INDEX.md)

---

## Version History

- **3.0+ (Feb 8, 2026)** - Post Sprint 11: TUI fixes, comprehensive documentation, AI integration stable
- **2.5 (Jan 2026)** - Added neural engine and AI payloads
- **2.0 (Dec 2025)** - HITL strategic planning added
- **1.0 (Oct 2025)** - Initial release with core exploitation modules

---

**Next Step:** [Installation & Setup](01_INSTALLATION_SETUP.md) →

