package date

import "time"

// formats the given time as string
func DateTimeToString(timeVar time.Time) string {
	return timeVar.UTC().Format(time.RFC3339Nano)
}

// returns the current date time in UTC
func UtcDateTimeNow() string {
	return DateTimeToString(time.Now())
}
