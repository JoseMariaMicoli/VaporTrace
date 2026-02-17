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

# Neuro Quick Reference - What Works & What Needs Fixing

## 📋 Complete Command Reference

```bash
# Enable/Disable Neural Engine
neuro on                    # Activate neuro engine
neuro off                   # Deactivate neuro engine

# Configure AI Provider (CORRECTED SYNTAX)
neuro config groq llama-3.1-8b gsk_xxxxx
neuro config openai gpt-4o sk_proj_xxxxx
neuro config google gemini-1.5-flash AIza_xxxxx
neuro config ollama mistral
neuro config hybrid gpt-4o sk_proj_xxxxx

# Test Connectivity
test-neuro                  # Verify provider is reachable

# Generate Payloads
neuro-gen "SQL injection on /api/search" 10
neuro-gen "XSS in profile form"
neuro-gen "JWT authentication bypass" 15

# Ask Direct Questions
ask "How do I bypass JWT signature verification?"
ask "Generate authentication bypass payloads"
```

---

## ✅ What Currently Works

| Feature | Status | Details |
|---------|--------|---------|
| Enable/disable neuro | ✅ | `neuro on` / `neuro off` |
| Command dispatch | ✅ | Routes to correct handlers |
| Integration with attacks | ✅ | Neuro calls in pipeline/scan |
| Multi-provider support | ✅ | 5 providers implemented |
| Hybrid fallback | ✅ | Cloud + local working |
| Rate limiting | ⚠️ | Works but too aggressive (6s) |
| Error logging | ⚠️ | Logs but misses some cases |

---

## ❌ What Needs Fixing

### Issue #1: Command Configuration (CRITICAL)
```
CURRENT:  neuro config --provider groq --apikey gsk_xxxxx
CORRECT:  neuro config groq llama-3.1-8b gsk_xxxxx

STATUS: Documentation fixed ✅ | Code needs verification ⏳
```

### Issue #2: Missing Nil Checks (HIGH)
```
PROBLEM: LLM responses can be nil → code panics
IMPACT: Random crashes under load
ESTIMATE: 30 minutes to fix

AFFECTED LOCATIONS:
- pkg/engine/neuro_engine.go (~line 250)
- pkg/logic/neuro_engine.go (~line 180)
- pkg/ai/client.go (multiple locations)
```

### Issue #3: Rate Limiting Too Slow (HIGH)
```
CURRENT: 6 seconds between API calls
CORRECT: 2 seconds (industry standard)
IMPACT: Slows attacks by 3x unnecessarily
ESTIMATE: 5 minutes to fix

FILE: pkg/logic/neuro_engine.go (line ~425)
```

### Issue #4: Race Conditions (CRITICAL)
```
PROBLEM: GlobalNeuro singleton not thread-safe
SCENARIO: Multiple goroutines access simultaneously
RESULT: Data corruption or panic
ESTIMATE: 20 minutes to fix

FILE: pkg/logic/neuro_engine.go (line ~50-100)
```

### Issue #5: Silent Errors (HIGH)
```
PROBLEM: neuro-gen silently fails on invalid count
CURRENT: > neuro-gen "payload" abc
         [no error message]

CORRECT: > neuro-gen "payload" abc
         [red]ERROR: Invalid count parameter: abc

FILE: pkg/engine/core.go (line ~280-295)
```

### Issue #6: Confidence Scoring (MEDIUM)
```
PROBLEM: Undefined status codes return confidence=0
AFFECTED: pkg/engine/neuro_engine.go (line ~320-350)
FIX: Add default case to switch statement
```

---

## 🎯 Usage Examples - All Features Together

### Example 1: Quick Setup (5 minutes)
```bash
> neuro on
> neuro config groq llama-3.1-8b gsk_xxxxx
> test-neuro
✓ Ready to use
```

### Example 2: Generate Custom Payloads (10 minutes)
```bash
> neuro config openai gpt-4o sk_proj_xxxxx
> neuro-gen "Authentication bypass for OAuth2" 20
[AI generates 20 unique payloads]
> (Copy payloads to your favorite testing tool)
```

### Example 3: Ask LLM Questions (5 minutes)
```bash
> neuro on
> ask "How do I test for BOLA vulnerabilities?"
[AI explains techniques]
> ask "Generate JWT forgery payloads"
[AI generates payloads]
```

### Example 4: Full Attack Workflow (30 minutes)
```bash
1. > target https://api.example.com
2. > neuro on
3. > neuro config hybrid gpt-4o sk_proj_xxxxx
4. > map                    # Discover endpoints
5. > neuro-gen "API auth bypass" 15
6. > bfla                   # Test with generated payloads
7. > ask "Summarize vulnerabilities"
8. > report                 # Generate report
```

---

## 📊 Implementation Roadmap

```
Phase 1: Command Parsing (10 min)
├─ Verify args array construction
└─ Confirm indices are correct

Phase 2: Race Conditions (20 min)
├─ Add sync.RWMutex to GlobalNeuro
└─ Implement double-checked locking

Phase 3: Nil Checks (30 min)
├─ Audit all LLM provider calls
├─ Add error handling
└─ Test with null responses

Phase 4: Error Handling (20 min)
├─ Add input validation
├─ Improve error messages
└─ Test edge cases

Phase 5: Performance (5 min)
├─ Reduce rate limit 6s → 2s
└─ Verify no throttling

Phase 6: Confidence Logic (15 min)
├─ Add default cases
└─ Test all status codes

TOTAL TIME: ~1.5 hours
```

---

## 🚀 Getting Started

### For Users: Use the Quick Usage Guide
→ Read: [NEURO_QUICK_USAGE_GUIDE.md](../manuals/NEURO_QUICK_USAGE_GUIDE.md)

**Quick Start:**
```bash
> neuro on
> neuro config groq llama-3.1-8b <YOUR_API_KEY>
> test-neuro
> neuro-gen "SQL injection" 5
```

### For Developers: Implement the Fixes
→ Read: [NEURO_SOURCE_FIXES_NEEDED.md](NEURO_SOURCE_FIXES_NEEDED.md)

**Priority:**
1. Fix command parsing verification
2. Fix race conditions
3. Add nil checks
4. Improve error handling
5. Reduce rate limiting
6. Fix confidence logic

---

## 📁 Documentation Structure

```
docs/
├── manuals/
│   ├── 07_AI_NEURO_ENGINE.md              ✅ Updated
│   └── NEURO_QUICK_USAGE_GUIDE.md         ✅ New
├── dev-logs/
│   ├── NEURO_AUDIT_REPORT.md              ✅ All 15 issues
│   ├── NEURO_AUDIT_EXECUTIVE_SUMMARY.md   ✅ Quick overview
│   ├── NEURO_AUDIT_VISUAL_MAP.md          ✅ Diagrams
│   ├── NEURO_SOURCE_FIXES_NEEDED.md       ✅ This one
│   └── AUDIT_SUMMARY_FOR_TEAM.md          ✅ Team briefing
```

---

## 🧪 Testing Commands

After implementing fixes, test with:

```bash
# Test all commands
> neuro on
> neuro config groq llama-3.1-8b gsk_xxxxx
> test-neuro
> neuro-gen "test payload" 3
> ask "Test query"

# Test error handling
> neuro-gen "test" abc              # Should error
> neuro config invalid-provider x   # Should error
> ask                               # Should error (missing prompt)

# Test under load (concurrent access)
> (Run multiple commands simultaneously)
# Should not panic or corrupt data
```

---

## ✨ When All Fixes Are Done

✅ All 8 documented commands work  
✅ All errors have helpful messages  
✅ No race conditions under load  
✅ Rate limiting is reasonable (2s)  
✅ AI provider responses are validated  
✅ Confidence scoring is reliable  
✅ Full production-ready neuro engine  

**Estimated total implementation time: 1.5-2 hours**

---

**Last Updated:** February 10, 2026  
**Status:** Ready for implementation
