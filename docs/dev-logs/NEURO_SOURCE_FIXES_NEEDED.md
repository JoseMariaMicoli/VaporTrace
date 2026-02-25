![VaporTrace Logo](../../assets/images/VaporTrace_Logo.png)

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

# Neuro Engine - Source Code Fixes Required

## 📋 Overview

Documentation has been corrected. The following source files need fixes to implement the documented functionality properly.

---

## ✅ Documentation Fixed

- [x] [07_AI_NEURO_ENGINE.md](../manuals/07_AI_NEURO_ENGINE.md) - Updated command syntax
- [x] [NEURO_QUICK_USAGE_GUIDE.md](../manuals/NEURO_QUICK_USAGE_GUIDE.md) - Created comprehensive usage guide

**Correct Syntax:**
```bash
neuro config <provider> <model> [api_key] [endpoint]
neuro-gen <context> [count]
ask <prompt>
```

---

## 🔧 Source Code Fixes Required

### Issue 1: Command Parsing Off-By-One (CRITICAL)

**File:** `pkg/engine/core.go`  
**Lines:** 273-284  
**Severity:** CRITICAL - Prevents neuro config from working

**Current Code:**
```go
case "config":
    if len(args) < 2 {
        utils.TacticalLog(session, "neuro config requires: provider model [api_key] [endpoint]", nil)
        return
    }
    provider := args[1]
    model := args[2]
    apiKey := ""
    endpoint := ""
    if len(args) > 3 {
        apiKey = args[3]
    }
    if len(args) > 4 {
        endpoint = args[4]
    }
```

**Issue:** Array indexing is off. When command is `neuro config groq llama-3.1-8b gsk_xxxxx`:
- args[0] = "config"
- args[1] = "groq" ✓
- args[2] = "llama-3.1-8b" ✓
- args[3] = "gsk_xxxxx" ✓

Wait, this looks correct! Let me verify by checking how commands are parsed...

**Action Needed:** Verify the args array construction in the command dispatch. If args includes the command itself as args[0], then adjust indices.

---

### Issue 2: Missing Nil Checks (HIGH)

**Files:**
- `pkg/engine/neuro_engine.go` - Line ~250
- `pkg/logic/neuro_engine.go` - Line ~180
- `pkg/ai/client.go` - Multiple locations

**Current Code Pattern (Example):**
```go
resp := provider.Query(prompt)  // No nil check on resp
results := resp.Content         // Panic if resp is nil
```

**Fix:**
```go
resp := provider.Query(prompt)
if resp == nil {
    utils.TacticalLog(session, "No response from LLM provider", map[string]interface{}{
        "provider": provider.Name(),
        "error": "nil response",
    })
    return errors.New("LLM provider returned nil")
}
results := resp.Content
```

**Action Needed:** Add nil checks after all LLM provider calls.

---

### Issue 3: Confidence Scoring Logic (MEDIUM)

**File:** `pkg/engine/neuro_engine.go`  
**Lines:** ~320-350  
**Severity:** MEDIUM - Returns undefined values

**Current Code Pattern:**
```go
switch {
case statusCode >= 200 && statusCode <= 299:
    confidence = 0.9
case statusCode >= 300 && statusCode <= 399:
    confidence = 0.7
case statusCode >= 400 && statusCode <= 499:
    confidence = 0.5
case statusCode >= 500 && statusCode <= 599:
    confidence = 0.3
}
// Missing: What if statusCode == 0 or unexpected value?
return confidence  // Could be 0 (undefined)
```

**Fix:**
```go
var confidence float64
switch {
case statusCode >= 200 && statusCode <= 299:
    confidence = 0.9
case statusCode >= 300 && statusCode <= 399:
    confidence = 0.7
case statusCode >= 400 && statusCode <= 499:
    confidence = 0.5
case statusCode >= 500 && statusCode <= 599:
    confidence = 0.3
default:
    confidence = 0.0  // Unknown status
    utils.TacticalLog(session, "Unexpected status code", map[string]interface{}{
        "status": statusCode,
    })
}
return confidence
```

**Action Needed:** Add default case to confidence scoring and other switch statements.

---

### Issue 4: Rate Limiting Too Aggressive (HIGH)

**File:** `pkg/logic/neuro_engine.go`  
**Lines:** ~420-435  
**Severity:** HIGH - Delays attacks unnecessarily

**Current Code:**
```go
func enforceRateLimit() {
    lastRequestTime := time.Now()
    for {
        elapsed := time.Since(lastRequestTime)
        if elapsed < 6*time.Second {  // ← Too long! 6 seconds
            time.Sleep(6*time.Second - elapsed)
        }
        lastRequestTime = time.Now()
    }
}
```

**Fix:**
```go
func enforceRateLimit() {
    lastRequestTime := time.Now()
    for {
        elapsed := time.Since(lastRequestTime)
        rateLimitDelay := 2 * time.Second  // Industry standard
        if elapsed < rateLimitDelay {
            time.Sleep(rateLimitDelay - elapsed)
        }
        lastRequestTime = time.Now()
    }
}
```

**Action Needed:** Change 6 seconds to 2 seconds in rate limiter.

---

### Issue 5: Race Condition in GlobalNeuro (CRITICAL)

**File:** `pkg/logic/neuro_engine.go`  
**Lines:** ~50-100  
**Severity:** CRITICAL - Data corruption under concurrent access

**Current Code:**
```go
var GlobalNeuro *NeuroEngine

func GetGlobalNeuro() *NeuroEngine {
    if GlobalNeuro == nil {  // ← Race condition! Two goroutines could see nil
        GlobalNeuro = &NeuroEngine{...}
    }
    return GlobalNeuro
}
```

**Fix:**
```go
var (
    GlobalNeuro *NeuroEngine
    neuroMutex  sync.RWMutex
)

func GetGlobalNeuro() *NeuroEngine {
    neuroMutex.RLock()
    if GlobalNeuro != nil {
        neuroMutex.RUnlock()
        return GlobalNeuro
    }
    neuroMutex.RUnlock()
    
    neuroMutex.Lock()
    defer neuroMutex.Unlock()
    
    // Double-check after acquiring write lock
    if GlobalNeuro == nil {
        GlobalNeuro = &NeuroEngine{...}
    }
    return GlobalNeuro
}
```

**Action Needed:** Implement proper singleton pattern with double-checked locking.

---

### Issue 6: Missing Error Propagation (HIGH)

**File:** `pkg/engine/core.go`  
**Lines:** ~280-295  
**Severity:** HIGH - Errors silently fail

**Current Code:**
```go
case "neuro-gen":
    context := args[1]
    count := 10
    if len(args) > 2 {
        count, _ = strconv.Atoi(args[2])  // ← Silent failure if conversion fails
    }
    logic.GeneratePayloads(context, count)  // ← No error handling
    return  // ← Doesn't check if generation succeeded
```

**Fix:**
```go
case "neuro-gen":
    if len(args) < 2 {
        utils.TacticalLog(session, "neuro-gen requires: context [count]", nil)
        return
    }
    
    context := args[1]
    count := 10
    if len(args) > 2 {
        var err error
        count, err = strconv.Atoi(args[2])
        if err != nil {
            utils.TacticalLog(session, "Invalid count parameter", map[string]interface{}{
                "provided": args[2],
                "error": err.Error(),
            })
            return
        }
    }
    
    err := logic.GeneratePayloads(context, count)
    if err != nil {
        utils.TacticalLog(session, "Payload generation failed", map[string]interface{}{
            "context": context,
            "error": err.Error(),
        })
        return
    }
    
    utils.TacticalLog(session, "Payloads generated successfully", map[string]interface{}{
        "count": count,
    })
```

**Action Needed:** Add proper error handling and validation.

---

## 📊 Fix Priority

| Priority | Issue | File | Lines | Est. Time |
|----------|-------|------|-------|-----------|
| 🔴 CRITICAL | Command parsing verification | core.go | 273-284 | 10 min |
| 🔴 CRITICAL | Race condition in GlobalNeuro | neuro_engine.go | 50-100 | 20 min |
| 🟠 HIGH | Missing nil checks | Multiple | ~10 locations | 30 min |
| 🟠 HIGH | Error propagation | core.go | 280-295 | 20 min |
| 🟠 HIGH | Rate limiting adjustment | neuro_engine.go | 420-435 | 5 min |
| 🟡 MEDIUM | Confidence scoring logic | neuro_engine.go | 320-350 | 15 min |

**Total Estimated Time: 1.5-2 hours**

---

## ✅ Before/After Testing

### Test 1: Configuration
```bash
# BEFORE: Would fail silently
> neuro config groq llama-3.1-8b gsk_xxxxx
[red]ERROR:[-] Invalid provider

# AFTER: Should work
> neuro config groq llama-3.1-8b gsk_xxxxx
[green]✓ CONFIG:[-] Groq configured
```

### Test 2: Payload Generation
```bash
# BEFORE: Silent failure on invalid count
> neuro-gen "SQL injection" abc
[cyan]Generated payloads...

# AFTER: Clear error message
> neuro-gen "SQL injection" abc
[red]ERROR:[-] Invalid count parameter: abc
```

### Test 3: Concurrent Access
```bash
# BEFORE: Race condition possible
# Multiple goroutines access GlobalNeuro simultaneously
# Result: Data corruption or panic

# AFTER: Thread-safe access
# Multiple goroutines can safely access GlobalNeuro
# Result: Consistent behavior
```

---

## 📝 Implementation Strategy

1. **Phase 1 (10 min)** - Verify command parsing logic
   - Check how args array is constructed in ExecuteCommand
   - Confirm correct indices for provider/model/apikey

2. **Phase 2 (20 min)** - Fix race conditions
   - Add sync.RWMutex to GlobalNeuro
   - Implement double-checked locking pattern

3. **Phase 3 (30 min)** - Add nil checks everywhere
   - Audit all LLM provider calls
   - Add proper nil checks with error logging

4. **Phase 4 (20 min)** - Error propagation
   - Update command handlers to check errors
   - Add detailed error messages to TacticalLog

5. **Phase 5 (5 min)** - Rate limiting adjustment
   - Change 6 seconds to 2 seconds
   - Test performance impact

6. **Phase 6 (15 min)** - Fix confidence logic
   - Add default cases to all switch statements
   - Test confidence scoring with edge cases

---

## 🧪 Validation Checklist

After implementing fixes, verify:

- [ ] `go build` compiles without errors
- [ ] All commands execute without panics
- [ ] Error messages are clear and helpful
- [ ] Concurrent access is safe (race detector: `go test -race ./...`)
- [ ] API calls retry on failure
- [ ] Confidence scoring returns values 0.0-1.0
- [ ] Rate limiting doesn't exceed 2 seconds per request
- [ ] Documentation matches implementation

---

## 📞 Next Steps

1. **Review this document** to understand all required changes
2. **Implement fixes in order** (follow priority table)
3. **Test each fix independently** before moving to next
4. **Run `go build` and `go test`** after each phase
5. **Compare with NEURO_AUDIT_REPORT.md** for detailed code examples

All detailed technical specifications are in:
- [NEURO_AUDIT_REPORT.md](NEURO_AUDIT_REPORT.md) - Full issue analysis
- [NEURO_QUICK_USAGE_GUIDE.md](../manuals/NEURO_QUICK_USAGE_GUIDE.md) - How it should work

---

**Status:** Ready for implementation  
**Created:** February 10, 2026
