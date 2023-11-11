package helpers

import (
	"strconv"
	"time"
)

// ParseDuration parse a duration string
func ParseDuration(str string, defaultDuration time.Duration) time.Duration {
	if str == "" {
		return defaultDuration
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		return defaultDuration
	}

	return duration
}

// ParseDurationToString parse a duration string
func ParseDurationToString(duration time.Duration) string {
	return duration.String()
}

// ParseInt64 parse string into int64
func ParseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}
