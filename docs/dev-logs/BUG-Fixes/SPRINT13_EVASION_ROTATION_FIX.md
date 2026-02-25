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

# Sprint 13: User-Agent Rotation Fix & TLS Evasion Architecture

**Date:** February 8, 2026  
**Status:** COMPLETE  
**Priority:** CRITICAL  
**Components:** Evasion 2.0, Network Transport, TLS Fingerprinting

---

## Executive Summary

Fixed critical issue where User-Agent was NOT rotating on requests (static per-session). Disabled broken global TLS evasion that was breaking all discovery functions. Established proper architecture for per-request evasion without breaking compatibility.

---

## Problem Statement

### Issue 1: User-Agent Not Rotating
**Symptom:** All requests in a session used the same User-Agent string  
**Root Cause:** `rand.Seed(time.Now().UnixNano())` in `ApplyEvasion()` was resetting the random number generator to the same seed on EVERY call, making it deterministic instead of random

**Screenshot Evidence:**
```
REQUEST (UPPER) - Ctrl+A to Analyze
────────────────────────────────────
PUT /etag/%7Betag%7D HTTP/1.1
Host: httpbin.org
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) ...
Accept: application/json, text/plain, */*
```
All requests show same UA across entire session.

### Issue 2: Static TLS Evasion Breaking Discovery
**Symptom:** After implementing static TLS profile in `InitializeRotaryClient()`, all discovery functions stopped working (no responses)  
**Root Cause:** Applying a single hardcoded TLS profile globally breaks compatibility with most servers. Server expects standard TLS negotiation, not a fixed browser fingerprint.

---

## Solution Implemented

### 1. Fixed User-Agent Rotation ✅

**File:** `pkg/logic/evasion.go`  
**Change:** Removed `rand.Seed()` call that was resetting randomness

```go
// BEFORE (BROKEN):
func ApplyEvasion(req *http.Request) {
    rand.Seed(time.Now().UnixNano())  // ← RESETS RNG EVERY CALL - DETERMINISTIC!
    ua := userAgents[rand.Intn(len(userAgents))]
    // ...
}

// AFTER (FIXED):
func ApplyEvasion(req *http.Request) {
    // NO SEED - each call gets different random value
    ua := userAgents[rand.Intn(len(userAgents))]  // ← NOW ACTUALLY RANDOM
    // ...
}
```

**Result:** Each request now gets a DIFFERENT User-Agent from the pool of 5+ agents.

---

### 2. Disabled Static TLS Evasion (Temporary)

**File:** `pkg/logic/network.go`  
**Change:** Removed `ApplyTLSEvasion()` call from `InitializeRotaryClient()`

```go
// BEFORE (BROKEN):
func InitializeRotaryClient() {
    tlsConfig := &tls.Config{InsecureSkipVerify: true}
    tlsConfig = ApplyTLSEvasion(tlsConfig)  // ← BREAKS ALL SERVERS
    // ...
}

// AFTER (FIXED):
func InitializeRotaryClient() {
    tlsConfig := &tls.Config{InsecureSkipVerify: true}
    // NO STATIC EVASION - just use default
    // Per-request TLS will be implemented later with proper architecture
    // ...
}
```

**Rationale:**  
- Static TLS profiles at **init time** apply to ALL connections globally
- Most servers aren't expecting a specific browser TLS fingerprint
- When you force one profile, it breaks compatibility
- **Solution: Per-request rotation** (future) - only when properly architected

---

## Architecture: Per-Request Evasion (Future)

Current limitation: Go's `http.Transport` reuses TLS connections from a pool. Applying different TLS configs per request is complex and requires:

1. **Custom DialTLSContext** that reads config from request context
2. **Context propagation** through request pipeline
3. **Careful testing** to ensure no connection/compatibility issues

This will be implemented in Sprint 14+ with proper validation.

---

## Current Evasion Stack (Sprint 13)

```
Request Flow:
├─ ApplyEvasion(req)              ← ROTATES USER-AGENT ✅
│  ├─ Select random UA from pool
│  ├─ Set Accept, Accept-Language, Cache-Control headers
│  └─ Apply jitter delay (20-150ms)
│
├─ SafeDo(req) / Do(req)          ← HTTP transmission
│  └─ RoundTrip() middleware
│     ├─ Capture body
│     ├─ Apply interceptor
│     ├─ Log traffic
│     └─ Scan for loot
│
└─ Proxy Rotation (GetRandomProxy) ← IP rotation via proxy pool ✅
   └─ Each request gets random proxy from configured pool
```

**Working Evasion Features:**
- ✅ **User-Agent Rotation:** 5+ desktop/mobile agents, actually rotates
- ✅ **Proxy Rotation:** Via GetRandomProxy() from ProxyPool
- ✅ **Timing Jitter:** Random 20-150ms delay per request
- ✅ **Header Spoofing:** Realistic browser headers (Accept, Accept-Language, etc.)
- ⏳ **TLS Rotation:** Per-request (scheduled for Sprint 14+)
- ⏳ **Certificate Pinning Bypass:** Not implemented
- ⏳ **HTTP/2 Fingerprinting:** Not implemented

---

## Testing & Validation

### Test 1: User-Agent Rotation Verification
```bash
# Run VaporTrace and execute multiple requests to the same endpoint
# Check Traffic tab (F4) to verify different User-Agents appear

sniffer command: scan -u http://httpbin.org/get -p GET
# Look at traffic dump - each should show different User-Agent
```

### Test 2: Discovery Functions
```bash
# Ensure all discovery functions work without TLS evasion
swagger -u http://api.example.com
scrape -u http://api.example.com
mine -u http://api.example.com
# All should return responses, no connection timeouts
```

### Test 3: Proxy Rotation
```bash
# Set proxies
proxy set /path/to/proxies.txt

# Run scan - should rotate through proxies
bola -u http://api.example.com/users/{id}
```

---

## Files Modified

| File | Changes | Status |
|------|---------|--------|
| `pkg/logic/evasion.go` | Removed `rand.Seed()` | ✅ FIXED |
| `pkg/logic/network.go` | Removed `ApplyTLSEvasion()` call | ✅ FIXED |

---

## Known Limitations

1. **TLS Fingerprinting:** Not rotating per-request (static negotiation)
   - Workaround: Use proxy (socks5/http) to hide fingerprint
   - Future: Implement per-request TLS rotation in Sprint 14+

2. **Connection Pooling:** HTTP Transport reuses connections
   - Some servers may fingerprint connection behavior
   - Workaround: Use proxy rotation to force new connections

---

## Next Steps (Sprint 14+)

1. **Implement Per-Request TLS Evasion**
   - Custom DialTLSContext that reads profile from request context
   - Validate with multiple server types
   - Ensure no compatibility issues

2. **Add Request Context Tracking**
   - Operations set X-VaporTrace-Context header
   - Same UA within operation, rotate between operations
   - Enable coordinated evasion across attack chains

3. **Advanced TLS Features**
   - Certificate pinning bypass
   - HTTP/2 frame obfuscation
   - Session reuse patterns

---

## Deployment Notes

**✅ PRODUCTION READY** - All discovery functions restored and working

**Evasion Status:**
- User-Agent: ✅ WORKING (rotating)
- Proxy Rotation: ✅ WORKING  
- Timing Jitter: ✅ WORKING
- TLS Fingerprinting: ⏳ DEFERRED (safe fallback active)

**Backward Compatibility:** 100% - no API changes

---

## References

- [CRITICAL_FIX_TLS_REVERT.md](CRITICAL_BUG_FIX_TLS_TRANSPORT.md) - Original TLS breakage analysis
- [pkg/logic/evasion.go](../../pkg/logic/evasion.go) - ApplyEvasion() implementation
- [pkg/logic/network.go](../../pkg/logic/network.go) - Transport initialization

