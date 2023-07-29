package period

import "time"

// One Month Period
type OneMonthPeriod struct{}

func (omp OneMonthPeriod) GetMatchingTimestamps(t1, t2 time.Time) []string {
	var ptlist []string

	// Round the start time to hour.
	t1 = t1.Round(60 * time.Minute)

	// Generate the periodic timestamps for one month
	for t := t1; t.Before(t2); t = t.AddDate(0, 1, 0) {
		_t := lastDateOfMonth(t.Year(), t.Month(), t.Hour())
		if _t.After(t2) {
			break
		}

		ptlist = append(ptlist, _t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}
