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

# Knowledge Base (KB) - Tier 4 Day 3: Institutional Memory

**Status:** ✅ Production Ready | **Sprint:** 20 | **Feature:** Institutional Memory  
**Objective:** Transform VaporTrace into a learning platform by recording successful attack vectors and feeding them back into the Neural Engine for continuous improvement.

---

## Table of Contents

1. [Overview](#overview)
2. [Architectural Design](#architectural-design)
3. [Command Reference](#command-reference)
4. [Vector Types](#vector-types)
5. [Workflows & Examples](#workflows--examples)
6. [Neural Engine Integration](#neural-engine-integration)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)

---

## Overview

### What is the Knowledge Base?

The **Knowledge Base (KB)** is VaporTrace's institutional memory system. It records every successful exploitation technique, creating a persistent library that:

- **Records** successful attack vectors (BOLA, BFLA, SSRF, etc.)
- **Feeds** these patterns back into the Neural Engine for AI mutation
- **Improves** over time as more exploits are documented
- **Scales** reconnaissance across multiple targets using learned patterns

### Key Features

| Feature | Description |
|---------|-------------|
| **Persistent Storage** | KB entries survive VaporTrace restarts |
| **AI Learning** | Neural Engine mutates payloads based on KB patterns |
| **Cross-Domain Reuse** | Exploit pattern from one target applies to another |
| **Export Capability** | Share KB entries across teams (JSON, CSV) |
| **Search & Discovery** | Quickly find relevant vectors by endpoint, type, or method |

### Why It Matters

- **Human-in-the-Loop Learning:** Each successful exploit improves the AI's future mutations
- **Institutional Knowledge:** Team accumulates domain-specific attack patterns
- **Faster Exploitation:** Pre-populated KB accelerates testing on similar targets
- **Compliance Reporting:** Track which attack vectors were used and their success rates

---

## Architectural Design

### Data Model

Each KB entry captures:

```
┌─────────────────────────────────────────┐
│      Knowledge Base Vector Entry       │
├─────────────────────────────────────────┤
│ Vector Type     │ BOLA, BFLA, SSRF ... │
│ Endpoint        │ /api/users/{id}      │
│ Method          │ GET, POST, DELETE    │
│ Payload         │ Exploitation payload │
│ Target Domain   │ api.example.com      │
│ Success Rate    │ 85% (tracked)        │
│ Date Recorded   │ 2026-02-11 12:30:15  │
│ AI Mutations    │ 12 variants learned  │
└─────────────────────────────────────────┘
```

### KB ↔ Neural Engine Pipeline

```
┌──────────────────────┐
│  Execute Attack      │  → User runs: bola https://api.target.com/users/1
│  (e.g., BOLA)        │
└──────────────────┬───┘
                   │
                   ↓
        ┌──────────────────────┐
        │  Record in KB         │  → kb add BOLA /users/{id} GET id=1
        │  (Vector Stored)      │
        └──────────────────┬───┘
                           │
                           ↓
        ┌──────────────────────────────────┐
        │  Neural Engine Ingests Pattern   │  → Learns: BOLA works on ID params
        │  (AI Learning Phase)             │
        └──────────────────┬───────────────┘
                           │
                           ↓
        ┌──────────────────────────────────┐
        │  Future Attack Mutations         │  → neuro-gen BOLA 10
        │  (Enhanced Payloads)             │  → Generates 10 variants
        └──────────────────────────────────┘
```

### Storage Backend

| Component | Implementation | Details |
|-----------|----------------|---------|
| **Storage** | SQLite3 (`kb.db`) | Persistent across sessions |
| **Schema** | Normalized tables | `vectors`, `metadata`, `performance` |
| **Indexing** | Hash + B-tree | Fast lookup by type/endpoint/method |
| **Encryption** | AES-256-GCM | Optional for sensitive payloads |
| **Sync** | Write-through | Immediate disk persistence |

---

## Command Reference

### KB List - Display All Vectors

Display all recorded attack vectors in the Knowledge Base.

```bash
kb list
```

**Output Example:**
```
[magenta]KNOWLEDGE BASE - Stored Attack Vectors:[-]
[cyan]=== Institutional Memory ===[-]
[yellow]Total Entries:[-] 3

[1] BOLA       | /api/users/{id}          | GET    | 85% success | 2026-02-11
[2] SSRF       | /api/webhook/notify      | POST   | 70% success | 2026-02-10
[3] MISCONFIG  | /api/health              | GET    | 100% success| 2026-02-09
```

**Fields:**
- **Type:** Vulnerability classification (BOLA, BFLA, SSRF, etc.)
- **Endpoint:** API path pattern
- **Method:** HTTP verb used
- **Success Rate:** Historical success percentage
- **Date:** When vector was first recorded

---

### KB Add - Record a Successful Vector

Record a new attack vector after successful exploitation.

```bash
kb add <vector_type> <endpoint> <method> <payload>
```

**Parameters:**
- `<vector_type>`: Attack type (required)
- `<endpoint>`: API path (required)
- `<method>`: HTTP method (required)
- `<payload>`: Attack payload or description (required)

**Examples:**

```bash
# Record a successful BOLA attack
kb add BOLA /api/users/{id} GET id=1

# Record SSRF vector
kb add SSRF /api/webhook/notify POST url=169.254.169.254/latest/meta-data

# Record mass assignment vulnerability
kb add BOPLA /api/profile PATCH is_admin=true

# Record injection point
kb add INJECTION /api/search QUERY q='; DROP TABLE users; --
```

**Output:**
```
[green]✓ Knowledge Vector Recorded:[-] BOLA on /api/users/{id} (GET)
[cyan]>>> KB ENTRY:[-] Type=BOLA, Endpoint=/api/users/{id}, Method=GET
[yellow]Payload:[-] id=1
[magenta]>>> KB Entry committed. Feeds back to Neural Engine for future mutations.[-]
```

**Workflow Integration:**
1. Execute attack: `bola https://api.target.com/users/1` ✓ Success
2. Record vector: `kb add BOLA /users/{id} GET id=1`
3. On next target: Neural Engine auto-applies similar patterns

---

### KB Search - Find Relevant Vectors

Search the Knowledge Base by endpoint, vector type, or method.

```bash
kb search <query>
```

**Parameters:**
- `<query>`: Search term (endpoint, type, or method)

**Examples:**

```bash
# Find all BOLA vectors
kb search BOLA

# Find all attacks on /users endpoint
kb search /users

# Find all POST attacks
kb search POST

# Find all injection vectors
kb search INJECTION
```

**Output:**
```
[blue]Searching KB for: /users[-]
[yellow]Found 3 matching vectors:[-]

[1] BOLA       | /users/{id}        | GET    | Success: 85%
[2] BOPLA      | /users/profile     | PATCH  | Success: 60%
[3] INJECTION  | /users/search      | GET    | Success: 90%
```

---

### KB Export - Share Knowledge Base

Export KB entries in standardized format for team sharing or reporting.

```bash
kb export [json|csv]
```

**Parameters:**
- `[json|csv]`: Export format (default: json)

**Examples:**

```bash
# Export as JSON
kb export json

# Export as CSV for spreadsheet
kb export csv

# Default format (JSON)
kb export
```

**Output:**
```
[blue]Exporting Knowledge Base as json...[-]
[green]Export saved to: ./kb_export_20260211_123045.json[-]
```

**JSON Format:**
```json
{
  "kb_version": "1.0",
  "export_date": "2026-02-11T12:30:45Z",
  "vector_count": 3,
  "vectors": [
    {
      "id": 1,
      "type": "BOLA",
      "endpoint": "/api/users/{id}",
      "method": "GET",
      "payload": "id=1",
      "success_rate": 0.85,
      "first_recorded": "2026-02-11",
      "usage_count": 12,
      "ai_mutations": 8
    }
  ]
}
```

---

### KB Clear - Delete All Entries

Permanently delete all Knowledge Base entries (requires confirmation).

```bash
kb clear --confirm
```

**Warning:**
```
[yellow]WARNING:[-] This will permanently delete all Knowledge Base entries.
[red]Confirmation Required. Use: kb clear --confirm[-]
```

**Backup Before Clearing:**
```bash
# Always export before clearing
kb export json
kb clear --confirm
```

---

## Vector Types

### Supported Vulnerability Classes

| Type | Description | HTTP Method | Target Pattern |
|------|-------------|-------------|-----------------|
| **BOLA** | Broken Object Level Auth | GET, DELETE | `/{resource}/{id}` |
| **BFLA** | Broken Function Level Auth | DELETE, PUT, PATCH | Any endpoint |
| **BOPLA** | Broken Object Property Auth | PATCH, PUT | JSON fields |
| **SSRF** | Server-Side Request Forgery | GET, POST | URL parameters |
| **EXHAUST** | Resource Exhaustion | GET | Pagination params |
| **MISCONFIG** | Security Misconfiguration | GET, HEAD | `/health`, `/status` |
| **INJECTION** | Code/SQL Injection | GET, POST | Query/body params |
| **CRYPTO** | Cryptographic Weakness | POST | Auth tokens |
| **XXEA** | XXE / XML Injection | POST | XML body |
| **AUTHZ** | Generic Auth Bypass | Any | API endpoints |

---

## Workflows & Examples

### Workflow 1: Record & Reuse Pattern

**Scenario:** You discover BOLA vulnerability on Target A, want to quickly test Target B.

```bash
# TARGET A: Initial discovery
target https://api.example-a.com
bola /users/999                              # ✓ Success! Found user enumeration

# Record the pattern
kb add BOLA /users/{id} GET id=999

# TARGET B: Reuse knowledge
target https://api.example-b.com
kb search BOLA                               # Show all BOLA patterns

# Neural Engine mutates based on KB
neuro on
neuro-gen BOLA 5                             # Generate 5 BOLA variants for target B
bola /users/999                              # ✓ Success again!
```

**Result:** 
- First target: Manual testing (30 min)
- Second target: Auto-mutation + quick test (5 min)
- **Productivity gain: 6x faster**

---

### Workflow 2: Team Knowledge Sharing

**Scenario:** Red team discovers complex attack chain, wants to share with others.

```bash
# On Attacker A's machine
kb search SSRF                               # Find all SSRF vectors
kb export json > ssrf_vectors.json

# Share file with team via secure channel

# On Attacker B's machine
kb import ssrf_vectors.json                  # Load team's knowledge

# Reuse patterns immediately
kb list                                      # See all imported vectors
target https://new-target.com
kb search SSRF                               # Apply learned patterns
```

---

### Workflow 3: Compliance Reporting

**Scenario:** Generate report showing all attack vectors used during penetration test.

```bash
# During engagement
kb add BOLA /api/users/{id} GET ...
kb add SSRF /webhook/notify POST ...
kb add INJECTION /search QUERY ...

# After engagement
kb export json > engagement_vectors.json

# Import into reporting tool
# Report shows: "3 vectors tested, 2 successful, success rate: 67%"
```

---

## Neural Engine Integration

### How KB Feeds the AI

**Step 1: Vector Ingestion**
```
KB stores: BOLA on /users/{id} with payload "id=1"
```

**Step 2: Pattern Recognition**
```
Neural Engine analyzes:
- Endpoint pattern: /{resource}/{id}
- Method: GET (retrieves by ID)
- Parameter name: "id"
- Value type: Numeric
```

**Step 3: Mutation Generation**
```
Neural Engine creates variants:
1. id=2, id=3, id=999
2. id=%s, id=%x (format string patterns)
3. id=../../../admin
4. id=-1
5. id=0
```

**Step 4: Future Application**
```
On new target:
neuro-gen BOLA 10
→ Generates 10 smart mutations based on KB patterns
→ Tests them automatically
```

### Configuration

Enable Neural Engine for KB:

```bash
# Start Neural Engine
neuro on

# Configure KB learning mode
neuro config kb-learning on

# Set mutation strategy
neuro config mutation-strategy adaptive

# View Neural Engine status
tasks                                        # Shows "Neural Engine: ACTIVE (Hybrid Mode)"
```

---

## Best Practices

### 1. Record Immediately After Success

✅ **DO:** Record vector right after exploitation succeeds
```bash
bola https://api.target.com/users/1         # ✓ Success
kb add BOLA /users/{id} GET id=1            # Record immediately
```

❌ **DON'T:** Forget to record and lose institutional knowledge

### 2. Use Consistent Endpoint Patterns

✅ **DO:** Normalize endpoints with `{id}` placeholders
```bash
kb add BOLA /users/{id} GET ...
kb add BOLA /posts/{id} GET ...
kb add BOLA /comments/{id} GET ...
```

❌ **DON'T:** Record as `/users/123` (not generalized)

### 3. Include Full Payload Context

✅ **DO:** Capture the complete exploitation payload
```bash
kb add SSRF /webhook/callback POST url=http://169.254.169.254/latest/meta-data
```

❌ **DON'T:** Just record "SSRF on webhook" without payload

### 4. Regularly Export & Backup

✅ **DO:** Export KB monthly for backup and team sharing
```bash
kb export json > kb_backup_$(date +%Y%m%d).json
```

❌ **DON'T:** Rely on single KB instance without backup

### 5. Clean Up Duplicate Vectors

✅ **DO:** Merge similar vectors to avoid clutter
```bash
kb search BOLA                               # Review all BOLA entries
# Manually combine similar patterns
```

❌ **DON'T:** Let KB grow with duplicate entries

### 6. Document Attack Context

Use extended KB entries to capture context:

```bash
# Standard entry
kb add BOLA /users/{id} GET id=1

# With context (in payload field)
kb add BOLA /users/{id} GET id=1 [admin_only_from_header=X-Admin]
```

---

## Troubleshooting

### Issue: KB is Empty After Restart

**Cause:** KB entries were not persisted to disk.

**Solution:**
```bash
# Check if KB database exists
ls -la kb.db

# If missing, reinitialize
init_db                                      # Create new KB database

# Re-import from backup if available
kb import kb_backup_20260211.json
```

### Issue: KB Search Returns No Results

**Cause:** Search query doesn't match stored vectors exactly.

**Solution:**
```bash
# List all to see exact names
kb list

# Try partial search
kb search user                                # Instead of kb search /users/{id}
```

### Issue: Neural Engine Not Learning KB Patterns

**Cause:** Neural Engine may be disabled or not configured.

**Solution:**
```bash
# Check status
tasks

# Enable Neural Engine
neuro on

# Configure KB learning
neuro config kb-learning on

# Verify integration
neuro-gen BOLA 5                             # Should now mutate based on KB
```

### Issue: KB Export File Too Large

**Cause:** KB has accumulated too many entries.

**Solution:**
```bash
# Export only specific type
kb search BOLA > /dev/null                   # Count BOLA entries

# Export and filter manually
kb export json | grep BOLA > bola_vectors.json

# Or clear old entries
kb clear --confirm
```

### Issue: Cannot Import KB from Team

**Cause:** File format incompatibility or permission issues.

**Solution:**
```bash
# Verify file format
file kb_vectors.json                         # Should show JSON

# Check permissions
ls -l kb_vectors.json

# Manually add vectors if import fails
cat kb_vectors.json | grep endpoint
# Then: kb add TYPE ENDPOINT METHOD PAYLOAD
```

---

## Advanced Topics

### Custom Vector Attributes

Record additional metadata in payload:

```bash
kb add BOLA /users/{id} GET id=1 [requires_auth=false][response_time=250ms]
```

### Encrypted KB Storage

For sensitive environments:

```bash
# Enable AES-256 encryption
neuro config kb-encryption on

# Export with encryption
kb export json --encrypted
```

### Distributed KB Synchronization

Share KB across team:

```bash
# Export from machine A
kb export json > ./shared_kb.json

# Sync to shared storage
gsutil cp ./shared_kb.json gs://team-vt-kb/

# Import on machine B
kb import gs://team-vt-kb/shared_kb.json
```

---

## Summary

| Task | Command | Purpose |
|------|---------|---------|
| View all vectors | `kb list` | Inventory institutional memory |
| Record new vector | `kb add TYPE ENDPOINT METHOD PAYLOAD` | Learn from exploit |
| Search vectors | `kb search QUERY` | Find relevant patterns |
| Share with team | `kb export json` | Distribute knowledge |
| AI learns patterns | `neuro-gen TYPE 10` | Generate mutations from KB |

**Tier 4 Day 3** transforms VaporTrace into a **learning platform** where:
- Every successful exploit improves AI mutations
- Team knowledge compounds over time
- Red teams share institutional memory
- Future engagements are faster and smarter

---

**Next Steps:**
1. ✅ Add kb command to all tools
2. ✅ Integrate with Neural Engine for AI learning
3. 📋 Expand KB storage to support cloud sync (Sprint 21)
4. 📋 Add KB analytics dashboard (Sprint 22)

