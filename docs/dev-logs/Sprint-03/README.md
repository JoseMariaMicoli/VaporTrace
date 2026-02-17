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

# Sprint 3: Authentication & Authorization - BOLA, BFLA, BOPLA

**Status:** ✅ COMPLETE | **Version:** v1.2-AuthZ | **Released:** Q4 2025

---

## 🎯 Sprint Overview

Sprint 3 implements the core exploitation engines for OWASP API Top 10 authorization vulnerabilities. This sprint delivers the BOLA (Broken Object Level Authorization), BFLA (Broken Function Level Authorization), and BOPLA (Broken Object Property Level Authorization) engines with tactical ID-swapping and mass assignment fuzzing.

**Slogan:** "Breaking Authorization at Every Level"

---

## 📋 Deliverables

### 3.1: BOLA Prober (API1 - Broken Object Level Authorization) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/logic/bola.go`

**Features Delivered:**
- **ID-Swapping Engine** - Systematic object ID enumeration and testing
- **Session Store Management** - Multi-user context handling
- **Response Diffing** - Automated false positive elimination
- **Concurrency Support** - Multi-threaded ID probing
- **Result Aggregation** - Confidence scoring for findings
- **Exploitation Verification** - Confirm authorization bypass

**BOLA Testing Workflow:**
```go
type BOLAContext struct {
    TargetURL      string              // Base endpoint (e.g., /api/users/)
    TestIDs        []string            // IDs to test (1, 100, admin, etc)
    AuthTokens     map[string]string   // User -> Token mapping
    Wordlist       []string            // Common ID patterns
    ResponseCache  map[string]string   // Store baseline responses
    Results        []BOLAFinding       // Positive results
}

type BOLAFinding struct {
    UserID        string              // Authenticated as
    TargetID      string              // Object ID tested
    AccessGranted bool                // Authorization bypass confirmed
    ResponseDiff  string              // Difference from expected
    Severity      string              // CRITICAL, HIGH, MEDIUM, LOW
}
```

**ID Testing Patterns:**
```bash
Test IDs:
  Numeric: 1, 2, 10, 100, 1000, admin_id
  String-based: "admin", "root", "test", "demo"
  UUID-like: [random-uuid]
  Sequential: enumeration of numeric ranges
```

**Example Exploitation:**
```bash
> bola https://api.example.com/api/users/
[cyan]BOLA:[-] Starting object enumeration...
[yellow]Testing with authenticated token...
[08:30:00] User 1: 200 OK (owned by attacker)
[08:30:01] User 2: 403 Forbidden (not owned)
[08:30:02] User 100: 200 OK (UNAUTHORIZED ACCESS!)
[green]CRITICAL FINDING:[-] User 100 accessible without authorization
[yellow]Confidence:[-] 95% (response contains PII: email, phone)
```

**Status:** ✅ Production-ready with false positive detection

---

### 3.2: BOPLA/Mass Assignment (API3 - Broken Object Property Level Authorization) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/logic/bopla.go`

**Features Delivered:**
- **Property Fuzzing Engine** - Systematic JSON property injection
- **Hidden Property Detection** - Discover non-documented fields
- **Privilege Escalation Testing** - Role/permission property injection
- **Mass Assignment Exploitation** - Modify properties in bulk
- **Response Analysis** - Detect successful modifications
- **Concurrent Testing** - High-speed property enumeration

**BOPLA Testing Approach:**
```go
type BOPLAContext struct {
    TargetURL      string              // PATCH/PUT endpoint
    Wordlist       []string            // Property names to try
    PayloadTemplate map[string]interface{} // Base JSON object
    TestProperties []string            // [admin, role, owner_id, verified]
    Results        []BOPLAFinding       // Modified properties
}

type BOPLAFinding struct {
    Property      string              // e.g., "is_admin"
    OriginalValue interface{}         // Original value
    ModifiedValue interface{}         // Value we injected
    Accepted      bool                // Server accepted modification
    Severity      string              // Impact level
}
```

**Property Dictionary:**
```go
commonProperties := []string{
    // Privilege/Role
    "admin", "is_admin", "role", "roles", "permission", "permissions",
    "owner_id", "owner", "created_by", "modified_by",
    
    // Sensitive Data
    "verified", "active", "disabled", "deleted", "archived",
    "email_verified", "phone_verified", "2fa_enabled",
    
    // Pricing/Limits
    "price", "cost", "discount", "plan", "subscription",
    "api_quota", "rate_limit", "concurrent_sessions",
    
    // Tracking/Metadata
    "user_id", "account_id", "organization_id", "tenant_id",
}
```

**Example Attack:**
```bash
> bopla https://api.example.com/api/profile
[cyan]BOPLA:[-] Starting property injection fuzzing...
[yellow]Base Request:[-]
  {
    "name": "Attacker",
    "email": "attacker@test.com"
  }

[yellow]Injection Attempts:[-]
[08:35:00] Trying: "admin": true → 200 OK (ACCEPTED!)
[08:35:01] Trying: "role": "superuser" → 200 OK (ACCEPTED!)
[08:35:02] Trying: "is_verified": true → 200 OK (ACCEPTED!)

[green]CRITICAL FINDINGS:[-]
  - Can self-promote to admin
  - Can self-verify email
  - Can modify account metadata
```

**Status:** ✅ Production-ready with high detection rate

---

### 3.3: BFLA Module (API5 - Broken Function Level Authorization) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/logic/bfla.go`

**Features Delivered:**
- **HTTP Method Tampering** - Test GET->POST, DELETE->PUT transformations
- **Hierarchical Access Testing** - Admin function access from user accounts
- **Action Bypass Detection** - Privilege escalation via method override
- **Endpoint Mutation** - Test alternative paths and methods
- **Result Aggregation** - Compile method-override findings
- **Concurrency Engine** - Multi-threaded method probing

**BFLA Testing Methodology:**
```go
type BFLAContext struct {
    TargetURL      string              // Base endpoint
    Methods        []string            // [GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS]
    UserTokens     []string            // Different auth levels
    Results        []BFLAFinding       // Method override findings
}

type BFLAFinding struct {
    Endpoint       string              // e.g., /api/users
    Method         string              // e.g., DELETE
    Result         string              // 200, 401, 403, 500
    RequiredRole   string              // Expected role
    AccessGranted  bool                // Unexpected access
    Severity       string              // Impact
}
```

**HTTP Method Matrix:**
```bash
Test Matrix:
  GET     /api/users/{id}        → Expected: 200, Test others
  POST    /api/users             → Expected: 403, Test others
  PUT     /api/users/{id}        → Expected: 403, Test GET/DELETE
  DELETE  /api/users/{id}        → Expected: 403, Test PUT/POST
  PATCH   /api/users/{id}        → Expected: 403, Try others
  OPTIONS /api/users             → Information disclosure
  HEAD    /api/users/{id}        → Returns headers only
```

**Example Exploitation:**
```bash
> bfla https://api.example.com/api/admin/users
[cyan]BFLA:[-] Testing method-level authorization...
[yellow]Base User Token:[-] user_token_123
[yellow]Testing Admin Function:[-] DELETE /api/admin/users/99

Method Tests:
[08:40:00] GET    /api/admin/users/99 → 403 Forbidden (expected)
[08:40:01] POST   /api/admin/users/99 → 403 Forbidden (expected)
[08:40:02] DELETE /api/admin/users/99 → 403 Forbidden (expected)
[08:40:03] PUT    /api/admin/users/99 → 200 OK (UNAUTHORIZED!)
[08:40:04] PATCH  /api/admin/users/99 → 405 Method Not Allowed

[green]HIGH FINDING:[-] PUT method bypasses DELETE authorization
Access Level: User can DELETE via PUT method override
Impact: Unauthorized data deletion via alternative HTTP method
```

**Status:** ✅ Production-ready with method matrix testing

---

## 🔄 Current Status by Sub-Phase

| Sub-Phase | Deliverable | Status | Completion |
|-----------|-------------|--------|------------|
| **3.1** | BOLA Prober | ✅ DONE | 100% |
| **3.2** | BOPLA/Mass Assignment | ✅ DONE | 100% |
| **3.3** | BFLA Module | ✅ DONE | 100% |

---

## 📊 Code Metrics

| Metric | Value |
|--------|-------|
| **New Files** | 3 modules (bola, bopla, bfla) |
| **Lines of Code** | ~1800 LOC |
| **Test IDs** | 50+ patterns |
| **HTTP Methods** | 7 tested (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS) |
| **Concurrency** | Worker pool support |

---

## 🎓 Architecture Decisions

### BOLA ID Testing Strategy
- Sequential enumeration combined with dictionary-based testing
- Response diffing to eliminate false positives
- Per-user testing to identify cross-user authorization bypass
- Confidence scoring based on response analysis

### BOPLA Property Fuzzing
- Dictionary-driven approach covering common exploitation patterns
- Request mutation at the JSON property level
- Response-based verification (did server accept modification?)
- No destructive modifications to production data

### BFLA Method Override Testing
- Systematic HTTP method matrix traversal
- Tests both expected and unexpected method combinations
- Hierarchical testing (user permissions on admin endpoints)
- Identifies common method-override vulnerabilities

---

## 🚀 Next Steps

Sprint 4 implements injection and resource exhaustion engines:
- SSRF (Server-Side Request Forgery) testing
- Resource exhaustion (pagination/limit probing)
- Security misconfiguration auditing
- Unsafe consumption in webhooks

---

## 📚 References

- **OWASP API Security Top 10:** https://owasp.org/www-project-api-security/
- **BOLA/IDOR:** https://owasp.org/www-project-web-security-testing-guide/
- **Go Concurrency:** https://pkg.go.dev/sync
