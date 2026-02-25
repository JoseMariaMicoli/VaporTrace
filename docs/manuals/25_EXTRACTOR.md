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

# 25. Value Extractor - Tier 4 Day 2

**Version:** 1.0  
**Tier:** Tier 4 (Advanced)  
**Status:** Sprint 20  
**Feature:** Flexible data extraction from HTTP responses for chaining and variable management

---

## Overview

The **Value Extractor** is a powerful data extraction engine that pulls structured and unstructured data from HTTP responses. Extracted values are stored as variables that can be injected into subsequent requests, making it essential for stateful attack chains.

### Key Features

- **Multiple Extraction Methods** - JSON, Regex, Cookies, Headers
- **Flexible Selectors** - JSONPath, regex patterns, header names
- **Variable Storage** - Automatic storage as `{{var_name}}`
- **Chaining Support** - Seamless integration with Chain Reactor
- **Error Handling** - Graceful failures with logging

---

## Extraction Types

### JSON Extraction

Extract values from JSON responses using JSONPath expressions.

**Syntax:**
```
extract config json <jsonpath>
```

**JSONPath Syntax:**

```
$.field              - Root-level field
$.field.nested       - Nested field
$.array[0]           - First element of array
$.array[*].id        - All IDs from array of objects
$.data.*.email       - All email fields under data
```

**Examples:**

```bash
# Simple field
extract config json $.access_token
# Response: {"access_token": "abc123"}
# Result: {{access_token}} = "abc123"

# Nested field
extract config json $.data.user.id
# Response: {"data": {"user": {"id": 42}}}
# Result: {{user_id}} = "42"

# Array element
extract config json $.tokens[0]
# Response: {"tokens": ["xyz", "abc"]}
# Result: {{first_token}} = "xyz"

# Multiple fields
extract config json $.response.auth.token
extract config json $.response.auth.refresh
# Extracts both token and refresh as separate variables
```

**Use Cases:**
- API responses with complex JSON structures
- OAuth/JWT token extraction
- User IDs, session IDs, resource IDs
- Nested configuration data

---

### Regex Extraction

Extract data using regular expressions for unstructured or HTML responses.

**Syntax:**
```
extract config regex <pattern>
```

**Pattern Groups:**

Use capturing groups `()` to extract specific parts:

```
pattern           - Entire match
(group1)(group2)  - Multiple captures
(.*?)             - Non-greedy capture
(?P<name>...)     - Named group
```

**Examples:**

```bash
# Extract user ID from HTML
extract config regex 'user[_]?id["\s:=]+(\d+)'
# HTML: <input type="hidden" name="user_id" value="12345">
# Result: {{user_id}} = "12345"

# Extract CSRF token
extract config regex 'name=["\']_token["\'][^>]*value=["\']([^"\']+)["\']'
# HTML: <input name="_token" value="abc123def456">
# Result: {{csrf_token}} = "abc123def456"

# Extract from URL in response
extract config regex 'redirect=([^&]+)'
# URL: https://app.com/callback?redirect=https%3A%2F%2Ftarget.com
# Result: {{redirect_url}} = "https%3A%2F%2Ftarget.com"

# Extract multiple values
extract config regex '(?P<token>\w+)|(?P<id>\d+)'
# Extracts both token and id into separate variables
```

**Use Cases:**
- HTML form data extraction
- Hidden input values
- Embedded URLs
- API keys in script tags
- Error messages containing data

---

### Cookie Extraction

Extract values from `Set-Cookie` headers.

**Syntax:**
```
extract config cookie <cookie_name>
```

**Cookie Format:**
```
Set-Cookie: name=value; Path=/; HttpOnly; Secure
```

**Examples:**

```bash
# Extract session cookie
extract config cookie PHPSESSID
# Header: Set-Cookie: PHPSESSID=abc123def456; Path=/
# Result: {{phpsessid}} = "abc123def456"

# Extract CSRF token from cookie
extract config cookie XSRF-TOKEN
# Header: Set-Cookie: XSRF-TOKEN=xyz789; Path=/api; Secure
# Result: {{csrf}} = "xyz789"

# Extract auth token
extract config cookie auth_token
# Header: Set-Cookie: auth_token=bearer_abc123
# Result: {{auth_token}} = "bearer_abc123"
```

**Use Cases:**
- Session cookies
- CSRF tokens
- Auth tokens
- Tracking cookies
- Any value in Set-Cookie headers

---

### Header Extraction

Extract values from HTTP response headers (non-cookie).

**Syntax:**
```
extract config header <header_name>
```

**Examples:**

```bash
# Extract authentication token from header
extract config header X-Auth-Token
# Response: X-Auth-Token: abc123xyz789
# Result: {{auth_token}} = "abc123xyz789"

# Extract location for redirect
extract config header Location
# Response: Location: https://api.example.com/callback?code=auth123
# Result: {{redirect}} = "https://api.example.com/callback?code=auth123"

# Extract server info
extract config header Server
# Response: Server: nginx/1.21.0
# Result: {{server}} = "nginx/1.21.0"
```

**Use Cases:**
- Custom authentication headers
- Redirect URLs
- Server/technology identification
- API versioning headers
- Rate limit information

---

## Command Reference

### extract list

List all available extraction methods and examples.

**Syntax:**
```
extract list
```

**Output:**
```
[magenta]ACTIVE EXTRACTORS:[-]
  JSON Path Extractor: Extract values from JSON responses
  Regex Extractor: Extract using regular expressions
  Cookie Extractor: Capture cookies from Set-Cookie headers
  Header Extractor: Extract custom HTTP headers
```

---

### extract config

Configure an extractor with a specific pattern.

**Syntax:**
```
extract config <type> <pattern>
```

**Examples:**
```bash
extract config json $.access_token
extract config regex 'user_id["\s:=]+(\d+)'
extract config cookie session_id
extract config header Authorization
```

**Output:**
```
[green]Extractor configured: json[-]
```

---

### extract run

Execute the configured extractor on the last response.

**Syntax:**
```
extract run
```

**Output:**
```
[yellow]Running extraction...[-]
[green]✓ Extraction complete. Results stored in context.[-]
Variable: {{access_token}} = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Integration with Chain Reactor

### Using Extractors in Chains

```bash
# 1. Create chain
chain create auth_chain

# 2. Add login step
chain add auth_chain POST https://api.example.com/login '{"user":"admin","pass":"test"}'

# 3. Configure extractor for JSON
extract config json $.token

# 4. Execute chain
chain run auth_chain

# 5. Extracted token available as {{token}} for next step
chain header auth_chain 2 Authorization "Bearer {{token}}"
```

### Chaining Multiple Extractors

```bash
# Step 1: Get auth token
chain add flow POST /auth '{"creds"}'
chain extract flow 1 token json $.access_token

# Step 2: Get user ID
chain add flow GET /user
chain header flow 2 Authorization "Bearer {{token}}"
chain extract flow 2 user_id json $.id

# Step 3: Modify user with both variables
chain add flow PUT /users/{{user_id}} '{"admin":true}'
chain header flow 3 Authorization "Bearer {{token}}"

# Execute complete chain
chain run flow
```

---

## Advanced Extraction Patterns

### Extract from Nested Arrays

```bash
# Response: {"users": [{"id": 1, "token": "a"}, {"id": 2, "token": "b"}]}

# Get all tokens
extract config json $.users[*].token

# Get first user's token
extract config json $.users[0].token

# Get user with specific ID
extract config json $.users[?(@.id==1)].token
```

### Extract from Malformed JSON

Use regex fallback for incomplete/malformed JSON:

```bash
# Response may be incomplete: {"token": "abc123...
# Use regex instead
extract config regex '"token"\s*:\s*"([^"]+)'
```

### Multi-Step Extraction

Extract data in steps, building on previous extractions:

```bash
# Step 1: Extract raw token
extract config regex 'token["\s:=]+([a-zA-Z0-9]+)'

# Step 2: Later, extract specific claim from token
# (after decoding JWT)
extract config json $.claims.user_id
```

---

## Common Extraction Scenarios

### OAuth/OIDC Flow

```bash
# 1. Get authorization code
chain add oauth GET https://auth.example.com/authorize?...
chain extract oauth 1 auth_code regex 'code=([a-zA-Z0-9-_]+)'

# 2. Exchange code for token
chain add oauth POST https://auth.example.com/token '{"code":"{{auth_code}}"}'
chain extract oauth 2 access_token json $.access_token
chain extract oauth 2 refresh_token json $.refresh_token

# 3. Use token for API access
chain add oauth GET https://api.example.com/me
chain header oauth 3 Authorization "Bearer {{access_token}}"
```

### Multi-Step API Testing

```bash
# 1. Create resource
chain add api POST /resources '{"name":"test"}'
chain extract api 1 resource_id regex '"id":\s*(\d+)'

# 2. Get resource details
chain add api GET /resources/{{resource_id}}
chain extract api 2 status json $.status

# 3. Update resource
chain add api PUT /resources/{{resource_id}} '{"status":"active"}'
chain header api 3 X-Request-ID "{{resource_id}}"

# 4. Delete resource
chain add api DELETE /resources/{{resource_id}}
```

### CSRF Token Rotation

```bash
# 1. Get initial CSRF token
chain add csrf_test GET https://app.example.com/form
chain extract csrf_test 1 csrf_token regex 'name=["\']_token["\'][^>]*value=["\']([^"\']+)["\']'

# 2. First form submission
chain add csrf_test POST https://app.example.com/action '{"_token":"{{csrf_token}}"}'

# 3. Get new CSRF token for second action
chain add csrf_test GET https://app.example.com/form2
chain extract csrf_test 3 csrf_token2 regex 'name=["\']_token["\'][^>]*value=["\']([^"\']+)["\']'

# 4. Second form submission
chain add csrf_test POST https://app.example.com/action2 '{"_token":"{{csrf_token2}}"}'
```

---

## Error Handling

### Extraction Failures

**Issue:** Pattern not found in response
```
Action: Check response body structure
Action: Verify pattern is correct for response format
Action: Test pattern independently with sample response
```

**Issue:** Multiple matches in response
```
Action: Refine pattern to be more specific
Action: Use first match if all values are equivalent
Action: Configure separate extractors for different patterns
```

**Issue:** Extracted value is empty/null
```
Action: Verify field exists in response
Action: Check for special characters or encoding
Action: Try alternative extraction type (JSON vs Regex)
```

---

## Performance Tips

1. **Use Specific JSONPath** - `$.user.id` faster than `$..*`
2. **Compile Regex Patterns** - Complex patterns evaluated once per response
3. **Limit Captures** - Each capture group adds processing
4. **Cache Patterns** - Store extractor configs for reuse
5. **Test Extraction** - Verify patterns work before chaining

---

## See Also

- [24_CHAIN_REACTOR.md](24_CHAIN_REACTOR.md) - Using extractors with chains
- [23_INTEL_OSINT.md](23_INTEL_OSINT.md) - Extracting data from OSINT results
- [09_ATTACK_CHAINS.md](09_ATTACK_CHAINS.md) - Building attack chains
- [06_EXPLOITATION.md](06_EXPLOITATION.md) - Individual attack modules
