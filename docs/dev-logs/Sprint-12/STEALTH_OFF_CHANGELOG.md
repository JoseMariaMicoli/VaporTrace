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

# Stealth Off Feature - Implementation Changelog

**Date:** February 2026  
**Feature:** "stealth off" Kill Switch Command  
**Purpose:** Easy one-command deactivation of all evasion techniques

---

## What's New

### 1. **Stealth Off Command** 
Added `stealth off` subcommand to disable all evasion techniques immediately.

**Usage:**
```bash
> stealth off
```

**Effects:**
- Disables all 5 evasion techniques: Jitter, Thinking, Backoff, Obfuscation, Encoding
- Resets Global Multiplier to 1.0x
- Running in fastest/most aggressive mode
- No delays, no obfuscation, no encoding

**Output:**
```
[green]STEALTH MODE: OFF[-]
[yellow]All evasion techniques disabled. Running in fastest/most aggressive mode.[-]
```

### 2. **Updated Usage Documentation**
Added `stealth off` to both `usage` and `usage 2` output:

**Location:** `pkg/engine/core.go` - printUsagePage2() function

```
[yellow]stealth off[-]      Kill Switch            Disable all evasion (fastest/most aggressive mode).
```

### 3. **Updated Help Documentation**
Added comprehensive help entry for `stealth off`:

**Location:** `pkg/engine/core.go` - printHelp() function

```bash
case "stealth off":
    Display stealth off documentation
    - Kill switch explanation
    - Effects (disables all 5 techniques + resets multiplier)
    - Usage scenarios
    - Warning about higher detection risk
```

### 4. **Manual Documentation Update**
Updated command reference manual with detailed `stealth off` documentation:

**Location:** `docs/manuals/18_COMMAND_REFERENCE.md`

**Contents:**
- Full command syntax
- Description and effects
- When to use (speed testing, aggressive probing, quick mapping)
- Warnings (detection risk on WAF targets)
- Real usage examples

---

## Implementation Details

### Code Changes

**File:** `pkg/engine/core.go`

**1. Stealth Command Handler (lines 649-707)**
```go
case "off":
    logic.SetEvasionToggle("jitter", false)
    logic.SetEvasionToggle("thinking", false)
    logic.SetEvasionToggle("backoff", false)
    logic.SetEvasionToggle("obfuscation", false)
    logic.SetEvasionToggle("encoding", false)
    logic.SetGlobalMultiplier(1.0)
    utils.TacticalLog("[green]STEALTH MODE: OFF[-]")
    utils.TacticalLog("[yellow]All evasion techniques disabled. Running in fastest/most aggressive mode.[-]")
```

**2. Updated Command Routing**
- Added "off" to stealth subcommand switch statement
- Proper arg validation and error handling
- Clear user feedback through TacticalLog

---

## Functionality Summary

| Feature | Status | Location |
|---------|--------|----------|
| `stealth off` command | ✅ Implemented | pkg/engine/core.go |
| Disables all 5 techniques | ✅ Implemented | pkg/logic/evasion.go |
| Resets multiplier to 1.0x | ✅ Implemented | pkg/logic/evasion.go |
| Usage documentation | ✅ Updated | pkg/engine/core.go - printUsagePage2() |
| Help entry | ✅ Added | pkg/engine/core.go - printHelp() |
| Manual documentation | ✅ Updated | docs/manuals/18_COMMAND_REFERENCE.md |
| Dashboard indicator update | ✅ Works | pkg/ui/dashboard.go (auto-updates) |

---

## Testing Scenarios

### Scenario 1: Quick Mode Switching
```bash
> stealth silent
[green]Stealth Mode:[-] silent
[yellow]All evasion techniques enabled...

> stealth off
[green]STEALTH MODE: OFF
[yellow]All evasion techniques disabled. Running in fastest/most aggressive mode.

> bola https://target.com/api/users/[id]
(Rapid exploitation without evasion delays)
```

### Scenario 2: Dashboard Indicator Update
- When `stealth off` is executed, dashboard status bar immediately updates
- Evasion indicator shows all toggles as darkgray (OFF)
- Mode indicator remains visible but indicates no evasion active

### Scenario 3: Speed Comparison
```bash
> stealth silent
> time bola https://target.com/api/users/[id]   (slow, with delays)

> stealth off
> time bola https://target.com/api/users/[id]   (fast, no delays)
```

---

## User Benefits

1. **Quick Deactivation:** Single command instead of manually toggling 5 features
2. **Maximum Speed:** Fastest possible mode for time-sensitive scenarios
3. **Easy Toggle:** Simple `stealth off` vs. `stealth aggressive` for mode switching
4. **Clear Intent:** Explicit "off" state is immediately clear to operators
5. **Dashboard Feedback:** Real-time indicator shows all evasion disabled

---

## Backward Compatibility

✅ **Fully Compatible** - All existing stealth modes and commands unchanged:
- `stealth aggressive|fast|silent|debug` - Still work as before
- `stealth status` - Still displays current state
- `stealth toggle` - Still works for individual features
- `stealth multiplier` - Still works for global scaling

---

## Auto-completion

`stealth off` automatically added to knownCommands[] for CLI auto-completion.

---

## Documentation Files Modified

1. **pkg/engine/core.go**
   - Updated stealth command help text (line 653)
   - Updated usage documentation (line 1524)
   - Added "stealth off" help case (new)
   - Added "off" case to stealth command switch

2. **docs/manuals/18_COMMAND_REFERENCE.md**
   - Added complete "stealth off" command reference section
   - Included usage, effects, scenarios, warnings, examples

---

## Related Features

- **Stealth Modes:** stealth aggressive|fast|silent|debug
- **Evasion Toggle:** stealth toggle <feature> on|off
- **Multiplier Scaling:** stealth multiplier <0.1-5.0>
- **Status Display:** stealth status
- **Dashboard Indicator:** Real-time evasion status in F1 footer

---

## Next Steps (Optional Enhancements)

- Alias: `stealth 0` as shorthand for `stealth off`
- Config file option: `default_stealth_mode: off`
- Keyboard shortcut: Consider Ctrl+Alt+O for "stealth off" toggle
