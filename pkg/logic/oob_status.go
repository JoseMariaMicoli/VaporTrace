package logic

import (
	"fmt"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// OOBChannelStatus represents the status of each OOB channel
type OOBChannelStatus struct {
	ChannelType string
	Status      string
	Description string
	Ready       bool
}

// OOBStatusManager manages OOB channel status and statistics
type OOBStatusManager struct {
	mu               sync.RWMutex
	Channels         []OOBChannelStatus
	TotalQueued      int
	TotalSent        int
	TotalFailed      int
	EncryptionMethod string
}

var globalOOBStatus = &OOBStatusManager{
	EncryptionMethod: "AES-256-GCM",
	Channels: []OOBChannelStatus{
		{
			ChannelType: "TCP",
			Status:      "READY",
			Description: "Custom TCP protocol to OOB receiver",
			Ready:       true,
		},
		{
			ChannelType: "DNS",
			Status:      "READY",
			Description: "Covert subdomain encoding",
			Ready:       true,
		},
		{
			ChannelType: "ICMP",
			Status:      "READY",
			Description: "Firewall evasion via Echo",
			Ready:       true,
		},
	},
}

// GetOOBStatus returns the current OOB channel status
func GetOOBStatus() *OOBStatusManager {
	globalOOBStatus.mu.RLock()
	defer globalOOBStatus.mu.RUnlock()

	// Return a copy to avoid race conditions
	copy := &OOBStatusManager{
		TotalQueued:      globalOOBStatus.TotalQueued,
		TotalSent:        globalOOBStatus.TotalSent,
		TotalFailed:      globalOOBStatus.TotalFailed,
		EncryptionMethod: globalOOBStatus.EncryptionMethod,
		Channels:         make([]OOBChannelStatus, len(globalOOBStatus.Channels)),
	}
	copy.Channels = append([]OOBChannelStatus{}, globalOOBStatus.Channels...)
	return copy
}

// UpdateOOBStats updates OOB transmission statistics
func UpdateOOBStats(queued, sent, failed int) {
	globalOOBStatus.mu.Lock()
	defer globalOOBStatus.mu.Unlock()

	globalOOBStatus.TotalQueued += queued
	globalOOBStatus.TotalSent += sent
	globalOOBStatus.TotalFailed += failed
}

// ReportOOBStatus outputs comprehensive OOB channel status
func ReportOOBStatus() {
	status := GetOOBStatus()

	utils.TacticalLog("[magenta::b]PHASE 7.2: OUT-OF-BAND EXFILTRATION CHANNEL[-:-:-]")
	utils.TacticalLog("[cyan]Status:[-] [green]CONFIGURED & MONITORING")
	utils.TacticalLog("")

	// Channel details
	utils.TacticalLog("[yellow]Available Channels:[-]")
	for _, ch := range status.Channels {
		statusColor := "[green]"
		if !ch.Ready {
			statusColor = "[red]"
		}
		utils.TacticalLog(fmt.Sprintf("  [blue]•[-] %s: %s%s[-] (%s)", ch.ChannelType, statusColor, ch.Status, ch.Description))
	}

	utils.TacticalLog("")
	utils.TacticalLog(fmt.Sprintf("[yellow]Encryption:[-] %s", status.EncryptionMethod))
	utils.TacticalLog("[yellow]Authentication:[-] [green]GCM AEAD (Authenticated Encryption with Associated Data)")

	// Statistics
	if status.TotalQueued > 0 || status.TotalSent > 0 || status.TotalFailed > 0 {
		utils.TacticalLog("")
		utils.TacticalLog("[yellow]Statistics:[-]")
		utils.TacticalLog(fmt.Sprintf("  Queued: [cyan]%d[-] | Sent: [green]%d[-] | Failed: [red]%d", status.TotalQueued, status.TotalSent, status.TotalFailed))
	}

	utils.TacticalLog("")
	utils.TacticalLog("[cyan]Integration:[-] Works seamlessly with SSRF, BOLA, and other attack vectors")
	utils.TacticalLog("[cyan]Workflow:[-] Automatically queues captured sensitive data for encrypted transmission")
}

// SetChannelStatus updates the status of a specific OOB channel
func SetChannelStatus(channelType, status string) {
	globalOOBStatus.mu.Lock()
	defer globalOOBStatus.mu.Unlock()

	for i, ch := range globalOOBStatus.Channels {
		if ch.ChannelType == channelType {
			globalOOBStatus.Channels[i].Status = status
			globalOOBStatus.Channels[i].Ready = status == "READY"
			break
		}
	}
}
