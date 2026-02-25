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
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// === PRIORITY BETA: PATH & PARAMETER OBFUSCATION ===

// PathObfuscationTechnique defines different ways to obfuscate API paths
type PathObfuscationTechnique int

const (
	CacheBusters PathObfuscationTechnique = iota
	PathParameters
	DoubleEncoding
	URLFragments
)

// ObfuscatePath adds "noise" parameters and path segments that WAFs must process
// but the server ignores, making automated patterns harder to detect
func ObfuscatePath(originalPath string, technique PathObfuscationTechnique) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	switch technique {
	case CacheBusters:
		return obfuscateWithCacheBusters(originalPath, seed)
	case PathParameters:
		return obfuscateWithPathParameters(originalPath, seed)
	case DoubleEncoding:
		return obfuscateWithDoubleEncoding(originalPath, seed)
	case URLFragments:
		return obfuscateWithFragments(originalPath, seed)
	default:
		return originalPath
	}
}

// obfuscateWithCacheBusters adds random query parameters that the server ignores
// Example: /api/v1/users?_debug=0&_ref=home&_t=1234567890
func obfuscateWithCacheBusters(path string, seed *rand.Rand) string {
	// Random noise parameters that common frameworks ignore
	noiseParams := []string{
		fmt.Sprintf("_debug=%d", seed.Intn(2)),
		fmt.Sprintf("_ref=%s", randomString(seed, 4)),
		fmt.Sprintf("_t=%d", time.Now().UnixNano()%100000),
		fmt.Sprintf("_v=%d", seed.Intn(999)),
		fmt.Sprintf("_cache=%s", randomString(seed, 8)),
		fmt.Sprintf("_rnd=%d", seed.Intn(999999)),
	}

	// Randomly select 1-3 parameters
	paramCount := seed.Intn(3) + 1
	var params []string
	for i := 0; i < paramCount; i++ {
		params = append(params, noiseParams[seed.Intn(len(noiseParams))])
	}

	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}

	obfuscated := path + separator + strings.Join(params, "&")
	utils.TacticalLog(fmt.Sprintf("[cyan]OBFUSCATION:[-] Cache-buster applied: %s", obfuscated))
	return obfuscated
}

// obfuscateWithPathParameters adds path parameters that servers typically ignore
// Example: /api/v1/users;v=1.0;x=y/profile
// These confuse regex-based WAF rules
func obfuscateWithPathParameters(path string, seed *rand.Rand) string {
	pathParams := []string{
		";v=1.0",
		";x=y",
		";client=web",
		";format=json",
		";nocache=true",
	}

	// Insert before the last path segment
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		lastIdx := len(parts) - 1
		param := pathParams[seed.Intn(len(pathParams))]
		// Insert parameter before last segment
		parts[lastIdx] = parts[lastIdx] + param
	}

	obfuscated := strings.Join(parts, "/")
	utils.TacticalLog(fmt.Sprintf("[cyan]OBFUSCATION:[-] Path parameters applied: %s", obfuscated))
	return obfuscated
}

// obfuscateWithDoubleEncoding double-encodes path components
// Example: /api/v1/users → /api/v1/%75sers (u encoded twice)
func obfuscateWithDoubleEncoding(path string, seed *rand.Rand) string {
	// Randomly select segments to encode
	parts := strings.Split(path, "/")
	encodeIndices := seed.Intn(len(parts))

	if encodeIndices >= 1 && encodeIndices < len(parts) {
		// Encode selected segment
		segment := parts[encodeIndices]
		encoded := url.QueryEscape(url.QueryEscape(segment))
		parts[encodeIndices] = encoded
	}

	obfuscated := strings.Join(parts, "/")
	utils.TacticalLog(fmt.Sprintf("[cyan]OBFUSCATION:[-] Double-encoding applied: %s", obfuscated))
	return obfuscated
}

// obfuscateWithFragments adds URL fragments that the server ignores but WAF must process
// Example: /api/v1/users#section1
func obfuscateWithFragments(path string, seed *rand.Rand) string {
	fragments := []string{"#api", "#v1", "#data", "#json", "#response"}
	fragment := fragments[seed.Intn(len(fragments))]

	obfuscated := path + fragment
	utils.TacticalLog(fmt.Sprintf("[cyan]OBFUSCATION:[-] Fragment added: %s", obfuscated))
	return obfuscated
}

// randomString generates a random alphanumeric string of length n
func randomString(seed *rand.Rand, n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[seed.Intn(len(letters))]
	}
	return string(b)
}

// SelectObfuscationStrategy picks a random obfuscation technique to rotate
func SelectObfuscationStrategy() PathObfuscationTechnique {
	techniques := []PathObfuscationTechnique{CacheBusters, PathParameters, URLFragments}
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	return techniques[seed.Intn(len(techniques))]
}
