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
	"os"
	"os/exec"
)

// GhostMasquerade implements Phase 9.11: Process Renaming
// Adopted from Ghost-Pipeline IoC: kworker_system_auth
func GhostMasquerade(binaryPath string) error {
	masqueradeName := "kworker_system_auth"

	// Atomic rename and execute
	cmd := exec.Command("cp", binaryPath, masqueradeName)
	if err := cmd.Run(); err != nil {
		return err
	}

	os.Chmod(masqueradeName, 0755)
	return nil
}
