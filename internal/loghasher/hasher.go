package loghasher

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

type HasherInterface interface {
	Hash(log string) string
}

type logHasher struct {
	timestampPattern *regexp.Regexp
}

func NewHasher() HasherInterface {
	return &logHasher{
		timestampPattern: regexp.MustCompile(
			// Patterns for the most common timestamp formats:
			// - 2025-01-17T14:32:11Z
			// - 2025-01-17T14:32:11+02:00
			// - 2025-01-17 14:32:11
			// - 14:32:11
			// - Jan 17 14:32:11
			`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z)|` +
				`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:[+-]\d{2}:\d{2}))|` +
				`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})|` +
				`(\b\d{2}:\d{2}:\d{2}\b)|` +
				`([A-Za-z]{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})|` +
				// IP address patterns (IPv4 + IPv6):
				`\b(?:\d{1,3}\.){3}\d{1,3}\b|` +
				`\b(?:[0-9a-fA-F]{0,4}:){2,7}[0-9a-fA-F]{0,4}\b|` +
				// UUID patterns:
				`\b[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\b|` +
				// Request id patterns:
				`\b[0-9a-fA-F]+/[A-Za-z0-9]+-\d+\b|` +
				// 32 character hash patterns:
				`\b[0-9a-fA-F]{32}\b`,
		),
	}
}

func (h *logHasher) Hash(log string) string {
	norimalizedLog := h.timestampPattern.ReplaceAllString(log, "")
	sum := sha256.Sum256([]byte(norimalizedLog))
	return hex.EncodeToString(sum[:])
}
