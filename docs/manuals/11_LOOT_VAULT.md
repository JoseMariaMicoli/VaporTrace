# 11 - Loot Vault: Secret & Credential Management

> **For:** All operators  
> **Read Time:** 15 minutes  
> **Difficulty:** ⭐⭐ Beginner  
> **Last Updated:** February 8, 2026

---

## Overview

> <span style="background-color: #059669; color: white; padding: 2px 6px; border-radius: 3px;">**VAULT**</span> **Loot Vault** automatically captures, classifies, and manages all secrets discovered during testing.

### Loot Categories

```
🔐 CREDENTIALS
  ├── API Keys (AWS, Azure, Groq)
  ├── Bearer Tokens (JWT, OAuth)
  ├── Basic Auth (username:password)
  └── Session Tokens

💾 DATA
  ├── User Emails (PII)
  ├── Financial Data (SSN, CC)
  └── Internal Docs

🔑 ENCRYPTION
  ├── Private Keys (RSA, EC)
  ├── Certificates (SSL, client)
  └── Secrets (API secrets)

📋 METADATA
  ├── Internal IPs
  ├── Hostnames
  └── Service Info
```

---

## Commands

### View Loot

```bash
> loot
[cyan]VAULT:[-] Total items: 247
[green]✓ API Keys: 12
[green]✓ Tokens: 45
[green]✓ Emails: 890
[green]✓ IPs: 23
```

### Export Loot

```bash
> loot export --format json --output loot.json
[green]✓ EXPORTED:[-] 247 items to loot.json

> loot export --format csv --output loot.csv
[green]✓ EXPORTED:[-] 247 items to loot.csv
```

### Filter Loot

```bash
> loot filter --type api_key
[cyan]VAULT:[-] API Keys: 12

> loot filter --type email
[cyan]VAULT:[-] Emails: 890
```

---

## Best Practices

✅ **DO:**
- Review before reporting
- Mask sensitive data
- Export for storage
- Categorize findings

❌ **DON'T:**
- Share unfiltered dumps
- Leave credentials in memory
- Report PII unnecessarily

---

**Last Updated:** February 8, 2026  
**Version:** 1.0 - Production Ready
