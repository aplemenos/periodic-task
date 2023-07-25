package timestamp

import "time"

// One Day Timestamp
type OneDayTimestamp struct{}

func (odt OneDayTimestamp) GetMatchingTimestamps(t1, t2 time.Time) []string {
	// Adjust t1 and t2 to the start and end of the day respectively
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour()+1, 0, 0, 0, time.UTC)
	t2 = t2.Truncate(time.Second)

	// Generate the periodic timestamps for one day
	var ptlist []string
	for t := t1; t.Before(t2); t = t.AddDate(0, 0, 1) {
		ptlist = append(ptlist, t.Format("20060102T150405Z"))
	}
	return ptlist
}
