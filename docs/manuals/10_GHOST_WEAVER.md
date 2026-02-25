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

# 10 - Ghost Weaver: Evasion & Detection Bypass

> **For:** Advanced operators  
> **Read Time:** 20 minutes  
> **Difficulty:** ⭐⭐⭐⭐⭐ Expert  
> **Last Updated:** February 8, 2026

---

## Overview

> <span style="background-color: #6b21a8; color: white; padding: 2px 6px; border-radius: 3px;">**STEALTH**</span> **Ghost Weaver** implements advanced evasion techniques to bypass WAF/IDS/IPS detection systems while maintaining attack effectiveness.

### Evasion Pipeline

```mermaid
graph LR
    A["💬 Attack Payload"] -->|Obfuscate| B["🔀 Mutation"]
    B -->|Encode| C["🔐 Encoding"]
    C -->|Timing| D["⏱️ Delays"]
    D -->|Headers| E["📝 Randomize"]
    E -->|Route| F["🛣️ Proxy Chain"]
    F -->|Signature| G["🎭 Fingerprint<br/>Spoof"]
    G -->|Send| H["🎯 Target"]
    
    style A fill:#dc2626,color:#fff
    style H fill:#10b981,color:#fff
```

---

## Core Evasion Techniques

### Technique 1: Payload Obfuscation

```bash
# Original (easily detected)
SELECT * FROM users WHERE id=1' OR '1'='1

# Ghost Weaver obfuscated variants:
SELECT /**/\* FROM users WHERE id=1' OR '1'='1
SELE/**/CT * FROM users WHERE id=1' OR '1'='1
SELECT * FROM /*!50000users*/ WHERE id=1' OR '1'='1
```

### Technique 2: HTTP/2 Pseudo-Header Randomization

```
Standard HTTP/1.1 signature:
  :method: POST
  :path: /api/search
  
Ghost Weaver variation:
  :method: posT       (case variation)
  :path: /api/../search  (path normalization)
```

---

## Configuration

```bash
> stealth silent
[cyan]STEALTH:[-] Maximum evasion, 2x slower

> multiplier 1.5
[cyan]EVASION:[-] All delays scaled 1.5x
```

---

**Last Updated:** February 8, 2026  
**Version:** 1.0
