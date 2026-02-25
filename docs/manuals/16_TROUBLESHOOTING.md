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

# 16 - Troubleshooting & Common Issues

> **For:** All operators  
> **Read Time:** 25 minutes  
> **Difficulty:** ⭐⭐ Beginner-Intermediate  
> **Last Updated:** February 8, 2026

---

## Overview

> <span style="background-color: #ef4444; color: white; padding: 2px 6px; border-radius: 3px;">**ISSUE**</span> Diagnosis and solutions for common VaporTrace problems.

---

## Quick Diagnosis Tree

```mermaid
graph TD
    A["🚨 Problem?"] --> B{Type?}
    B -->|Connection| C["Network Issues"]
    B -->|API Issues| D["API Discovery"]
    B -->|Too Slow| E["Performance"]
    B -->|Errors| F["Crashes/Errors"]
    B -->|Detection| G["WAF/IDS Detection"]
    
    C -->|Timeout| C1["Increase timeout"]
    C -->|SSL Error| C2["Disable SSL verify"]
    
    D -->|No endpoints| D1["Enable verbose mode"]
    D -->|Wrong data| D2["Check filters"]
    
    E -->|Hung process| E1["Reduce concurrency"]
    E -->|Memory leak| E2["Restart process"]
    
    F -->|Segfault| F1["Check config"]
    F -->|Panic| F2["Review logs"]
    
    G -->|Blocked| G1["Enable evasion"]
    G -->|Slow| G2["Increase delays"]
```

---

## Network Issues

### Connection Timeout

**Problem:** `ERROR: Connection timeout after 30s`

**Causes:**
- Target is slow or unreachable
- Network path blocked
- Firewall rules

**Solutions:**

```bash
# Increase timeout
> config set network.timeout 60

# Or per-target
> scan https://example.com --timeout 120

# Test connectivity
> network diagnose https://example.com
[green]✓ REACHABLE:[-] Target online
[green]✓ DNS:[-] Resolves correctly
[green]✓ PORT:[-] 443 open
[red]✗ SSL:[-] Certificate invalid (ignore with --no-verify)
```

### SSL Certificate Error

**Problem:** `ERROR: certificate verify failed`

**Causes:**
- Self-signed certificate
- Outdated CA bundle
- MITM proxy present

**Solutions:**

```bash
# Disable SSL verification (risky, use carefully)
> scan https://example.com --no-verify
[yellow]⚠ WARNING:[-] SSL verification disabled

# Or globally
> config set http.verify_ssl false

# Check certificate
> network cert-info https://example.com
[cyan]CERTIFICATE:[-]
  Subject: CN=example.com
  Issuer: CN=example.com (self-signed)
  Expires: 2025-12-31
```

### DNS Resolution Failed

**Problem:** `ERROR: Failed to resolve example.com`

**Causes:**
- DNS server down
- Domain doesn't exist
- Network isolation

**Solutions:**

```bash
# Use different DNS
> config set network.dns_server 8.8.8.8

# Test DNS
> network dns-check example.com
[green]✓ RESOLVES:[-] 192.0.2.1

# Try IP directly
> scan https://192.0.2.1 --host example.com
```

### Proxy Connection Failed

**Problem:** `ERROR: Unable to connect to proxy`

**Causes:**
- Proxy URL incorrect
- Proxy server down
- Authentication failed

**Solutions:**

```bash
# Verify proxy
> network proxy-check http://proxy.example.com:8080
[red]✗ UNREACHABLE:[-] Connection refused

# Test with auth
> config set proxy.url http://user:pass@proxy.example.com:8080

# Disable proxy
> config set proxy.enabled false
```

---

## API Discovery Issues

### No Endpoints Found

**Problem:** `WARNING: No endpoints discovered`

**Causes:**
- Wrong target URL
- API not spidered correctly
- Endpoints hidden

**Solutions:**

```bash
# Enable verbose mode
> scan https://api.example.com -vvv

[cyan]VERBOSE:[-] Trying spider...
[cyan]VERBOSE:[-] Found 500 endpoints
[cyan]VERBOSE:[-] Parsing Swagger...

# Manually specify endpoints
> scan https://api.example.com \
    --endpoints /api/v1/users,/api/v1/products

# Check robots.txt and sitemap
> network fetch-urls https://example.com/robots.txt
[green]✓ FOUND:[-]
  /api/v1/*
  /admin/*
```

### Wrong Data Being Extracted

**Problem:** `Endpoints found but data looks incomplete`

**Causes:**
- Filters too restrictive
- Authentication required
- Wrong content-type

**Solutions:**

```bash
# Remove filters
> scan https://api.example.com --no-filters

# Add authentication
> scan https://api.example.com \
    --header "Authorization: Bearer token123"

# Check content-type
> network fetch https://api.example.com/api/v1 -H "Accept: application/json"
```

### Swagger/OpenAPI Not Parsing

**Problem:** `ERROR: Failed to parse Swagger spec`

**Causes:**
- Invalid YAML/JSON
- Version unsupported
- Format issue

**Solutions:**

```bash
# Check Swagger URL
> network fetch https://api.example.com/swagger.json | jq .

# Validate format
> discovery validate-swagger swagger.json
[red]✗ ERROR:[-] Invalid JSON on line 42

# Try alternative URL
> scan https://api.example.com/v1/swagger.json
```

---

## Performance Issues

### Process Hung/Frozen

**Problem:** VaporTrace appears frozen, no output

**Causes:**
- Too many concurrent requests
- Large endpoint list
- Deadlock

**Solutions:**

```bash
# Check status in another terminal
> status

# Reduce concurrency
> config set network.concurrent_requests 5

# Force kill and restart
> kill-scan <scan_id>
> scan https://example.com --workers 2

# Check memory
> system status
[cyan]MEMORY:[-] 85% (2.1GB/2.5GB)
```

### Memory Leak / Growing Memory

**Problem:** Memory usage keeps increasing

**Causes:**
- Caching issue
- Connection pool leak
- Large result set

**Solutions:**

```bash
# Disable cache
> config set performance.cache_enabled false

# Reduce batch size
> config set performance.batch_size 100

# Restart daemon
> daemon restart

# Monitor memory
> system monitor --watch
[cyan]MEMORY:[-] 512MB (↓ decreasing)
```

### Very Slow / Hanging on Specific Target

**Problem:** Scan works on other targets but slow on this one

**Causes:**
- Target rate limiting
- WAF throttling
- Network path issues

**Solutions:**

```bash
# Enable evasion
> config set evasion.multiplier 2.0

# Add delays
> config set evasion.jitter.min 100
> config set evasion.jitter.max 500

# Test manually
> network fetch https://target.example.com -v

# Check target status
> network ping-test https://target.example.com
[cyan]LATENCY:[-] 2500ms (very slow)
[yellow]⚠ WARNING:[-] Target may be rate-limiting
```

---

## Error Messages & Solutions

### Segmentation Fault

```
[FATAL] Segmentation fault (core dumped)
```

**Solutions:**
```bash
# Check configuration
> config view
> config validate

# Reset and retry
> config reset
> scan https://example.com

# Report bug with logs
> logs export --last-crash
```

### Out of Memory

```
ERROR: Cannot allocate memory
```

**Solutions:**
```bash
# Reduce workers
> config set performance.workers 2

# Disable cache
> config set performance.cache_enabled false

# Kill other processes
> system processes

# Use system limit
> ulimit -m 1000000
```

### Panic: Index out of bounds

```
PANIC: index out of bounds [5] with length 4
```

**Solutions:**
```bash
# Usually a bug - update VaporTrace
> version
VaporTrace v3.2.1 (outdated)

> update
[green]✓ UPDATED:[-] v3.2.1 → v3.3.0

# Retry scan
> scan https://example.com
```

---

## WAF/IDS Detection Issues

### Getting Blocked by WAF

**Problem:** `ERROR: 403 Forbidden` or `ERROR: 429 Too Many Requests`

**Causes:**
- WAF signature match
- Rate limiting
- Behavioral detection

**Solutions:**

```bash
# Enable full evasion
> config set evasion.mode silent
> config set evasion.multiplier 3.0

# Enable obfuscation
> config set evasion.obfuscation.path true
> config set evasion.obfuscation.payload true

# Retry with delays
> config set evasion.jitter.min 500
> config set evasion.jitter.max 2000

# Check what's detected
> scan https://example.com -vvv
[yellow]⚠ WAF:[-] Detected payload pattern in parameter 'id'
[cyan]SOLUTION:[-] Try encoding with UTF-8

# Use Ghost Weaver
> ghost enable
> ghost payload-obfuscate
> scan https://example.com
```

### Getting Slow/Rate-Limited

**Problem:** Responses getting slower, 429 responses

**Causes:**
- Rate limiter activated
- Behavioral pattern detected
- Too many requests

**Solutions:**

```bash
# Increase delays significantly
> config set evasion.jitter.min 2000
> config set evasion.jitter.max 5000

# Spread over time
> scan https://example.com --timeout 600

# Use proxy rotation
> config set proxy.rotation random
> scan https://example.com --proxies proxy-list.txt

# Check headers for limits
> network fetch https://api.example.com \
    --show-headers | grep -i "rate\|limit\|retry"
```

---

## Help & Support

### Get Help for Specific Issue

```bash
> help troubleshoot <keyword>
[cyan]HELP:[-]
  timeout - Connection timeout solutions
  ssl - SSL certificate errors
  waf - WAF detection evasion
  memory - Memory leak diagnosis
  slow - Performance optimization

> help troubleshoot timeout
[Shows detailed timeout troubleshooting guide]
```

### Generate Diagnostic Report

```bash
> diagnostics generate
[green]✓ GENERATED:[-] diagnostics-2025-02-08.tar.gz

# Includes:
# - System information
# - Configuration
# - Recent logs
# - Network tests
# - Memory usage

# Send to support
> diagnostics send --email support@vaportrace.io
```

### View Recent Logs

```bash
# Last 50 lines
> logs tail -n 50

# Search for errors
> logs grep ERROR

# Export for analysis
> logs export --last-24h
```

---

## Performance Baselines

| Metric | Baseline | Good | Poor |
|--------|----------|------|------|
| **Endpoint Discovery** | 1-5s | <10s | >60s |
| **Memory Usage** | 100MB | <500MB | >2GB |
| **CPU Usage** | 20% | <60% | >90% |
| **Request Rate** | 100/s | 50-200/s | <10/s |
| **API Response Time** | <200ms | <500ms | >2s |

---

## Escalation

### When to Restart

- After changing major configurations
- After crashes
- After memory warnings
- Every 24 hours in production

```bash
> daemon restart
[green]✓ RESTARTED:[-] VaporTrace daemon
```

### When to Seek Support

- Persistent crashes
- Memory leaks not resolved
- WAF evasion not working
- Network errors not resolved

**Contact:** support@vaportrace.io  
**Docs:** https://docs.vaportrace.io  
**Community:** https://discord.gg/vaportrace

---

**Last Updated:** February 8, 2026  
**Version:** 1.0 - Production Ready
