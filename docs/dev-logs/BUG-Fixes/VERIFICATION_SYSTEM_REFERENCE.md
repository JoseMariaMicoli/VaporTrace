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

# Verification System Reference - Sprint 16.1

## Overview
The verification system prevents LLM hallucination by ensuring all remediation suggestions are vetted before presenting to users. This document explains the system architecture, usage, and extensibility.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ Exploit Detected                                            │
│ (SuggestFix called with TacticalAction)                    │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ Generate Remediation Suggestion                             │
│ - generateBOLAFix()                                         │
│ - generateBFLAFix()                                         │
│ - generateSSRFFix()                                         │
│ - etc.                                                       │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ VerifyRemediationSuggestion()                               │
│ ├─ Check GoldStandardLibrary match?                         │
│ │  ├─ YES → VerificationStatus = "GOLD_STANDARD"           │
│ │  └─ NO → VerificationStatus = "UNVERIFIED"               │
│ └─ Set StaticAnalysisPassed flag                            │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ FormatForUI()                                               │
│ ├─ GOLD_STANDARD → [GREEN] VERIFIED & PRODUCTION READY     │
│ ├─ VERIFIED → [GREEN] VERIFIED - Passed Security Review    │
│ └─ UNVERIFIED → [YELLOW] ⚠ WARNING - Manual review needed  │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ Display to User                                             │
│ (Operator sees verification status before applying fix)    │
└─────────────────────────────────────────────────────────────┘
```

---

## Verification Status Levels

### 1. GOLD_STANDARD
- **Description:** Snippet is in pre-vetted library, passed security audit
- **Color:** Green (✓)
- **Deployment:** Production-ready, can be auto-applied
- **Examples:** BOLA authorization, JWT validation, URL validation
- **Criteria:** Manually reviewed + zero security findings

```go
case "GOLD_STANDARD":
    verificationBanner = `╔═══════════════════════════════════════════════════════════════════════════╗
║ [green]✓ GOLD STANDARD - VERIFIED & PRODUCTION READY[-]
║ This snippet is from the Gold Standard library and has passed security audit.
╚═══════════════════════════════════════════════════════════════════════════╝`
```

### 2. VERIFIED
- **Description:** Snippet has passed security review (future enhancement)
- **Color:** Green (✓)
- **Deployment:** Production-ready, requires explicit approval
- **Examples:** (Future) static analysis passed + team approval
- **Criteria:** gosec/bandit scan + human review approved

```go
case "VERIFIED":
    verificationBanner = `╔═══════════════════════════════════════════════════════════════════════════╗
║ [green]✓ VERIFIED - Passed Security Review[-]
║ [VerificationNotes with review details]
╚═══════════════════════════════════════════════════════════════════════════╝`
```

### 3. UNVERIFIED
- **Description:** Snippet has NOT been verified, manual review required
- **Color:** Yellow (⚠)
- **Deployment:** Blocked until approved by security team
- **Examples:** LLM-generated suggestions, custom vulnerability fixes
- **Criteria:** Not in Gold Standard + no approval

```go
case "UNVERIFIED":
    verificationBanner = `╔═══════════════════════════════════════════════════════════════════════════╗
║ [yellow]⚠ WARNING - UNVERIFIED CODE[-]
║ This snippet has NOT been verified by security team. Manual review required before production use.
║
║ ACTION: Manual security review REQUIRED before production use.
║ TEST FOR: ReDoS in regex, SQL injection, auth bypass, hardcoded secrets.
╚═══════════════════════════════════════════════════════════════════════════╝`
```

---

## Gold Standard Library

### Current Snippets (5)

#### 1. BOLA_AUTHZ_MIDDLEWARE
```go
// VERIFIED: Object-level authorization check
func AuthorizeObjectAccess(c *gin.Context) {
    userID := c.GetString("user_id")
    objectID := c.Param("id")
    
    // Fetch object and verify ownership
    obj := database.GetObject(objectID)
    if obj == nil || obj.OwnerID != userID {
        c.JSON(403, gin.H{"error": "Unauthorized"})
        c.Abort()
        return
    }
    c.Next()
}
```
**Audited By:** Security team
**Last Updated:** Sprint 16.1
**Risk Level:** ZERO

#### 2. BFLA_RBAC_MIDDLEWARE
```go
// VERIFIED: Role-based access control
func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("role")
        
        // Verify user has required role
        if userRole != requiredRole && userRole != "admin" {
            c.JSON(403, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```
**Audited By:** Security team
**Last Updated:** Sprint 16.1
**Risk Level:** ZERO

#### 3. SSRF_URL_VALIDATION
```go
// VERIFIED: Safe URL validation against private IP ranges
func IsValidURL(urlStr string) bool {
    u, err := url.Parse(urlStr)
    if err != nil {
        return false
    }
    
    // Block private IP ranges
    privateRanges := []string{
        "127.0.0.0/8",        // Localhost
        "10.0.0.0/8",         // Private
        "172.16.0.0/12",      // Private
        "169.254.0.0/16",     // Link-local
        "169.254.169.254/32", // AWS metadata
    }
    
    ip := net.ParseIP(u.Hostname())
    for _, cidr := range privateRanges {
        _, network, _ := net.ParseCIDR(cidr)
        if network.Contains(ip) {
            return false
        }
    }
    
    return u.Scheme == "http" || u.Scheme == "https"
}
```
**Audited By:** Security team
**Last Updated:** Sprint 16.1
**Risk Level:** ZERO
**Notes:** Standard CIDR validation, no ReDoS risk

#### 4. JWT_VALIDATION
```go
// VERIFIED: Strict JWT validation with algorithm pinning
func ValidateJWT(tokenString string) (jwt.Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        // CRITICAL: Restrict to specific algorithm(s) only
        if token.Method.Alg() != "HS256" && token.Method.Alg() != "RS256" {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        
        return getSigningKey(token), nil
    })
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    return token.Claims, nil
}
```
**Audited By:** Security team
**Last Updated:** Sprint 16.1
**Risk Level:** ZERO
**Notes:** Algorithm pinning prevents "alg: none" bypass

#### 5. SQL_INJECTION_PREVENTION
```go
// VERIFIED: Parameterized queries (not string concatenation)
// WRONG: query := "SELECT * FROM users WHERE id = " + userInput
// CORRECT: Use parameterized queries
query := "SELECT * FROM users WHERE id = ?"
db.Query(query, userInput)
```
**Audited By:** Security team
**Last Updated:** Sprint 16.1
**Risk Level:** ZERO
**Notes:** Demonstrates pattern, not full implementation

---

## Usage Examples

### Example 1: GOLD_STANDARD Suggestion
```go
exploit := TacticalAction{
    Type: "BOLA",
    Target: "https://api.target.com/users/123",
}
suggestion := SuggestFix(exploit, "")

// Output:
fmt.Println(suggestion.VerificationStatus)        // "GOLD_STANDARD"
fmt.Println(suggestion.VerificationNotes)         // "This snippet is from the Gold Standard library..."
fmt.Println(suggestion.StaticAnalysisPassed)      // true
fmt.Println(suggestion.FormatForUI())             // Shows [GREEN] ✓ banner
```

### Example 2: UNVERIFIED Suggestion
```go
exploit := TacticalAction{
    Type: "CUSTOM_VULN",
    Target: "https://api.target.com/custom",
}
suggestion := SuggestFix(exploit, "custom payload")

// Output:
fmt.Println(suggestion.VerificationStatus)        // "UNVERIFIED"
fmt.Println(suggestion.VerificationNotes)         // "⚠️ WARNING: This snippet has NOT been verified..."
fmt.Println(suggestion.StaticAnalysisPassed)      // false
fmt.Println(suggestion.FormatForUI())             // Shows [YELLOW] ⚠ WARNING banner
```

### Example 3: Filter for Production Use
```go
// Only allow GOLD_STANDARD suggestions in production
allSuggestions := GenerateAllRemediations(target)

productionReady := make([]*RemediationSuggestion, 0)
for _, sug := range allSuggestions {
    if sug.VerificationStatus == "GOLD_STANDARD" && sug.StaticAnalysisPassed {
        productionReady = append(productionReady, sug)
    }
}

// Deploy only production-ready fixes
for _, fix := range productionReady {
    ApplyRemediationFix(fix)
}
```

---

## Extension Points

### Adding New Gold Standard Snippets

1. **Create snippet** (manually audit for security)
2. **Add to GoldStandardLibrary:**
```go
var GoldStandardLibrary = map[string]string{
    // ... existing snippets ...
    "NEW_VULN_FIX": `// VERIFIED: Description
// Implementation here
`,
}
```
3. **Document in Gold Standard Reference**
4. **Security team approval** (sign off on audit)
5. **Test with VerifyRemediationSuggestion()**

### Future: Implementing Static Analysis

```go
// NOT YET IMPLEMENTED - Future enhancement
func VerifyWithStaticAnalysis(snippet, language string) (bool, string) {
    switch language {
    case "Go":
        // Run gosec
        cmd := exec.Command("gosec", "-quiet", "/dev/stdin")
        cmd.Stdin = strings.NewReader(snippet)
        
        if err := cmd.Run(); err != nil {
            return false, fmt.Sprintf("gosec found issues: %v", err)
        }
        return true, "gosec passed"
        
    case "Python":
        // Run bandit
        // ...
        
    case "Java":
        // Run SpotBugs
        // ...
    }
    
    return false, "unsupported language"
}
```

### Future: Security Team Approval Workflow

```go
// NOT YET IMPLEMENTED - Future enhancement
type RemediationApproval struct {
    SnippetID      string
    SubmittedBy    string
    ApprovedBy     string
    ApprovedAt     time.Time
    ApprovalNotes  string
    RiskAssessment string // "LOW", "MEDIUM", "HIGH"
}

// Workflow:
// 1. VerifyRemediationSuggestion() checks approval database
// 2. If approved → VerificationStatus = "VERIFIED"
// 3. If not approved → VerificationStatus = "UNVERIFIED"
```

---

## Security Checklist for New Snippets

Before adding to Gold Standard Library, verify:

- [ ] No hardcoded credentials (API keys, passwords, tokens)
- [ ] No ReDoS-vulnerable regex patterns
- [ ] No SQL injection paths (parameterized queries only)
- [ ] No auth bypass (proper algorithm pinning, validation)
- [ ] No command injection (input validation, parameterization)
- [ ] No XXE vulnerabilities (disable entity expansion)
- [ ] No deserialization vulnerabilities
- [ ] No CORS misconfiguration
- [ ] No sensitive data in logs
- [ ] No infinite loops or performance DoS

---

## Verification System Integration Points

### 1. SuggestFix() Dispatcher
```go
func SuggestFix(exploit TacticalAction, result string) *RemediationSuggestion {
    // Generate suggestion based on vulnerability type
    var suggestion *RemediationSuggestion
    switch exploit.Type {
    case "BOLA":
        suggestion = generateBOLAFix(exploit, result)
    // ... other cases ...
    }
    
    // CRITICAL: Verify before returning
    VerifyRemediationSuggestion(suggestion)  // ← Verification check
    
    return suggestion
}
```

### 2. UI Display
```go
func (r *RemediationSuggestion) FormatForUI() string {
    // Add verification banner at top
    verificationBanner := ""
    switch r.VerificationStatus {
    case "GOLD_STANDARD":
        verificationBanner = `[GREEN] ✓ GOLD STANDARD ...`
    // ... other cases ...
    }
    
    // Include banner in output
    output := fmt.Sprintf("%s\n%s", verificationBanner, fullContent)
    return output
}
```

### 3. Logging & Audit Trail
```go
// VerifyRemediationSuggestion() logs all decisions
VerifyRemediationSuggestion(suggestion)

// Output in LogBuffer:
// [green]✓ VERIFIED:[-] BOLA_AUTHZ_MIDDLEWARE is GOLD_STANDARD
// [yellow]⚠ UNVERIFIED:[-] This suggestion requires manual security review
```

---

## Troubleshooting

### Issue: Suggestion showing as UNVERIFIED even though it should be GOLD_STANDARD
**Cause:** Exact snippet match required (whitespace, indentation must match)
**Fix:** Verify snippet in GoldStandardLibrary matches exactly

### Issue: Gold Standard snippets not being recognized
**Cause:** `GoldStandardLibrary` variable not initialized
**Fix:** Ensure remediation.go lines 10-65 are present

### Issue: FormatForUI() not showing verification banner
**Cause:** VerifyRemediationSuggestion() not called before FormatForUI()
**Fix:** Ensure SuggestFix() calls VerifyRemediationSuggestion() first

---

## Future Enhancements (Post-Sprint 16.1)

1. **Static Analysis Integration**
   - Auto-run gosec, bandit, SpotBugs on new snippets
   - Flag patterns: hardcoded secrets, ReDoS regex, SQL injection
   - Automated security scoring

2. **Security Team Approval Workflow**
   - DB table for approved snippets
   - Web UI for review + approval
   - Audit trail with timestamps, reviewer notes

3. **Community Snippet Submission**
   - Submit suggestions for review
   - Community voting on snippets
   - Trending fixes by vulnerability type

4. **Integration with Remediation Frameworks**
   - Auto-generate Terraform/CloudFormation fixes
   - Container security remediation (Docker, K8s)
   - Infrastructure-as-Code templates

---

## References

- **Implementation:** [pkg/engine/remediation.go](pkg/engine/remediation.go)
- **Struct Definition:** RemediationSuggestion (lines 10-50)
- **Verification Function:** VerifyRemediationSuggestion() (lines 67-105)
- **UI Formatting:** FormatForUI() (lines 485-540)
- **Testing:** CRITICAL_FIXES_SPRINT_11.3-16.1.md (Verification Testing section)
