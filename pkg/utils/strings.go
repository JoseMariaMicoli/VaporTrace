package utils

import "regexp"

// StripANSI removes raw console color codes and tview markup to prevent TUI artifacts.
// This uses a robust regex to catch complex ANSI sequences and also removes tview tags.
func StripANSI(str string) string {
	// Remove raw ANSI escape sequences
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	re := regexp.MustCompile(ansi)
	str = re.ReplaceAllString(str, "")

	// Remove tview markup tags like [cyan], [red], [-], etc.
	tviewTags := regexp.MustCompile(`\[[a-zA-Z0-9\-:;]*\]`)
	str = tviewTags.ReplaceAllString(str, "")

	return str
}
