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
