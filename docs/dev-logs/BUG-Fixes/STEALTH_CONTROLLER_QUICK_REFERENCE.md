# Stealth Controller - Quick Reference

## Hotkey Summary

| Hotkey | Command | Effect |
|--------|---------|--------|
| `Ctrl+A` | `stealth aggressive` | Maximum evasion, 1x speed |
| `Ctrl+F` | `stealth fast` | Reduced delays, 0.5x speed |
| `Ctrl+S` | `stealth silent` | Maximum delays, 2x speed |
| `Ctrl+D` | `stealth debug` | All evasion disabled |

## Command Examples

### View Current Status
```
> status stealth
[cyan]STEALTH STATUS[-]
  Mode: Aggressive | Multiplier: 1.0x
  Jitter: [green]ON | Thinking: [green]ON | Backoff: [green]ON | 
  Obfuscation: [green]ON | Encoding: [green]ON
```

### Switch Modes
```
> stealth fast
[blue::b]STEALTH:[-] Mode set to [blue]FAST[-] (reduced delays, 0.5x speed)[-:-:-]

> stealth silent
[green::b]STEALTH:[-] Mode set to [green]SILENT[-] (maximum evasion, 2x speed)[-:-:-]

> stealth aggressive
[red::b]STEALTH:[-] Mode set to [red]AGGRESSIVE[-] (all evasion enabled, 1x speed)[-:-:-]

> stealth debug
[yellow::b]STEALTH:[-] Mode set to [yellow]DEBUG[-] (all evasion disabled)[-:-:-]
```

### Toggle Individual Features
```
> toggle jitter on
[cyan]STEALTH:[-] Jitter [green]ENABLED

> toggle thinking off
[cyan]STEALTH:[-] Thinking Time [red]DISABLED

> toggle backoff on
[cyan]STEALTH:[-] Rate-Limit Backoff [green]ENABLED

> toggle obfuscation off
[cyan]STEALTH:[-] Path Obfuscation [red]DISABLED

> toggle encoding on
[cyan]STEALTH:[-] Payload Encoding [green]ENABLED
```

### Adjust Speed Multiplier
```
> multiplier 0.1
[cyan]STEALTH:[-] Global evasion multiplier set to 0.1x

> multiplier 0.5
[cyan]STEALTH:[-] Global evasion multiplier set to 0.5x

> multiplier 1.0
[cyan]STEALTH:[-] Global evasion multiplier set to 1.0x

> multiplier 2.0
[cyan]STEALTH:[-] Global evasion multiplier set to 2.0x

> multiplier 5.0
[cyan]STEALTH:[-] Global evasion multiplier set to 5.0x
```

## Configuration Presets

### Aggressive (Default)
- **When**: First reconnaissance, lighter WAF targets
- **Command**: `stealth aggressive`
- **Behavior**: 
  - All evasion techniques active
  - Normal speed (1.0x multiplier)
  - Good balance of stealth and speed

### Fast
- **When**: Time-sensitive operations, lighter WAF
- **Command**: `stealth fast`
- **Behavior**: 
  - Payload encoding disabled
  - 0.5x multiplier (2x faster)
  - Quick reconnaissance phase

### Silent
- **When**: Hardened targets, long-duration scans
- **Command**: `stealth silent`
- **Behavior**: 
  - All evasion techniques active
  - 2.0x multiplier (2x slower)
  - Maximum stealth posture

### Debug
- **When**: Troubleshooting, testing
- **Command**: `stealth debug`
- **Behavior**: 
  - All evasion disabled
  - Normal speed (1.0x)
  - Direct baseline testing

## Real-World Scenarios

### Scenario 1: Increasing 429 Responses
```
Initial: stealth fast (running quick scan)
↓ (after seeing 429 rate-limit responses)
> stealth silent
> multiplier 2.0
Result: Slower requests, longer delays between retries
```

### Scenario 2: Deep Reconnaissance
```
Initial: stealth debug (verify target responds)
↓ (target confirmed live)
> stealth aggressive
↓ (after 5 minutes, no blocking)
> stealth silent
Result: Progressive hardening as scan progresses
```

### Scenario 3: Time-Constrained Assessment
```
> stealth aggressive
> multiplier 0.1
Result: Maximum speed while keeping evasion active
(Only use when WAF is weak or testing shows high tolerance)
```

### Scenario 4: Debugging Response Encoding
```
> stealth debug
(verify target responds without any evasion)
↓
> toggle encoding on
(test only payload encoding)
↓ 
> toggle thinking on
(add thinking time delays)
Result: Systematic testing of individual techniques
```

## Troubleshooting

### "Timeouts increasing - requests blocked"
```
Current: stealth silent (2.0x multiplier)
→ Try: multiplier 1.0
→ Or: toggle encoding off
→ Or: stealth fast (0.5x multiplier)
```

### "Too many 429 responses"
```
Current: stealth aggressive
→ Switch: stealth silent
→ Then: multiplier 2.0 or 3.0
→ Or: toggle backoff on (if off)
```

### "Requests completing too slowly"
```
Current: stealth silent (2.0x)
→ Try: multiplier 1.0
→ Then: stealth fast (0.5x)
→ Or: multiplier 0.1 (experimental)
```

### "Need to verify target accepts requests"
```
→ Command: stealth debug
→ This disables all evasion
→ Make baseline request
→ If succeeds: target accepts traffic
→ Then: stealth aggressive (re-enable evasion)
```

## Performance Tips

| Goal | Configuration |
|------|---------------|
| **Maximum Speed** | `stealth fast` + `multiplier 0.1` |
| **Balanced** | `stealth aggressive` + `multiplier 1.0` |
| **Maximum Stealth** | `stealth silent` + `multiplier 2.0` |
| **Testing** | `stealth debug` (all evasion off) |
| **Fine-tuning** | Keep stealth mode, adjust multiplier |

## Monitor During Execution

Watch the tactical log for:
- `[cyan]STEALTH:[-]` entries show toggle/multiplier changes
- `[yellow]STEALTH:[-]` entries show skipped sleeps
- `[red::b]STEALTH:[-]` entries show mode changes (red = aggressive)
- `[blue::b]STEALTH:[-]` entries show mode changes (blue = fast)
- `[green::b]STEALTH:[-]` entries show mode changes (green = silent)
- `[yellow::b]STEALTH:[-]` entries show mode changes (yellow = debug)

Example log flow:
```
[red::b]STEALTH:[-] Mode set to [red]AGGRESSIVE
[blue]THINKING:[-] Discovery request: Light jitter (10-50ms)
[blue]EVASION:[-] Applied stochastic jitter
[cyan]BEHAVIOR:[-] Contextual thinking time: 45ms
[yellow]PAUSE:[-] Human behavior simulation: sleeping for 2 seconds
```
