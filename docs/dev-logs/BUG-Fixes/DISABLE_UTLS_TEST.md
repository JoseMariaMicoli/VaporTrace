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

# Emergency: Disable uTLS Temporarily to Test

**If you just need to get it working RIGHT NOW to test findings logic:**

---

## 🚨 QUICK DISABLE (60 seconds)

Edit `pkg/logic/network.go` and find line ~286 where it says:

```go
DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
    // Try uTLS first
    conn, err := tlsTransport.DialTLSContext(ctx, network, addr)
```

**Replace the ENTIRE DialTLSContext function with this simple version:**

```go
DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
    // Temporary: Use only standard TLS for testing
    host, _, _ := net.SplitHostPort(addr)
    
    dialer := &net.Dialer{Timeout: 30 * time.Second}
    tcpConn, err := dialer.DialContext(ctx, "tcp", addr)
    if err != nil {
        utils.TacticalLog(fmt.Sprintf("[red]TCP DIAL ERROR:[-] %v", err))
        return nil, err
    }
    
    tlsConn := tls.Client(tcpConn, &tls.Config{
        ServerName:         host,
        InsecureSkipVerify: true,
    })
    
    err = tlsConn.Handshake()
    if err != nil {
        tcpConn.Close()
        utils.TacticalLog(fmt.Sprintf("[red]TLS HANDSHAKE ERROR:[-] %v", err))
        return nil, err
    }
    
    utils.TacticalLog(fmt.Sprintf("[yellow]STANDARD TLS:[-] Connected to %s", addr))
    return tlsConn, nil
},
```

---

## ⚡ BUILD & TEST

```bash
cd /home/xoce/Workspace/VaporTrace

# Rebuild
go build ./...

# If errors, check import section has:
# "crypto/tls"
# "fmt" 
# (they should already be there)

# Run VaporTrace
./VaporTrace  # or however you start it

# In UI:
VAPOR/INT> target https://httpbin.org
VAPOR/INT> map

# Wait 30 seconds...
# Check if MAP tab populates now
```

---

## ✅ TEST RESULTS

### If findings appear NOW:
- ✅ Issue IS the uTLS code
- The standard TLS works fine
- Problem: uTLS configuration or library issue

### If STILL no findings:
- ❌ Issue is NOT uTLS
- Problem: Response handling, findings extraction, or other logic

---

## 🔍 WHAT TO LOOK FOR IN LOGS

**With standard TLS enabled, you should see:**

```
[yellow]STANDARD TLS:[-] Connected to httpbin.org:443
[yellow]STANDARD TLS:[-] Connected to httpbin.org:443
[yellow]STANDARD TLS:[-] Connected to httpbin.org:443
... (multiple connections)

[green]✓ MAP:[-] Found 5 endpoints
[green]✓ LOOT:[-] Captured credentials
```

**If still nothing:**
```
[yellow]STANDARD TLS:[-] Connected...
[blue]i[-] No endpoints found
```

→ Then it's NOT a TLS problem, it's response parsing

---

## 🔄 RESTORE uTLS (When ready)

When you're done testing and want to go back to uTLS:

```bash
# Revert the file:
git checkout pkg/logic/network.go

# Or manually change back to original code

# Rebuild:
go build ./...
```

---

## 📊 COMPARISON TABLE

| Setting | Requests? | Findings? | Time | Security |
|---------|-----------|-----------|------|----------|
| uTLS enabled | ? | ? | ? | ✅ High |
| Standard TLS | Should work | Should work | Fast | ⚠️ Standard |
| No TLS at all | Would fail | N/A | N/A | ❌ None |

---

## 💡 THIS TELLS US

- **If standard TLS works** → uTLS has a bug we need to fix
- **If standard TLS fails** → Completely different issue
- **If partially works** → Partial uTLS issue

---

**Do this now, report back if findings appear!**
