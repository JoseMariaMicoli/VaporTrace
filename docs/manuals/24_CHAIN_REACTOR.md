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

# 24. Chain Reactor - Tier 4 Day 2

**Version:** 1.0  
**Tier:** Tier 4 (Advanced)  
**Status:** Sprint 20  
**Feature:** Stateful multi-step attack orchestration with variable extraction

---

## Overview

The **Chain Reactor** is a sophisticated workflow engine enabling multi-step attack sequences where each step depends on data extracted from the previous step. It's the foundation of advanced exploitation requiring authentication flows, state management, and conditional logic.

### Key Features

- **Stateful Execution** - Maintains context across steps
- **Variable Extraction** - Extract and reuse data between steps
- **Header Injection** - Automatically inject extracted values into headers
- **Flexible HTTP Methods** - Support for GET, POST, PUT, PATCH, DELETE
- **Real-time Monitoring** - Watch execution progress in F1 logs

---

## Core Concepts

### Chains

A chain is a named sequence of HTTP steps with state management. Each step can:
- Execute an HTTP request
- Capture data from the response
- Pass extracted data to subsequent steps

### Variables

Variables are stored as `{{var_name}}` and automatically injected into:
- Request headers
- Request body
- URL parameters

### Steps

Each step in a chain is a complete HTTP request with:
- HTTP method (GET, POST, etc.)
- URL
- Body (optional)
- Headers (optional)
- Extraction rules (optional)

---

## Commands

### chain create

Create a new, empty chain definition.

**Syntax:**
```
chain create <name>
```

**Examples:**
```bash
chain create authentication_flow
chain create csrf_bypass
chain create privilege_escalation
```

**Output:**
```
[blue]Creating chain: authentication_flow[-]
[yellow]Chain created. Use 'chain add' to define steps.[-]
```

---

### chain add

Add an HTTP step to an existing chain.

**Syntax:**
```
chain add <chain_name> <method> <url> [body]
```

**Parameters:**
- `chain_name` - Name of the chain (must exist)
- `method` - HTTP method: GET, POST, PUT, PATCH, DELETE
- `url` - Full URL (can include variables like {{token}})
- `body` - (Optional) Request body for POST/PUT/PATCH

**Examples:**

```bash
# Step 1: Login
chain add login_flow POST https://api.example.com/login '{"user":"admin","pass":"1234"}'

# Step 2: Get profile (uses token from Step 1)
chain add login_flow GET https://api.example.com/profile

# Step 3: Change password (injects token into header)
chain add login_flow PUT https://api.example.com/profile/password '{"new":"newpass123"}'
```

**Output:**
```
[green]Added step to chain login_flow: POST https://api.example.com/login[-]
```

---

### chain extract

Configure data extraction from a step's response.

**Syntax:**
```
chain extract <chain_name> <step_index> <var_name> <type> <selector>
```

**Parameters:**
- `chain_name` - Name of the chain
- `step_index` - Step number (1-based)
- `var_name` - Variable name to store as (e.g., `token`)
- `type` - Extraction type: `json`, `regex`, `header`, `cookie`
- `selector` - Pattern/path for extraction

**Examples:**

```bash
# Extract access_token from JSON response
chain extract login_flow 1 token json $.access_token

# Extract CSRF token from Set-Cookie header
chain extract login_flow 1 csrf cookie XSRF-TOKEN

# Extract ID using regex
chain extract login_flow 2 user_id regex 'user[_]?id["\s:=]+(\d+)'

# Extract from response header
chain extract login_flow 1 auth_header header Authorization
```

**Stored Variables:**
```
Step 1 response: {"access_token": "abc123", "user_id": 42}

After extraction:
{{token}}     = "abc123"
{{user_id}}   = "42"
```

---

### chain header

Add or modify HTTP headers for a specific step.

**Syntax:**
```
chain header <chain_name> <step_index> <key> <value>
```

**Parameters:**
- `chain_name` - Name of the chain
- `step_index` - Step number (1-based)
- `key` - Header name (e.g., Authorization, Content-Type)
- `value` - Header value (supports variables like {{token}})

**Examples:**

```bash
# Add Authorization header with extracted token
chain header login_flow 2 Authorization "Bearer {{token}}"

# Add custom header
chain header login_flow 3 X-CSRF-Token "{{csrf}}"

# Add Admin header
chain header login_flow 3 X-Admin-Mode "true"
```

---

### chain run

Execute the chain in sequence.

**Syntax:**
```
chain run <chain_name>
```

**Examples:**
```bash
chain run authentication_flow
chain run csrf_bypass
```

**Execution Flow:**
```
1. Execute Step 1 (POST /login)
   ↓ Response: {"access_token": "xyz"}
   ✓ Extract {{token}} = "xyz"
   
2. Execute Step 2 (GET /profile)
   ↑ Inject header: Authorization: Bearer {{token}}
   ↓ Response: {"user_id": 42, "permissions": [...]}
   ✓ Extract {{user_id}} = "42"
   
3. Execute Step 3 (PUT /admin)
   ↑ Inject: Authorization, X-Admin-Mode headers
   ↓ Response: {"status": "success"}
```

**Output:**
```
[blue]Executing chain: authentication_flow[-]
[yellow]Chain execution started in background.[-]
[Step 1] POST /login → 200 OK
[Step 2] GET /profile → 200 OK, extracted {{user_id}}=42
[Step 3] PUT /admin → 403 Forbidden (escalation failed)
```

---

### chain list

List all defined chains and their steps.

**Syntax:**
```
chain list
```

**Output:**
```
[magenta]ACTIVE CHAINS:[-]
  - authentication_flow (3 steps)
  - csrf_bypass (2 steps)
  - privilege_escalation (4 steps)
```

---

## Complete Workflow Example

### Authentication Flow

```bash
# 1. Create chain
chain create auth_flow

# 2. Add login step
chain add auth_flow POST https://api.example.com/auth/login '{"email":"test@example.com","password":"SecurePass123"}'

# 3. Extract JWT token from response
chain extract auth_flow 1 access_token json $.data.access_token

# 4. Add step to access protected endpoint
chain add auth_flow GET https://api.example.com/api/users/me

# 5. Inject token into Authorization header
chain header auth_flow 2 Authorization "Bearer {{access_token}}"

# 6. Execute the chain
chain run auth_flow
```

### CSRF Bypass Flow

```bash
# 1. Create chain
chain create csrf_bypass

# 2. Get CSRF token from login page
chain add csrf_bypass GET https://app.example.com/login

# 3. Extract CSRF token from HTML form
chain extract csrf_bypass 1 csrf_token regex '<input[^>]*name=["\']_token["\'][^>]*value=["\']([^"\']+)["\']'

# 4. Extract session cookie
chain extract csrf_bypass 1 session_id cookie PHPSESSID

# 5. Submit form with CSRF token
chain add csrf_bypass POST https://app.example.com/action '{"_token":"{{csrf_token}}","action":"delete_user"}'

# 6. Inject session cookie
chain header csrf_bypass 2 Cookie "PHPSESSID={{session_id}}"

# 7. Execute
chain run csrf_bypass
```

### Privilege Escalation Flow

```bash
# 1. Create chain
chain create priv_esc

# 2. Login as regular user
chain add priv_esc POST https://api.example.com/login '{"user":"regular","pass":"pass123"}'
chain extract priv_esc 1 user_token json $.token

# 3. Get user info
chain add priv_esc GET https://api.example.com/users/me
chain header priv_esc 2 Authorization "Bearer {{user_token}}"
chain extract priv_esc 2 user_id json $.id

# 4. Attempt to modify admin field (BOPLA)
chain add priv_esc PUT https://api.example.com/users/{{user_id}} '{"is_admin":true,"role":"administrator"}'
chain header priv_esc 3 Authorization "Bearer {{user_token}}"

# 5. Get admin actions endpoint
chain add priv_esc GET https://api.example.com/admin/actions
chain header priv_esc 4 Authorization "Bearer {{user_token}}"

# 6. Execute
chain run priv_esc
```

---

## Integration with Tier 4 Features

### Using Extractors with Chains

```bash
# Configure extractor for complex data
extract config json $.nested.data.token

# Chain uses extracted values automatically
chain add my_chain POST url '{"token":"{{extracted_value}}"}'
```

### Using with OSINT Results

```bash
# Discover endpoint via OSINT
intel wayback api.example.com
# Result: /api/v1/admin exists in 2021 snapshot

# Create chain to test legacy endpoint
chain create legacy_test
chain add legacy_test GET https://api.example.com/api/v1/admin
chain run legacy_test
```

---

## Advanced Features

### Conditional Steps (Future)

While not yet implemented, future releases will support:

```bash
# Run Step 3 only if Step 2 returned 200
chain add flow GET /endpoint1
chain add flow GET /endpoint2 --if "status=200"
chain add flow POST /endpoint3
```

### Looping (Future)

```bash
# Repeat step with different variables
chain add enumeration GET /users/{{user_id}} --for "user_id in 1..100"
```

### Response Assertions (Future)

```bash
# Fail chain if condition not met
chain add test POST /login --assert "status=200"
chain add test GET /admin --assert "contains='admin'"
```

---

## Troubleshooting

### Common Issues

**Issue:** "Chain not found"
```
Solution: Use 'chain list' to see defined chains
Solution: Chain names are case-sensitive
```

**Issue:** Variable not injected
```
Solution: Ensure variable name matches exactly: {{exact_name}}
Solution: Verify extraction was successful in Step N
```

**Issue:** Step returns 401 Unauthorized
```
Solution: Verify token extraction is working
Solution: Check Authorization header format: "Bearer {{token}}"
```

**Issue:** Response extraction fails
```
Solution: Test selector/regex independently
Solution: Log response body to verify structure
Solution: Use different extraction type (json vs regex)
```

---

## Best Practices

1. **Test Each Step** - Verify each request independently first
2. **Name Variables Clearly** - Use descriptive names like `access_token`, `user_id`, not just `token1`
3. **Verify Extraction** - Ensure extracted values are actually needed
4. **Document Chains** - Comment what each step does
5. **Test Edge Cases** - What if token expires mid-chain?
6. **Use Unique Chain Names** - Avoid conflicts with descriptive naming

---

## Performance Considerations

- **Parallel Chains** - Create multiple chains to test different paths
- **Batch Testing** - Use chains for repeated patterns (user enumeration, etc.)
- **Timeout Management** - Chains continue even if a step times out
- **Resource Limits** - Consider Tier 3 race conditions for concurrent chains

---

## See Also

- [25_EXTRACTOR.md](25_EXTRACTOR.md) - Data extraction details
- [23_INTEL_OSINT.md](23_INTEL_OSINT.md) - Finding endpoints to chain
- [09_ATTACK_CHAINS.md](09_ATTACK_CHAINS.md) - Orchestrated attack sequences
- [06_EXPLOITATION.md](06_EXPLOITATION.md) - Individual attack modules
