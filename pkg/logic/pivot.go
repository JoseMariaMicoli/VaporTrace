package logic

import (
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

func TriggerCloudPivot(url string) {
	if strings.Contains(url, "127.0.0.1") || strings.Contains(url, "169.254.169.254") {
		utils.TacticalLog("[cyan]PIVOT:[-] Intercepted Metadata Target: " + url)

		target := "127.0.0.1"
		if strings.Contains(url, "169.254.169.254") {
			target = "169.254.169.254"
		}
		go ExecutePivot(target, url)
	}
}
