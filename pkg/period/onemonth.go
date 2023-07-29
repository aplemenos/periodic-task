package period

import (
	"log"
	"time"
)

// One Month Period
type OneMonthPeriod struct{}

func (omp OneMonthPeriod) GetMatchingTimestamps(t1, t2 time.Time, tz *time.Location) []string {
	var ptlist []string

	// Get the time zone offset of the start time based on the requested timezone.
	_, offsetSecs := t1.In(tz).Zone()
	log.Println("Offset ", time.Duration(offsetSecs)*time.Second)

	// Generate the periodic timestamps for one month
	for t := t1; t.Before(t2); t = t.AddDate(0, 1, 0) {
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
