package period

import (
	"time"
)

// One Day Period
type OneDayPeriod struct{}

func (odp OneDayPeriod) GetMatchingTimestamps(t1, t2 time.Time) []string {
	var ptlist []string

	// Round the start time to hour.
	t1 = t1.Round(60 * time.Minute)

	// Generate the periodic timestamps for one day
	for t := t1; t.Before(t2); t = t.AddDate(0, 0, 1) {
		ptlist = append(ptlist, t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
