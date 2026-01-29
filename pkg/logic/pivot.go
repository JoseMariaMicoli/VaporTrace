package logic

import (
	"strings"
	"github.com/pterm/pterm"
)

func TriggerCloudPivot(url string) {
	if strings.Contains(url, "127.0.0.1") || strings.Contains(url, "169.254.169.254") {
		pterm.Info.WithPrefix(pterm.Prefix{Text: "PIVOT"}).Printfln("Intercepted Metadata Target: %s", url)
		
		target := "127.0.0.1"
		if strings.Contains(url, "169.254.169.254") {
			target = "169.254.169.254"
		}
		go ExecutePivot(target, url)
	}
}