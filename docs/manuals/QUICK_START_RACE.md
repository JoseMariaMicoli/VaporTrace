# Quick Start: Race Condition Testing

VaporTrace now includes a dedicated synchronization engine to test for logic flaws like:
- Coupon reuse (double spending).
- Bypassing transfer limits.
- Creating multiple resources when only 1 is allowed.

## How to use it

### 1. Identify a Target
Find an endpoint that changes state (e.g., `POST /api/coupons/redeem` or `GET /api/gift?claim=1`).

### 2. Run the Race
In the VaporTrace dashboard:

```bash
# Basic test (20 threads)
race https://target.com/api/claim?code=WINNER

# High intensity (50 threads)
race https://target.com/api/claim?code=WINNER 50
```

### 3. Analyze Results
Watch the **F1 Logs**.
- **Gray:** "No variances detected" (Secure / Failed).
- **Red:** "⚠ RACE ANOMALY DETECTED".
  - Example: `Code 200: 5 responses`, `Code 400: 15 responses`.
  - This means 5 requests "won" the race successfully.

### 4. Verify
Check the **F7 Report** tab. A CRITICAL finding will be logged if status codes were mixed.

## Technical Details

### The Synchronization Gate
The race engine uses a channel-based barrier pattern:

```
1. Spawn N goroutines (workers)
2. All workers pre-build their requests and wait on startGate channel
3. Main thread closes startGate
4. All workers fire simultaneously (nanosecond precision)
5. Collect results and analyze for variance
```

### Detection Logic
- **Status Code Variance:** Mixed 200/400/500 responses indicate some requests succeeded while others failed → **CRITICAL**
- **Response Length Variance:** Different body sizes suggest one request succeeded differently → **HIGH**
- **Identical Responses:** All requests returned the same status/length → **SECURE** (no race found)

### Common TOCTOU Vulnerabilities Detected
| Scenario | Expected Behavior | Race Exploit |
|----------|-------------------|--------------|
| Gift card redemption | Only 1 claim per code | 50 threads, 5 succeed |
| Bank transfer limit bypass | Max 2 transfers/day | Race bypasses counter |
| Coupon double-spend | Single use only | N simultaneous claims |
| Resource creation | Only 1 allowed | Race creates N duplicates |

## Troubleshooting

### "No variances detected"
- Target endpoint may not be vulnerable to races
- Try increasing thread count: `race <url> 100`
- Check if endpoint requires authentication: `auth attacker <token>` first

### High false positives
- Some endpoints naturally return varied responses (randomization, load balancing)
- Verify manually by checking database/state changes, not just HTTP responses

### Timeout errors
- Increase timeout in `pkg/attack/race.go` (currently 10 seconds)
- Target may be rate-limiting or dropping connections
- Try reducing thread count or adding delays between requests

## Advanced Usage

### With F5 Planner Integration
The Neuro Engine can automatically suggest race tests:

```bash
analyze          # AI examines target
list-plan        # Review suggested actions (including RACE_CONDITION tasks)
commit           # Execute all planned tests (including races)
```

### With Proxy Integration
Send race traffic through your proxy:

```bash
proxy 127.0.0.1:8080
race https://target.com/api/claim?code=ABC 30
```

## Reporting
All race findings are logged to the database and appear in F7 reports:
- **Phase:** PHASE III: RACE CONDITION
- **Severity:** CRITICAL (CVSS 8.5+) or HIGH
- **Action:** **ARCHITECTURAL FIX REQ** (not a simple patch)

Findings include:
- Thread count used
- Status code distribution
- Response length variation
- Proof of simultaneous execution
