![VaporTrace Logo](assets/images/VaporTrace_Logo.png)

# VaporTrace v3.2-Hydra

Enterprise-grade API security testing platform for red teams, penetration testers, and security engineers.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/License-Business_Source_1.1-red)
![Status](https://img.shields.io/badge/Status-Tier_4_Operational-brightgreen?style=flat-square)
![OWASP](https://img.shields.io/badge/OWASP-API%20Top%2010-blue?style=flat-square)
![MITRE](https://img.shields.io/badge/MITRE-ATT%26CK-orange?style=flat-square)
![NIST](https://img.shields.io/badge/NIST-CSF%20v2.0-purple?style=flat-square)
![AI](https://img.shields.io/badge/AI-Model_Agnostic-ff69b4?style=flat-square)

## License Update Notice
Starting from version 3.1.1, VaporTrace transitioned to **Business Source License 1.1 (BSL)**.

Reason:
To sustain development, support, and long-term maintenance while keeping source-visible transparency.

Previous versions remain under their original license terms.

---

## Table of Contents
- [What VaporTrace Is](#what-vaportrace-is)
- [Core Capabilities](#core-capabilities)
- [Architecture Overview](#architecture-overview)
- [Quick Start](#quick-start)
- [Operator Workflow](#operator-workflow)
- [Dashboard Tabs and Hotkeys](#dashboard-tabs-and-hotkeys)
- [Command Map](#command-map)
- [Plan Lifecycle and Persistence](#plan-lifecycle-and-persistence)
- [Neural Engine (F6)](#neural-engine-f6)
- [Reporting and Compliance](#reporting-and-compliance)
- [Project Structure](#project-structure)
- [Documentation](#documentation)
- [Roadmap](#roadmap)
- [Security and Usage Policy](#security-and-usage-policy)
- [Contributing](#contributing)
- [License](#license)

---

## What VaporTrace Is
VaporTrace is a tactical API offensive security platform that combines:
- API attack-surface discovery
- OWASP API Top 10 exploitation modules
- Human-in-the-loop strategic planning
- AI-assisted analysis and payload generation
- compliance-ready evidence reporting

It is built for engagements where speed, explainability, and operational control all matter.

---

## Core Capabilities

### Discovery and Reconnaissance
- OpenAPI/Swagger parsing (v2 and v3)
- JavaScript endpoint extraction
- endpoint/parameter discovery (`map`, `scrape`, `mine`, `spider`, `fuzz`)
- discovery-to-pipeline tagging for exploit routing

### Exploitation Engines
- BOLA, BFLA, BOPLA
- SSRF, resource exhaustion, misconfiguration audit
- integration/webhook probing
- intruder sniper mode and race condition testing

### Tactical Orchestration
- HITL action buffer (`analyze`, `list-plan`, `edit`, `drop`, `commit`)
- status-aware execution (`RUNNING`, `SUCCESS`, `FAILED`, `DROPPED`)
- auto follow-up plan generation after commit
- persisted tactical plan state (`list-plan` survives restart)

### AI and Neuro Layer
- model-agnostic provider strategy (cloud + local fallback)
- AI recommendation pass during strategic analysis
- payload mutation support and contextual analysis
- operator-visible neuro output in F6

### Evasion and Operational Stealth
- UA/header rotation and traffic mimicry
- jitter, thinking-time, rate-limit backoff
- path/payload obfuscation controls
- proxy and rotating pool support

### Evidence and Reporting
- SQLite-backed findings and context store
- markdown report generation
- OWASP, MITRE ATT&CK, NIST CSF alignment
- CVSS scoring support

---

## Architecture Overview

### Runtime layers
1. **UI Layer**: Hydra TUI (F1-F8) + shell mode.
2. **Engine Layer**: strategic planner, command router, action buffer.
3. **Logic Layer**: discovery, exploit modules, network middleware, loot scanning.
4. **Data Layer**: SQLite mission database (`findings`, `context_store`, `tactical_actions`, etc.).
5. **AI Layer**: neuro provider abstraction and tactical recommendation pipeline.

### Data flow (high level)
`Discovery -> Analyze -> Plan Buffer -> Commit -> Findings/Telemetry -> Next Plan`

---

## Quick Start

### Prerequisites
- Go 1.21+
- Linux/macOS (or WSL)
- Git

### Build and run
```bash
git clone https://github.com/JoseMariaMicoli/VaporTrace.git
cd VaporTrace
go mod tidy
go build -o VaporTrace .
./VaporTrace
```

### First commands
```bash
init_db
target https://api.example.com
map
analyze
list-plan
commit
report
```

---

## Operator Workflow

### 1) Initialize and scope
```bash
init_db
target https://api.example.com
```

### 2) Discover attack surface
```bash
map
swagger https://api.example.com/openapi.json
scrape https://api.example.com
```

### 3) Build and refine tactical plan
```bash
analyze
list-plan
edit 1 "ID: 9001"
drop 3
```

### 4) Execute and review
```bash
commit
list-plan
report
```

### 5) Reset mission state (single pass)
```bash
reset_db
```
This now clears database state and volatile runtime state in one run.

---

## Dashboard Tabs and Hotkeys

### Tabs
- **F1 LOGS**: tactical feed and module output
- **F2 MAP**: discovered endpoints and attack surface
- **F3 LOOT**: extracted secrets/artifacts
- **F4 TRAFFIC**: request/response telemetry
- **F5 PLAN**: strategic actions and statuses
- **F6 NEURO**: AI engine output and recommendations
- **F7 REPORT**: report editing/export flow
- **F8 HISTORY**: traffic and operational history

### Core shortcuts
- `F1..F8`: tab navigation
- `Ctrl+I`: interceptor toggle
- `Ctrl+H`: keybind help
- `Esc`: graceful exit

Full reference: [docs/manuals/17_KEYBOARD_SHORTCUTS.md](docs/manuals/17_KEYBOARD_SHORTCUTS.md)

---

## Command Map

### Discovery
- `target <url>`
- `map [url]`
- `swagger <url>`
- `scrape <url>`
- `mine <url> <endpoint>`
- `spider <url> [depth]`
- `fuzz <url> [params|paths]`

### Exploitation
- `bola <url> <id>` or pipeline mode
- `bfla`
- `bopla`
- `ssrf <url> <param> <callback>`
- `exhaust <url> <param>`
- `audit <url>`
- `probe <url> [type]`
- `intruder sniper <url> <param> <wordlist>`
- `race <url> [threads]`

### Planning
- `analyze`
- `list-plan`
- `edit <id> <new_payload>`
- `drop <id>`
- `commit`

### AI and Neuro
- `neuro on|off`
- `neuro-gen <context> <count>`
- `test-neuro`
- `ask <prompt>`

### Intel and Tier-4
- `intel wayback <domain>`
- `intel shodan <host>`
- `chain ...`
- `extract ...`
- `kb ...`

### System
- `init_db`
- `reset_db`
- `loot list|clear`
- `proxy <url|off>`
- `proxies load <file>|reset`
- `report`

Complete command details: [docs/manuals/18_COMMAND_REFERENCE.md](docs/manuals/18_COMMAND_REFERENCE.md)

---

## Plan Lifecycle and Persistence

### How it works now
- `analyze` generates actions and persists them into `tactical_actions`.
- `list-plan` loads from in-memory buffer, and falls back to persisted plan when needed.
- `edit`, `drop`, and `commit` update both runtime and DB-backed plan state.

### Why this matters
- tactical state survives process restart
- F5 reflects durable status transitions
- operational continuity improves across long sessions

---

## Neural Engine (F6)

F6 is the AI decision surface for operators. Current role:
- visualize neuro recommendations and analysis output
- track AI-assisted tactical context during engagements
- monitor recommendation-to-execution flow

Upgrade roadmap (F6 + MCP + persistence):
- [docs/manuals/27_NEURO_F6_MCP_UPGRADE_PLAN.md](docs/manuals/27_NEURO_F6_MCP_UPGRADE_PLAN.md)

---

## Reporting and Compliance

### Reporting outputs
- mission findings in SQLite
- markdown report generation from aggregated findings and telemetry

### Compliance mapping
- OWASP API Top 10 categories
- MITRE ATT&CK references
- NIST CSF control tags
- CVSS scoring support

---

## Project Structure

```text
VaporTrace/
├── cmd/                  # CLI entrypoints
├── pkg/
│   ├── engine/           # command engine + strategic planner
│   ├── logic/            # exploitation/discovery/network modules
│   ├── discovery/        # recon parsers and crawlers
│   ├── ui/               # Hydra TUI and shell interfaces
│   ├── db/               # SQLite schema and persistence
│   └── report/           # report generation
├── docs/
│   ├── manuals/          # operator and admin manuals
│   └── dev-logs/         # implementation logs by sprint
└── assets/               # branding and static assets
```

---

## Documentation

### Start here
- [docs/manuals/INDEX.md](docs/manuals/INDEX.md)

### Recommended reading order
1. [01_INSTALLATION_SETUP.md](docs/manuals/01_INSTALLATION_SETUP.md)
2. [02_FIRST_RUN.md](docs/manuals/02_FIRST_RUN.md)
3. [03_UI_OVERVIEW.md](docs/manuals/03_UI_OVERVIEW.md)
4. [04_STRATEGIC_PLANNING.md](docs/manuals/04_STRATEGIC_PLANNING.md)
5. [05_RECONNAISSANCE.md](docs/manuals/05_RECONNAISSANCE.md)
6. [06_EXPLOITATION.md](docs/manuals/06_EXPLOITATION.md)
7. [07_AI_NEURO_ENGINE.md](docs/manuals/07_AI_NEURO_ENGINE.md)

### Advanced / Tier 4
- [23_INTEL_OSINT.md](docs/manuals/23_INTEL_OSINT.md)
- [24_CHAIN_REACTOR.md](docs/manuals/24_CHAIN_REACTOR.md)
- [25_EXTRACTOR.md](docs/manuals/25_EXTRACTOR.md)
- [26_KNOWLEDGE_BASE.md](docs/manuals/26_KNOWLEDGE_BASE.md)
- [27_NEURO_F6_MCP_UPGRADE_PLAN.md](docs/manuals/27_NEURO_F6_MCP_UPGRADE_PLAN.md)

---

## Roadmap

### Completed
- Core OWASP API exploitation engines
- Hydra TUI operational dashboard
- strategic planner and commit lifecycle
- Tier-4 intel/chain/kb features

### In progress / next
- F6 neuro decision cockpit upgrade
- MCP integration with policy-gated tools
- deeper telemetry explainability and action traceability

---

## Security and Usage Policy
- Use only on systems you own or are explicitly authorized to test.
- Enforce scope via `target` and operational controls.
- Do not run against production without formal authorization and change control.

---

## Contributing
This repository accepts issue reports and improvement proposals through the project workflow.
For architecture-level changes, include:
- affected modules
- threat model implications
- operator UX impact
- testing plan

---

## License
Business Source License 1.1 (BSL). See [LICENSE](LICENSE) for full terms.

