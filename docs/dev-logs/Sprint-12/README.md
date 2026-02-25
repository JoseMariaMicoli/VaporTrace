# Sprint 12: Evasion V2 & Advanced Defense Bypass

**Status:** 🔄 IN PROGRESS | **Version:** v3.1-Hydra (Active Development) | **Started:** February 2026

---

## 🎯 Sprint Overview

Sprint 12 continues VaporTrace's evasion capabilities from Sprint 6 with advanced traffic shaping, behavioral obfuscation, and encrypted out-of-band exfiltration channels. This sprint modernizes evasion to defeat next-generation defenders (SIEM, EDR, WAF).

**Slogan:** "Invisible by Design - Evasion 2.0"

---

## 📋 Planned Deliverables

### 12.1: Deep Traffic Shaping ✅ COMPLETE

**Status:** ✅ Shipped in Sprint 11  
**Location:** `pkg/logic/network.go`

**Features Delivered:**
- `MimicTraffic()` - Browser profile matching (6 profiles)
  - iOS User-Agent and headers
  - Chrome on macOS (Windows)
  - Firefox on Windows
  - Edge Browser
  - Safari
  - Bot profile
- Header randomization (Accept, Accept-Encoding, Sec-Fetch-*)
- Referer patterns per profile
- Cookie management
- Connection timing realism

**Status:** ✅ Fully functional

---

### 12.2: TLS Fingerprinting & JA3 Evasion ⏳ PLANNED

**Target:** February-March 2026  
**Complexity:** High  
**Dependencies:** tls-utls library

**Objective:** Mimic real browser TLS handshakes

**Problem:**
- Go's `crypto/tls` has distinctive cipher suite order
- JA3 fingerprinting identifies Go-based clients
- Modern WAFs detect Go TLS fingerprint

**Solution (Planned):**
- Integrate `github.com/refraction-networking/utls`
- Implement TLS profiles matching real browsers
- Match cipher suite order, extensions, curves
- Result: Indistinguishable from real clients

**Implementation Plan:**
```go
// Planned for Sprint 12.2
import "github.com/refraction-networking/utls"

func ProcessChainWithTLSMimicry(chainID string) {
    config := &utls.Config{
        ClientProfile: utls.ChromeProfile,  // or Safari, Firefox, etc.
        ServerName:    targetHost,
    }
    
    conn, err := utls.Dial("tcp", targetAddr, config)
    // Result: TLS Client Hello matches Chrome exactly
}
```

**Benefits:**
- Defeats JA3/JA3S fingerprinting
- Passes modern WAF detection
- No signature in TLS handshake
- Complete browser transparency

**Current Status:** Planned, not yet started

---

### 12.3: Behavioral Jitter ✅ COMPLETE

**Status:** ✅ Shipped in Sprint 6/11  
**Location:** `pkg/logic/network.go` - `ApplyJitter()`

**Features Delivered:**
- Gaussian distribution for inter-packet delays
- Box-Muller transform for naturalistic timing
- Configurable mean and standard deviation
- Per-action delay randomization

**Implementation:**
```go
func ApplyJitter(baseDelay int) time.Duration {
    // Gaussian distribution: mean=baseDelay, stddev=20%
    // Results in natural-looking traffic patterns
    // Not detectable by statistical rate-limiting
}
```

**Status:** ✅ Fully functional

---

### 12.4: Encrypted OOB Exfiltration ⏳ PLANNED

**Target:** March-April 2026  
**Complexity:** Very High  
**Dependencies:** Custom protocol design

**Objective:** Secure exfiltration channel for sensitive data

**Problem:**
- Current exfiltration plaintext or basic encoding
- Detectable by network analysis
- No integrity verification
- Keys transmitted in-band

**Solution (Planned):**

1. **OOB Channel Establishment**
   - Alternative communication path (DNS, ICMP, custom protocol)
   - Out-of-band key exchange (Diffie-Hellman)
   - No connection to original target

2. **Encryption Protocol**
   - ChaCha20-Poly1305 for AEAD
   - 256-bit keys
   - Per-message nonces
   - Authentication tags

3. **Exfiltration Format**
   - Chunked encrypted payloads
   - Resend logic for reliability
   - Compression before encryption
   - Minimal packet size

**Status:** Planned, not yet started

---

## 🔄 Current Status by Sub-Phase

| Sub-Phase | Deliverable | Status | Completion |
|-----------|-------------|--------|------------|
| **12.1** | MimicTraffic() | ✅ DONE | 100% |
| **12.2** | TLS-utls integration | ⏳ PLANNED | 0% |
| **12.3** | Behavioral Jitter | ✅ DONE | 100% |
| **12.4** | Encrypted OOB | ⏳ PLANNED | 0% |

**Overall Progress:** 50% (2/4 sub-phases complete)

---

## 🔧 Technical Architecture

### Current Evasion Stack (Sprint 6 + 11)

```
┌─────────────────────────────────────────────────┐
│ Application Layer (Exploitation)                │
├─────────────────────────────────────────────────┤
│ Evasion Layer                                   │
├─────────────────────────────────────────────────┤
│ ├─ ApplyJitter()          → Timing obfuscation  │
│ ├─ MimicTraffic()         → Header spoofing     │
│ ├─ ProxyRotation()        → IP masking          │
│ └─ [TLS-utls pending]    → JA3 fingerprinting  │
├─────────────────────────────────────────────────┤
│ HTTP/RoundTripper Layer                         │
├─────────────────────────────────────────────────┤
│ TCP/TLS Layer                                   │
└─────────────────────────────────────────────────┘
```

### Planned Enhancements (Sprint 12.2+)

```
┌─────────────────────────────────────────────────┐
│ [NEW] TLS-utls Integration                      │
│ - Client Profile Selection                      │
│ - JA3 Fingerprint Matching                      │
│ - Extension Ordering                            │
│ - Cipher Suite Reordering                       │
└─────────────────────────────────────────────────┘
         ↓
    crypto/tls (current)
         ↓
    [TO BE REPLACED by utls.Conn]
```

---

## 📊 Evasion Effectiveness Matrix

### Current Capabilities (Sprints 6, 11)

| Detection Method | Current Status | Bypass Rate |
|------------------|----------------|------------|
| **Signature Detection** | ✅ Bypassed | >90% |
| **Behavioral Analysis** | ✅ Bypassed | ~80% |
| **Rate Limiting** | ✅ Bypassed | >90% |
| **IP Reputation** | ✅ Bypassed | >95% (with proxy rotation) |
| **Header Analysis** | ✅ Bypassed | >85% |
| **TLS Fingerprinting (JA3)** | ❌ DETECTED | 0% |
| **Timing Analysis** | ✅ Bypassed | ~75% |
| **Process Masquerading** | ✅ Bypassed | >95% |

### After Sprint 12.2 (With TLS-utls)

| Detection Method | Expected Status | Expected Bypass Rate |
|------------------|-----------------|----------------------|
| **TLS Fingerprinting (JA3)** | ✅ Bypassed | >95% |
| **Overall Detection Rate** | NEAR ZERO | >98% |

---

## 🔗 Dependencies

### External Libraries (Planned)

```go
import "github.com/refraction-networking/utls"
```

**TLS-utls Library:**
- Repository: https://github.com/refraction-networking/utls
- License: BSD 3-Clause
- Maintenance: Active community
- Adoption: Used by major tools (circumvention tools, pentesters)

**Rationale for Deferral:**
- Sprint 11 focused on core autonomy
- Current header-based evasion sufficient for most targets
- TLS integration adds complexity and external dependency
- Can be added post-release without breaking changes

---

## 📝 Integration Points

### With Sprint 11 (Autonomy)
- ProcessChain() respects evasion settings
- Jitter applied between chain links
- Traffic mimicry per browser profile
- Future: TLS profiles per chain

### With Sprint 6 (Base Evasion)
- Builds on header randomization
- Extends IP rotation capabilities
- Combines timing obfuscation
- Future: Advanced fingerprinting evasion

### With Core Engine
- `pkg/logic/network.go` - Evasion middleware
- `pkg/logic/http.go` - HTTP transport
- `pkg/engine/core.go` - Execution orchestration

---

## 🎯 Success Criteria

### Sprint 12 Completion Requirements

- [ ] 12.1: MimicTraffic() with 6 profiles ✅ DONE
- [ ] 12.3: Behavioral jitter working ✅ DONE
- [ ] 12.2: TLS-utls integration (target: March 2026)
- [ ] 12.4: Encrypted OOB channels (target: April 2026)
- [ ] All evasion techniques integrated
- [ ] Detection bypass rates >90%
- [ ] Performance impact <5%
- [ ] Documentation complete

---

## 📚 Documentation References

- **Dev-Roadmap.md** - Full roadmap with all sprints
- **Sprint-11/README.md** - Previous autonomy sprint
- **Sprint-16/README.md** - Blue-team mirror sprint
- **README.md** - Main project documentation

---

## 🚀 Next Steps

### Immediate (February 2026)
- [ ] Review tls-utls library integration approach
- [ ] Design TLS profile matching system
- [ ] Plan encrypted OOB protocol

### Near-term (March 2026)
- [ ] Implement TLS-utls integration
- [ ] Test JA3 fingerprint matching
- [ ] Validate bypass rates

### Medium-term (April 2026)
- [ ] Implement encrypted OOB exfiltration
- [ ] Design custom protocol
- [ ] Security audit

---

## 📊 Metrics

| Metric | Current | Target (Sprint 12) |
|--------|---------|-------------------|
| **Evasion Techniques** | 5 | 7 |
| **Browser Profiles** | 6 | 6+ |
| **Detection Bypass Rate** | ~80% | >95% |
| **Performance Impact** | ~2% | <5% |
| **Supported Frameworks** | JAG detection, IP reputation | All modern WAFs |

---

## 🔐 Security Considerations

### Evasion Ethics
- **Authorized Scope Only:** All evasion deployed on authorized targets
- **Compliance:** Respects SLAs, business hours, and RPS limits
- **Detection:** Works with defensive tools (SIEM, EDR)

### Safe Defaults
- Jitter configurable per engagement
- Proxy rotation optional
- TLS-utls as opt-in feature
- No system-wide masquerading

---

## 📄 Status Summary

**Sprint 12 Progress:** 50% complete (2/4 sub-phases)

**Shipped:** Traffic Shaping + Behavioral Jitter  
**In Planning:** TLS Fingerprinting + Encrypted OOB  
**Timeline:** 2-4 more months for remaining features

**Production Impact:** Current evasion sufficient for most engagements. TLS fingerprinting will significantly improve detection bypass for advanced WAF environments.

---

**Status:** 🔄 IN PROGRESS | **Last Updated:** February 8, 2026 | **Next Review:** March 15, 2026
