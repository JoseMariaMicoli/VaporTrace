# VaporTrace Keybindings Quick Reference

**Last Updated:** February 8, 2026  
**Access:** Press `Ctrl+H` in TUI or run `help keys` in CLI

---

## Global Hotkeys (Work Everywhere)

| Key | Function | Notes |
|-----|----------|-------|
| **Ctrl + H** | 📋 Show Keybindings Modal | Opens this reference popup |
| **Ctrl + I** | 🚨 Toggle Interceptor | On/Off switch for request interception |
| **F1** | 📊 LOGS Tab | Tactical feed & system messages |
| **F2** | 🗺️ MAP Tab | Discovered endpoints & attack surface |
| **F3** | 💎 LOOT Tab | Captured secrets & credentials |
| **F4** | 📡 TRAFFIC Tab | HTTP requests & responses |
| **F5** | 🎯 PLAN Tab | Strategic actions & planner |
| **F6** | 🧠 NEURO Tab | AI engine output & analysis |
| **F7** | 📄 REPORT Tab | Findings export (Markdown/PDF) |
| **Esc** | 🚪 Exit | Terminate VaporTrace (with confirmation) |

---

## Modal-Only Hotkeys (When Interceptor Modal Open)

| Key | Function | Notes |
|-----|----------|-------|
| **Ctrl + F** | ✈️ Forward Packet | Send to network immediately |
| **Ctrl + D** | ✋ Drop Packet | Block and discard packet |
| **Ctrl + B** | 🧬 Neuro Brute | Generate AI payloads for field |
| **Ctrl + S** | 💾 Sync Loot | Save request to Loot DB |

---

## Tab-Specific Hotkeys

### F1 Tab (Logs)
| Key | Function |
|-----|----------|
| **Page Up** | Scroll up in logs |
| **Page Down** | Scroll down in logs |

### F4 Tab (Traffic)
| Key | Function |
|-----|----------|
| **Ctrl + A** | Analyze snapshot with AI Brain |

### F7 Tab (Report)
| Key | Function |
|-----|----------|
| **Ctrl + W** | Save report to disk |
| **Ctrl + X** | Delete session & clear report |

---

## Command Line Help

```bash
# View all available commands
> usage

# Get help on specific command
> help <command>

# Examples:
> help analyze       # Strategic planning
> help commit        # Execute tactical actions
> help keys          # Show this keybindings list
> help neuro         # Neural Engine configuration
```

---

## Command Categories

### 🎯 Strategic Planning (HITL - Human-in-the-Loop)
- `analyze` - Generate attack vectors
- `list-plan` - View pending actions
- `edit` - Override payloads
- `drop` - Reject actions
- `commit` - Execute all actions

### 🔍 Reconnaissance  
- `target` - Set scope URL
- `map` - Full recon pipeline
- `swagger` - Parse OpenAPI specs
- `scrape` - Mine JavaScript
- `mine` - Fuzz parameters
- `sessions` - Auth management

### ⚔️ Exploitation (OWASP Top 10)
- `bola` - ID enumeration
- `bfla` - Privilege escalation
- `bopla` - Mass assignment
- `ssrf` - Internal access
- `exhaust` - Resource DoS
- `audit` - Config checks
- `probe` - Webhook injection
- `flow` - Attack chains

### 🧠 AI & Neural Engine
- `neuro` - Enable/disable AI
- `neuro-gen` - Generate payloads
- `test-neuro` - Connectivity test
- `ask` - Direct LLM query

### 🔧 Infrastructure
- `proxy` - Set proxy
- `proxies` - Load proxy list
- `init_db` - Create database
- `seed_db` - Add test data
- `reset_db` - Clear database
- `report` - Export findings

### 🛠️ Utilities
- `tasks` - Engine status
- `loot` - Vault manager
- `clear` - Clear logs
- `usage` - Command reference ← You are here!
- `help` - Detailed command help
- `exit` - Graceful shutdown

---

## How to Use the Keybindings Popup

### Opening
1. **In TUI:** Press `Ctrl+H` at any time
2. **In CLI:** Run `help keys` to see text version

### Navigating the Modal
- Use arrow keys to navigate the table (optional)
- Read entries for scope (Global/Modal/Tab-specific)
- See function description in right column

### Closing the Popup
- Press **Esc** - Standard exit
- Press **Enter** - Accept & close
- Press **Ctrl+H** - Toggle (close or reopen)

### After Closing
- Focus returns to command input (`>` prompt)
- You can continue typing commands immediately

---

## Pro Tips

✅ **Ctrl+H is your friend** - Forget a hotkey? Press Ctrl+H!  
✅ **F1-F7 for everything** - Tab navigation is fast and visual  
✅ **Ctrl+I for interceptor** - Most powerful feature, toggle on demand  
✅ **help <cmd> for details** - Every command has comprehensive help  
✅ **SSRF/Commit with F6 open** - Watch real-time AI mutations  

---

## Troubleshooting

**Q: Ctrl+H not working?**  
A: Make sure you're in the TUI (not stuck in command input with a long command). Try pressing Escape first to clear focus.

**Q: Interceptor modal won't close?**  
A: Press Esc, then try again. If stuck, press Ctrl+D to drop the packet and close the modal.

**Q: Can't switch tabs?**  
A: Make sure you're not in the command input field. Click outside it or press Escape, then use F1-F7.

**Q: Report tab commands not working?**  
A: You must be ON the F7 tab. Switch with F7, then try Ctrl+W to save or Ctrl+X to delete.

---

## Next Steps

1. Try `Ctrl+H` now to see this list in a modal!
2. Run `usage` to see all commands
3. Run `help analyze` to learn tactical planning
4. Ready? Run `target https://target-api.com` to get started!

---

**Happy hacking!** 🚀

