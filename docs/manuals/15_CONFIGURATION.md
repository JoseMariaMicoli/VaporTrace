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

# 15 - Advanced Configuration & Tuning

> **For:** Advanced operators & DevOps  
> **Read Time:** 20 minutes  
> **Difficulty:** ⭐⭐⭐ Intermediate  
> **Last Updated:** February 8, 2026

---

## Overview

> <span style="background-color: #6366f1; color: white; padding: 2px 6px; border-radius: 3px;">**CONFIG**</span> Fine-tune VaporTrace for different target environments, network conditions, and security postures.

### Configuration File Location

```bash
~/.vapor/config.yaml      # User config (default)
/etc/vapor/config.yaml    # System config
./vapor.yaml              # Project config (overrides)
```

### Full Configuration Template

```yaml
# VaporTrace Configuration
vapor:
  # Network Settings
  network:
    timeout: 30
    retries: 3
    concurrent_requests: 10
    connection_pool_size: 100
    dns_cache: true
    
  # HTTP Settings
  http:
    user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
    follow_redirects: true
    max_redirects: 5
    verify_ssl: false
    
  # Proxy Settings
  proxy:
    enabled: false
    url: "http://proxy.example.com:8080"
    rotation: random
    
  # Evasion Settings
  evasion:
    enabled: true
    mode: aggressive        # aggressive|fast|silent|debug
    multiplier: 1.0
    
    jitter:
      enabled: true
      min: 10              # ms
      max: 150             # ms
      
    thinking_time:
      enabled: true
      
    rate_limit:
      backoff: true
      exponential: true
      
    obfuscation:
      path: true
      payload: true
      
    encoding:
      enabled: false
      types: ["gzip", "deflate"]
  
  # Database
  database:
    engine: sqlite
    path: /home/user/.vapor/vapor.db
    auto_backup: true
    
  # Logging
  logging:
    level: info             # debug|info|warn|error
    format: colored
    file: /var/log/vapor.log
    max_size: 100M
    retention: 30           # days
    
  # AI/LLM Settings
  ai:
    enabled: true
    provider: groq          # groq|openai|local|anthropic
    api_key: "${GROQ_API_KEY}"
    model: "llama2-70b"
    temperature: 0.7
    timeout: 30
    
  # Performance Tuning
  performance:
    batch_size: 1000
    cache_enabled: true
    cache_ttl: 3600
    workers: 10
    
  # Security
  security:
    tls_version: "1.2"
    allowed_protocols: ["https"]
    certificate_verify: true
```

### Environment Variables

```bash
# Network
export VAPOR_TIMEOUT=30
export VAPOR_RETRIES=3
export VAPOR_WORKERS=10

# Evasion
export VAPOR_EVASION_MODE=silent
export VAPOR_MULTIPLIER=1.5

# AI
export GROQ_API_KEY=gsk_xxxxx
export VAPOR_AI_PROVIDER=groq

# Database
export VAPOR_DB_PATH=/home/user/vapor.db

# Logging
export VAPOR_LOG_LEVEL=debug
export VAPOR_LOG_FILE=/var/log/vapor.log
```

---

## Target-Specific Configurations

### Light Target (Simple API)

```yaml
network:
  timeout: 10
  concurrent_requests: 20
  
evasion:
  mode: aggressive
  multiplier: 0.5        # Fast
```

### Heavy Target (Enterprise API)

```yaml
network:
  timeout: 60
  concurrent_requests: 5  # Slower
  
evasion:
  mode: silent
  multiplier: 2.0        # Slower delays
```

### WAF-Protected Target

```yaml
evasion:
  mode: silent
  multiplier: 3.0
  jitter:
    min: 50
    max: 500
  thinking_time:
    enabled: true
  obfuscation:
    path: true
    payload: true
  encoding:
    enabled: true
```

---

## Performance Optimization

### Fast Mode (Quick Discovery)

```bash
> config set network.concurrent_requests 50
> config set evasion.mode fast
> config set evasion.multiplier 0.1

[green]✓ CONFIG:[-] Optimized for speed
[cyan]SPEED:[-] ~10x faster (less stealthy)
```

### Stealth Mode (Slow & Careful)

```bash
> config set network.concurrent_requests 2
> config set evasion.mode silent
> config set evasion.multiplier 5.0

[green]✓ CONFIG:[-] Optimized for stealth
[cyan]STEALTH:[-] Very slow but very evasive
```

---

## Advanced Tuning

### Batch Size Optimization

```bash
# For 1M endpoints
> config set performance.batch_size 10000

# For limited memory
> config set performance.batch_size 100
```

### Cache Configuration

```bash
# Enable aggressive caching
> config set performance.cache_enabled true
> config set performance.cache_ttl 7200  # 2 hours

# Disable cache (always fresh)
> config set performance.cache_enabled false
```

### Thread Pool Sizing

```bash
# Formula: workers = (core_count * 2) + 1

# 4-core machine
> config set performance.workers 9

# 16-core machine
> config set performance.workers 33
```

---

## Database Tuning

### SQLite Optimization

```yaml
database:
  engine: sqlite
  options:
    journal_mode: WAL          # Write-Ahead Logging
    synchronous: NORMAL
    cache_size: 10000
    temp_store: MEMORY
```

### Connection Pooling

```yaml
database:
  connection_pool:
    min_connections: 5
    max_connections: 50
    connection_timeout: 30
    idle_timeout: 300
```

---

## Logging Configuration

### Debug Logging

```bash
> config set logging.level debug
> config set logging.format json

[green]✓ CONFIG:[-] Debug logging enabled
[cyan]OUTPUT:[-] Detailed JSON logs to /var/log/vapor.log
```

### Production Logging

```bash
> config set logging.level error
> config set logging.format simple

[green]✓ CONFIG:[-] Production logging enabled
[cyan]OUTPUT:[-] Errors only (minimal overhead)
```

---

## Configuration Management

### View Current Configuration

```bash
> config view
[cyan]ACTIVE CONFIGURATION:[-]
  network.timeout: 30
  evasion.mode: aggressive
  evasion.multiplier: 1.0
  ...
```

### Reset to Defaults

```bash
> config reset
[green]✓ CONFIG:[-] Reset to defaults
```

### Export Configuration

```bash
> config export myconfig.yaml
[green]✓ EXPORTED:[-] Configuration saved

# Later, restore
> config import myconfig.yaml
```

---

## Best Practices

✅ **DO:**
- Start with defaults
- Adjust based on target behavior
- Save configurations by target type
- Test config changes on lab first
- Monitor performance metrics

❌ **DON'T:**
- Set timeout too high (hangs)
- Use too many workers (resource exhaustion)
- Disable all evasion on production tests
- Use debug logging in production
- Make all changes at once

---

**See Also:**
- [10_GHOST_WEAVER.md](10_GHOST_WEAVER.md) - Evasion tuning
- [14_ANALYTICS.md](14_ANALYTICS.md) - Performance monitoring

---

**Last Updated:** February 8, 2026  
**Version:** 1.0 - Production Ready
