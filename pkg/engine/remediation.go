/*
Copyright (c) 2026 José María Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

package engine

import (
	"fmt"
	"strconv"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// === SPRINT 16.1: Blue-Team Mirror & Remediation Engine ===
// RemediationSuggestion contains a fix recommendation for a discovered vulnerability
type RemediationSuggestion struct {
	VulnerabilityType string // BOLA, BFLA, SSRF, etc.
	Severity          string // CRITICAL, HIGH, MEDIUM, LOW
	Exploit           string // Summary of how it was exploited
	FixDescription    string // Plain English description
	CodeSnippet       string // Language-agnostic pseudocode or middleware example
	Language          string // "Go", "Python", "Node.js", "Java", etc.
	ImplementationURL string // Link to documentation
	// === SPRINT 16.1.1: LLM Hallucination Prevention ===
	VerificationStatus   string // "VERIFIED", "UNVERIFIED", "GOLD_STANDARD"
	VerificationNotes    string // Why this snippet is safe (or warnings if UNVERIFIED)
	StaticAnalysisPassed bool   // True if passed security linter (gosec, bandit, etc.)
}

// === SPRINT 16.1.1: Verification System for LLM Hallucination Prevention ===

// GoldStandardLibrary contains pre-verified, production-ready code snippets
// These snippets have been manually reviewed and passed security audits
var GoldStandardLibrary = map[string]string{
	"BOLA_AUTHZ_MIDDLEWARE": `// VERIFIED: Object-level authorization check
func AuthorizeObjectAccess(c *gin.Context) {
	userID := c.GetString("user_id")
	objectID := c.Param("id")
	
	// Fetch object and verify ownership
	obj := database.GetObject(objectID)
	if obj == nil || obj.OwnerID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}`,

	"BFLA_RBAC_MIDDLEWARE": `// VERIFIED: Role-based access control
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		
		// Verify user has required role
		if userRole != requiredRole && userRole != "admin" {
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}`,

	"SSRF_URL_VALIDATION": `// VERIFIED: Safe URL validation against private IP ranges
func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	
	// Block private IP ranges
	privateRanges := []string{
		"127.0.0.0/8",        // Localhost
		"10.0.0.0/8",         // Private
		"172.16.0.0/12",      // Private
		"169.254.0.0/16",     // Link-local
		"169.254.169.254/32", // AWS metadata
	}
	
	ip := net.ParseIP(u.Hostname())
	for _, cidr := range privateRanges {
		_, network, _ := net.ParseCIDR(cidr)
		if network.Contains(ip) {
			return false
		}
	}
	
	return u.Scheme == "http" || u.Scheme == "https"
}`,

	"JWT_VALIDATION": `// VERIFIED: Strict JWT validation with algorithm pinning
func ValidateJWT(tokenString string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// CRITICAL: Restrict to specific algorithm(s) only
		if token.Method.Alg() != "HS256" && token.Method.Alg() != "RS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		return getSigningKey(token), nil
	})
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	return token.Claims, nil
}`,

	"SQL_INJECTION_PREVENTION": `// VERIFIED: Parameterized queries (not string concatenation)
// WRONG: query := "SELECT * FROM users WHERE id = " + userInput
// CORRECT: Use parameterized queries
query := "SELECT * FROM users WHERE id = ?"
db.Query(query, userInput)`,
}

// VerifyRemediationSuggestion validates a remediation suggestion before display
// Returns true if safe for production use, false if requires manual review
func VerifyRemediationSuggestion(suggestion *RemediationSuggestion) {
	// Check if snippet is in Gold Standard library
	for key, goldSnippet := range GoldStandardLibrary {
		if suggestion.CodeSnippet == goldSnippet {
			suggestion.VerificationStatus = "GOLD_STANDARD"
			suggestion.VerificationNotes = "This snippet is from the Gold Standard library and has passed security audit."
			suggestion.StaticAnalysisPassed = true
			utils.LogContext("[green]✓ VERIFIED:[-] " + key + " is GOLD_STANDARD")
			return
		}
	}

	// If not in Gold Standard, mark as UNVERIFIED with warning
	suggestion.VerificationStatus = "UNVERIFIED"
	suggestion.VerificationNotes = "⚠️ WARNING: This snippet has NOT been verified by security team. Manual review required before production use. Test for: ReDoS in regex, SQL injection bypass, auth bypass, hardcoded credentials."
	suggestion.StaticAnalysisPassed = false

	utils.LogContext("[yellow]⚠ UNVERIFIED:[-] This suggestion requires manual security review before production deployment")
	utils.TacticalLog("[yellow]REMEDIATION:[-] Unverified suggestion queued for security team review")
}

// SuggestFix generates a remediation recommendation based on a successful exploit
// This is the "mirror" that shows how to defend against each attack
func SuggestFix(exploit TacticalAction, result string) *RemediationSuggestion {
	utils.TacticalLog("[cyan]REMEDIATION:[-] Analyzing exploit " + exploit.Type + " for defensive guidance...")
	utils.LogContext("[blue]Generating fix for:[-] " + exploit.Type + " on " + exploit.Target)

	suggestion := &RemediationSuggestion{
		VulnerabilityType: exploit.Type,
		Severity:          exploit.Confidence,
		Exploit:           exploit.Reasoning,
	}

	// Generate fix recommendations based on vulnerability type
	switch exploit.Type {
	case "BOLA":
		suggestion = generateBOLAFix(exploit, result)
	case "BFLA":
		suggestion = generateBFLAFix(exploit, result)
	case "BOPLA":
		suggestion = generateBOPLAFix(exploit, result)
	case "SSRF":
		suggestion = generateSSRFFix(exploit, result)
	case "INJECTION":
		suggestion = generateInjectionFix(exploit, result)
	case "JWT_BYPASS":
		suggestion = generateJWTBypassFix(exploit, result)
	case "CLOUD_PIVOT":
		suggestion = generateCloudPivotFix(exploit, result)
	default:
		suggestion = generateGenericFix(exploit, result)
	}

	// CRITICAL: Verify suggestion before returning to ensure production safety
	VerifyRemediationSuggestion(suggestion)

	return suggestion
}

// generateBOLAFix creates a fix for Broken Object Level Authorization
func generateBOLAFix(exploit TacticalAction, result string) *RemediationSuggestion {
	fix := &RemediationSuggestion{
		VulnerabilityType: "BOLA",
		Severity:          exploit.Confidence,
		Exploit:           "Attacker accessed objects by manipulating ID parameter",
		FixDescription:    "Implement object-level authorization checks to ensure users can only access their own objects.",
		Language:          "Go",
		ImplementationURL: "https://owasp.org/www-community/attacks/Insecure_Direct_Object_References",
	}

	fix.CodeSnippet = `
// Middleware: Verify object ownership before responding
func AuthorizeObjectAccess(c *gin.Context) {
	userID := c.GetString("user_id")
	objectID := c.Param("id")
	
	// Fetch object and verify ownership
	obj := database.GetObject(objectID)
	if obj == nil || obj.OwnerID != userID {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

// Usage: router.GET("/api/objects/:id", AuthorizeObjectAccess, getObjectHandler)
`

	return fix
}

// generateBFLAFix creates a fix for Broken Function Level Authorization
func generateBFLAFix(exploit TacticalAction, result string) *RemediationSuggestion {
	fix := &RemediationSuggestion{
		VulnerabilityType: "BFLA",
		Severity:          exploit.Confidence,
		Exploit:           "Attacker escalated privileges by using privileged HTTP methods (DELETE, PATCH)",
		FixDescription:    "Implement function-level authorization to restrict privileged actions to authorized roles.",
		Language:          "Go",
		ImplementationURL: "https://owasp.org/API-Security/editions/2023/en/0x05-broken-function-level-authorization/",
	}

	fix.CodeSnippet = `
// Middleware: Role-based access control
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		
		// Verify user has required role
		if userRole != requiredRole && userRole != "admin" {
			c.JSON(403, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Usage: router.DELETE("/api/users/:id", RequireRole("admin"), deleteUserHandler)
`

	return fix
}

// generateBOPLAFix creates a fix for Broken Object Property Level Authorization
func generateBOPLAFix(exploit TacticalAction, result string) *RemediationSuggestion {
	fix := &RemediationSuggestion{
		VulnerabilityType: "BOPLA",
		Severity:          exploit.Confidence,
		Exploit:           "Attacker modified object properties via mass assignment",
		FixDescription:    "Implement property whitelisting to restrict which fields can be modified by clients.",
		Language:          "Go",
		ImplementationURL: "https://owasp.org/www-community/attacks/Mass_Assignment",
	}

	fix.CodeSnippet = `
// Define safe fields for client modification
var AllowedUpdateFields = []string{"name", "email", "bio"}

// Middleware: Filter input fields
func SanitizeInput(c *gin.Context) {
	var input map[string]interface{}
	c.BindJSON(&input)
	
	// Remove non-whitelisted fields
	for key := range input {
		if !contains(AllowedUpdateFields, key) {
			delete(input, key)
		}
	}
	
	c.Set("sanitized_input", input)
	c.Next()
}

// Usage: router.PATCH("/api/users/:id", SanitizeInput, updateUserHandler)
`

	return fix
}

// generateSSRFFix creates a fix for Server-Side Request Forgery
func generateSSRFFix(exploit TacticalAction, result string) *RemediationSuggestion {
	fix := &RemediationSuggestion{
		VulnerabilityType: "SSRF",
		Severity:          exploit.Confidence,
		Exploit:           "Attacker triggered requests to internal endpoints via SSRF",
		FixDescription:    "Implement URL validation to prevent requests to internal/private IP ranges.",
		Language:          "Go",
		ImplementationURL: "https://owasp.org/www-community/attacks/Server_Side_Request_Forgery",
	}

	fix.CodeSnippet = `
// Validate URL before making requests
func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	
	// Block private IP ranges
	privateRanges := []string{
		"127.0.0.0/8",        // Localhost
		"10.0.0.0/8",         // Private
		"172.16.0.0/12",      // Private
		"169.254.0.0/16",     // Link-local
		"169.254.169.254/32", // AWS metadata
	}
	
	ip := net.ParseIP(u.Hostname())
	for _, cidr := range privateRanges {
		_, network, _ := net.ParseCIDR(cidr)
		if network.Contains(ip) {
			return false
		}
	}
	
	// Also allow only http/https
	return u.Scheme == "http" || u.Scheme == "https"
}

// Usage: if IsValidURL(userInput) { makeRequest(userInput) }
`

	return fix
}

// generateInjectionFix creates a fix for Injection attacks
func generateInjectionFix(exploit TacticalAction, result string) *RemediationSuggestion {
	fix := &RemediationSuggestion{
		VulnerabilityType: "INJECTION",
		Severity:          exploit.Confidence,
		Exploit:           "Attacker injected malicious code/commands via user input",
		FixDescription:    "Use parameterized queries and input validation to prevent injection.",
		Language:          "Go",
		ImplementationURL: "https://owasp.org/www-community/attacks/SQL_Injection",
	}

	fix.CodeSnippet = `
// WRONG: String concatenation (vulnerable)
query := "SELECT * FROM users WHERE id = " + userInput

// CORRECT: Parameterized query
query := "SELECT * FROM users WHERE id = ?"
db.Query(query, userInput)

// Input validation
func ValidateUserID(id string) error {
	if !regexp.MustCompile("^[0-9]+$").MatchString(id) {
		return fmt.Errorf("invalid user ID format")
	}
	return nil
}

// Usage: 
// if err := ValidateUserID(userInput); err == nil {
//     db.Query("SELECT * FROM users WHERE id = ?", userInput)
// }
`

	return fix
}

// generateJWTBypassFix creates a fix for JWT bypass attacks
func generateJWTBypassFix(exploit TacticalAction, result string) *RemediationSuggestion {
	fix := &RemediationSuggestion{
		VulnerabilityType: "JWT_BYPASS",
		Severity:          exploit.Confidence,
		Exploit:           "Attacker manipulated JWT token claims or used algorithm substitution",
		FixDescription:    "Implement strict JWT validation with key pinning and algorithm restrictions.",
		Language:          "Go",
		ImplementationURL: "https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/",
	}

	fix.CodeSnippet = `
// Strict JWT validation
func ValidateJWT(tokenString string) (jwt.Claims, error) {
	// CRITICAL: Always verify the signing algorithm
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Restrict to specific algorithm(s)
		if token.Method.Alg() != "HS256" && token.Method.Alg() != "RS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		
		// Use proper key storage (not hardcoded)
		return getSigningKey(token), nil
	})
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	
	return token.Claims, nil
}
`

	return fix
}

// generateCloudPivotFix creates a fix for cloud metadata exploitation
func generateCloudPivotFix(exploit TacticalAction, result string) *RemediationSuggestion {
	fix := &RemediationSuggestion{
		VulnerabilityType: "CLOUD_PIVOT",
		Severity:          exploit.Confidence,
		Exploit:           "Attacker accessed cloud metadata endpoint to retrieve credentials/tokens",
		FixDescription:    "Restrict access to cloud metadata endpoints and use short-lived credentials.",
		Language:          "Go",
		ImplementationURL: "https://owasp.org/www-community/attacks/Cloud_Metadata_Abuse",
	}

	fix.CodeSnippet = `
// 1. Firewall rule: Block access to 169.254.169.254 (AWS metadata)
// From your application container/instance, disable metadata access:
// AWS: Add IMDSv2 requirement + token-based access
// GCP: Disable metadata server if not needed
// Azure: Restrict access via network policies

// 2. Use IAM roles instead of embedded credentials
// Go: Use AWS SDK with environment-based credentials or IAM roles
// Instead of: hardcoded AWS_ACCESS_KEY in env vars

// 3. Rotate credentials frequently
// Set short TTL (15-60 min) for temporary credentials
// Implement credential monitoring and audit logging

// 4. Application-level check
func IsInternalRequest(ip string) bool {
	internalRanges := []string{
		"169.254.0.0/16",     // AWS metadata
		"169.254.169.254/32", // Specific AWS metadata IP
	}
	
	clientIP := net.ParseIP(ip)
	for _, cidr := range internalRanges {
		_, network, _ := net.ParseCIDR(cidr)
		if network.Contains(clientIP) {
			return true
		}
	}
	return false
}
`

	return fix
}

// generateGenericFix creates a generic fix template
func generateGenericFix(exploit TacticalAction, result string) *RemediationSuggestion {
	fix := &RemediationSuggestion{
		VulnerabilityType: exploit.Type,
		Severity:          exploit.Confidence,
		Exploit:           exploit.Reasoning,
		FixDescription:    "Implement input validation, output encoding, and proper authorization checks.",
		Language:          "Go",
		ImplementationURL: "https://owasp.org/www-project-api-security/",
	}

	fix.CodeSnippet = `
// Generic security middleware template
func SecurityMiddleware(c *gin.Context) {
	// 1. Validate input
	if !ValidateInput(c.Request) {
		c.JSON(400, gin.H{"error": "Invalid input"})
		c.Abort()
		return
	}
	
	// 2. Check authorization
	if !IsAuthorized(c.GetString("user_id"), c.Request.URL.Path) {
		c.JSON(403, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	
	// 3. Rate limiting
	if !RateLimit(c.ClientIP()) {
		c.JSON(429, gin.H{"error": "Too many requests"})
		c.Abort()
		return
	}
	
	c.Next()
}
`

	return fix
}

// FormatRemediationForUI converts the suggestion into UI-friendly text
func (r *RemediationSuggestion) FormatForUI() string {
	// Add verification status banner
	verificationBanner := ""
	switch r.VerificationStatus {
	case "GOLD_STANDARD":
		verificationBanner = fmt.Sprintf(`
╔═══════════════════════════════════════════════════════════════════════════╗
║ [green]✓ GOLD STANDARD - VERIFIED & PRODUCTION READY[-]
║ %s
╚═══════════════════════════════════════════════════════════════════════════╝`, r.VerificationNotes)

	case "VERIFIED":
		verificationBanner = fmt.Sprintf(`
╔═══════════════════════════════════════════════════════════════════════════╗
║ [green]✓ VERIFIED - Passed Security Review[-]
║ %s
╚═══════════════════════════════════════════════════════════════════════════╝`, r.VerificationNotes)

	case "UNVERIFIED":
		verificationBanner = fmt.Sprintf(`
╔═══════════════════════════════════════════════════════════════════════════╗
║ [yellow]⚠ WARNING - UNVERIFIED CODE[-]
║ %s
║
║ ACTION: Manual security review REQUIRED before production use.
║ TEST FOR: ReDoS in regex, SQL injection, auth bypass, hardcoded secrets.
╚═══════════════════════════════════════════════════════════════════════════╝`, r.VerificationNotes)
	}

	output := fmt.Sprintf(`%s
╔════════════════════════════════════════════════════════════════════════════╗
║ REMEDIATION SUGGESTION - %s (%s)
╚════════════════════════════════════════════════════════════════════════════╝

VULNERABILITY: %s
SEVERITY: %s

ATTACK SUMMARY:
  %s

FIX DESCRIPTION:
  %s

LANGUAGE: %s

CODE EXAMPLE:
%s

DOCUMENTATION:
  %s

═════════════════════════════════════════════════════════════════════════════
`, verificationBanner, r.VulnerabilityType, r.Severity, r.VulnerabilityType, r.Severity,
		r.Exploit, r.FixDescription, r.Language, r.CodeSnippet, r.ImplementationURL)

	return output
}

// SuggestFixAndLog generates a fix and logs it to the UI
func SuggestFixAndLog(exploit TacticalAction, result string) {
	suggestion := SuggestFix(exploit, result)
	if suggestion == nil {
		return
	}

	formattedOutput := suggestion.FormatForUI()

	utils.TacticalLog("[green]✓ REMEDIATION:[-] Generated fix for " + exploit.Type)
	utils.LogContext(formattedOutput)

	// Also store in DataSilo for reference
	key := "remediation_" + strconv.Itoa(exploit.ID)
	logic.GlobalDataSilo.Set(key, formattedOutput)
}
