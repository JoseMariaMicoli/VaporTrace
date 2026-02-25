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

package logic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson" // Assumes already in go.mod, standard for Go JSON parsing
)

// ExtractValue pulls a value from a response body/header based on a selector strategy
func ExtractValue(source string, strategy string, selector string) (string, error) {
	switch strings.ToLower(strategy) {
	case "json":
		// Uses GJSON syntax (e.g., "auth.access_token" or "users.0.id")
		result := gjson.Get(source, selector)
		if !result.Exists() {
			return "", fmt.Errorf("json key '%s' not found", selector)
		}
		return result.String(), nil

	case "regex":
		// Uses Go Regex (e.g., `token":"(.*?)"`)
		re, err := regexp.Compile(selector)
		if err != nil {
			return "", fmt.Errorf("invalid regex: %v", err)
		}
		matches := re.FindStringSubmatch(source)
		if len(matches) < 2 {
			return "", fmt.Errorf("regex match failed for '%s'", selector)
		}
		return matches[1], nil // Return the first capture group

	case "header":
		// Expects source to be the raw header dump
		// Selector is the header name (e.g., "Set-Cookie")
		lines := strings.Split(source, "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.ToLower(line), strings.ToLower(selector)+":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
		return "", fmt.Errorf("header '%s' not found", selector)

	default:
		return "", fmt.Errorf("unknown strategy: %s", strategy)
	}
}

// InjectVariables replaces {{key}} placeholders with values from the context map
func InjectVariables(input string, context map[string]string) string {
	result := input
	for key, val := range context {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, val)
	}
	return result
}
