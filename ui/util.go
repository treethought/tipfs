package ui

import (
	"fmt"
	"strings"
)

func truncateMiddle(s string, length int) string {
	if len(s) < length*2+3 {
		return s
	}
	// midLen := len(s) - 8
	suffixStart := len(s) - 4

	remove := s[4:suffixStart]
	return strings.ReplaceAll(s, remove, "...")

}

func byteCount(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
