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

# 13 - Professional Reporting & Delivery

> **For:** All operators  
> **Read Time:** 20 minutes  
> **Difficulty:** ⭐⭐ Beginner  
> **Last Updated:** February 8, 2026

---

## Overview

> <span style="background-color: #7c2d12; color: white; padding: 2px 6px; border-radius: 3px;">**REPORTING**</span> Generate professional penetration test reports from VaporTrace findings. Customizable formats for different stakeholders.

### Report Generation

```bash
> report
[cyan]REPORT:[-] Generating report...
[blue]ANALYSIS:[-] Processing 47 vulnerabilities
[blue]METRICS:[-] 5 critical, 12 high, 18 medium
[green]✓ REPORT:[-] Generated: VaporTrace_PenTest_20260208.md
[cyan]FILE:[-] /reports/VaporTrace_PenTest_20260208.md
```

### Export Formats

```bash
# Markdown (default, GitHub-friendly)
> report --format markdown
[green]✓ Saved:[-] report.md

# PDF (professional documents)
> report --format pdf
[green]✓ Saved:[-] report.pdf

# HTML (interactive viewer)
> report --format html
[green]✓ Saved:[-] report.html

# JSON (programmatic access)
> report --format json
[green]✓ Saved:[-] report.json

# SARIF (security analysis exchange)
> report --format sarif
[green]✓ Saved:[-] report.sarif
```

### Report Customization

```bash
# Executive summary only
> report --template executive
[cyan]REPORT:[-] Executive summary format

# Detailed technical
> report --template technical
[cyan]REPORT:[-] Technical deep-dive format

# Remediation focused
> report --template remediation
[cyan]REPORT:[-] Remediation guidance format

# CVSS scoring detailed
> report --template cvss
[cyan]REPORT:[-] CVSS analysis detailed
```

### Report Sections

```
1. Executive Summary
   ├── Overall Risk Rating
   ├── Key Findings
   └── Recommendations

2. Vulnerability Details
   ├── ID, Title, Severity
   ├── Description
   ├── Impact Analysis
   ├── Proof of Concept
   └── Remediation Steps

3. Metrics & Analytics
   ├── Vulnerability Distribution
   ├── CVSS Score Breakdown
   └── Test Coverage

4. Appendices
   ├── Tested Endpoints
   ├── Tools Used
   └── Timeline
```

### Filtering Findings

```bash
# Only critical findings
> report --severity critical

# Exclude informational
> report --min-severity high

# Specific vulnerability types
> report --type BOLA,SSRF

# Date range
> report --since 2026-02-01 --until 2026-02-08
```

### Exfiltration & Redaction

```bash
# Mask sensitive data
> report --redact-pii
[green]✓ Report:[-] PII masked (emails, IPs)

# Exclude credentials from report
> report --exclude-loot
[green]✓ Report:[-] Credentials excluded

# Custom redaction patterns
> report --redact 'password|api_key|secret'
[green]✓ Report:[-] Custom patterns redacted
```

---

## Report Templates

### Template: Executive Summary

**Best For:** C-level stakeholders

```markdown
# VaporTrace Penetration Test Report
## Executive Summary

**Test Date:** February 1-8, 2026  
**Target:** api.vulnerable.com  
**Overall Risk:** 🔴 CRITICAL

### Key Findings
- 5 Critical vulnerabilities discovered
- Direct administrative access possible
- Customer data exposure
- Remediation priority: IMMEDIATE

### Recommendations
1. Emergency patching (within 24 hours)
2. Credential rotation
3. Security audit
```

### Template: Technical Detail

**Best For:** Security teams

```markdown
## Vulnerability: BOLA in /api/users/{id}

**Severity:** Critical (CVSS 9.8)  
**CWE:** CWE-639 (Authorization Bypass)

### Technical Details
- Endpoint: GET /api/users/{id}/profile
- Method: Direct object reference (no auth check)
- Impact: Access to all user profiles

### Proof of Concept
1. Authenticate as User A
2. Request: GET /api/users/999/profile
3. Response: User 999's data without authorization
4. Verified: Can access any user (1-100000)

### Remediation
Implement user-context verification:
```go
if user.ID != requestedID {
    return 403 Forbidden
}
```
```

---

## Best Practices

✅ **DO:**
- Review loot before including
- Use professional templates
- Include CVSS scores
- Provide remediation guidance
- Add timeline of activities

❌ **DON'T:**
- Include unverified findings
- Expose client credentials
- Use inflammatory language
- Report duplicates
- Make unsupported claims

---

## Integration Workflow

```bash
# Complete testing workflow
> target https://api.vulnerable.com
> map
> bola
> bfla
> ssrf
> report --format pdf
> loot export --format csv

[green]✓ COMPLETE:[-] Report + Loot exported
```

---

**See Also:**
- [04_STRATEGIC_PLANNING.md](04_STRATEGIC_PLANNING.md) - Planning phase
- [18_COMMAND_REFERENCE.md](18_COMMAND_REFERENCE.md) - Report commands

---

**Last Updated:** February 8, 2026  
**Version:** 1.0 - Production Ready
