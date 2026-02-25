![VaporTrace Logo](../../../assets/images/VaporTrace_Logo.png)

/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

# Tier 4 Day 1: Strategic Intelligence (OSINT)

**Status:** ✅ COMPLETE
**Date:** February 11, 2026
**Focus:** Passive Reconnaissance & Infrastructure Analysis

---

## 🧠 The "Platform" Shift

Tier 4 moves VaporTrace from "scanning what we see" to "finding what was hidden". By integrating external intelligence sources, we can identify targets that aren't linked in the current application HTML but technically still exist (and are often vulnerable).

---

## 📦 What Was Built

### 1. Wayback Machine Integration (`pkg/intel/wayback.go`)
- **Capability:** Queries the Internet Archive's CDX API
- **Filtering:** Automatically removes static assets (images, CSS) to focus on API endpoints
- **Integration:** Results are fed directly into the **F2 Map** and **Database**
- **Value:** Finds "Ghost Endpoints" (e.g., old `/v1/` APIs) that developers forgot to disable

### 2. Shodan Integration (`pkg/intel/shodan.go`)
- **Capability:** Resolves domain to IP and queries Shodan.io for open ports
- **Output:** Lists ports, services, and banners directly in VaporTrace logs
- **Value:** Identifies potential origin servers or exposed administrative panels on non-standard ports

### 3. Intel Command Integration (`pkg/engine/core.go`)
- **New command structure:**
  - `intel wayback <domain>` - Fetch historical endpoints
  - `intel shodan <ip/domain>` - Query Shodan for ports/services
  - `intel config shodan <key>` - Configure Shodan API key

---

## 🧪 How to Test It

### 1. **Rebuild & Run**
```bash
go build
./VaporTrace
```

### 2. **Ghost Hunting (No API Key needed)**
```bash
intel wayback tesla.com
```
Watch F1 Logs and F2 Map populate with thousands of historical endpoints.

### 3. **Infrastructure Scan (API Key needed)**
```bash
intel config shodan YOUR_API_KEY
intel shodan 8.8.8.8
```

---

## 🚀 Impact on Workflow

1. **Run intel wayback** to discover historical endpoints
2. **Wait for F2 Map** to populate with "ghost endpoints"
3. **Run analyze** (Tier 1) to generate attack vectors
4. **Run commit** (Tier 3) to fuzz these forgotten endpoints with **Intruder/Race**

### Example Attack Chain
```
intel wayback api.example.com
  ↓
F2 Map discovers: /api/v1/users, /api/v1/admin, /api/v2/profile
  ↓
analyze
  ↓
commit
  ↓
Intruder/Race engines automatically test these endpoints
```

---

## 🔧 Implementation Details

### Import Added
```go
"github.com/JoseMariaMicoli/VaporTrace/pkg/intel"
```

### Files Updated
- `pkg/engine/core.go` - Added intel command + help + usage
- `GetAvailableCommands()` - Includes "intel"
- `GetCommandSyntax()` - Maps syntax
- `printHelp()` - Added help case for intel
- `printUsage()` - Added to discovery section

### Architecture
```
External Sources
├─ Wayback Machine CDX API
└─ Shodan.io API

        ↓
    pkg/intel/
    ├─ config.go (API key management)
    ├─ wayback.go (CDX client + filtering)
    └─ shodan.go (Host API client)

        ↓
pkg/engine/core.go (intel command)

        ↓
GlobalDiscovery (F2 Map) ← Results populated
        ↓
Future: Tier 3 Intruder/Race testing
```

---

## ✅ Tier 4 Day 1 Checklist

- [x] **Package Structure:** Created `pkg/intel/` foundation
- [x] **Wayback Module:** Implemented CDX API client + filtering logic
- [x] **Shodan Module:** Implemented Host API client
- [x] **CLI Integration:** Wired intel command into `core.go`
- [x] **Help System:** Added intel to help, usage, and autocomplete
- [x] **Build Verification:** Code compiles without errors

---

## 📋 Next Steps (Tier 4 - Day 2 & 3)

### Day 2: The Chain Reactor
- Build stateful multi-request engine
- Variable extraction (Regex/JSONPath)
- Automated login → token → attack chains

### Day 3: Knowledge Base
- Store successful attack patterns
- Feed proven payloads into AI
- Implement "Institutional Memory"

---

## 🎯 Tier 4 Strategic Value

| Feature | Value | Use Case |
|---------|-------|----------|
| **Wayback Machine** | Find forgotten endpoints | Penetration testing old APIs |
| **Shodan** | Identify exposed infrastructure | Bypass WAF via origin IP |
| **Chaining** | Multi-step exploitation | Automate login → attack sequences |
| **Knowledge Base** | Learn from successful attacks | Build custom AI training data |

---

**Status:** ✅ Day 1 Complete - Ready for Day 2 (Chain Reactor)
