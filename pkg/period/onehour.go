package period

import "time"

// One Hour Period
type OneHourPeriod struct{}

func (ohp OneHourPeriod) GetMatchingTimestamps(t1, t2 time.Time, tz *time.Location) []string {
	var ptlist []string

	// Get the time zone offset of the start time based on the requested timezone.
	_, offsetSecs := t1.In(tz).Zone()

	// Round the start time to hour.
	t1 = t1.Round(60 * time.Minute)

	// Generate the periodic timestamps for one hour
	for t := t1; t.Before(t2); t = t.Add(time.Hour) {
		// Get the time zone offset of the current time based on the requested timezone.
		_, newOffsetSecs := t.In(tz).Zone()
		_t := t
		// Recalculate the offset if different from the start time's offset
		// Remove the daylight saving time
		if offsetSecs != newOffsetSecs {
			_t = _t.Add(time.Duration(offsetSecs-newOffsetSecs) * time.Second)
		}

		ptlist = append(ptlist, _t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
