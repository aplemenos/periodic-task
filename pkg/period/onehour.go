package period

import "time"

// One Hour Period
type OneHourPeriod struct{}

func (ohp OneHourPeriod) GetMatchingTimestamps(t1, t2 time.Time) []string {
	var ptlist []string

	// Round the start time to hour.
	t1 = t1.Round(60 * time.Minute)

	// Generate the periodic timestamps for one hour
	for t := t1; t.Before(t2); t = t.Add(time.Hour) {
		ptlist = append(ptlist, t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
