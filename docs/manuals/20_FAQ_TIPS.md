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

# 20 - FAQ & Best Practices

> **For:** All users  
> **Read Time:** 15 minutes  
> **Difficulty:** ⭐ Beginner  
> **Last Updated:** February 8, 2026

---

## Frequently Asked Questions

> <span style="background-color: #14b8a6; color: white; padding: 2px 6px; border-radius: 3px;">**FAQ**</span> Common questions about VaporTrace usage, licensing, and capabilities.

---

## General Questions

### Q: What is VaporTrace?

**A:** VaporTrace is an enterprise-grade API penetration testing framework with:
- **Automated discovery** of API endpoints and Swagger specs
- **Advanced exploitation** of OWASP API Top 10 vulnerabilities
- **AI-powered payloads** using Groq LLM
- **Real-time interception** of HTTP requests/responses
- **Professional reporting** with executive summaries
- **Stealth evasion** against WAF/IDS/IPS systems

### Q: Is VaporTrace free?

**A:** VaporTrace is open-source (MIT License) and completely free. Enterprise support available via GitHub Sponsors.

### Q: What languages does VaporTrace support?

**A:** VaporTrace is built in **Go** for performance and compiled into a single binary. It can test APIs in any language.

### Q: Can I use VaporTrace commercially?

**A:** Yes! MIT License allows commercial use. Just include license attribution.

### Q: What's the system requirement?

**A:** 
- **CPU:** 2+ cores recommended
- **RAM:** 512MB minimum, 2GB+ recommended
- **Disk:** 100MB for binary + logs
- **Network:** Internet connection
- **OS:** Linux, macOS, Windows (WSL2)

---

## Installation & Setup

### Q: How do I install VaporTrace?

**A:** 
```bash
# Clone repo
git clone https://github.com/vaportrace/vaportrace.git

# Build
cd vaportrace
go build -o vapor ./cmd

# Run
./vapor scan https://api.example.com
```

Or download pre-compiled binary from releases.

### Q: Do I need Go installed?

**A:** No - download the pre-compiled binary. Go only needed if building from source.

### Q: How do I update VaporTrace?

**A:**
```bash
> update check
> update install
```

Or pull latest and rebuild:
```bash
git pull origin main
go build -o vapor ./cmd
```

### Q: Can I use VaporTrace on Windows?

**A:** Yes, using **WSL2** (Windows Subsystem for Linux 2). Native Windows port in development.

### Q: How do I uninstall?

**A:** Just delete the binary and config directory:
```bash
rm ~/vapor
rm -rf ~/.vapor
```

---

## Usage & Scanning

### Q: How long does a typical scan take?

**A:**
- **Light API:** 2-5 minutes (10-50 endpoints)
- **Medium API:** 15-30 minutes (100-500 endpoints)
- **Heavy API:** 1-4 hours (1000+ endpoints)

Use `--fast` mode for quick scans: `-50% time, -50% coverage`

### Q: How many endpoints can VaporTrace handle?

**A:** Tested up to 10,000+ endpoints. Performance depends on:
- Network speed (most critical)
- Server response time
- Local disk speed
- Configured concurrency

### Q: What's the difference between `scan` and `test`?

**A:**
| Command | Purpose | Speed | Coverage |
|---------|---------|-------|----------|
| **scan** | Full discovery + testing | Slow | 100% |
| **test** | Test specific endpoint | Fast | Manual |
| **quick** | Fast scan + critical tests | 10x faster | ~60% |

### Q: How do I test authenticated APIs?

**A:**
```bash
# Bearer token
> scan https://api.example.com \
    --header "Authorization: Bearer token123"

# API key
> scan https://api.example.com \
    --header "X-API-Key: key123"

# Basic auth
> scan https://api.example.com \
    --auth user:password

# OAuth2 (token endpoint)
> scan https://api.example.com \
    --oauth2 https://auth.example.com/token \
    --client-id CLIENT_ID \
    --client-secret SECRET
```

### Q: Can I scan private/internal APIs?

**A:** Yes! Use VPN or proxy:
```bash
# Via VPN
> scan https://internal-api.local --vpn

# Via proxy
> config set proxy.url http://proxy.example.com:8080
> scan https://internal-api.local
```

### Q: How do I exclude certain endpoints?

**A:**
```bash
# By pattern
> scan https://api.example.com --exclude "/admin,/internal"

# By tag
> scan https://api.example.com --exclude-tag deprecated

# By method
> scan https://api.example.com --exclude-method DELETE
```

---

## Evasion & WAF Bypass

### Q: Why am I getting blocked by WAF?

**A:** Common reasons:
- Too many requests (rate limiting)
- Payload pattern detected
- Behavioral pattern detected
- User-Agent suspicious

**Solutions:**
```bash
# Enable Ghost Weaver evasion
> ghost enable
> config set evasion.mode silent

# Increase delays
> config set evasion.jitter.min 500
> config set evasion.jitter.max 2000

# Rotate User-Agent
> config set http.user_agent "Mozilla/5.0..."

# Use proxy rotation
> config set proxy.rotation random
```

### Q: What's the "Ghost Weaver"?

**A:** Advanced WAF evasion technique using:
- Payload obfuscation
- Protocol-level evasion
- Encoding tricks (UTF-8, double encoding, etc)
- Timing manipulation
- Behavioral obfuscation

### Q: How long does evasion take?

**A:** ~3-10x slower but much more stealthy:
- **No evasion:** 100 requests/sec
- **Standard:** 50 requests/sec
- **Aggressive:** 10 requests/sec
- **Silent:** 2 requests/sec

### Q: Can VaporTrace bypass Cloudflare?

**A:** Partially:
- ✅ Bypass basic rate limiting
- ✅ Bypass simple WAF rules
- ⚠️ Challenge pages may block
- ❌ Cannot bypass bot detection

Use legitimate access with Cloudflare bypass service for best results.

---

## AI & Payload Generation

### Q: Do I need a Groq API key?

**A:** Optional. If not provided:
- Use basic payloads (no AI)
- Limited customization
- Slower exploitation

Get free Groq API key: https://console.groq.com

### Q: How much does Groq cost?

**A:** Free tier:
- 30 requests/minute
- Perfect for most scans
- No credit card needed

Paid plans start at $0.50/1M tokens.

### Q: Can I use OpenAI instead of Groq?

**A:** Yes:
```bash
> config set ai.provider openai
> config set ai.api_key sk-xxxxx
```

Supported providers:
- Groq (free, recommended)
- OpenAI (paid, more powerful)
- Anthropic Claude (paid)
- Local Ollama (free, requires setup)

### Q: What payloads does VaporTrace generate?

**A:** Context-aware payloads for:
- SQL Injection (Union, Error-based, Blind, Time-based)
- NoSQL Injection
- Command Injection
- XSS (Stored, Reflected, DOM)
- XXE
- SSRF
- LDAP Injection
- Template Injection
- Path Traversal
- And more...

### Q: Can I use custom payloads?

**A:** Yes:
```bash
# Custom wordlist
> scan https://api.example.com --payloads custom-payloads.txt

# Custom generator
> scan https://api.example.com --gen custom-generator.py
```

---

## Reporting

### Q: What report formats are available?

**A:**
- **Markdown** - For collaboration
- **PDF** - For printing/delivery
- **HTML** - For sharing
- **JSON** - For parsing
- **CSV** - For spreadsheets

```bash
> report generate --format pdf --template executive
```

### Q: Can I customize the report?

**A:** Yes, multiple templates:
- `executive` - High-level summary (C-suite)
- `technical` - Deep technical detail
- `audit` - Compliance-focused
- `remediation` - Fix-focused
- `minimal` - Quick summary

### Q: How do I share findings with developers?

**A:**
```bash
# Export as JSON
> report export findings.json

# Open in VS Code
code findings.json

# Or HTML dashboard
> report export findings.html
# Share via web link
```

### Q: Can I automate reporting?

**A:** Yes:
```bash
# Generate and send
> scan https://api.example.com --report pdf --email admin@example.com

# Or webhook
> config set reporting.webhook https://slack.com/hook
```

---

## Performance & Optimization

### Q: VaporTrace is too slow. What can I do?

**A:**
```bash
# Fastest mode
> scan https://api.example.com --fast

# Manual tuning
> config set network.concurrent_requests 100
> config set performance.batch_size 5000
> config set evasion.multiplier 0.1
```

**Warning:** Fast mode reduces coverage and evasion.

### Q: How do I monitor performance?

**A:**
```bash
# Real-time dashboard
> dashboard

# Export metrics
> metrics export prometheus

# Check resource usage
> system status
```

### Q: Can I run multiple scans simultaneously?

**A:** Not recommended (resource intensive), but:
```bash
# Background scan
> scan https://api1.example.com &
> scan https://api2.example.com &

# Monitor
> jobs list
```

---

## Security & Privacy

### Q: Is my data safe with VaporTrace?

**A:** All data stored locally:
- Database: `~/.vapor/vapor.db`
- Logs: `~/.vapor/logs/`
- No cloud upload
- No telemetry (unless enabled)

### Q: Can I use VaporTrace in production?

**A:** Not recommended - even with evasion:
- ✅ For authorized testing only
- ✅ With explicit written permission
- ✅ On non-production systems
- ✅ During designated test windows

### Q: Does VaporTrace cause damage?

**A:** No - it's read-only testing:
- Doesn't modify data
- Doesn't delete anything
- Doesn't install backdoors
- Exploitation is safe/reversible

### Q: How do I report security issues?

**A:** Email: security@vaportrace.io  
**No public disclosure** until patch available (90-day window)

---

## Legal & Licensing

### Q: Do I need permission to test an API?

**A:** **YES** - explicit written authorization required:
- OWASP WASC guidelines
- NIST standards
- Legal liability

Always get written permission before testing.

### Q: Is VaporTrace legal to use?

**A:** Yes, if:
- ✅ Testing with explicit permission
- ✅ Used on authorized systems
- ✅ For security improvement

Illegal if:
- ❌ Testing without permission (unauthorized access)
- ❌ Testing competitor systems
- ❌ Malicious intent

### Q: What's the license?

**A:** MIT License:
- Free for commercial use
- Modify and redistribute
- Include license attribution
- No warranty/liability

---

## Troubleshooting

### Q: VaporTrace keeps crashing

**A:** Try in order:
1. `> daemon restart`
2. `> config reset`
3. Update: `git pull && go build`
4. Check logs: `> logs tail -n 100`
5. Report issue with logs

### Q: Can't connect to target

**A:**
```bash
# Test connectivity
> network diagnose https://target.example.com

# Try without SSL verify
> scan https://target.example.com --no-verify

# Try with different DNS
> config set network.dns_server 8.8.8.8
```

### Q: No vulnerabilities found

**A:** Try:
1. Lower sensitivity: `--sensitivity high`
2. Enable AI: Set Groq API key
3. Check endpoints: `--list-endpoints`
4. Manual testing: `> test /endpoint --method POST`

### Q: Memory keeps growing

**A:**
```bash
> config set performance.cache_enabled false
> daemon restart
```

---

## Best Practices

> <span style="background-color: #06b6d4; color: white; padding: 2px 6px; border-radius: 3px;">**TIP**</span>

### ✅ DO:

- **Get authorization** first (written consent)
- **Test in lab** before production
- **Enable evasion** on real targets
- **Monitor logs** during scans
- **Use proxies** for anonymity
- **Document findings** thoroughly
- **Follow responsible disclosure**
- **Update regularly** for security

### ❌ DON'T:

- **Test without permission** (criminal)
- **Use default configs** on real targets
- **Disable SSL verification** permanently
- **Ignore rate limits** (causes denial)
- **Modify data** during tests
- **Use on competitor systems** (illegal)
- **Share credentials** in reports
- **Leave scans running** unattended

---

## Getting Help

### Documentation
- 📖 **Full Manuals:** `/docs/manuals/`
- 🔍 **API Reference:** [19_API_MODULES.md](19_API_MODULES.md)
- 🐛 **Troubleshooting:** [16_TROUBLESHOOTING.md](16_TROUBLESHOOTING.md)
- ⚙️ **Configuration:** [15_CONFIGURATION.md](15_CONFIGURATION.md)

### Community
- 💬 **GitHub Discussions:** https://github.com/vaportrace/vaportrace/discussions
- 🐛 **Report Bug:** https://github.com/vaportrace/vaportrace/issues
- 📧 **Email Support:** support@vaportrace.io
- 🤝 **Contributing:** Contributing guidelines in README

### Premium Support
- 🚀 **GitHub Sponsor:** https://github.com/sponsors/vaportrace
- 🏢 **Enterprise:** enterprise@vaportrace.io
- 📞 **Priority Support:** Available to sponsors

---

## Quick Reference

### Common Commands

```bash
# Start scan
> scan https://api.example.com

# With options
> scan https://api.example.com --fast --no-verify

# List endpoints
> scan https://api.example.com --list

# Test specific endpoint
> test /api/v1/users --method GET

# Generate report
> report generate --format pdf

# View dashboard
> dashboard

# Show config
> config view
```

### Common Issues

| Issue | Solution |
|-------|----------|
| Timeout | `config set network.timeout 60` |
| SSL Error | `scan --no-verify` |
| Blocked | `config set evasion.mode silent` |
| Slow | `scan --fast` |
| Memory | `config set performance.cache_enabled false` |

---

**Last Updated:** February 8, 2026  
**Version:** 1.0 - Current Release

**Need Help?** → [support@vaportrace.io](mailto:support@vaportrace.io)
