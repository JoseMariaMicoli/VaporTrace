# VaporTrace Technical Architecture & Deep Dive

**Document:** dev-logs/INDEX.md  
**Version:** 3.0+  
**Audience:** Developers, Security Engineers, Contributors  
**Last Updated:** February 8, 2026

---

## Overview

VaporTrace is a Go-based API penetration testing framework with real-time TUI dashboard, AI-driven payload generation, and human-in-the-loop attack orchestration. This document provides comprehensive technical documentation of the architecture, implementation, and data flows.

---

## Quick Navigation

### Release & Status (Latest)
- ⭐ [Sprint-20/README.md](Sprint-20/README.md) - **Tier 4: Intelligence, Chain Reactor, Extractor** (Feb 11, 2026) ✅ NEW
- ⭐ [TIER_3_IMPLEMENTATION_SUMMARY.md](TIER_3_IMPLEMENTATION_SUMMARY.md) - **Tier 3 Complete Overview** (Feb 11, 2026) ✅
- ⭐ [YOUR_ACTION_ITEMS.md](YOUR_ACTION_ITEMS.md) - **Tier 3 Checklist & Deployment Status** ✅
- [Dev-Roadmap.md](Dev-Roadmap.md) - Strategic roadmap and upcoming work
- [00_START_HERE.md](00_START_HERE.md) - Starting point for new developers

### Architecture & Design
- [01_ARCHITECTURE.md](01_ARCHITECTURE.md) - System design, components, modules
- [DIAGRAMS/](DIAGRAMS/) - ASCII flow diagrams and architecture visualizations

### Implementation Details
- [02_IMPLEMENTATION.md](02_IMPLEMENTATION.md) - Code walkthrough, key functions
- [03_MODULES_DETAILED.md](03_MODULES_DETAILED.md) - Each exploitation module explained
- [04_DATA_FLOW.md](04_DATA_FLOW.md) - How data flows through the system
- [05_AI_INTEGRATION.md](05_AI_INTEGRATION.md) - Neural engine implementation

### System Internals
- [06_TUI_RENDERING.md](06_TUI_RENDERING.md) - Dashboard rendering, batch updates
- [07_DATABASE.md](07_DATABASE.md) - SQLite schema, persistence layer
- [08_CHANNEL_SYSTEM.md](08_CHANNEL_SYSTEM.md) - Go channels for telemetry
- [09_PERFORMANCE.md](09_PERFORMANCE.md) - Optimization notes, metrics

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    VAPORTRACE 3.0 ARCHITECTURE                 │
└─────────────────────────────────────────────────────────────────┘

┌──────────────────── USER INTERFACE LAYER ─────────────────────┐
│                                                                │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Tactical TUI Dashboard (tview)             │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │  │
│  │  │ F1: Logs │ │ F2: Map  │ │ F3: Loot │ │F4:Traffic│  │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐              │  │
│  │  │F5: Plan  │ │F6: Neuro │ │F7: Report│              │  │
│  │  └──────────┘ └──────────┘ └──────────┘              │  │
│  │  Left Panel: Pipeline Status, Target, Auth, Proxies   │  │
│  │  Bottom: Command Input, Status Bar                    │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                │
└────────────────────────────────────────────────────────────────┘
                              ▼
┌──────────────────── COMMAND ENGINE LAYER ──────────────────────┐
│                                                                │
│  Command Interpreter (ExecuteCommand in core.go)              │
│  ├─ Strategic Planning: analyze, list-plan, edit, drop        │
│  ├─ Reconnaissance: target, map, swagger, scrape, mine        │
│  ├─ Exploitation: bola, bfla, bopla, ssrf, exhaust, ...       │
│  ├─ AI/Neuro: neuro, neuro-gen, ask, test-neuro              │
│  └─ Infrastructure: proxy, init_db, report, sessions          │
│                                                                │
└────────────────────────────────────────────────────────────────┘
                              ▼
┌──────────────────── EXECUTION ENGINE LAYER ────────────────────┐
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │          Core Engine (pkg/engine/core.go)               │ │
│  │                                                          │ │
│  │  ┌─────────────────────────────────────────────────────┐│ │
│  │  │         Strategic Aggregator (HITL Planner)        ││ │
│  │  │  • Analyze: Extract attack vectors from discovery   ││ │
│  │  │  • Plan: Generate tactical actions with confidence  ││ │
│  │  │  • Edit: Allow operator override                    ││ │
│  │  │  • Commit: Execute all actions in parallel          ││ │
│  │  └─────────────────────────────────────────────────────┘│ │
│  │                                                          │ │
│  │  ┌─────────────────────────────────────────────────────┐│ │
│  │  │      Module Orchestrator (pkg/logic/*.go)          ││ │
│  │  │  Routes commands to appropriate exploitation module ││ │
│  │  └─────────────────────────────────────────────────────┘│ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
└────────────────────────────────────────────────────────────────┘
                              ▼
┌──────────────────── EXPLOITATION MODULES ──────────────────────┐
│                                                                │
│  Discovery        Exploitation      Evasion                  │
│  ├─ Scraper      ├─ BOLA          ├─ Ghost Weaver          │
│  ├─ Miner        ├─ BFLA          ├─ Token Forge           │
│  ├─ Spider       ├─ BOPLA         └─ Data Masking          │
│  ├─ Swagger      ├─ SSRF                                    │
│  └─ Sessions     ├─ Exhaust       AI Augmentation           │
│                  ├─ Audit         ├─ Neural Engine          │
│                  ├─ Probe         └─ Payload Mutation       │
│                  └─ Flow                                     │
│                                                                │
└────────────────────────────────────────────────────────────────┘
                              ▼
┌──────────────────── NETWORK LAYER ─────────────────────────────┐
│                                                                │
│  HTTP Client (pkg/logic/network.go)                            │
│  ├─ Request interception & modification                       │
│  ├─ Request body capture (RequestBodyBytes)                   │
│  ├─ Response parsing                                          │
│  ├─ Proxy routing                                             │
│  └─ SSL/TLS handling                                          │
│                                                                │
│  Interceptor Modal (pkg/ui/interceptor.go)                    │
│  ├─ Real-time request viewing                                 │
│  ├─ Payload editing                                           │
│  ├─ Forward/Drop decisions                                    │
│  └─ Loot capture from intercepts                              │
│                                                                │
└────────────────────────────────────────────────────────────────┘
                              ▼
┌──────────────────── AI ENGINE LAYER ───────────────────────────┐
│                                                                │
│  Neural Engine (pkg/engine/neuro_engine.go)                   │
│  ├─ Provider abstraction (Groq, OpenAI, Claude, etc)          │
│  ├─ Prompt templates for:                                     │
│  │  ├─ Payload generation                                     │
│  │  ├─ Mutation strategies                                    │
│  │  └─ Evasion techniques                                     │
│  ├─ Rate limiting & caching                                   │
│  └─ Latency measurement                                       │
│                                                                │
└────────────────────────────────────────────────────────────────┘
                              ▼
┌──────────────────── PERSISTENCE LAYER ────────────────────────┐
│                                                                │
│  Database Manager (pkg/db/manager.go)                         │
│  ├─ Endpoints table (discovered URLs)                         │
│  ├─ Findings table (vulnerability results)                    │
│  ├─ Loot table (captured secrets)                             │
│  ├─ Findings table (audit results)                            │
│  └─ SQLite3 backend                                           │
│                                                                │
│  File System                                                   │
│  ├─ ~/.vaportrace/config.yaml                                 │
│  ├─ ~/.vaportrace/vaportrace.db                               │
│  └─ ./reports/ (generated Markdown/PDF reports)               │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

---

## Data Flow: A Complete Attack Scenario

```
User Input: > target https://api.example.com
    │
    ▼
Command: TARGET
    │
    ▼
Set Global Target URL
Update Pipeline Panel (F5 left side)
Log to Tactical Feed (F1)
    │
    ├─────────────────────────────────────────┐
    │                                         │
    ▼                                         ▼
User: > map                      [Meanwhile: Batch Ticker (200ms)]
    │                            Collects telemetry from all channels
    ▼                            │
Command: MAP                     ▼
    │                            Flushes to UI in single render
    ▼                            (Prevents cascading collapse)
Module: Discovery
├─ Spider: Crawl endpoints
├─ Swagger: Parse OpenAPI specs
└─ Scrape: Extract JS routes
    │
    ▼
Send to MapDataChan
    │
    ▼
UI: F2 MAP tab updates with discovered endpoints
Log: F1 LOGS shows progress
    │
    │
    ├─────────────────────────────────────────┐
    │                                         │
    ▼                                         ▼
User: > analyze                  DataSilo aggregates F1-F4 data
    │                            ├─ Endpoints (discovered)
    ▼                            ├─ Loot (captured secrets)
Strategic Planner                ├─ Traffic patterns
├─ Analyze each endpoint        └─ Context (auth, targets)
├─ Generate attack scenarios        │
├─ Assign confidence scores         ▼
└─ Create tactical actions      AI (if enabled)
    │                            └─ Analyze patterns
    ▼                               Generate insights
ActionBuffer populated
Log to F5 PLAN tab
    │
    ├─────────────────────────────────────────┐
    │                                         │
    ▼                                         ▼
User: > list-plan               User: > analyze
    │                            [Already shows in F5]
    ▼
Display all pending actions
(Show ID, Type, Target, Payload, Confidence, Status)
    │
    │
    ├─────────────────────────────────────────┐
    │                                         │
    ▼                                         ▼
User: > edit 1 /api/users/999   User: > commit
    │                            │
    ▼                            ▼
Update Action #1 payload        Execute all PENDING actions
    │                            │
    │                            ▼
    │                        For each action:
    │                        ├─ Execute exploitation module
    │                        ├─ Collect results
    │                        ├─ Update ActionBuffer status
    │                        ├─ Capture findings
    │                        ├─ Extract loot (if any)
    │                        └─ Log execution
    │                            │
    │                            ▼
    │                        Send findings to Database
    │                            │
    │                            ▼
    │                        UI Updates:
    │                        ├─ F5: Action status → SUCCESS/FAILED
    │                        ├─ F3: Loot table populated
    │                        ├─ F1: Execution logs
    │                        └─ DB: Persistent storage
    │                            │
    ▼                            ▼
User: > report              Action execution complete
    │                       (All parallel, real-time feedback)
    ▼
Generate findings report
├─ Read from database
├─ Format as Markdown
├─ Include screenshots (optional)
└─ Save to ./reports/
    │
    ▼
F7 REPORT tab displays
    │
    ▼
Export complete!
(User can view, edit, and save report)
```

---

## Key Architectural Decisions

### 1. **TUI-First Design**
- All output goes to TUI dashboard (tabbed interface)
- Real-time feedback for long-running operations
- Batch rendering (200ms ticker) prevents cascading collapse

### 2. **Human-in-the-Loop (HITL) Strategy**
- Operator approves/modifies actions before execution
- Tactical planner generates confidence scores
- Operators can edit payloads or drop actions
- Reduces false positives and improves precision

### 3. **Channel-Based Telemetry**
- Go channels for async event streaming
- Multiple channels for different data types:
  - `UI_Log_Chan` - Tactical feed messages
  - `MapDataChan` - Endpoint discoveries
  - `LootDataChan` - Secret captures
  - `TrafficChan` - HTTP traffic
- LogBuffer system batches updates

### 4. **Modular Exploitation**
- Each attack vector is a separate module
- Easy to add new exploit types
- Independent parameter handling
- Pluggable into orchestrator

### 5. **AI Augmentation (Optional)**
- Neural engine as optional layer
- Provider-agnostic (supports Groq, OpenAI, etc)
- Payload mutation and generation
- Integrated with all exploitation modules

---

## Component Interaction Diagram

```
┌────────────────────────────────────────────────────────────┐
│                  USER INTERFACE LAYER                      │
│  Dashboard (tview) with 7 tabs + left panel + cmd input    │
└──────────────────┬─────────────────────────────────────────┘
                   │ ShowHelpModal()
                   │ UpdateTabs()
                   │ SetInputCapture()
                   ▼
┌────────────────────────────────────────────────────────────┐
│                  COMMAND ENGINE LAYER                      │
│  ExecuteCommand() - Switch statement routing               │
│  40+ command cases (analyze, commit, ssrf, etc)            │
└──────────────────┬─────────────────────────────────────────┘
                   │
        ┌──────────┼──────────┬────────────┬─────────────┐
        │          │          │            │             │
        ▼          ▼          ▼            ▼             ▼
     PLAN      DISCOVERY   EXPLOIT     NEURO        INFRA
     ├─List    ├─Map       ├─BOLA      ├─Enable    ├─Proxy
     ├─Edit    ├─Swagger   ├─BFLA      ├─Ask       ├─Report
     ├─Drop    ├─Scrape    ├─BOPLA     ├─Neuro-gen ├─DB
     └─Commit  ├─Mine      ├─SSRF      └─Test      └─Sessions
              └─Sessions  └─Exhaust
                │          │
                └──────────┼─────────────────────────┐
                           ▼                         ▼
        ┌─────────────────────────────┐  ┌──────────────────┐
        │  Exploitation Modules        │  │ Neural Engine    │
        │  (pkg/logic/*.go)            │  │ (neuro_engine.go)│
        │  - Network I/O               │  │ - LLM queries    │
        │  - Pattern matching          │  │ - Mutation logic │
        │  - Payload generation        │  │ - Caching        │
        │  - Response parsing          │  └──────────────────┘
        └────────┬────────────────────┘
                 │ HTTP calls, Response capture
                 ▼
        ┌─────────────────────────────┐
        │  HTTP Network Layer          │
        │  (network.go)                │
        │  - Interceptor MITM          │
        │  - Request/response capture  │
        │  - Proxy routing             │
        │  - SSL/TLS handling          │
        └────────┬────────────────────┘
                 │ HTTP traffic
                 ▼
        [TARGET API SERVER]
                 │ Response
                 ▼
        ┌─────────────────────────────┐
        │  Telemetry Router            │
        │  (dashboard.go, channels)    │
        │  - LogBuffer batching        │
        │  - Channel listeners         │
        │  - UI update coordination    │
        └────────┬────────────────────┘
                 │
      ┌──────────┼──────────┬────────────┐
      ▼          ▼          ▼            ▼
    F1 LOGS   F3 LOOT   F5 PLAN    DATABASE
    Tactical  Captured  Actions    Persistent
    Feed      Secrets   Status     Storage
```

---

## Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **UI** | tview, tcell | Terminal UI rendering |
| **Language** | Go 1.25+ | Core implementation |
| **CLI** | spf13/cobra | Command parsing |
| **Database** | SQLite3 | Persistent storage |
| **HTTP** | net/http | Network requests |
| **AI** | Provider API (Groq/OpenAI) | Payload generation |
| **Config** | YAML | Settings management |
| **Reporting** | Markdown | Findings export |

---

## File Structure

```
VaporTrace/
├── main.go                          # Entry point
├── go.mod, go.sum                   # Dependencies
│
├── pkg/
│   ├── ui/                          # User interface
│   │   ├── dashboard.go             # Main TUI (756 lines)
│   │   ├── interceptor.go           # MITM modal
│   │   ├── help.go                  # Keybindings modal
│   │   ├── shell.go                 # Shell-like helpers
│   │   ├── report_tab.go            # Report editor
│   │   └── report_manager.go        # Report generation
│   │
│   ├── engine/                      # Execution & AI
│   │   ├── core.go                  # Command interpreter (1483 lines)
│   │   ├── neuro_engine.go          # AI/ML integration
│   │   └── heuristics.go            # Pattern analysis
│   │
│   ├── logic/                       # Exploitation modules
│   │   ├── network.go               # HTTP + Interceptor
│   │   ├── discovery.go             # Endpoint discovery
│   │   ├── scraper.go               # JS parsing
│   │   ├── miner.go                 # Parameter fuzzing
│   │   ├── bola.go                  # BOLA exploitation
│   │   ├── bfla.go                  # BFLA exploitation
│   │   ├── bopla.go                 # BOPLA exploitation
│   │   ├── ssrf.go                  # SSRF exploitation
│   │   ├── exhaust.go               # Resource DoS
│   │   ├── audit.go                 # Config audit
│   │   ├── probe.go                 # Webhook injection
│   │   ├── flow.go                  # Attack chains
│   │   ├── loot.go                  # Secret capture
│   │   ├── evasion.go               # General evasion
│   │   ├── ghost_weaver.go          # Token forgery
│   │   ├── pipeline.go              # Action orchestration
│   │   ├── auth_store.go            # Auth management
│   │   └── *.go                     # (20+ modules total)
│   │
│   ├── db/                          # Persistence
│   │   └── manager.go               # SQLite manager
│   │
│   ├── utils/                       # Helpers
│   │   ├── logger.go                # Logging system
│   │   ├── channels.go              # Channel definitions
│   │   ├── http.go                  # HTTP helpers
│   │   └── strings.go               # String utilities
│   │
│   ├── discovery/                   # Reconnaissance
│   │   ├── discovery.go             # URL discovery
│   │   ├── scraper.go               # JS endpoint extraction
│   │   ├── miner.go                 # Param fuzzing
│   │   └── swagger.go               # OpenAPI parsing
│   │
│   ├── enrichment/                  # Data analysis
│   │   └── tagging.go               # Response tagging
│   │
│   ├── ai/                          # AI providers
│   │   ├── client.go                # LLM client
│   │   └── prompts.go               # Prompt templates
│   │
│   └── report/                      # Report generation
│       └── generator.go             # Markdown export
│
├── cmd/                             # CLI subcommands
│   ├── root.go                      # Root command
│   ├── map.go                       # map command
│   ├── scan.go                      # scan command
│   └── *.go                         # Other CLI commands
│
├── docs/                            # Documentation
│   ├── manuals/
│   │   ├── INDEX.md
│   │   ├── 01_INSTALLATION_SETUP.md
│   │   ├── 02_FIRST_RUN.md
│   │   ├── 03_UI_OVERVIEW.md
│   │   └── ... (20 total)
│   │
│   └── dev-logs/
│       ├── INDEX.md (you are here)
│       ├── 01_ARCHITECTURE.md
│       ├── 02_IMPLEMENTATION.md
│       └── DIAGRAMS/
│
└── README.md
```

---

## Next Steps

1. **System Design:** → [01_ARCHITECTURE.md](01_ARCHITECTURE.md)
2. **Code Walkthrough:** → [02_IMPLEMENTATION.md](02_IMPLEMENTATION.md)
3. **Module Details:** → [03_MODULES_DETAILED.md](03_MODULES_DETAILED.md)
4. **Data Flows:** → [04_DATA_FLOW.md](04_DATA_FLOW.md)
5. **AI Integration:** → [05_AI_INTEGRATION.md](05_AI_INTEGRATION.md)

---

**For User Documentation:** See [../manuals/INDEX.md](../manuals/INDEX.md)

