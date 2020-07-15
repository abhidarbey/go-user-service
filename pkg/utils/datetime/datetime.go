package datetime

import (
	"time"
)

const (
	dateTimeLayout = "2006-01-02T15:04:05Z"
)

// GetDateTime function
func GetDateTime() time.Time {
	return time.Now()
}

// GetDateTimeString function
func GetDateTimeString() string {
	return GetDateTime().Format(dateTimeLayout)
}
