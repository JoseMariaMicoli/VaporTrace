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

# Network Diagnostic Guide - Burp Suite & Wireshark

**Purpose:** Capture and analyze VaporTrace traffic to diagnose TLS/connection issues

---

## 📊 QUICK DIAGNOSIS FLOWCHART

```
Start VaporTrace → Run Map/Scan → Check:
├─ Are requests being sent?        (Burp/Wireshark)
├─ What TLS version/cipher?        (Wireshark)
├─ Are responses coming back?      (Burp)
├─ Are responses valid?            (Burp)
└─ Why aren't findings populating? (VaporTrace logs)
```

---

## 🔧 METHOD 1: BURP SUITE (Easiest - See Requests/Responses)

### Setup Burp

**Step 1: Start Burp Suite**
```bash
# On your machine (not in VaporTrace container)
burpsuite &  # or use your installed copy
```

**Step 2: Configure Listener**
- Go to: `Proxy` → `Settings` → `Proxy Listeners`
- Click `Add`
- Set to:
  - Bind to address: `127.0.0.1` 
  - Port: `8080`
  - Click `OK`

**Step 3: Accept CA Certificate**
- Go to: `Proxy` → `HTTP history`
- Make a simple HTTP request
- Burp will ask to accept certificate
- Accept it

**Step 4: Configure VaporTrace to Use Proxy**

In VaporTrace, before running commands:
```
VAPOR/INT> set-proxy http://127.0.0.1:8080

# Verify it's set
VAPOR/INT> keys
# Should show: PROXY (STAT): http://127.0.0.1:8080
```

### Run Test & Capture

**Step 5: Run Discovery Map**
```
VAPOR/INT> target https://target-api.com
VAPOR/INT> map
# Wait for scan to complete
```

**Step 6: Check Burp HTTP History**

In Burp:
- Click `Proxy` → `HTTP history`
- Look for requests FROM VaporTrace

**What to look for:**

```
Request Details:
✓ Host/Port correct?
✓ URL path correct?
✓ User-Agent header present?
✓ Headers look realistic?

Response Details:
✓ Status code (200, 401, 403, 500)?
✓ Response body not empty?
✓ Response headers present?
```

### Interpret Results

**If you see requests in Burp:**
- ✅ Requests ARE being sent
- ✅ TLS connection successful
- Check response for errors

**If you DON'T see requests:**
- ❌ VaporTrace not reaching Burp
- Check proxy setting
- Check firewall/network

**If you see requests but no responses:**
- ⚠️ Server not responding
- Target might be blocking
- Check target availability

---

## 🔬 METHOD 2: WIRESHARK (See Raw Traffic & TLS Details)

### Setup Wireshark

**Step 1: Install Wireshark**
```bash
# Ubuntu/Debian
sudo apt-get install wireshark

# macOS
brew install wireshark
```

**Step 2: Start Capture**
```bash
# Launch Wireshark GUI
wireshark &

# OR command line (if running in container)
sudo tcpdump -i any -w vaporTrace_capture.pcap 'tcp port 443 or tcp port 80'
```

**Step 3: Select Interface**
- In Wireshark, select interface (usually `eth0` or `docker0`)
- Click capture button

### Run Test

**Step 4: In another terminal, run VaporTrace without proxy**
```bash
cd /home/xoce/Workspace/VaporTrace
./VaporTrace  # or run it normally

# In UI:
VAPOR/INT> target https://ja3er.com
VAPOR/INT> map
```

**Step 5: Let it run and capture traffic**
- Let Wireshark capture for 30 seconds
- Or until you see activity

**Step 6: Stop capture in Wireshark**
- Click stop button

### Analyze Capture

**Step 7: Filter for TLS Traffic**
```
# In Wireshark's filter box, type:
tls.handshake
```

**This shows all TLS handshake packets**

**Step 8: Expand a TLS Handshake packet**
- Double-click on `TLS` or `SSL` row
- Look for: `Client Hello`
- Expand it

**What to look for in Client Hello:**

```
✓ TLS Version: Should be TLS 1.2 or 1.3
✓ Cipher Suites: Should see multiple
✓ Extensions: Should be > 5
✓ SNI: Should show target hostname
```

**Example Output:**
```
TLS Record Layer: Handshake Protocol: Client Hello
    Version: TLS 1.2 (0x0303)
    Session ID: (empty)
    Cipher Suites (20 suites):
        0x1301 - TLS_AES_128_GCM_SHA256
        0x1302 - TLS_AES_256_GCM_SHA384
        0x1303 - TLS_CHACHA20_POLY1305_SHA256
        ...
    Extensions:
        server_name (SNI): ja3er.com
        supported_groups: x25519, secp256r1, secp384r1
        supported_versions: TLS 1.3, TLS 1.2
        ...
```

**Step 9: Look for Server Response**
- Find `Server Hello` in same packet sequence
- Check TLS version negotiated
- Check cipher selected

---

## 📊 METHOD 3: COMBINED APPROACH (Burp + Wireshark)

**Best for Complete Picture:**

### Setup
```
┌─────────────────────────────────┐
│     VaporTrace                  │
│  (with proxy = 127.0.0.1:8080) │
└──────────────┬──────────────────┘
               │ HTTP (port 8080)
               ↓
        ┌──────────────┐
        │ Burp Proxy   │
        │ (8080)       │
        └──────────────┘
               │ HTTPS (port 443)
               ↓
    ┌──────────────────┐
    │ Wireshark        │
    │ (capture)        │
    └──────────────────┘
               │
               ↓
         Target Server
```

### Execute

**Terminal 1: Start Wireshark**
```bash
sudo wireshark &
# Select interface, start capture
```

**Terminal 2: Start Burp**
```bash
burpsuite &
# Configure proxy as above
```

**Terminal 3: Run VaporTrace**
```bash
cd /home/xoce/Workspace/VaporTrace
# Set proxy
set-proxy http://127.0.0.1:8080

# Run test
target https://ja3er.com
map
```

### Analyze Both Simultaneously

**In Burp:**
- See what VaporTrace THINKS it's sending
- See responses from target

**In Wireshark:**
- See actual TLS handshake
- See cipher negotiation
- Confirm uTLS vs standard TLS

---

## 🔍 DIAGNOSTIC CHECKLIST

### Basic Connectivity

```bash
# In terminal, before running VaporTrace
# Test basic connectivity to target:

curl -v https://ja3er.com
# Should work and show TLS details

# Test through Burp proxy:
curl -v -x http://127.0.0.1:8080 https://ja3er.com
# Should show in Burp
```

### TLS-Specific Checks in Wireshark

**Look for these Red Flags:**

```
❌ "TLS Alert: Handshake Failure"
   → Certificate/cipher mismatch

❌ "Packet too short" 
   → Fragmentation issue

❌ No SNI in Client Hello
   → Server Name Indication missing

❌ Very few cipher suites (< 5)
   → Looks like Go standard library

✅ Many cipher suites (10+)
   → Looks like uTLS (good!)

✅ SNI present
   → Proper configuration

✅ TLS 1.3
   → Modern connection
```

---

## 🎯 WHAT TO CAPTURE & SEND ME

### For Diagnosis, Capture:

**1. VaporTrace Logs (Full Output)**
```bash
# Run VaporTrace with full logging
# Copy ALL terminal output when running:
target https://ja3er.com
map

# Screenshot or save output
```

**2. Burp HTTP History**
```bash
# In Burp: Proxy → HTTP history
# Right-click on request → "Copy all columns"
# Paste into file or screenshot

# Or export: Proxy → Options → "Request/Response copy"
```

**3. Wireshark Capture (If possible)**
```bash
# Save capture
# File → Export Packets As → vaporTrace_capture.pcapng
# Send the .pcapng file or describe what you see
```

**4. Specific Information Needed:**

```
From VaporTrace logs:
- "[TLS:[-]" messages - any errors?
- "[cyan]TLS PROFILE:[-]" - what profile selected?
- "[yellow]WARNING:[-]" - any fallback messages?
- "[red]ERROR:[-]" - any error messages?

From Burp:
- Status code of requests (200, 401, 403, 500)?
- Response body (truncated or empty)?
- Content-Type header?
- Any error messages in response?

From Wireshark:
- TLS version (1.2 or 1.3)?
- Cipher suite count (many or few)?
- Any "Alert" packets?
- SNI present or missing?
```

---

## 🛠️ QUICK FIXES BASED ON FINDINGS

### If No Requests in Burp

**Problem:** Proxy not working
```bash
# Check VaporTrace proxy setting:
keys
# Should show PROXY set

# Test manually:
curl -v -x http://127.0.0.1:8080 http://httpbin.org/get
# Should appear in Burp
```

**Fix:**
```bash
# In VaporTrace:
set-proxy http://127.0.0.1:8080

# Verify:
keys | grep PROXY
```

### If Requests in Burp But No Findings in VaporTrace

**Problem:** Responses not being processed
```bash
# Check response in Burp:
- Is status code 200?
- Is body not empty?
- Is Content-Type correct?
```

**Likely Cause:** Response parser issue, not TLS

### If TLS Handshake Failure in Wireshark

**Problem:** Certificate/cipher issue
```bash
# In Wireshark, look for:
- "TLS Alert: Handshake Failure"
- Handshake packet sequence incomplete
```

**Fix:** May need to disable uTLS
```go
// In tls_evasion.go, change:
EnableJitter: true  // Try disabling to test
// Or check certificate error
```

---

## 🚀 STEP-BY-STEP TEST CASE

### Simplest Test to Start

**Step 1: Test with HTTP (no TLS)**
```bash
# Use simple HTTP target
VAPOR/INT> target http://httpbin.org
VAPOR/INT> map
# This should work regardless of TLS
# If findings appear → TLS is the issue
# If still no findings → Different issue
```

**Step 2: Test with HTTPS + Burp**
```bash
VAPOR/INT> set-proxy http://127.0.0.1:8080
VAPOR/INT> target https://httpbin.org
VAPOR/INT> map
# Check Burp for requests
# If requests appear → Network works, findings issue elsewhere
# If no requests → Proxy/TLS issue
```

**Step 3: Test with HTTPS + Wireshark**
```bash
# No proxy, just Wireshark capture
VAPOR/INT> target https://httpbin.org
VAPOR/INT> map
# Look in Wireshark for TLS traffic
# If TLS handshakes → Connection successful
# If no TLS → Network/firewall issue
```

**Step 4: Analyze Results**
```bash
# Based on what you find:
- HTTP works → TLS issue
- HTTP + Burp works → uTLS issue specifically
- Wireshark shows TLS → Connection success, response handling issue
- No traffic anywhere → Network blocked
```

---

## 📋 REPORTING TEMPLATE

When you capture the data, report in this format:

```
=== VaporTrace Test Report ===

Target: [URL tested]
Command: [map/scan/etc]
Time: [duration]

=== VaporTrace Logs ===
[Paste key TLS-related messages]

=== Burp HTTP History ===
- Request sent? [Yes/No]
- Status code: [200/401/500/etc]
- Response body: [Present/Empty/Error]
- Header issues: [None/Missing/Malformed]

=== Wireshark Analysis ===
- TLS handshake: [Successful/Failed/No packets]
- TLS version: [1.2/1.3/Unknown]
- Cipher suite count: [Number]
- SNI present: [Yes/No]
- Alerts: [None/List any]

=== Findings ===
- Endpoints found: [Number]
- Loot captured: [Yes/No/Number]
- Errors in logs: [List]

=== My Assessment ===
[What you think the issue is]
```

---

## 🎓 EXAMPLE: GOOD vs BAD Output

### ✅ GOOD Output (Everything Working)

**Burp:**
```
Request: POST /api/users
Status: 200
Response: {"users": [{...}, {...}]}
Content-Type: application/json
```

**Wireshark:**
```
TLS Client Hello with 12 cipher suites
TLS Server Hello - TLS 1.3 selected
Cipher: TLS_AES_128_GCM_SHA256
SNI: target.com present
```

**VaporTrace:**
```
[cyan]TLS:[-] Dialing api.target.com:443 with profile chrome-windows
[cyan]TLS:[-] Connected with chrome-windows profile
[cyan]MAP:[-] Found 25 endpoints
[green]LOOT:[-] Captured 3 credentials
```

### ❌ BAD Output (Problem)

**Burp:**
```
No requests appear
```

**Wireshark:**
```
No TLS packets captured
```

**VaporTrace:**
```
[yellow]WARNING:[-] uTLS failed
[yellow]FALLBACK:[-] Using standard TLS
[blue]MAP:[-] Found 0 endpoints
```

---

## 📞 Next Steps

1. **Choose one method** (Burp is easiest to start)
2. **Run the test** following the steps
3. **Capture the output** (logs, Burp history, or Wireshark)
4. **Report findings** using the template
5. **I'll analyze** and provide specific fix

---

**Ready to diagnose? Pick METHOD 1 (Burp) and report back!**
