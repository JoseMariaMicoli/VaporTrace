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

# 23. Intelligence Layer (OSINT) - Tier 4 Day 1

**Version:** 1.0  
**Tier:** Tier 4 (Advanced)  
**Status:** Sprint 20  
**Feature:** External intelligence gathering and ghost endpoint discovery

---

## Overview

The **Intelligence Layer** (Tier 4 - Day 1) provides reconnaissance capabilities beyond the target domain itself. It leverages external data sources to discover "ghost endpoints" - forgotten, deprecated, or staging versions of APIs that may still be live and exploitable.

### Key Capabilities

- **Wayback Machine Integration** - Retrieve historical URLs and archived versions
- **Shodan Integration** - Discover infrastructure, open ports, and services
- **Configuration Management** - Store and manage API keys securely
- **Automated Data Collection** - Results feed directly into the F2 Map

---

## Commands

### intel wayback

Query the Internet Archive's Wayback Machine for historical URLs of your target.

**Syntax:**
```
intel wayback <domain>
```

**Examples:**
```bash
intel wayback api.example.com
intel wayback beta.api.example.com
intel wayback admin.example.com
```

**Use Cases:**

1. **Legacy API Discovery** - Find old `/v1/`, `/beta/`, or `/staging/` endpoints
2. **Deprecated Endpoints** - Identify endpoints removed from production but still active
3. **Archive Exploration** - Retrieve snapshots from different time periods
4. **Forgotten Configurations** - Discover debug endpoints or test APIs

**Output:**
- Results automatically added to F2 Map
- Logged with timestamp and snapshot date
- Filtered for unique endpoints
- Status codes for each discovered endpoint

**Example Workflow:**
```
target https://api.example.com
intel wayback example.com
# Results show: /api/v1/users (2023-06-15), /admin/debug (2022-11-20)
# All added to F2 Map for testing
```

---

### intel shodan

Query Shodan for infrastructure information about your target.

**Syntax:**
```
intel shodan <ip|domain>
```

**Examples:**
```bash
intel shodan 192.168.1.100
intel shodan example.com
intel shodan api.example.com
```

**What Shodan Returns:**

- Open ports and services
- Web server version/technology stack
- Open buckets (S3, Azure Storage)
- API endpoints exposed
- Certificates and SSL/TLS info
- Geographic location
- ISP and hosting provider

**Use Cases:**

1. **Infrastructure Mapping** - Understand the target's tech stack
2. **Port Enumeration** - Discover non-standard API ports (8080, 8443, 9000, etc.)
3. **Service Identification** - Find database ports, management interfaces
4. **Supply Chain Discovery** - Identify CDN, WAF, and load balancer info
5. **Historical Data** - Compare current vs. past infrastructure

**Example Workflow:**
```
intel config shodan YOUR_SHODAN_API_KEY
intel shodan example.com
# Results show: Port 8443 (Internal API), Port 5432 (PostgreSQL), etc.
```

---

### intel config

Configure API keys and authentication for external services.

**Syntax:**
```
intel config <service> <api_key>
```

**Supported Services:**

```
shodan     - Shodan.io API key
wayback    - Wayback Machine (usually free, no key needed)
```

**Examples:**
```bash
intel config shodan xxxxxxxxxxxxxxxxxxxxxxxx
intel config shodan $(cat ~/.shodan_key)
```

**API Key Management:**

1. **Shodan API Key:**
   - Sign up at https://account.shodan.io/register
   - Get free key (limited queries) or premium (unlimited)
   - API panel: https://account.shodan.io/

2. **Storing Keys Securely:**
   ```bash
   # Option 1: Store in environment variable
   export SHODAN_KEY="your_key_here"
   intel config shodan $SHODAN_KEY

   # Option 2: Store in config file (chmod 600 for security)
   intel config shodan $(cat ~/.vaportrace/shodan.key)
   ```

**Security Note:** API keys are stored in VaporTrace's configuration. Never commit keys to version control.

---

## Tier 4 Day 1 Workflow

### Complete OSINT Reconnaissance Flow

```bash
# 1. Set target
target https://api.example.com

# 2. Configure API keys
intel config shodan YOUR_SHODAN_KEY

# 3. Query Wayback Machine for historical URLs
intel wayback example.com
# ✓ Found 47 unique URLs across snapshots

# 4. Query Shodan for infrastructure
intel shodan example.com
# ✓ Discovered 5 open ports, 3 web services

# 5. Check F2 Map for all discovered endpoints
# UI shows merged results from both sources

# 6. Proceed with exploitation testing
bola https://api.example.com/v1/users/1
ssrf https://api.example.com
```

---

## Integration with Attack Chains

### Using OSINT Findings with Chains

Once you've discovered endpoints via OSINT, integrate them into attack chains:

```bash
# 1. Discover endpoint via Wayback
intel wayback old-api.example.com
# Shows: /api/admin/debug endpoint (2021 snapshot)

# 2. Create chain to test the old endpoint
chain create legacy_test
chain add legacy_test GET https://old-api.example.com/api/admin/debug
chain run legacy_test

# 3. Extract sensitive data if endpoint is still active
extract config json $.config
extract run legacy_test
```

---

## Data Flow

```
┌─────────────────────┐
│  intel wayback      │  → Internet Archive API
│  intel shodan       │  → Shodan.io API
└──────────┬──────────┘
           │
           ▼
    ┌─────────────┐
    │ Aggregation │
    │   Engine    │
    └──────┬──────┘
           │
           ▼
    ┌─────────────────┐
    │ F2 Map (UI)     │
    │ All Endpoints   │
    └─────────────────┘
           │
           ▼
    ┌──────────────────┐
    │  Attack Testing  │
    │ BOLA, SSRF, etc. │
    └──────────────────┘
```

---

## Advanced Scenarios

### Scenario 1: Finding Forgotten Admin Panels

```bash
intel wayback admin.example.com
# Snapshot from 2020 shows /admin/debug endpoint
# Endpoint removed from main app but still in infrastructure
# Test with: admin https://admin.example.com/admin/debug
```

### Scenario 2: Discovering Staging API

```bash
intel shodan example.com
# Shodan reveals: staging.api.example.com port 8443
# Add to target: target https://staging.api.example.com:8443
# Staging typically has weaker security, credentials exposed
```

### Scenario 3: Version-Specific Vulnerabilities

```bash
intel wayback api.example.com
# Discover: /api/v1/ (deprecated in v3)
# v1 may have known vulnerabilities not fixed in v3
# bola https://api.example.com/api/v1/users/1
```

---

## Troubleshooting

### Shodan API Key Issues

**Error:** "Invalid API key"
```
Solution: Verify key via https://account.shodan.io/
```

**Error:** "Query limit exceeded"
```
Solution: Upgrade to premium or wait for quota reset
```

### Wayback Machine Issues

**Error:** "No snapshots found"
```
Solution: Try broader domain or different time period
```

**Error:** "Slow response"
```
Solution: Wayback can be slow; try alternative domains
```

---

## Best Practices

1. **Prioritize Findings** - Tier discoveries by recency and functionality
2. **Cross-Reference** - Compare Wayback and Shodan results for comprehensive map
3. **Test Cautiously** - Legacy endpoints may have stricter logging
4. **Store Configuration** - Use `intel config` to avoid re-entry
5. **Document Sources** - Note where each endpoint was discovered (useful for reporting)

---

## Integration with Other Tier 4 Features

- **Chain Reactor:** Use discovered endpoints with chains for state-aware testing
- **Extractor:** Extract credentials/tokens from legacy endpoint responses
- **Attack Pipeline:** Automatically route legacy endpoints to BOLA/SSRF tests

---

## Next Steps

After OSINT reconnaissance:
1. Review consolidated endpoint list in F2 Map
2. Prioritize endpoints by criticality and recency
3. Create attack chains for multi-step exploitation
4. Use extractors to capture credentials from legacy systems
5. Document findings for final reporting

---

## See Also

- [09_ATTACK_CHAINS.md](09_ATTACK_CHAINS.md) - Build workflows with discovered endpoints
- [24_CHAIN_REACTOR.md](24_CHAIN_REACTOR.md) - Orchestrate complex attacks
- [25_EXTRACTOR.md](25_EXTRACTOR.md) - Extract values from responses
- [22_DISCOVERY_GUIDE.md](22_DISCOVERY_GUIDE.md) - Additional discovery techniques
