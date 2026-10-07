package utils

import (
	"fmt"
	"time"
)

// FormatTime 统一 UTC -> 前端可解析 RFC3339；空值返回空串。
func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// MinutesText 延误/超时分钟文案。
func MinutesText(minutes int) string {
	if minutes <= 0 {
		return "0 分钟"
	}
	return fmt.Sprintf("%d 分钟", minutes)
}

func OneLine(value string) string {
	out := make([]rune, 0, len(value))
	for _, r := range value {
		if r == '\n' || r == '\r' {
			out = append(out, ' ')
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
