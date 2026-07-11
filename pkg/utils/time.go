package utils

import "time"

// FormatTimePtr 将 *time.Time 格式化为 "2006-01-02 15:04:05" 字符串指针
func FormatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	str := t.Format("2006-01-02 15:04:05")
	return &str
}

func FormatCreatedAt(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
