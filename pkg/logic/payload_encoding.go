package logic

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// === PRIORITY EPSILON: PAYLOAD ENCODING & CASE RANDOMIZATION ===

// EncodingType defines payload encoding strategy
type EncodingType int

const (
	NoEncoding EncodingType = iota
	GzipEncoding
	DeflateEncoding
	ChunkedEncoding
)

// PayloadTransformation applies encoding/case variations to POST bodies
type PayloadTransformation struct {
	Encoding      EncodingType
	CaseVariation bool
	Whitespace    bool
	KeyOrder      bool
}

// TransformPayload applies multiple encoding techniques to obscure JSON/form payloads
func TransformPayload(payload []byte, technique PayloadTransformation) ([]byte, string) {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Step 1: Try to parse as JSON and apply transformations
	if isJSON(payload) {
		payload = transformJSON(payload, seed, technique)
	}

	// Step 2: Apply content encoding
	encodingHeader := "identity"
	switch technique.Encoding {
	case GzipEncoding:
		payload, encodingHeader = encodeGzip(payload)
	case DeflateEncoding:
		payload, encodingHeader = encodeDeflate(payload)
	case ChunkedEncoding:
		// Chunked will be handled at HTTP level
		encodingHeader = "chunked"
	}

	return payload, encodingHeader
}

// transformJSON applies case and whitespace variations to JSON
func transformJSON(payload []byte, seed *rand.Rand, tech PayloadTransformation) []byte {
	var obj interface{}
	if err := json.Unmarshal(payload, &obj); err != nil {
		return payload
	}

	// Apply transformations
	if tech.CaseVariation {
		payload = randomizeCaseJSON(payload, seed)
	}

	if tech.Whitespace {
		payload = randomizeWhitespaceJSON(payload, seed)
	}

	if tech.KeyOrder {
		// This would require reordering map keys
		payload = randomizeKeyOrderJSON(obj, seed)
	}

	utils.TacticalLog("[green]✓ PAYLOAD:[-] Applied JSON transformations")
	return payload
}

// randomizeCaseJSON changes case of non-critical parts
// Example: {"userId": 1} → {"userId": 1} (keys stay same for accuracy)
func randomizeCaseJSON(payload []byte, seed *rand.Rand) []byte {
	// For now, maintain strict JSON validity
	// This is a placeholder for more aggressive case variation
	return payload
}

// randomizeWhitespaceJSON adds/removes whitespace around colons and commas
// {"id":1} → { "id" : 1 } or {  "id":1  }
func randomizeWhitespaceJSON(payload []byte, seed *rand.Rand) []byte {
	str := string(payload)

	// Randomly add spaces after commas
	if seed.Float64() > 0.5 {
		str = strings.ReplaceAll(str, ",", ", ")
	}

	// Randomly add spaces around colons
	if seed.Float64() > 0.5 {
		str = strings.ReplaceAll(str, ":", ": ")
	}

	// Randomly add newlines
	if seed.Float64() > 0.3 {
		str = strings.ReplaceAll(str, ",", ",\n")
		str = strings.ReplaceAll(str, "{", "{\n")
		str = strings.ReplaceAll(str, "}", "\n}")
	}

	utils.TacticalLog("[cyan]ENCODING:[-] Whitespace randomization applied")
	return []byte(str)
}

// randomizeKeyOrderJSON reorders JSON keys (requires map reconstruction)
func randomizeKeyOrderJSON(obj interface{}, seed *rand.Rand) []byte {
	// Note: Go's json.Marshal doesn't preserve key order from maps
	// For true reordering, we'd need a custom JSON encoder
	// For now, marshal normally and randomize in transformJSON
	data, _ := json.Marshal(obj)
	return data
}

// encodeGzip compresses payload with gzip
func encodeGzip(payload []byte) ([]byte, string) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(payload)
	writer.Close()

	utils.TacticalLog(fmt.Sprintf("[cyan]ENCODING:[-] Gzip applied (size: %d → %d bytes)", len(payload), buf.Len()))
	return buf.Bytes(), "gzip"
}

// encodeDeflate compresses payload with deflate
func encodeDeflate(payload []byte) ([]byte, string) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	writer.Write(payload)
	writer.Close()

	utils.TacticalLog(fmt.Sprintf("[cyan]ENCODING:[-] Deflate applied (size: %d → %d bytes)", len(payload), buf.Len()))
	return buf.Bytes(), "deflate"
}

// isJSON checks if payload is valid JSON
func isJSON(payload []byte) bool {
	payload = bytes.TrimSpace(payload)
	return (bytes.HasPrefix(payload, []byte("{")) || bytes.HasPrefix(payload, []byte("[")))
}

// SelectRandomEncoding picks a random encoding strategy
func SelectRandomEncoding() PayloadTransformation {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	encodings := []EncodingType{NoEncoding, GzipEncoding, DeflateEncoding}
	encoding := encodings[seed.Intn(len(encodings))]

	return PayloadTransformation{
		Encoding:      encoding,
		CaseVariation: seed.Float64() > 0.5,
		Whitespace:    seed.Float64() > 0.5,
		KeyOrder:      seed.Float64() > 0.7,
	}
}

// DecodePayload handles decompression for testing
func DecodePayload(payload []byte, encoding string) ([]byte, error) {
	switch encoding {
	case "gzip":
		reader, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return io.ReadAll(reader)

	case "deflate":
		reader, err := zlib.NewReader(bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return io.ReadAll(reader)

	default:
		return payload, nil
	}
}
