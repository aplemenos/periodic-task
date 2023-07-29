package period

import (
	"time"
)

// One Year Period
type OneYearPeriod struct{}

func (oyp OneYearPeriod) GetMatchingTimestamps(t1, t2 time.Time, tz *time.Location) []string {
	// Get the time zone offset of the start time based on the requested timezone.
	_, offsetSecs := t1.In(tz).Zone()

	// Adjust t1 at the end of the year
	t1 = time.Date(t1.Year()+1, 1, 0, t1.Hour(), t1.Minute(), t1.Second(),
		0, time.UTC)

	// Generate the periodic timestamps for one year
	var ptlist []string
	for t := t1; t.Before(t2); t = t.AddDate(1, 0, 0) {
		_t := lastDateOfMonth(t)
		if _t.After(t2) {
			break
		}

		// Get the time zone offset of the current time based on the requested timezone.
		_, newOffsetSecs := _t.In(tz).Zone()
		if newOffsetSecs != offsetSecs {
			// Add offset of start time (removing the daylight saving time)
			_t = _t.Add(time.Duration(offsetSecs-newOffsetSecs) * time.Second)
		}

		// Round the hour
		_t = _t.Round(60 * time.Minute)
		// Append the time to the list
		ptlist = append(ptlist, _t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
