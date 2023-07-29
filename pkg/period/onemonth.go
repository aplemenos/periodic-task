package period

import (
	"time"
)

// One Month Period
type OneMonthPeriod struct{}

func (omp OneMonthPeriod) GetMatchingTimestamps(t1, t2 time.Time) []string {
	var ptlist []string

	utc, _ := time.LoadLocation("UTC")

	// Get the time zone offset of the start time.
	_, offsetSecs := t1.Zone()

	// Generate the periodic timestamps for one month
	for t := t1; t.Before(t2); t = t.AddDate(0, 1, 0) {
		_t := lastDateOfMonth(t)
		if _t.After(t2) {
			break
		}

		// Load time in UTC
		_t = _t.In(utc)
		// Add offset of start time (removing the daylight saving time)
		_t = _t.Add(time.Duration(offsetSecs) * time.Second)

		// Round the hour
		_t = _t.Round(60 * time.Minute)
		// Append the time to the list
		ptlist = append(ptlist, _t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
