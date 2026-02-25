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
	"sync"
)

// ScopeManager handles the inclusion and exclusion rules for target auditing.
// It uses RWMutex to allow high-concurrency reads during fuzzing.
type ScopeManager struct {
	mu           sync.RWMutex
	IncludeRegex []*regexp.Regexp
	ExcludeRegex []*regexp.Regexp
	RawInclude   []string
	RawExclude   []string
}

// GlobalScope is the singleton instance used across the engine
var GlobalScope = &ScopeManager{
	RawInclude: []string{},
	RawExclude: []string{},
}

// UpdateRules rebuilds the regex objects from raw strings. Thread-safe.
func (sm *ScopeManager) UpdateRules(include, exclude []string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.RawInclude = include
	sm.RawExclude = exclude
	sm.IncludeRegex = make([]*regexp.Regexp, 0, len(include))
	sm.ExcludeRegex = make([]*regexp.Regexp, 0, len(exclude))

	// Compile Includes
	for _, s := range include {
		if s == "" {
			continue
		}
		re, err := regexp.Compile(s)
		if err != nil {
			return fmt.Errorf("invalid include regex '%s': %v", s, err)
		}
		sm.IncludeRegex = append(sm.IncludeRegex, re)
	}

	// Compile Excludes
	for _, s := range exclude {
		if s == "" {
			continue
		}
		re, err := regexp.Compile(s)
		if err != nil {
			return fmt.Errorf("invalid exclude regex '%s': %v", s, err)
		}
		sm.ExcludeRegex = append(sm.ExcludeRegex, re)
	}

	return nil
}

// IsInScope checks if a target URL is allowed.
// Logic:
// 1. If Exclude matches, return FALSE immediately.
// 2. If Include is empty, return TRUE (default allow if no whitelist defined).
// 3. If Include is set, return TRUE only if it matches one of the Include regexes.
func (sm *ScopeManager) IsInScope(targetURL string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 1. Check Exclusions first (Safety First)
	for _, re := range sm.ExcludeRegex {
		if re.MatchString(targetURL) {
			return false
		}
	}

	// 2. If no whitelist defined, allow everything not excluded
	if len(sm.IncludeRegex) == 0 {
		return true
	}

	// 3. Check Inclusions
	for _, re := range sm.IncludeRegex {
		if re.MatchString(targetURL) {
			return true
		}
	}

	return false
}

// GetRawRules returns copies of the current rule strings for the UI.
func (sm *ScopeManager) GetRawRules() ([]string, []string) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	inc := make([]string, len(sm.RawInclude))
	copy(inc, sm.RawInclude)

	exc := make([]string, len(sm.RawExclude))
	copy(exc, sm.RawExclude)

	return inc, exc
}
