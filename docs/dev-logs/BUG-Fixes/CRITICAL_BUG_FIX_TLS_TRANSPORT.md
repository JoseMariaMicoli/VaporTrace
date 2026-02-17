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

# ⚠️ CRITICAL BUG FIX REPORT

**Date:** February 8, 2026  
**Status:** 🚨 EMERGENCY FIX APPLIED  
**Severity:** CRITICAL

---

## Issue Identified

The TLS browser fingerprint rotation feature was **breaking the entire HTTP/2 negotiation**, causing:

1. ❌ **All discovery functions failing** - audit, mining, swagger couldn't connect
2. ❌ **Sniffer not capturing traffic** - F4 TRAFFIC tab empty
3. ❌ **Malformed HTTP responses** - "HTTP/1.x transport connection broken"
4. ❌ **Generic garbage data** - `\x00\x00\x12\x04` (HTTP/2 frames misinterpreted)

## Root Cause

The `ApplyTLSEvasion()` function was being called during **client initialization** and replacing the entire TLS config. This had two critical problems:

1. **Loss of HTTP/2 negotiation** - The new config didn't preserve proper NextProtos settings
2. **Timing issue** - TLS profiles need per-request selection, not global assignment

## Solution Applied

### Immediate Fix (Emergency)
Reverted `InitializeRotaryClient()` to use standard TLS config:

```go
// CRITICAL: Use standard config with HTTP/2 support
// TLS profiles are applied per-request, NOT during client init
baseTransport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    // ← No TLS profile override (preserves HTTP/2)
}
```

**Result:** ✅ All systems restored to working state

### Proper Implementation (Next Sprint)

The TLS browser fingerprint rotation **must be implemented differently**:

1. **Per-Request Application** - Apply profiles in RoundTrip(), not during init
2. **Proper HTTP/2 Negotiation** - Preserve NextProtos settings
3. **Custom Dialer** - Use net.Dialer with custom TLS handshake
4. **No Transport Replacement** - Work within existing middleware chain

## Current Status

✅ **BUILD:** SUCCESS  
✅ **Sniffer:** CAPTURING TRAFFIC  
✅ **Discovery:** FUNCTIONS WORKING  
✅ **HTTP/2:** NEGOTIATION RESTORED

All systems operational. TLS profile feature **disabled pending proper implementation**.

---

## Lessons Learned

1. **HTTP/2 negotiation requires careful handling** - Can't replace transport mid-flight
2. **Global TLS config too broad** - Need per-request granularity
3. **Test with real targets early** - HTTP/2 issues would have been caught immediately

## Next Steps

- [ ] Implement TLS profiles via custom net.Dialer
- [ ] Apply per-request in RoundTrip()
- [ ] Test with HTTP/2 targets (httpbin.org, etc.)
- [ ] Preserve backward compatibility
- [ ] Add integration tests

---

**Severity Fixed:** CRITICAL → RESOLVED ✅

