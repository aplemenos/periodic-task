package timestamp

import "time"

// One Hour Timestamp
type OneHourTimestamp struct{}

func (oht OneHourTimestamp) GetMatchingTimestamps(t1, t2 time.Time) []string {
	// Adjust t1 and t2 to the start and end of the hour respectively
	t1 = t1.Truncate(time.Hour).Add(time.Hour + 1)
	t2 = t2.Truncate(time.Hour).Add(time.Hour - 1)

	// Generate the periodic timestamps for one hour
	var ptlist []string
	for t := t1; t.Before(t2); t = t.Add(time.Hour) {
		ptlist = append(ptlist, t.Format("20060102T150405Z"))
	}
	return ptlist
}
