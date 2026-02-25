![VaporTrace Logo](../../../assets/images/VaporTrace_Logo.png)

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

# Quick Troubleshooting - No Findings Issue

**Quick Test Before Using Burp/Wireshark**

---

## 🔧 STEP 1: Test Basic Connectivity

```bash
# In VaporTrace terminal:

# First, test HTTP (no TLS)
VAPOR/INT> target http://httpbin.org
VAPOR/INT> map

# Wait 10 seconds...

# Check results:
# - Did MAP tab populate?
# - Did you see endpoints?
```

### Results Interpretation:

**✅ If HTTP works (findings appear):**
- Network: ✅ OK
- VaporTrace: ✅ OK  
- Issue: TLS/HTTPS only

→ Go to STEP 2

**❌ If HTTP also doesn't work:**
- Network: ❌ PROBLEM
- Issue: Not TLS-related

→ Skip to STEP 3 (Network Check)

---

## 🔧 STEP 2: Test HTTPS with Burp

```bash
# Terminal 1: Start Burp on port 8080
# (burpsuite should be running)

# Terminal 2: VaporTrace
VAPOR/INT> set-proxy http://127.0.0.1:8080
VAPOR/INT> keys | grep PROXY
# Should show: PROXY (STAT): http://127.0.0.1:8080

# Terminal 2: Run test
VAPOR/INT> target https://httpbin.org
VAPOR/INT> map

# Wait 10 seconds...
```

### Results Interpretation:

**✅ Requests appear in Burp:**
- Proxy: ✅ OK
- TLS Connection: ✅ OK
- Issue: Response handling or findings logic

→ Check VaporTrace logs for errors

**❌ No requests in Burp:**
- Proxy: ❌ NOT SET
- Issue: Either proxy not working or TLS failing

→ Disable TLS temporarily to test

---

## 🔧 STEP 3: Bypass TLS to Isolate Issue

**Create a simple wrapper to disable uTLS:**

Edit `pkg/logic/network.go` - Find `InitializeRotaryClient()`:

```go
// Find this section:
DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
    // Try uTLS first
    conn, err := tlsTransport.DialTLSContext(ctx, network, addr)
    ...
}

// Temporarily replace with STANDARD TLS ONLY:
DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
    host, port, _ := net.SplitHostPort(addr)
    
    dialer := &net.Dialer{Timeout: 30 * time.Second}
    tcpConn, err := dialer.DialContext(ctx, "tcp", addr)
    if err != nil {
        return nil, err
    }
    
    tlsConn := tls.Client(tcpConn, &tls.Config{
        ServerName:         host,
        InsecureSkipVerify: true,
    })
    
    tlsConn.Handshake()
    return tlsConn, nil
},
```

**Then rebuild and test:**
```bash
go build ./...
# Run VaporTrace again
VAPOR/INT> target https://httpbin.org
VAPOR/INT> map
```

### Results:

**✅ Works with standard TLS:**
- Issue: uTLS implementation
- → uTLS library issue, not core logic

**❌ Still fails:**
- Issue: Not uTLS
- → Core response handling or findings logic

---

## 🔧 STEP 4: Check VaporTrace Logs

While running map/scan, look for these messages:

**Good signs (request being sent):**
```
[cyan]TLS:[-] Dialing...
[cyan]TLS:[-] Connected with...
[blue]MAP:[-] Discovered endpoint...
```

**Bad signs (something's wrong):**
```
[red]ERROR:[-] ...
[yellow]WARNING:[-] ...
PANIC: ...
```

**If you see errors:**
- Copy the exact error message
- This tells us what's failing

---

## 🔧 STEP 5: Check Response Handling

In VaporTrace logs, look for:

```bash
# Good:
[green]✓ MAP:[-] Found 15 endpoints
[yellow]LOOT:[-] Captured credential...
[cyan]TRAFFIC:[-] Logged request...

# Bad:
[blue]i[-] No endpoints found
[blue]i[-] Empty response body
[yellow]WARNING:[-] Invalid JSON response
```

If seeing "no endpoints" with no errors:
- Response received but not parsed correctly
- Issue: JSON parsing or endpoint extraction logic

---

## 📋 DECISION TREE

```
NO FINDINGS
├─ Test HTTP
│  ├─ HTTP works → Issue is HTTPS/TLS
│  │  ├─ Test with Burp
│  │  │  ├─ Requests in Burp → uTLS issue
│  │  │  └─ No requests → Proxy/network issue
│  │  └─ Test standard TLS only
│  │     ├─ Works → uTLS problem
│  │     └─ Fails → Response handling problem
│  │
│  └─ HTTP fails → Issue is not TLS
│     ├─ Check logs for errors
│     ├─ Check network connectivity
│     └─ Check response parsing logic
```

---

## 🎯 WHAT EACH TEST TELLS YOU

| Test | Result | Means |
|------|--------|-------|
| HTTP works | ✅ | Network OK, issue is HTTPS |
| HTTP fails | ❌ | Network problem OR response parsing |
| HTTPS + Burp: requests seen | ✅ | Connection OK, findings issue |
| HTTPS + Burp: no requests | ❌ | TLS/proxy/network issue |
| Standard TLS works | ✅ | uTLS is the problem |
| Standard TLS fails | ❌ | Not uTLS, core issue |

---

## 🚨 COMMON ISSUES & QUICK FIXES

### Issue: "PROXY not set" in keys

**Fix:**
```bash
VAPOR/INT> set-proxy http://127.0.0.1:8080
VAPOR/INT> keys  # Verify
```

### Issue: Burp shows "HTTP 403 Forbidden"

**Means:** Server is blocking something (User-Agent? Headers?)

**Fix:**
```bash
# Check request headers in Burp
# Look for: User-Agent, Accept, Authorization
# May need to adjust in evasion.go
```

### Issue: Wireshark shows "TLS Alert: Handshake Failure"

**Means:** Certificate or cipher mismatch

**Fix:**
```bash
# Try standard TLS first
# If standard TLS works → uTLS issue
# If standard TLS fails → Cert issue
```

### Issue: "No endpoints found" but no errors

**Means:** Response received but not parsed

**Fix:**
```bash
# Check response format in Burp
# Is it valid JSON?
# Is it the expected structure?
# Check endpoint extraction regex
```

---

## 📊 REAL EXAMPLE: WORKING vs NOT WORKING

### ❌ NOT WORKING (What you're seeing now):

```
VaporTrace logs:
[cyan]TLS:[-] Dialing api.example.com:443...
[blue]i[-] No endpoints found
(Nothing in MAP tab)

Burp:
(No requests)

Analysis: 
→ Requests not being sent
→ TLS connection failing silently
→ Fallback to standard TLS not happening?
```

### ✅ WORKING (What we want):

```
VaporTrace logs:
[cyan]TLS:[-] Dialing api.example.com:443 with profile chrome-windows
[cyan]TLS:[-] Connected with chrome-windows profile
[green]✓ MAP:[-] Found 25 endpoints
(MAP tab populated with endpoints)

Burp:
POST /api/users - 200 OK - {"users": [...]}
POST /api/products - 200 OK - {"products": [...]}
... (multiple requests)

Analysis:
→ Requests being sent
→ Responses received
→ Endpoints extracted and displayed
```

---

## 🚀 ACTION PLAN

1. **Run HTTP test** (httpbin.org) - takes 30 seconds
   - If works → Go to Burp test
   - If fails → Check network/firewall

2. **If HTTP works, run HTTPS with Burp** - takes 1 minute
   - If requests show → uTLS issue
   - If no requests → Proxy/network issue

3. **If HTTPS fails, test standard TLS** - takes 2 minutes
   - If works → uTLS definitely the problem
   - If fails → Something else is wrong

4. **Report findings** with test results
   - I'll know exactly where the problem is
   - Can provide targeted fix

---

**Start with: Test HTTP first**

Then report back what you see!
