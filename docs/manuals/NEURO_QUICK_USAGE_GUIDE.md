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

# VaporTrace Neuro Engine - Quick Usage Guide

## 🚀 Complete Neuro Workflow

This guide shows how to use all neuro functionalities together for a complete attack flow.

---

## 📋 Available Commands

```
neuro on                                     - Enable neural engine
neuro off                                    - Disable neural engine
neuro config <provider> <model> [api_key] [endpoint]  - Configure AI provider
neuro-gen <context> [count]                 - Generate attack payloads
test-neuro                                  - Test neural engine connectivity
ask <prompt>                                 - Direct LLM query
```

---

## 🎯 Complete Workflow Example

### Step 1: Enable Neural Engine

```bash
> neuro on
[green]Neural Engine Activated.[-]
```

### Step 2: Configure AI Provider (Choose One)

#### Option A: Groq (Fastest - Recommended)
```bash
> neuro config groq llama-3.1-8b gsk_xxxxxxxxxxxxx
[cyan]NEURO:[-] Groq LPU Cloud Selected.
[green]✓ CONFIG:[-] Groq API configured
```

#### Option B: OpenAI (Most Capable)
```bash
> neuro config openai gpt-4o sk_proj_xxxxxxxxxxxxx
[cyan]NEURO:[-] OpenAI Cloud Selected.
[green]✓ CONFIG:[-] OpenAI API configured
```

#### Option C: Google Gemini
```bash
> neuro config google gemini-1.5-flash AIzaxxxxxxxxxxxxx
[cyan]NEURO:[-] Google Gemini Selected (gemini-1.5-flash).
[green]✓ CONFIG:[-] Gemini API configured
```

#### Option D: Local Ollama (No API Key Needed)
```bash
> neuro config ollama mistral
[yellow]NEURO:[-] Warning - Local Inference uses significant RAM. Ensure 8GB+ avail.
[green]✓ CONFIG:[-] Ollama local configured
```

#### Option E: Hybrid Mode (Cloud Primary + Local Fallback)
```bash
> neuro config hybrid gpt-4o sk_proj_xxxxxxxxxxxxx
[magenta]NEURO:[-] Hybrid Brain Activated. Primary: OpenAI | Fallback: Ollama.
[green]✓ CONFIG:[-] Hybrid mode configured
```

**Syntax:** `neuro config <provider> <model> [api_key] [endpoint]`

**Parameters:**
- `provider` - One of: `groq`, `openai`, `google`/`gemini`, `ollama`, `hybrid`
- `model` - Model name (defaults to fastest for each provider)
- `api_key` - (Optional) API key from provider
- `endpoint` - (Optional) Custom endpoint URL

---

### Step 3: Test Connectivity

```bash
> test-neuro
[blue]NEURO:[-] Sending heartbeat packet to Brain...
[green]CONNECTIVITY CHECK:[-] Pong!
[green]NEURO ONLINE:[-] Check Neural Tab.
```

This verifies:
- ✓ Primary provider is reachable
- ✓ API key is valid
- ✓ Fallback mechanism works (if one fails)

---

### Step 4: Generate Attack Payloads

#### Basic Payload Generation
```bash
> neuro-gen "SQL injection on /api/search parameter"
[cyan]PAYLOAD_GEN:[-] Requesting AI payloads...
[blue]GENERATING:[-] 10 unique SQL injection payloads...
[yellow]PAYLOAD 1:[-] 1' UNION SELECT NULL,username,password FROM users--
[yellow]PAYLOAD 2:[-] 1' AND 1=1 UNION ALL SELECT 1,2,3,4,5--
[yellow]PAYLOAD 3:[-] 1' AND (SELECT COUNT(*) FROM information_schema.tables)>0--
...
[cyan]PAYLOADS:[-] Generated 10 payloads ready for testing.
```

#### With Custom Count
```bash
> neuro-gen "XSS in user profile" 15
[cyan]PAYLOAD_GEN:[-] Generating 15 unique XSS payloads...
```

---

### Step 5: Interactive LLM Query

```bash
> ask "How do I bypass JWT signature verification?"
[cyan]AI:[-] Let me analyze your request...
[blue]RESPONSE:[-]
"To bypass JWT signature verification:
1. Token manipulation (change claims)
2. JWKS key confusion
3. Algorithm swap (HS256 vs RS256)
4. Expired token bypass (remove exp claim)
5. Weak signature key brute force

I recommend testing JWT signature verification first.
Would you like me to generate JWT manipulation payloads?"
```

---

### Step 6: Automated Analysis in Scan

When running other commands like `scan`, `pipeline`, or `bfla`, neuro automatically:

```bash
> scan https://target.com

[cyan]TARGET:[-] https://target.com
[blue]⠋[-] Engaging NEURO engine for analysis...

[yellow]PAYLOAD 1:[-] Testing /api/users/1 with BOLA payload
  Response: 200 OK - User data returned
  [red]✗ VULNERABLE:[-] BOLA detected (unauthorized resource access)

[yellow]PAYLOAD 2:[-] Testing /api/users/2 with modified role
  Response: 403 Forbidden
  [green]✓ PROTECTED:[-] Input validation working

[cyan]NEURO ANALYSIS:[-] 8 potential vulnerabilities identified
[magenta]AI RECOMMENDATION:[-] Focus on role-based authorization bypass
```

---

## 🔄 Common Workflows

### Workflow 1: Quick Setup & Test

```bash
> neuro on
> neuro config groq llama-3.1-8b gsk_xxxxx
> test-neuro
[green]✓ Ready to analyze
```

### Workflow 2: Generate Custom Payloads

```bash
> neuro on
> neuro config openai gpt-4o sk_xxxxx
> neuro-gen "Authentication bypass for API key validation" 20
> (Copy payloads to attack tools)
```

### Workflow 3: Interactive Problem Solving

```bash
> neuro on
> ask "What are the best techniques to bypass rate limiting?"
[AI Response with multiple techniques]
> ask "Generate payloads for HTTP/2 connection reset"
[AI Response with specific payloads]
```

### Workflow 4: Full Attack Analysis

```bash
> target https://target.com
> neuro on
> neuro config hybrid gpt-4o sk_xxxxx
> map
[Discovery runs with neuro analysis]
> pipeline
[Automatically uses neuro for attack generation and analysis]
> analyze
[Strategic planning with AI insights]
> commit
[Execute with AI-optimized payloads]
```

---

## 🧠 Provider Comparison

| Provider | Speed | Cost | Quality | Best For |
|----------|-------|------|---------|----------|
| **Groq** | ⭐⭐⭐⭐⭐ | 💰 Free | Very Good | Quick tests, free tier |
| **OpenAI** | ⭐⭐⭐ | 💰💰💰 | Excellent | Production, complex analysis |
| **Gemini** | ⭐⭐⭐⭐ | 💰💰 | Very Good | Good balance |
| **Ollama** | ⭐⭐ | 💰 Free | Good | Privacy, no internet |
| **Hybrid** | ⭐⭐⭐⭐ | 💰 Varies | Excellent | Reliability, fallback |

---

## 💡 Tips & Best Practices

### Tip 1: Use Hybrid Mode for Reliability
```bash
> neuro config hybrid gpt-4o sk_xxxxx
```
If OpenAI API has issues, automatically falls back to local Ollama.

### Tip 2: Test Connectivity After Configuration
```bash
> neuro config groq llama-3.1-8b gsk_xxxxx
> test-neuro
# Verify API key is correct before using in attacks
```

### Tip 3: Generate Context-Specific Payloads
```bash
> neuro-gen "Authentication bypass on OAuth2 implicit flow" 25
# Specific context = better payloads
```

### Tip 4: Use Ask for Research
```bash
> ask "What are the latest vulnerabilities in Spring Boot?"
[AI provides current information]
> ask "Generate payloads for CVE-2024-xxxxx"
[AI generates specific exploits]
```

### Tip 5: Combine with Other Commands
```bash
> map          # Find endpoints
> neuro-gen "..." # Generate payloads for found endpoints
> bfla         # Test with generated payloads
> report       # AI summarizes findings
```

---

## 🐛 Troubleshooting

### Issue: "Neuro Engine Inactive"
```
❌ [yellow]NEURO:[-] Engine inactive. Auto-starting Hybrid mode...
✓ Fixed: Run `neuro on` first
```

### Issue: "Cloud Brain Quota/Rate-Limit (429) Hit"
```
❌ [red]NEURO:[-] Cloud Brain Quota/Rate-Limit (429) Hit.
✓ Solution: 
  - Switch to Ollama (local): `neuro config ollama mistral`
  - Use different provider: `neuro config gemini gemini-1.5-flash AIza...`
  - Wait 1 minute and retry
```

### Issue: "Invalid provider"
```
❌ [red]ERROR:[-] Invalid provider specified
✓ Check: Use one of: groq, openai, google, gemini, ollama, hybrid
> neuro config groq llama-3.1-8b gsk_xxxxx
```

### Issue: "Connection failed"
```
❌ [red]NEURO:[-] connection failed: connection refused
✓ For Ollama: Ensure running `ollama serve` in another terminal
✓ For Cloud: Check internet connection and API key
```

---

## 📊 Example: Complete Attack with Neuro

```bash
# 1. Setup
> target https://api.example.com
> neuro on
> neuro config groq llama-3.1-8b gsk_xxxxx

# 2. Discovery
> map
[Discovers 45 endpoints using swagger, scraping, etc.]

# 3. Generate Payloads
> neuro-gen "API authentication bypass" 20
[20 unique payloads generated by AI]

# 4. Strategic Planning
> analyze
[AI analyzes endpoints and creates attack plan]

# 5. Execute Pipeline
> pipeline
[Runs BOLA, BFLA, BOPLA with AI-generated payloads]

# 6. Get Insights
> ask "What are the top 3 vulnerabilities found?"
[AI summarizes findings and recommendations]

# 7. Report
> report
[Generates report with AI-powered analysis]
```

---

## ⚙️ Configuration Details

### Groq Setup
```bash
1. Visit: https://console.groq.com
2. Create API key
3. Run: neuro config groq llama-3.1-8b gsk_xxxxx
```

### OpenAI Setup
```bash
1. Visit: https://platform.openai.com/api-keys
2. Create API key
3. Run: neuro config openai gpt-4o sk_proj_xxxxx
```

### Google Gemini Setup
```bash
1. Visit: https://makersuite.google.com/app/apikey
2. Create API key
3. Run: neuro config google gemini-1.5-flash AIza_xxxxx
```

### Ollama Setup (Local)
```bash
1. Install: https://ollama.ai
2. Run in terminal: ollama serve
3. In VaporTrace: neuro config ollama mistral
4. (No API key needed - runs locally)
```

---

## 🎓 Learning Path

**Beginner:**
1. `neuro on`
2. `neuro config ollama mistral`
3. `test-neuro`
4. `ask "What is SQL injection?"`

**Intermediate:**
1. Setup Groq API (free)
2. `neuro-gen "authentication bypass"` 
3. Combine with `scan` command
4. Review payloads in Neural Tab

**Advanced:**
1. Use `hybrid` mode for redundancy
2. Combine with `pipeline` for automated attacks
3. Use `ask` for advanced technique research
4. Create custom workflows combining multiple engines

---

**Ready to use!** Start with: `neuro on` → `neuro config groq llama-3.1-8b <YOUR_API_KEY>`
