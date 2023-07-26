package timestamp

import "time"

// One Hour Timestamp
type OneHourTimestamp struct{}

func (oht OneHourTimestamp) GetMatchingTimestamps(t1, t2 time.Time) []string {
	// Adjust t1 and t2 to the start and end of the hour respectively
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(),
		0, time.UTC)
	t2 = t2.Truncate(time.Second)

	// Round the start time to hour.
	t1 = t1.Round(60 * time.Minute)

	// Generate the periodic timestamps for one hour
	var ptlist []string
	for t := t1; t.Before(t2); t = t.Add(time.Hour) {
		ptlist = append(ptlist, t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
