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

# Sprint 2: Reconnaissance - Spec Parsing & Endpoint Mining

**Status:** ✅ COMPLETE | **Version:** v1.1-Recon | **Released:** Q4 2025

---

## 🎯 Sprint Overview

Sprint 2 implements comprehensive reconnaissance capabilities for API discovery and mapping. This sprint builds the attack surface foundation by ingesting API specifications, extracting endpoints from JavaScript bundles, and mining hidden parameters. All reconnaissance features are designed to be fully automated or manual.

**Slogan:** "Mapping the API Attack Surface"

---

## 📋 Deliverables

### 2.1: Specification Ingestion (Swagger & OpenAPI) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/discovery/swagger.go`

**Features Delivered:**
- **Swagger v2 Parsing** - Full OpenAPI 2.0 specification support
- **OpenAPI v3 Parsing** - OpenAPI 3.0+ specification support
- **Endpoint Extraction** - Automatic path + method enumeration
- **Parameter Detection** - Query, path, body parameter identification
- **Authentication Scheme Recognition** - OAuth, API Key, JWT detection
- **Request/Response Schema Mapping** - JSON schema parsing and analysis

**Implementation Highlights:**
```go
// Swagger v2/v3 parser
type SpecParser struct {
    SpecURL  string        // URL to Swagger/OpenAPI spec
    Spec     interface{}   // Parsed spec object
    Endpoints []APIEndpoint // Extracted endpoints
}

type APIEndpoint struct {
    Path        string       // e.g., /api/users/{id}
    Method      string       // GET, POST, DELETE, etc.
    Description string       // From spec
    Parameters  []Parameter  // Query, path, header
    RequestBody interface{}  // JSON schema
    Responses   []Response   // Status codes + schemas
}
```

**Supported Formats:**
- ✅ Swagger 2.0 (JSON and YAML)
- ✅ OpenAPI 3.0 (JSON and YAML)
- ✅ Local file paths (file://)
- ✅ HTTP URLs (https://)

**Example Usage:**
```bash
> swagger https://api.example.com/swagger.json
[cyan]SWAGGER:[-] Parsing OpenAPI 3.0 specification...
[green]Found 47 endpoints[-]
[yellow]Methods:[-] GET (18), POST (12), PUT (10), DELETE (7)
[yellow]Authentication:[-] OAuth 2.0, API Key, JWT
```

**Status:** ✅ Production-ready with multi-format support

---

### 2.2: JavaScript Route Scraper ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/discovery/scraper.go`

**Features Delivered:**
- **Regex-Based Endpoint Detection** - Pattern matching for API routes
- **Fetch/Axios Call Parsing** - Modern JavaScript API patterns
- **jQuery AJAX Detection** - Legacy jQuery $.ajax() patterns
- **String Literal Extraction** - Hardcoded endpoints in code
- **URL Parameter Analysis** - Query string parameter discovery
- **Relative Path Resolution** - Convert relative paths to absolute

**Regex Patterns Deployed:**
```go
// API route patterns
patterns := []string{
    `/api/[a-z0-9/_-]+`,           // /api/v1/users, /api/auth/login
    `fetch\(['"](.*?)['"]\)`,       // fetch('/api/users')
    `axios\.\w+\(['"](.*?)['"]\)`, // axios.get('/api/users')
    `\$.ajax\({.*?url:\s*['"](.*?)['"]\)`, // $.ajax({url: '/api/data'})
    `http[s]?://[a-z0-9.-]+/[a-z0-9/_-]+`,  // Full URLs
}
```

**Extraction Process:**
1. Download JavaScript bundle from target
2. Apply multi-pattern regex scanning
3. Deduplicate discovered routes
4. Resolve relative paths to base URL
5. Classify as API, static, or unknown

**Example Output:**
```bash
> scrape https://example.com/app.js
[cyan]SCRAPER:[-] Downloading app.js (245KB)...
[yellow]Discovered Routes:[-]
  /api/v1/users
  /api/v1/products
  /api/v1/orders
  /api/v1/admin/dashboard
  /static/js/vendor.js
  /css/styles.css
Total: 24 unique routes
```

**Status:** ✅ Production-ready with high accuracy rate

---

### 2.3: Version Walker (API Version Detection) ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/discovery/discovery.go`

**Features Delivered:**
- **Version String Detection** - Identify /v1/, /v2/, /beta/, etc.
- **Deprecated API Testing** - Probe for outdated endpoints
- **API Migration Tracking** - Correlate versions across endpoints
- **Version-Specific Behavior** - Test differences between versions
- **Backward Compatibility Testing** - Legacy version exploitability

**Version Detection Logic:**
```go
// Common API version patterns
versionPatterns := []string{
    `/v\d+/`,         // /v1/, /v2/, /v3/
    `/v\d+\.\d+/`,    // /v1.0/, /v2.1/
    `/beta/`,         // /beta/
    `/alpha/`,        // /alpha/
    `/staging/`,      // /staging/
    `/preview/`,      // /preview/
    `api_version=`,   // Query parameter
    `X-API-Version:`, // Header parameter
}
```

**Version Testing:**
```bash
> target https://api.example.com
> map
[yellow]Versions Detected:[-]
  /api/v1/  (deprecated, 2020)
  /api/v2/  (current, active)
  /api/v3/  (beta, bleeding-edge)
  
Testing endpoints across versions...
[green]v1 /users:[-] 200 OK (outdated auth)
[green]v2 /users:[-] 200 OK (modern auth)
[green]v3 /users:[-] 404 Not Found (pre-release)
```

**Status:** ✅ Fully implemented with version correlation

---

### 2.4: Parameter Miner ✅ COMPLETE

**Status:** ✅ Shipped  
**Location:** `pkg/discovery/miner.go`

**Features Delivered:**
- **Common Parameter Fuzzing** - Brute-force hidden parameters
- **Context-Aware Suggestions** - Based on endpoint name
- **HTTP Method-Specific** - Different params for GET vs POST
- **Header Parameter Discovery** - X-* and standard headers
- **Request Body Property Mining** - Hidden JSON properties

**Parameter Wordlist:**
```go
commonParams := []string{
    // Debug/Admin
    "debug", "admin", "test", "internal", "backdoor",
    
    // Filtering
    "filter", "sort", "order", "search", "query",
    
    // Pagination
    "page", "offset", "limit", "size", "count",
    
    // Authentication
    "token", "api_key", "secret", "password", "auth",
    
    // Bypass
    "bypass", "override", "force", "skip", "validate",
    
    // Hidden Features
    "include", "exclude", "expand", "fields", "select",
}
```

**Mining Process:**
1. Enumerate GET parameters (query string)
2. Enumerate POST parameters (request body)
3. Enumerate headers (HTTP headers)
4. Test for status changes with parameter presence
5. Flag anomalies as potential hidden functionality

**Example:**
```bash
> mine https://api.example.com/api/users
[cyan]MINER:[-] Probing 45 common parameters...
[yellow]Discovery:[-]
  ✓ page         - Controls pagination
  ✓ include      - Returns nested resources
  ✓ admin        - Reveals admin fields (HIDDEN)
  ✓ bypass       - Skips validation (CRITICAL)
  ✓ internal     - Returns internal metadata (HIGH)
```

**Status:** ✅ Production-ready with high discovery rate

---

## 🔄 Current Status by Sub-Phase

| Sub-Phase | Deliverable | Status | Completion |
|-----------|-------------|--------|------------|
| **2.1** | Swagger/OpenAPI Parser | ✅ DONE | 100% |
| **2.2** | JavaScript Route Scraper | ✅ DONE | 100% |
| **2.3** | Version Walker | ✅ DONE | 100% |
| **2.4** | Parameter Miner | ✅ DONE | 100% |

---

## 📊 Code Metrics

| Metric | Value |
|--------|-------|
| **New Files** | 4 modules (swagger, scraper, discovery, miner) |
| **Lines of Code** | ~1200 LOC |
| **Patterns/Signatures** | 40+ regex patterns |
| **API Formats Supported** | 4 (Swagger 2.0, OpenAPI 3.0, JS, custom) |
| **Test Coverage** | Spec parsing + regex tests |

---

## 🎓 Architecture Decisions

### Why Regex-Based JS Scraping?
- No JavaScript execution required
- Fast pattern matching across large bundles
- Reduces false positives with context-aware patterns
- Can be extended with custom patterns

### Version Detection Approach
- Probes multiple version paths simultaneously
- Identifies deprecated versions with known vulns
- Enables cross-version comparison testing
- Helps identify API migration patterns

### Parameter Mining Strategy
- Uses common parameter dictionary
- Context-sensitive based on endpoint structure
- Detects both expected and anomalous behavior
- Identifies debug/admin parameters automatically

---

## 🚀 Next Steps

Sprint 3 implements authentication and BOLA/BFLA engines:
- Session management and token storage
- Object-level authorization testing
- Function-level authorization bypassing
- Role-based access control probing

---

## 📚 References

- **OpenAPI Specification:** https://spec.openapis.org/
- **Swagger Tools:** https://swagger.io/tools/
- **Go Regexp:** https://pkg.go.dev/regexp
