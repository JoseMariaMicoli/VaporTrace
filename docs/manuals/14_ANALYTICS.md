# 14 - Analytics & Metrics Dashboard

> **For:** All operators  
> **Read Time:** 15 minutes  
> **Difficulty:** ⭐⭐ Beginner  
> **Last Updated:** February 8, 2026

---

## Overview

> <span style="background-color: #1e40af; color: white; padding: 2px 6px; border-radius: 3px;">**METRICS**</span> Real-time metrics and analytics dashboard for penetration tests. Monitor performance, success rates, and vulnerability distribution.

### Accessing Analytics

```bash
Press: [F5]  (Metrics Tab)
```

### Dashboard Layout

```
┌─────────────────────────────────────────────────────────────────┐
│ VaporTrace - Metrics Dashboard                     [F1-F7 Tabs] │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  📊 Test Statistics          🎯 Vulnerability Distribution      │
│  ├─ Total Requests: 15,234   ├─ Critical: 5 (🔴)               │
│  ├─ Success Rate: 92.3%      ├─ High: 12 (🟠)                  │
│  ├─ Avg Response: 245ms      ├─ Medium: 18 (🟡)                │
│  ├─ Data Transfer: 2.3GB     └─ Low: 34 (🟢)                   │
│  └─ Duration: 4h 32m         Total: 69 vulnerabilities         │
│                                                                 │
│  📈 Performance Metrics      ⏱️  Timeline                        │
│  ├─ Req/Sec: 1.2             ├─ Start: 2026-02-08 10:00        │
│  ├─ Error Rate: 7.7%         ├─ End: 2026-02-08 14:32          │
│  ├─ Avg Latency: 245ms       └─ Runtime: 4h 32m                │
│  └─ Peak Throughput: 5.2MB/s                                   │
│                                                                 │
│  🎯 Module Performance       📍 Attack Coverage                 │
│  ├─ BOLA: 8/10 targets      ├─ Endpoints: 340/500 (68%)        │
│  ├─ BFLA: 5/10 targets      ├─ Parameters: 2300/3000 (77%)     │
│  ├─ SSRF: 3/10 targets      ├─ Auth Methods: 5/5 (100%)        │
│  └─ Other: 12/20 targets    └─ Coverage: 81% (Good)            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
[Command Input]> _
```

### Real-Time Metrics

```bash
# View live metrics
> metrics
[cyan]METRICS:[-] Live dashboard
[blue]REQUESTS:[-] 1,234 so far
[blue]RATE:[-] 1.2 req/sec
[blue]SUCCESS:[-] 92.3%
[blue]ERRORS:[-] 7.7%

# Export metrics
> metrics export --format json
[green]✓ EXPORTED:[-] metrics.json
```

### Vulnerability Distribution

```
Severity Breakdown:
🔴 CRITICAL: 5 (7.2%)
  • BOLA in /api/users
  • SSRF in /api/fetch
  • BFLA in /api/admin

🟠 HIGH: 12 (17.4%)
  • Path traversal
  • SQL injection
  • Command injection

🟡 MEDIUM: 18 (26.1%)
  • Information disclosure
  • Weak authentication
  • Rate limiting bypass

🟢 LOW: 34 (49.3%)
  • Misconfiguration
  • Default credentials
  • Outdated dependencies
```

### Performance Timeline

```mermaid
graph LR
    A["Start<br/>10:00"] -->|1h| B["Discovery<br/>Phase"]
    B -->|1.5h| C["Exploitation<br/>Phase"]
    C -->|1.5h| D["Verification<br/>Phase"]
    D -->|30m| E["Reporting<br/>14:32"]
    
    style A fill:#2563eb,color:#fff
    style E fill:#10b981,color:#fff
```

### Custom Dashboards

```bash
# Create custom metric view
> metrics custom --name "executive"
  --include severity,count,remediation
  --exclude findings-detail

[green]✓ DASHBOARD:[-] "executive" created

# View custom dashboard
> metrics view executive
```

---

## Metrics Interpretation

| Metric | Interpretation |
|--------|-----------------|
| **Success Rate > 95%** | Efficient testing, good coverage |
| **Success Rate < 80%** | Network issues or WAF blocking |
| **Error Rate > 10%** | Target instability or evasion enabled |
| **Avg Latency > 500ms** | Slow target or heavy WAF |
| **Coverage > 80%** | Comprehensive testing achieved |

---

## Exports

```bash
# Export all metrics
> metrics export --format csv --output metrics.csv

# Chart-ready format
> metrics export --format json-chart

# Stakeholder summary
> metrics export --format summary
```

---

**See Also:**
- [13_REPORTING.md](13_REPORTING.md) - Formal reporting
- [04_STRATEGIC_PLANNING.md](04_STRATEGIC_PLANNING.md) - Planning integration

---

**Last Updated:** February 8, 2026  
**Version:** 1.0 - Production Ready
