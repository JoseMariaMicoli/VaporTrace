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

# Advanced Discovery Guide - Spider & Fuzz Techniques

**VaporTrace Version:** Post-Sprint-18  
**Last Updated:** February 11, 2026  
**Difficulty:** Intermediate → Advanced  

---

## Table of Contents

1. [Overview](#overview)
2. [Spider Command](#spider-command)
3. [Fuzz Command](#fuzz-command)
4. [Discovery Pipeline](#discovery-pipeline)
5. [Performance Optimization](#performance-optimization)
6. [WAF Evasion](#waf-evasion)
7. [Real-World Workflows](#real-world-workflows)
8. [Troubleshooting](#troubleshooting)

---

## Overview

VaporTrace provides three complementary discovery methods:

| Method | Speed | Coverage | Use Case |
|--------|-------|----------|----------|
| **spider** | Medium | High | Domain crawl, link extraction |
| **swagger** | Fast | Medium | API documentation parsing |
| **scrape** | Medium | High | JavaScript endpoint mining |
| **fuzz-params** | Slow | Medium | Hidden parameter discovery |
| **fuzz-paths** | Medium | Medium | Hidden path/admin panel discovery |
| **mine** | Medium | Medium | Quick parameter mining |

**Best Practice:** Combine all methods for maximum coverage.

---

## Spider Command

### Basic Usage

```bash
# Crawl with default depth (2)
spider https://target.com

# Crawl to depth 5
spider https://target.com 5

# Crawl with global target set
target https://target.com
spider                      # Uses global target, default depth
```

### How Spider Works

1. **URL Extraction**
   - Fetches target URL
   - Extracts all `href="..."` attributes
   - Extracts all `src="..."` attributes
   - Normalizes relative URLs to absolute

2. **Recursive Traversal**
   - Visits each discovered link
   - Repeats extraction for each page
   - Continues until max depth reached

3. **Deduplication**
   - Uses thread-safe map to avoid revisiting URLs
   - Checks domain scope to prevent external crawling

4. **Output**
   - Logs all findings to F1 (Tactical Log)
   - Adds endpoints to F2 (Map)
   - Records in database for reporting

### Example: Crawling httpbin.org

```bash
# Command
spider https://httpbin.org 3

# Output
[green]SPIDER:[-] Initializing crawler on https://httpbin.org (Depth: 3, Scope: httpbin.org)
[green]SPIDER:[-] Crawling: https://httpbin.org/
[green]FOUND:[-] /status [200]
[green]FOUND:[-] /get [200]
[green]FOUND:[-] /post [405]
[green]FOUND:[-] /put [405]
[green]FOUND:[-] /delete [405]
[green]FOUND:[-] /patch [405]
[green]FOUND:[-] /anything [200]
[green]SPIDER:[-] Crawling: https://httpbin.org/anything
[green]FOUND:[-] /anything/thing [200]
✓ SPIDER: Crawl complete. Check F2 Map.

# Time: ~10 seconds
# Endpoints found: 10-20
```

### Spider Configuration

**Control via Stealth Settings:**

```bash
# For fast crawling
stealth off
spider https://target.com

# For WAF-heavy targets
stealth silent
spider https://target.com 5

# For balanced approach
stealth fast
spider https://target.com 3
```

**Depth Guidelines:**

| Depth | Time | Coverage | Recommendation |
|-------|------|----------|-----------------|
| 1 | 5-10s | 20% | Quick assessment |
| 2 | 10-20s | 50% | Default, balanced |
| 3 | 20-40s | 70% | Thorough |
| 4-5 | 40-120s | 85% | Comprehensive |
| 6+ | 2+ min | 90%+ | Only if necessary |

### Advanced Spider Techniques

**Technique 1: Incremental Crawling**
```bash
# First pass: shallow
spider https://target.com 2

# Analyze findings
analyze

# Second pass: deeper
spider https://target.com 4
```

**Technique 2: Post-Analysis Crawling**
```bash
# After analyzing discovered endpoints
analyze

# Get more context by crawling interesting endpoints
spider https://target.com/api 3
```

**Technique 3: Subdomain Crawling**
```bash
# Crawl main domain
spider https://target.com 3

# Crawl API subdomain
spider https://api.target.com 3

# Crawl admin panel
spider https://admin.target.com 2
```

---

## Fuzz Command

### Basic Usage

```bash
# Fuzz paths (find hidden admin panels, APIs, etc)
fuzz https://target.com paths

# Fuzz parameters (find hidden query params)
fuzz https://api.target.com params

# Use global target
target https://target.com
fuzz paths          # Uses global target

# Combine both
fuzz https://target.com paths
fuzz https://target.com params
```

### Mode 1: Path Fuzzing

**What it does:**
- Tries 100 common administrative paths
- Detects any response that is NOT 404
- Returns findings with status codes

**Common Paths Tested:**
```
admin, administrator, login, register, signin
api, v1, v2, api/v1, api/v2, dashboard, console
manage, internal, intranet, private, test
debug, metrics, health, healthz, status, info
swagger, openapi, api-docs, swagger.json
config, setup, install, backup, .env
payment, billing, invoice, shop, store
webhook, callback, notification, secrets, tokens
```

**Example:**
```bash
fuzz https://example.com paths

# Output
[cyan]FUZZ:[-] Path enumeration starting on https://example.com (100 payloads)
[green]FOUND:[-] /admin [301]
[green]FOUND:[-] /api [200]
[green]FOUND:[-] /api/v1 [200]
[green]FOUND:[-] /swagger [200]
[green]FOUND:[-] /health [200]
[green]FOUND:[-] /config [403]
[green]FOUND:[-] /backup.sql [403]
[green]FOUND:[-] /payment [200]
✓ FUZZ: Path enumeration complete.
```

### Mode 2: Parameter Fuzzing

**What it does:**
- Gets baseline response (normal request)
- Tries 100 common parameter names
- Detects anomalies:
  - Status code different from baseline
  - Response size differs by > 100 bytes
- Returns findings with reason

**Common Parameters Tested:**
```
id, user, username, password, admin, debug, test
token, auth, key, secret, api_key, session
filter, limit, offset, page, q, search, query
file, upload, download, path, dir, folder
role, permissions, group, email, grant
action, method, format, version, source, config
env, details, info, csrf, xsrf, payment, amount
```

**Example:**
```bash
fuzz https://api.example.com/users params

# Output
[cyan]FUZZ:[-] Parameter mining on https://api.example.com/users (100 payloads)
Baseline: Status=200, Size=1024 bytes

[green]FOUND:[-] Parameter 'id' (Status Code Anomaly: 404 vs baseline 200)
[green]FOUND:[-] Parameter 'debug' (Size Anomaly: +500 bytes)
[green]FOUND:[-] Parameter 'admin' (Status Code Anomaly: 403 vs baseline 200)
[green]FOUND:[-] Parameter 'token' (Size Anomaly: -200 bytes)
✓ FUZZ: Parameter mining complete.
```

### Fuzz Configuration

**Speed vs Stealth:**

```bash
# Fast fuzzing (aggressive)
stealth off
fuzz https://target.com paths   # 30-60 seconds, high detection risk

# Balanced
stealth fast
fuzz https://target.com paths   # 60-90 seconds, moderate detection

# Stealth-focused
stealth silent
fuzz https://target.com paths   # 2-3 minutes, low detection risk
```

**Concurrent Workers:**

Currently hardcoded to 5 workers per mode. To increase (requires code mod):

```go
// In pkg/discovery/fuzzer.go, change:
for i := 0; i < 5; i++ {  // ← Change 5 to desired concurrency
    wg.Add(1)
    go func() { ... }()
}
```

---

## Discovery Pipeline

### The Complete Workflow

```
STEP 1: Reconnaissance
└─ set target
└─ map https://target.com
   ├─ swagger https://target.com/api/docs (if exists)
   ├─ scrape https://target.com/assets/app.js
   └─ mine https://target.com/api

STEP 2: Advanced Discovery
├─ spider https://target.com 3
└─ fuzz https://target.com paths
└─ fuzz https://target.com params

STEP 3: Analysis
└─ analyze (AI analysis of all endpoints)

STEP 4: Exploitation
└─ commit (Execute 30-50 tactical actions)

STEP 5: Reporting
└─ report (Generate finding summary)
```

### Real Example: Discovering httpbin.org

**Phase 1: API Documentation**
```bash
target https://httpbin.org
swagger https://httpbin.org/docs  # If available
```

**Phase 2: JavaScript Analysis**
```bash
scrape https://httpbin.org/assets/app.js
```

**Phase 3: Domain Crawling**
```bash
spider https://httpbin.org 3
```

**Phase 4: Brute Force**
```bash
fuzz https://httpbin.org paths
fuzz https://httpbin.org params
```

**Phase 5: Mine Specific Endpoint**
```bash
mine https://httpbin.org /get
```

**Phase 6: Analyze All Findings**
```bash
analyze
```

**Result: F2 Map**
```
/status [200]
/get [200]
/post [200]
/put [200]
/delete [200]
/anything [200]
/anything/* [200]
... (50+ more endpoints discovered)
```

---

## Performance Optimization

### Tuning for Speed

```bash
# Kill all evasion
stealth off

# Use default depth
spider https://target.com      # depth 2 by default

# Fuzz in parallel
fuzz https://target.com paths &
fuzz https://target.com params
wait
```

**Expected time:** 5-10 minutes for 500+ endpoints

### Tuning for Stealth

```bash
# Maximum evasion
stealth silent
stealth multiplier 5.0

# Shallow crawl
spider https://target.com 1

# Slower fuzzing
fuzz https://target.com paths
fuzz https://target.com params
```

**Expected time:** 30-60 minutes, minimal WAF triggers

### Memory Optimization

For large targets (1000+ endpoints):

```bash
# Process in stages
map https://target.com              # Stage 1: Initial discovery
sleep 5
analyze                            # Stage 2: Analyze findings
sleep 5
commit -batch 10                   # Stage 3: Execute actions (10 at a time)
```

---

## WAF Evasion

### Detection Indicators

```bash
# Check current WAF status
waf detect

# Sample output
[cyan]WAF DETECTION ENGINE[-]
Rate Limit (429): 5 blocks
WAF Blocks (403): 2 blocks
Redirects (30x): 3 redirects
Server Errors (50x): 0 errors
✓ No active WAF detection patterns observed
```

### Evasion Techniques for Discovery

**Technique 1: Delay-Based**
```bash
stealth silent           # Adds 3.0x delay multiplier
spider https://target.com 2
fuzz https://target.com params
```

**Technique 2: User-Agent Rotation**
```bash
# Built into stealth mode
stealth fast
spider https://target.com
# Automatically rotates User-Agent between requests
```

**Technique 3: Request Spacing**
```bash
# Increase multiplier for extreme WAFs
stealth multiplier 10.0
spider https://target.com 1   # Very slow but stealthy
```

**Technique 4: Path Obfuscation**
```bash
# Enable path obfuscation
stealth toggle obfuscation on
fuzz https://target.com paths
# Adds cache busters: ?v=random, ;foo=bar, etc
```

### When to Stop

Signs you're being blocked:
- Multiple 429 (Rate Limit) responses
- Sudden 403/401 errors on previously accessible endpoints
- Consistent 500 errors
- Redirects to captcha

**Response:**
```bash
# Increase stealth
stealth multiplier 5.0
stealth toggle encoding on

# OR accept detection and switch to exploitation
# (you already have 500+ endpoints, that's enough)
```

---

## Real-World Workflows

### Scenario 1: Quick Assessment (15 minutes)

```bash
# Goal: Quick endpoint discovery for known target

target https://api.example.com

# Fast spider
stealth off
spider https://api.example.com 2   # 10-20 seconds

# Quick fuzz
fuzz https://api.example.com paths # 30-60 seconds

# Analyze
analyze

# Results: 50-100 endpoints, ready for exploitation
```

### Scenario 2: Enterprise Target (2+ hours)

```bash
# Goal: Comprehensive discovery for large target

target https://enterprise.example.com

# Set stealth mode (avoid detection)
stealth fast

# Phase 1: Documentation
swagger https://enterprise.example.com/api/docs

# Phase 2: Code extraction
scrape https://enterprise.example.com/assets/app.js
scrape https://enterprise.example.com/assets/vendor.js

# Phase 3: Crawling (30 minutes)
spider https://enterprise.example.com 4
spider https://api.enterprise.example.com 3
spider https://admin.enterprise.example.com 2

# Phase 4: Fuzzing (45 minutes)
fuzz https://enterprise.example.com paths
fuzz https://enterprise.example.com params
fuzz https://api.enterprise.example.com paths
fuzz https://api.enterprise.example.com params

# Phase 5: Analysis & Exploitation
analyze
# Results: 500-1000+ endpoints
commit
```

### Scenario 3: WAF-Protected Target (with evasion)

```bash
# Goal: Discovery without triggering WAF

target https://protected.example.com

# Maximum stealth
stealth silent
stealth multiplier 5.0

# Shallow crawl (minimal requests)
spider https://protected.example.com 1

# Fuzzing with long delays
fuzz https://protected.example.com paths
# (This takes 10+ minutes but won't trigger WAF)

# Check detection status
waf detect

# If still detected, switch to Interceptor + manual analysis
# (VaporTrace already found 100+ endpoints, that's good enough)
```

### Scenario 4: Incremental Discovery

```bash
# Goal: Continuous discovery as new endpoints emerge

target https://api.example.com
map https://api.example.com      # First pass

# Wait, analyze what you found
analyze
commit

# Come back later with deeper crawl
spider https://api.example.com 4

# Analyze again with new findings
analyze
commit
```

---

## Troubleshooting

### Issue: Spider returns 0 endpoints

**Possible causes:**
1. Target URL is unreachable
2. Target blocks without User-Agent
3. Target returns 404 for everything
4. Depth too shallow

**Solutions:**
```bash
# Check connectivity
target https://example.com
audit                      # Tests basic connectivity

# Check target validity
# Use Burp Suite or curl to verify target is accessible
curl -H "User-Agent: Mozilla/5.0" https://example.com

# Increase depth
spider https://example.com 5

# Disable stealth to remove User-Agent rotation
stealth off
spider https://example.com
```

### Issue: Fuzz-params returns 0 anomalies

**Possible causes:**
1. API doesn't respond to unknown parameters (secure)
2. Baseline response includes error handling
3. Parameter anomaly threshold too strict (> 100 bytes)

**Solutions:**
```bash
# Check baseline manually
curl "https://api.example.com/endpoint?unknownparam=test"

# If baseline is stable, target may be secure (that's good!)

# Try fuzz-paths instead (usually more effective)
fuzz https://api.example.com paths
```

### Issue: "Too many requests" errors

**Cause:** WAF rate limiting detected

**Solutions:**
```bash
# Increase stealth
stealth silent
stealth multiplier 10.0

# Reduce concurrency (requires code change)
# OR wait and retry later

# OR check if you already have enough endpoints
# (500+ endpoints is plenty for analysis)
list-plan  # See if analyze already found 30-50 actions
```

### Issue: Spider takes too long

**Possible causes:**
1. Depth too deep
2. Target has many endpoints
3. Server responding slowly

**Solutions:**
```bash
# Kill current spider (Ctrl+C in terminal)
# Increase stealth multiplier to 0.1x (fast)
stealth multiplier 0.1

# Or reduce depth
spider https://target.com 1
```

### Issue: Discovered endpoints seem wrong

**Example:** Finds `/api/v1/users/../../admin` instead of `/admin`

**Cause:** URL normalization issue in spider

**Solution:**
```bash
# This is harmless - the endpoint still works
# You can 'edit' the action in the buffer to fix the payload
# OR use 'fuzz paths' which explicitly tests for /admin
```

---

## Best Practices

### Discovery Phase

✅ **DO:**
- Run map first (combines swagger + scrape)
- Use spider for domain structure
- Use fuzz for hidden paths/params
- Set target before running commands
- Use `stealth fast` on production targets
- Document which commands found which endpoints

❌ **DON'T:**
- Run fuzz-params on every endpoint (too slow)
- Ignore 403/401 responses (they indicate real endpoints!)
- Use `stealth off` on protected targets
- Set depth > 4 on large sites (diminishing returns)
- Expect 100% coverage (some endpoints hidden intentionally)

### Analysis Phase

✅ **DO:**
- Run 'analyze' after discovery completes
- Review generated actions in F5
- Edit/drop actions that seem wrong
- Check confidence scores
- Use 'commit' to execute top-priority actions

❌ **DON'T:**
- Run 'commit' without reviewing (may be noisy)
- Expect AI to find everything (augment with manual testing)
- Ignore rate limiting (stop if getting many 429s)

---

## Advanced Configuration

### Custom Wordlists (Tier 3 feature)

Currently not supported. When available:

```bash
# Load custom wordlist
fuzz https://target.com paths --wordlist /path/to/wordlist.txt

# Or edit pkg/discovery/wordlists.go and rebuild
# Add to Top100Paths or Top100Params slices
```

### Proxy Integration

```bash
# Route all discovery through Burp Suite
proxy 127.0.0.1:8080
spider https://target.com 3    # Goes through Burp

# Check Burp for all requests
```

### Database Storage

All findings automatically stored:
```bash
# View stored findings
report                         # Generates full report with all discovered endpoints

# Or check database directly (if you know SQL)
# Database: ~/.VaporTrace/findings.db
```

---

## References

- [05_RECONNAISSANCE.md](./05_RECONNAISSANCE.md) - General reconnaissance
- [cmd/spider.go](../cmd/spider.go) - Spider source code
- [cmd/fuzz.go](../cmd/fuzz.go) - Fuzz source code
- [pkg/discovery/spider.go](../pkg/discovery/spider.go) - Spider implementation
- [pkg/discovery/fuzzer.go](../pkg/discovery/fuzzer.go) - Fuzz implementation
- [pkg/discovery/wordlists.go](../pkg/discovery/wordlists.go) - Embedded wordlists

---

## Conclusion

Spider and Fuzz commands are your primary tools for automated discovery. Use them in combination with Swagger parsing and JS scraping to build a comprehensive picture of your target's attack surface. Combined with VaporTrace's AI-powered analysis, this creates an efficient reconnaissance → analysis → exploitation workflow.

**Next Steps:**
- Run `help spider` for quick syntax
- Run `help fuzz` for parameter fuzzing details
- Try `map https://your-target.com` for full pipeline
- Check results in F2 (Map) and F3 (Loot) tabs
