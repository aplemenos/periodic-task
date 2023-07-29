package period

import (
	"time"
)

// One Year Period
type OneYearPeriod struct{}

func (oyp OneYearPeriod) GetMatchingTimestamps(t1, t2 time.Time) []string {
	// Get the time zone offset of the start time.
	_, offsetSecs := t1.Zone()

	// Adjust t1 at the end of the year
	t1 = time.Date(t1.Year()+1, 1, 0, t1.Hour(), t1.Minute(), t1.Second(),
		0, t1.Location())

	// Load time in UTC
	utc, _ := time.LoadLocation("UTC")
	t1 = t1.In(utc)

	// Generate the periodic timestamps for one year
	var ptlist []string
	for t := t1; t.Before(t2); t = t.AddDate(1, 0, 0) {
		_t := lastDateOfMonth(t)
		if _t.After(t2) {
			break
		}

		// Add offset of start time (removing the daylight saving time)
		_t = _t.Add(time.Duration(offsetSecs) * time.Second)

		// Round the hour
		_t = _t.Round(60 * time.Minute)
		// Append the time to the list
		ptlist = append(ptlist, _t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
