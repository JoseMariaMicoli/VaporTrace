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

# Installation & Setup Guide

**Document:** 01_INSTALLATION_SETUP.md  
**Version:** 3.0+  
**Time to Complete:** 10-15 minutes

---

## Table of Contents
1. [System Requirements](#system-requirements)
2. [Dependencies](#dependencies)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Verification](#verification)
6. [Troubleshooting](#troubleshooting)

---

## System Requirements

### Minimum
- **OS:** Linux (Ubuntu 18.04+, Debian 10+, CentOS 7+)
- **RAM:** 2GB (4GB recommended)
- **Disk:** 500MB free space
- **CPU:** 2 cores minimum

### Recommended
- **OS:** Ubuntu 20.04+ or Debian 11+
- **RAM:** 8GB+
- **Disk:** 2GB+ SSD
- **CPU:** 4+ cores
- **Network:** Stable internet connection (for AI features)

### Supported Platforms
- ✅ Linux (primary)
- ✅ WSL2 (Windows Subsystem for Linux)
- ✅ Docker (containerized)
- ✅ macOS (Homebrew, experimental)

---

## Dependencies

### System Packages (Ubuntu/Debian)

```bash
# Update package manager
sudo apt update
sudo apt upgrade -y

# Install Go (1.25+)
wget https://go.dev/dl/go1.25.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install build tools
sudo apt install -y build-essential git curl wget

# Install optional dependencies
sudo apt install -y sqlite3           # Database
sudo apt install -y postgresql        # Advanced DB (optional)
```

### Verify Go Installation

```bash
go version
# Expected: go version go1.25.6 linux/amd64
```

---

## Installation

### Method 1: From Source (Recommended)

```bash
# 1. Clone repository
git clone https://github.com/JoseMariaMicoli/VaporTrace.git
cd VaporTrace

# 2. Install Go dependencies
go mod download

# 3. Build binary
go build -o ./VaporTrace ./main.go

# 4. Verify build
./VaporTrace --version
# or
./VaporTrace --help
```

### Method 2: Pre-built Binary

```bash
# Download latest release
cd /tmp
wget https://github.com/JoseMariaMicoli/VaporTrace/releases/download/v3.0/VaporTrace-linux-amd64.tar.gz
tar -xzf VaporTrace-linux-amd64.tar.gz
sudo mv VaporTrace /usr/local/bin/

# Verify
VaporTrace --version
```

### Method 3: Docker

```bash
# Build container
docker build -t vaportrace:latest .

# Run container
docker run -it --network host vaportrace:latest

# Or with volume for data persistence
docker run -it -v $(pwd)/data:/app/data vaportrace:latest
```

---

## Configuration

### Initial Setup

```bash
# 1. Create config directory
mkdir -p ~/.vaportrace

# 2. Initialize database
VaporTrace
> init_db

# 3. Seed test data (optional - for demo)
> seed_db

# 4. Exit
> exit
```

### Configuration File

**Location:** `~/.vaportrace/config.yaml`

```yaml
# VaporTrace Configuration

# Target scope
target:
  url: ""                           # Default target URL
  timeout: 30                       # Request timeout (seconds)
  max_redirects: 10                 # Follow redirects

# Proxy configuration
proxy:
  enabled: false                    # Enable upstream proxy
  host: "localhost"                 # Proxy host
  port: 8080                        # Proxy port
  username: ""                      # Proxy auth username
  password: ""                      # Proxy auth password

# Authentication
auth:
  cookies_enabled: true             # Cookie support
  token_header: "Authorization"     # Token header name
  session_timeout: 3600             # Session timeout (seconds)

# AI/Neural Engine
neuro:
  enabled: false                    # Enable AI features
  provider: "groq"                  # AI provider (groq, openai, etc)
  model: "llama-3.1-8b"             # Model selection
  api_key: ""                       # API key (set via env: NEURO_API_KEY)
  temperature: 0.7                  # Creativity (0.0-1.0)
  max_tokens: 1000                  # Response length

# Interceptor
interceptor:
  enabled: false                    # Enable by default
  listen_port: 8081                 # Interception port
  save_requests: true               # Save to loot DB

# Database
database:
  type: "sqlite"                    # sqlite or postgresql
  path: "~/.vaportrace/vaportrace.db" # SQLite path
  postgresql:
    host: "localhost"               # PostgreSQL host
    port: 5432                      # PostgreSQL port
    user: "vaportrace"              # PostgreSQL user
    password: ""                    # PostgreSQL password
    database: "vaportrace"          # PostgreSQL database

# Logging
logging:
  level: "info"                     # debug, info, warn, error
  file: "~/.vaportrace/vaportrace.log"
  max_size_mb: 100                  # Log rotation size

# Report generation
report:
  format: "markdown"                # markdown or pdf
  include_screenshots: false        # Include UI screenshots
  compress: true                    # Compress findings
```

### Environment Variables

```bash
# AI Provider Configuration
export NEURO_API_KEY="your-api-key-here"
export NEURO_PROVIDER="groq"           # or openai, claude, etc
export NEURO_MODEL="llama-3.1-8b-instant"

# Proxy Configuration
export VAPORTRACE_PROXY="http://proxy.example.com:8080"
export VAPORTRACE_PROXY_USER="username"
export VAPORTRACE_PROXY_PASS="password"

# Database
export VAPORTRACE_DB="sqlite"
export VAPORTRACE_DB_PATH="/path/to/vaportrace.db"

# Logging
export VAPORTRACE_LOG_LEVEL="debug"
```

---

## First Run

### Quick Start Script

```bash
#!/bin/bash

# Create config
mkdir -p ~/.vaportrace

# Initialize database
echo "init_db" | VaporTrace

# Run VaporTrace
VaporTrace
```

### Manual First Run

```bash
# 1. Launch VaporTrace
./VaporTrace

# 2. You should see:
#    ██╗   ██╗ █████╗ ██████╗  ██████╗ ██████╗ ████████╗██████╗  █████╗  ██████╗
#    ...
#    [F1] LOGS  [F2] MAP  [F3] LOOT  [F4] TRAFFIC  [F5] PLAN  [F6] NEURO  [F7] RPT

# 3. Initialize (if first run)
> init_db

# 4. Set target
> target https://api.example.com

# 5. Discover endpoints
> map

# See next document: 02_FIRST_RUN.md for hands-on walkthrough
```

---

## Verification

### Check Installation

```bash
# Verify Go
go version                              # Should be 1.25+

# Verify VaporTrace binary
./VaporTrace --version
./VaporTrace --help

# Verify can launch
./VaporTrace
> usage                                 # Should show all commands
> exit
```

### Verify Database

```bash
# Check database file
ls -la ~/.vaportrace/vaportrace.db

# Or check from SQLite
sqlite3 ~/.vaportrace/vaportrace.db ".tables"
```

### Test Network Connectivity

```bash
# Test DNS
nslookup api.example.com

# Test HTTP connectivity
curl -v https://api.example.com

# Test proxy (if configured)
curl --proxy http://localhost:8080 https://api.example.com
```

---

## Troubleshooting

### Issue: "Go not found"

**Solution:**
```bash
# Check Go installation
which go

# If not found, add to PATH
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

### Issue: "Permission denied"

**Solution:**
```bash
# Make binary executable
chmod +x ./VaporTrace

# Or use full path with ./
./VaporTrace
```

### Issue: "Cannot find package"

**Solution:**
```bash
# Download dependencies
go mod download
go mod tidy

# Then rebuild
go build -o ./VaporTrace ./main.go
```

### Issue: "Database locked"

**Solution:**
```bash
# Check for running instances
ps aux | grep VaporTrace

# Kill if needed
pkill -f VaporTrace

# Remove lock file
rm ~/.vaportrace/vaportrace.db-lock

# Restart
./VaporTrace
```

### Issue: "Port already in use"

**Solution:**
```bash
# Find process using port
sudo lsof -i :8081

# Kill process
sudo kill -9 <PID>

# Or change port in config.yaml
```

### Issue: "Network timeout"

**Solution:**
```bash
# Check network connectivity
ping api.example.com

# Increase timeout in config.yaml
timeout: 60

# Or via command
> target https://api.example.com
> timeout 60
```

---

## Next Steps

1. ✅ Installation complete
2. → Read [02_FIRST_RUN.md](02_FIRST_RUN.md) for hands-on walkthrough
3. → Read [03_UI_OVERVIEW.md](03_UI_OVERVIEW.md) for dashboard tour
4. → See [17_KEYBOARD_SHORTCUTS.md](17_KEYBOARD_SHORTCUTS.md) for hotkeys

---

## Quick Reference

| Task | Command |
|------|---------|
| Build from source | `go build -o ./VaporTrace ./main.go` |
| Launch | `./VaporTrace` |
| Init database | `> init_db` |
| View config | Edit `~/.vaportrace/config.yaml` |
| View logs | `tail -f ~/.vaportrace/vaportrace.log` |

---

**Issues?** See [16_TROUBLESHOOTING.md](16_TROUBLESHOOTING.md) for detailed help.

