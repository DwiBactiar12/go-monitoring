package utils

import "time"

// Format waktu ke string readable
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Parse string ke time.Time
func ParseTime(value string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", value)
}
