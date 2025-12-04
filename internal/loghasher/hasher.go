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
		// Patterns for the most common timestamp formats:
		// - 2025-01-17T14:32:11Z
		// - 2025-01-17T14:32:11+02:00
		// - 2025-01-17 14:32:11
		// - 14:32:11
		// - Jan 17 14:32:11
		timestampPattern: regexp.MustCompile(
			`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z)|` +
				`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:[+-]\d{2}:\d{2}))|` +
				`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})|` +
				`(\b\d{2}:\d{2}:\d{2}\b)|` +
				`([A-Za-z]{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})`,
		),
	}
}

func (h *logHasher) Hash(log string) string {
	norimalizedLog := h.timestampPattern.ReplaceAllString(log, "")
	sum := sha256.Sum256([]byte(norimalizedLog))
	return hex.EncodeToString(sum[:])
}
