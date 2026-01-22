package utils

import "regexp"

func StripANSI(str string) string {
    // Standard Regex for ANSI escape sequences
    const ansi = `\x1b\[[0-9;]*[a-zA-Z]`
    re := regexp.MustCompile(ansi)
    return re.ReplaceAllString(str, "")
}