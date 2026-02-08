# 12 - Proxy Network Configuration

> **For:** Advanced operators  
> **Read Time:** 15 minutes  
> **Difficulty:** ⭐⭐⭐ Intermediate  
> **Last Updated:** February 8, 2026

---

## Overview

> <span style="background-color: #0891b2; color: white; padding: 2px 6px; border-radius: 3px;">**PROXY**</span> Route traffic through proxy chains for anonymization, IP rotation, and geolocation spoofing.

### Proxy Setup

```bash
# Single proxy
> proxy http://10.0.0.1:8080
[green]✓ PROXY:[-] Configured

# Multiple proxies (rotation)
> proxies load /path/to/proxy-list.txt
[green]✓ PROXIES:[-] Loaded 50 proxies

# SOCKS5
> proxy socks5://proxy.example.com:1080
[green]✓ PROXY:[-] SOCKS5 configured
```

### Proxy List Format

```
http://proxy1.com:8080
https://proxy2.com:8080
socks5://proxy3.com:1080
http://user:pass@proxy4.com:8080
```

### Proxy Chain Configuration

```bash
> proxy chain
  proxy1: http://10.0.0.1:8080
  proxy2: socks5://10.0.0.2:1080
  rotation: sequential

[green]✓ CHAIN:[-] 2-proxy chain configured
[cyan]TRAFFIC:[-] Routes: proxy1 → proxy2 → target
```

### Geolocation Spoofing

```bash
> proxy --location US
[green]✓ PROXY:[-] US proxy selected
[cyan]IP:[-] 203.0.113.45 (New York, USA)

> proxy --location EU
[green]✓ PROXY:[-] EU proxy selected
[cyan]IP:[-] 198.51.100.12 (Frankfurt, Germany)

> proxy --location ASIA
[green]✓ PROXY:[-] ASIA proxy selected
[cyan]IP:[-] 192.0.2.88 (Tokyo, Japan)
```

### IP Rotation Strategy

```bash
# Sequential (predictable)
> rotation sequential

# Random (more stealthy)
> rotation random

# Geographic (appears natural)
> rotation geographic

# Time-based (human-like patterns)
> rotation time-based --interval 5m
```

---

## Best Practices

✅ **DO:**
- Use legitimate proxy services only
- Rotate proxies during long scans
- Test proxy connectivity first
- Mix proxy types for diversity

❌ **DON'T:**
- Use free proxies (often monitored)
- Keep same proxy too long
- Use proxies without rotation
- Assume anonymity is guaranteed

---

## Integration with Stealth

```bash
> stealth silent
[cyan]STEALTH:[-] Maximum evasion enabled

> proxy chain --with-stealth
  - Randomized headers
  - Delayed requests
  - Proxy rotation every 10 req

[green]✓ COMBINED:[-] Stealth + Proxy chain active
```

---

**See Also:**
- [10_GHOST_WEAVER.md](10_GHOST_WEAVER.md) - Evasion techniques
- [18_COMMAND_REFERENCE.md](18_COMMAND_REFERENCE.md) - Proxy commands

---

**Last Updated:** February 8, 2026  
**Version:** 1.0 - Production Ready
