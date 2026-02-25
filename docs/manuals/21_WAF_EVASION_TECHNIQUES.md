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

# WAF EVASION TECHNIQUES - QUICK REFERENCE

## Overview
VaporTrace now includes 5 coordinated WAF evasion techniques that work together to bypass modern WAF detection.

---

## 1. HTTP/2 Pseudo-Header Randomization

**Module**: `pkg/logic/http2_evasion.go`

### What It Does
- Rotates HTTP/2 ClientHello fingerprints to mimic different browsers
- Changes pseudo-header ordering based on User-Agent
- Adds Sec-Fetch-* headers for realistic browser behavior

### Usage
```go
profile := GetHTTP2Profile(userAgent)
ApplyHTTP2Evasion(req, profile)
```

### Profiles Available
- Chrome/Windows (standard order)
- Firefox/Windows (reordered)
- Safari/macOS (aggressive)
- Brave/Linux (standard)

### Bypass Target
- Cloudflare ClientHello fingerprinting
- Akamai bot detection
- Custom HTTP/2 analyzers

---

## 2. Path & Parameter Obfuscation

**Module**: `pkg/logic/path_obfuscation.go`

### What It Does
- Appends noise parameters that servers ignore
- Inserts path parameters (`;v=1.0;x=y`)
- Uses cache-buster patterns
- Double-encodes path segments

### Usage
```go
obfuscated := ObfuscatePath("/api/v1/users", CacheBusters)
// Result: /api/v1/users?_debug=0&_ref=home&_t=1234567890
```

### Techniques
| Technique | Example |
|---|---|
| Cache-Busters | `?_debug=0&_t=123` |
| Path Parameters | `;v=1.0;x=y` |
| Double Encoding | `%75sers` (u encoded) |
| Fragments | `#section1` |

### Bypass Target
- Regex-based WAF rules
- ModSecurity path matching
- Signature-based detection

---

## 3. Contextual "Thinking Time"

**Module**: `pkg/logic/thinking_time.go`

### What It Does
- Applies request-type-specific delays
- Simulates human "thinking" between actions
- Detects and warns about bot patterns

### Usage
```go
ApplyContextualBehavior(req)  // Auto-detects method and applies delay
```

### Delay Profiles
| Request Type | Delay | Example |
|---|---|---|
| GET (Discovery) | 10-50ms | Fetching list of users |
| HEAD/OPTIONS | 50-300ms | Checking API capabilities |
| POST/PUT/DELETE | 800-3000ms | Filling form, then submitting |
| Static/Pixel | 0-5ms | Loading image |

### Bypass Target
- DataDome bot scoring
- Imperva behavioral analysis
- ML-based detection (Cloudflare Bot Management)

---

## 4. Intelligent Rate-Limit (429) Backoff

**Module**: `pkg/logic/rate_limit_backoff.go`

### What It Does
- Automatically detects 429 (Too Many Requests) responses
- Applies exponential backoff with random jitter
- Rotates proxy and User-Agent after cooldown
- Prevents "hammering" patterns

### Usage
```go
// In SafeDo or your request handler:
if resp.StatusCode == 429 {
    delay := HandleRateLimit(resp.StatusCode, resp.Header)
    time.Sleep(delay)
}

// Check if backoff is active:
if IsBackoffActive() {
    waitTime := GetBackoffWaitTime()
    // Skip requests until waitTime expires
}
```

### Backoff Progression
```
1st 429: Wait 30-60 seconds (2^0 + jitter)
2nd 429: Wait 60-120 seconds (2^1 + jitter)
3rd 429: Wait 120 seconds max (2^2 capped)
```

### Post-Backoff Action
- Proxy switched to random from pool
- User-Agent rotated on next request
- Backoff resets on 200/201 response

### Bypass Target
- Rate-limit triggers
- Automated response patterns
- Continuous hammering detection

---

## 5. Payload Encoding & Case Randomization

**Module**: `pkg/logic/payload_encoding.go`

### What It Does
- Applies random encoding to POST bodies (gzip, deflate, identity)
- Randomizes JSON whitespace and formatting
- Varies Content-Encoding header
- Defeats signature-based detection

### Usage
```go
technique := SelectRandomEncoding()
encodedBody, contentEncoding := TransformPayload(jsonPayload, technique)

req.Header.Set("Content-Encoding", contentEncoding)
req.Body = io.NopCloser(bytes.NewReader(encodedBody))
```

### Encoding Strategies
| Strategy | Payload | Encoding Header |
|---|---|---|
| Identity | `{"id":1}` | `identity` |
| Gzip | `[compressed]` | `gzip` |
| Deflate | `[compressed]` | `deflate` |
| Whitespace | `{ "id" : 1 }` | `identity` |

### Bypass Target
- Signature-based payload detection
- ModSecurity payload rules
- Custom WAF regex patterns

---

## Integration Checklist

When integrating these into your exploit engines:

- [ ] Call `ApplyHTTP2Evasion()` before every request
- [ ] Call `ApplyContextualBehavior()` before `SafeDo()`
- [ ] Handle 429 with `HandleRateLimit()` in SafeDo response
- [ ] Apply `TransformPayload()` to all POST/PUT bodies
- [ ] Optionally call `ObfuscatePath()` to modify request paths

---

## Configuration & Tuning

### Disable Specific Techniques
In the respective modules, set to `false` or return early:

```go
// Disable HTTP/2 evasion:
// if req.URL.Scheme == "http" { return } // Skip for HTTP-only targets

// Disable path obfuscation for specific endpoints:
if strings.Contains(req.URL.Path, "/admin") { obfuscated = path }

// Disable thinking time:
// return 0 in ContextualThinkingTime()

// Disable backoff:
// Remove HandleRateLimit() call from SafeDo()

// Disable payload encoding:
// return payload, "identity" in TransformPayload()
```

### Performance Impact
- **HTTP/2 Evasion**: ~0ms (headers only)
- **Path Obfuscation**: ~1ms (string manipulation)
- **Thinking Time**: 10ms - 3000ms (intentional delay)
- **Rate-Limit Handling**: ~0ms (unless backoff active)
- **Payload Encoding**: 1-5ms (gzip faster on large payloads)

**Total per request**: 11ms - 3005ms depending on request type

---

## Testing & Validation

### Test HTTP/2 Evasion
```bash
# Verify Sec-Fetch headers are present
curl -H "User-Agent: Mozilla/5.0..." https://target.com/api
```

### Test Path Obfuscation
```bash
# Verify query parameters are added
GET /api/v1/users?_debug=0&_t=123 HTTP/1.1
```

### Test Rate-Limit Handling
```bash
# Generate 429 response and verify backoff triggers
# Check logs for "RATE_LIMIT: Backoff #1"
```

### Test Payload Encoding
```bash
# Verify Content-Encoding header changes
POST /api/v1/exploit HTTP/1.1
Content-Encoding: gzip
[compressed body]
```

---

## Estimated Effectiveness

Against different WAF solutions:

| WAF | Before | After | Bypass Rate |
|---|---|---|---|
| ModSecurity (basic) | 50% | 70% | Good |
| Cloudflare WAF | 10% | 20% | Partial |
| Akamai | 5% | 15% | Weak |
| DataDome | <5% | 10% | Experimental |

---

## Known Issues & Limitations

1. **HTTP/2 Key Reordering**: Go's stdlib doesn't expose pseudo-header control. Current implementation uses header hints only.

2. **JSON Key Ordering**: Not implemented (requires custom marshaling). Currently only whitespace variation.

3. **TLS Fingerprinting**: Still vulnerable to TLS ClientHello analysis. Consider external proxy for sensitive targets.

4. **Session Anomalies**: Cookie rotation per proxy may trigger session flags.

---

## Support & Troubleshooting

### Technique Not Applied?
- Check module is imported in your file
- Verify function called before `SafeDo()`
- Check log output for `[cyan]` colored messages

### Too Many Delays?
- Contextual thinking time adds 0-3s per POST
- For rapid scanning, disable in loop
- Use `PixelRequest` context for static resources

### Rate-Limit Triggered?
- Normal behavior - cooldown should resume automatically
- Check `IsBackoffActive()` in your loop
- Verify proxy pool is loaded (`LoadProxiesFromFile()`)

---

**Last Updated**: February 8, 2026  
**Sprint**: 17 WAF Evasion V2
