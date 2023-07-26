package timestamp

import "time"

// One Month Timestamp
type OneMonthTimestamp struct{}

func (omt OneMonthTimestamp) GetMatchingTimestamps(t1, t2 time.Time) []string {
	// Adjust t1 and t2 to the start and end of the month respectively
	t1 = time.Date(t1.Year(), t1.Month()+1, 0, t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	t2 = t2.Truncate(time.Second)

	// Round the start time to hour.
	t1 = t1.Round(60 * time.Minute)

	// Generate the periodic timestamps for one month
	var ptlist []string
	for t := t1; t.Before(t2); t = t.AddDate(0, 1, 0) {
		_t := lastDateOfMonth(t.Year(), t.Month(), t.Hour())
		ptlist = append(ptlist, _t.Format("20060102T150405Z"))
	}
	return ptlist
}
