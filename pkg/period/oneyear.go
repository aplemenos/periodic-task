package period

import "time"

// One Year Period
type OneYearPeriod struct{}

func (oyp OneYearPeriod) GetMatchingTimestamps(t1, t2 time.Time) []string {
	// Adjust t1 and t2 to the start and end of the month respectively
	t1 = time.Date(t1.Year()+1, 1, 0, t1.Hour(), t1.Minute(), t1.Second(),
		0, time.UTC).Round(60 * time.Minute)

	// Generate the periodic timestamps for one year
	var ptlist []string
	for t := t1; t.Before(t2); t = t.AddDate(1, 0, 0) {
		_t := lastDateOfMonth(t.Year(), t.Month(), t.Hour())
		ptlist = append(ptlist, _t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
