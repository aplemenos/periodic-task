package period

import (
	"time"
)

// One Day Period
type OneDayPeriod struct{}

func (odp OneDayPeriod) GetMatchingTimestamps(t1, t2 time.Time) []string {
	var ptlist []string

	// Get the time zone offset of the start time.
	_, offsetSecs := t1.Zone()

	// Generate the periodic timestamps for one day,
	// taking daylight saving time changes into account.
	for t := t1; t.Before(t2); t = t.Add(24 * time.Hour) {
		// Get the time zone offset of the current time.
		_, newOffsetSecs := t.Zone()
		_t := t
		// Recalculate the offset if different from the start time's offset
		// Remove the daylight saving time
		if offsetSecs != newOffsetSecs {
			_t = _t.Add(time.Duration(newOffsetSecs) * time.Second)
		}

		// Round the hour
		_t = _t.Round(60 * time.Minute)
		// Append the time to the list
		ptlist = append(ptlist, _t.Format(SUPPORTEDFORMAT))
	}

	return ptlist
}
