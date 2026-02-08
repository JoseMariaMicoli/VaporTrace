# Interactive Keyboard Shortcuts

**Document:** manuals/17_KEYBOARD_SHORTCUTS.md  
**Version:** 3.0+  
**Reference:** Interactive UI Hotkeys (19 total)

---

## Hotkeys Guide

VaporTrace supports 19 keyboard shortcuts across three categories:
- **Global Hotkeys**: Work everywhere (UI, modals, command line)
- **Modal Hotkeys**: Work only inside modals (help, interceptor, analysis)
- **Tab Hotkeys**: Tab-specific navigation

---

## Global Hotkeys (Always Available)

### F1 - Toggle Logs Panel
**Scope:** Global  
**Action:** Show/hide F1 (LOGS) tab

```
Press: F1
Result: Toggle visibility of tactical feed logs
Context: Shows live attack execution logs
```

**When to use:** To monitor attack progress in real-time

### F2 - Toggle Map Panel
**Scope:** Global  
**Action:** Show/hide F2 (MAP) tab - discovered endpoints

```
Press: F2
Result: Toggle visibility of discovered endpoints table
Context: Shows all endpoints found by map/swagger/scrape/mine
```

**When to use:** Review your attack surface

### F3 - Toggle Loot Panel
**Scope:** Global  
**Action:** Show/hide F3 (LOOT) tab - captured secrets

```
Press: F3
Result: Toggle visibility of captured secrets
Context: Credentials, tokens, emails found during exploitation
```

**When to use:** Review captured data between attacks

### F4 - Toggle Analysis Panel
**Scope:** Global  
**Action:** Show/hide F4 (ANALYSIS) tab

```
Press: F4
Result: Toggle visibility of findings analysis
Context: Vulnerabilities, confidence scores, remediation
```

**When to use:** Review vulnerability analysis

### F5 - Toggle Plan Panel
**Scope:** Global  
**Action:** Show/hide F5 (PLAN) tab - tactical actions

```
Press: F5
Result: Toggle visibility of tactical action plan
Context: All pending/executed actions from 'analyze' and 'commit'
```

**When to use:** Monitor tactical action execution

### F6 - Toggle Neuro Panel
**Scope:** Global  
**Action:** Show/hide F6 (NEURO) tab - AI mutations

```
Press: F6
Result: Toggle visibility of neural engine panel
Context: Shows AI-generated payloads and mutations
```

**When to use:** See AI payload alternatives

### F7 - Toggle Settings Panel
**Scope:** Global  
**Action:** Show/hide F7 (SETTINGS) tab - configuration

```
Press: F7
Result: Toggle visibility of settings/configuration panel
Context: Runtime settings, API keys, proxy, modules
```

**When to use:** Adjust configuration mid-test

### Ctrl+H - Open Help Modal
**Scope:** Global  
**Action:** Display comprehensive help modal

```
Press: Ctrl+H
Result: Opens help modal showing all 19 hotkeys
Format: Organized table with Hotkey | Scope | Description
```

**When to use:** Need to quickly see hotkeys (inside dashboard)

### Ctrl+I - Open Interceptor
**Scope:** Global  
**Action:** Open MITM request interceptor modal

```
Press: Ctrl+I
Result: Opens interceptor modal showing pending requests
Context: Inspect/modify requests, forward/drop
Features: Edit payload, view response, modify headers
```

**When to use:** Intercept and modify requests in real-time

### Ctrl+F - Search
**Scope:** Global  
**Action:** Open search modal

```
Press: Ctrl+F
Result: Opens search modal to filter logs/endpoints/loot
```

**When to use:** Find specific endpoints or logs

### Ctrl+D - Debug Mode
**Scope:** Global  
**Action:** Toggle verbose debug output

```
Press: Ctrl+D
Result: Enable/disable debug logging (verbose)
Output: Additional technical details in logs panel
```

**When to use:** Troubleshooting issues

### Ctrl+B - Batch Operation
**Scope:** Global  
**Action:** Open batch operation dialog

```
Press: Ctrl+B
Result: Opens batch operations menu
```

**When to use:** Apply operations to multiple items

### Ctrl+S - Save Session
**Scope:** Global  
**Action:** Explicitly save current session

```
Press: Ctrl+S
Result: Saves session to database
Auto-saved: Every 30 seconds anyway
```

**When to use:** Ensure data saved before critical operations

### Ctrl+A - Select All
**Scope:** Global  
**Action:** Select all items in current view

```
Press: Ctrl+A
Result: Highlight all endpoints/loot/actions for batch ops
```

**When to use:** When performing bulk operations

### Ctrl+W - Close Tab
**Scope:** Global  
**Action:** Close current tab (if closeable)

```
Press: Ctrl+W
Result: Closes F6/F7 (NEURO/SETTINGS) tabs
Note: Main 5 tabs (F1-F5) cannot be closed
```

**When to use:** Minimize clutter

### Ctrl+X - Export
**Scope:** Global  
**Action:** Export current view to file

```
Press: Ctrl+X
Result: Export visible data (endpoints, loot, findings)
Format: CSV, JSON, or Markdown depending on context
```

**When to use:** Share data with team

---

## Modal-Specific Hotkeys

### Escape - Close/Cancel Modal
**Scope:** Modal context (Help, Interceptor, Search, etc)  
**Action:** Close modal and return to dashboard

```
Press: Esc
Result: Closes any open modal dialog
State: Returns to main dashboard
```

**When to use:** Exit modal without action

### Page Up - Scroll Up
**Scope:** Modal/List context  
**Action:** Scroll up in scrollable views

```
Press: Page Up
Result: Scroll up 10 items
Context: In modal lists, endpoints table, loot table
```

**When to use:** Review earlier items in list

### Page Down - Scroll Down
**Scope:** Modal/List context  
**Action:** Scroll down in scrollable views

```
Press: Page Down
Result: Scroll down 10 items
Context: In modal lists, endpoints table, loot table
```

**When to use:** Review later items in list

---

## Tab Navigation

### Tab Key - Next Tab
**Scope:** Dashboard  
**Action:** Navigate to next tab (circular)

```
Press: Tab
Result: Move focus to next tab
Order: F1 → F2 → F3 → F4 → F5 → F6 → F7 → F1
```

### Shift+Tab - Previous Tab
**Scope:** Dashboard  
**Action:** Navigate to previous tab

```
Press: Shift+Tab
Result: Move focus to previous tab
Order: F7 → F6 → F5 → F4 → F3 → F2 → F1 → F7
```

---

## Use Case Examples

### Example 1: Complete Reconnaissance
```
1. Press F2 to focus MAP tab
2. Type: target https://api.example.com
3. Type: map
4. F2 shows discovered endpoints
5. Press Page Down to scroll through endpoints
6. F3 to toggle LOOT visibility for comparison
7. Type: analyze
8. F5 to review generated tactical actions
```

### Example 2: Interactive Attack with Interception
```
1. Type: commit (start attack execution)
2. Press Ctrl+I to open interceptor
3. Interceptor shows pending requests
4. Edit payload in interceptor modal
5. Press Enter to forward modified request
6. Press Esc to close interceptor
7. F5 shows attack results in real-time
8. Press Page Down to see full results
```

### Example 3: Fast Vulnerability Review
```
1. Press Ctrl+H to see all hotkeys
2. Press Esc to close help
3. Press F1 to focus LOGS
4. Press Ctrl+F to search for "VULNERABLE"
5. Results show matching vulnerabilities
6. F3 to view associated loot
7. Ctrl+X to export findings
```

### Example 4: Multi-Step Exploitation
```
1. F7 to check SETTINGS
2. Verify proxy/API key configs
3. Ctrl+S to save settings
4. Ctrl+W to close SETTINGS tab
5. Type: analyze
6. F5 to review plan
7. Edit specific actions with arrow keys
8. F1 to monitor execution logs
9. Ctrl+D for debug details if needed
10. Ctrl+X to export results
```

### Example 5: Help and Documentation
```
1. At any time, press Ctrl+H for interactive help
2. Hotkeys table displays with descriptions
3. Press Page Down for more hotkeys
4. Type: help <command> in command line for command help
5. Type: help keys to see keys command info
6. Type: usage to see all commands
7. Type: help bola for exploitation module help
8. Press Esc to close any modal
```

---

## Hotkeys by Context

### In Logs Tab (F1)
- **Ctrl+F** - Search logs
- **Page Up/Down** - Scroll log history
- **Ctrl+A** - Select all logs
- **Ctrl+C** - Copy selected
- **Ctrl+X** - Export logs

### In Map Tab (F2)
- **Ctrl+F** - Find endpoints
- **Page Up/Down** - Scroll endpoints
- **Ctrl+A** - Select all endpoints
- **Ctrl+B** - Batch test selected endpoints

### In Loot Tab (F3)
- **Ctrl+F** - Search credentials
- **Page Up/Down** - Scroll loot items
- **Ctrl+A** - Select all loot
- **Ctrl+C** - Copy secret to clipboard
- **Ctrl+X** - Export secrets

### In Analysis Tab (F4)
- **Ctrl+F** - Search vulnerabilities
- **Page Up/Down** - Scroll findings
- **Ctrl+A** - Select all findings

### In Plan Tab (F5)
- **Ctrl+F** - Search actions
- **Page Up/Down** - Scroll actions
- **Ctrl+A** - Select all actions
- **Ctrl+B** - Batch modify selected

### In Neuro Tab (F6)
- **Page Up/Down** - Scroll mutations
- **Ctrl+C** - Copy payload
- **Ctrl+X** - Export payloads

### In Settings Tab (F7)
- **Ctrl+S** - Save settings

### Interceptor Modal (Ctrl+I)
- **Tab/Shift+Tab** - Cycle through requests
- **Page Up/Down** - Scroll request details
- **Enter** - Forward request
- **Delete** - Drop request
- **Esc** - Close interceptor

### Help Modal (Ctrl+H)
- **Page Up/Down** - Scroll hotkeys
- **Ctrl+F** - Search within help
- **Esc** - Close help

---

## Hotkeys Comparison with Popular Tools

| Action | VaporTrace | Burp | OWASP ZAP |
|--------|-----------|------|-----------|
| Help | Ctrl+H | F1 | F1 |
| Search | Ctrl+F | Ctrl+F | Ctrl+F |
| Interceptor | Ctrl+I | Custom | Custom |
| Export | Ctrl+X | - | - |
| Close Tab | Ctrl+W | Ctrl+W | - |
| Scroll | Page Up/Down | Scroll | Scroll |
| Select All | Ctrl+A | Ctrl+A | - |

---

## Keyboard Layout Diagram

```
┌─────────────────────────────────────────────────────┐
│ F1   F2   F3   F4   F5   F6   F7                   │
│ LOGS MAP  LOOT ANAL PLAN NEURO SET                  │
├─────────────────────────────────────────────────────┤
│ ESC  [F keys toggle tabs]                           │
│ CTRL+H = Help                                       │
│ CTRL+I = Interceptor                                │
│ CTRL+F = Search                                     │
│ CTRL+D = Debug                                      │
│ CTRL+B = Batch                                      │
│ CTRL+S = Save                                       │
│ CTRL+A = Select All                                 │
│ CTRL+W = Close Tab                                  │
│ CTRL+X = Export                                     │
│ PAGE UP/DOWN = Scroll                               │
│ TAB/SHIFT+TAB = Tab navigation                      │
└─────────────────────────────────────────────────────┘
```

---

## Tips & Tricks

### Rapid Navigation
- Use F1-F7 for direct tab access (don't tab through)
- Ctrl+H to remind yourself of all hotkeys
- Ctrl+F + search term for fast filtering

### Batch Operations
- Ctrl+A to select all items
- Ctrl+B for batch operations
- Apply same settings to multiple endpoints

### Real-time Monitoring
- Open F1 (LOGS) for execution logs
- F5 for tactical action status
- F3 for loot in real-time
- Ctrl+D for extra details

### Export & Sharing
- Ctrl+X from any tab to export
- Creates CSV/JSON files for team review
- Useful for compliance documentation

### Debugging
- Ctrl+D to enable debug mode
- F1 will show verbose logs
- Ctrl+I to inspect intercepted requests
- Ctrl+F to find specific debug entries

---

## Cheat Sheet (Print This!)

```
TAB TOGGLES (Press once to hide, again to show)
F1=LOGS | F2=MAP | F3=LOOT | F4=ANALYSIS | F5=PLAN | F6=NEURO | F7=SETTINGS

CORE OPERATIONS
Ctrl+H = Help (all hotkeys)      Ctrl+I = Interceptor
Ctrl+F = Search/Find              Ctrl+D = Debug mode
Ctrl+S = Save session             Ctrl+A = Select all
Ctrl+B = Batch operations         Ctrl+W = Close tab
Ctrl+X = Export                   Esc = Close modal

NAVIGATION
Tab = Next tab              Shift+Tab = Previous tab
Page Up = Scroll up         Page Down = Scroll down
```

---

See also: [18_COMMAND_REFERENCE.md](18_COMMAND_REFERENCE.md) for CLI commands, and [02_FIRST_RUN.md](02_FIRST_RUN.md) for step-by-step guide.

