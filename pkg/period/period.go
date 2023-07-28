package period

import "time"

// Constants for all supported periods
const (
	ONEHOUR  = "1h"
	ONEDAY   = "1d"
	ONEMONTH = "1mo"
	ONEYEAR  = "1y"
)

const SUPPORTEDFORMAT = "20060102T150405Z"

type Period struct {
	prd   string
	st    time.Time
	et    time.Time
	intvl time.Duration
}

// NewPeriod returns the period object at the runtime
// based on the requested period.
func NewPeriod(period string, startTime, endTime time.Time) *Period {
	et := endTime.Truncate(time.Second)
	switch period {
	case ONEHOUR:
		st := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startTime.Hour(),
			startTime.Minute(), startTime.Second(), 0, time.UTC)

		return &Period{
			prd:   period,
			st:    st,
			et:    et,
			intvl: time.Hour,
		}
	case ONEDAY:
		st := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startTime.Hour(),
			startTime.Minute(), startTime.Second(), 0, time.UTC)

		return &Period{
			prd:   period,
			st:    st,
			et:    et,
			intvl: 24 * time.Hour,
		}
	case ONEMONTH:
		st := time.Date(startTime.Year(), startTime.Month()+1, 0, startTime.Hour(),
			startTime.Minute(), startTime.Second(), 0, time.UTC)

		return &Period{
			prd:   period,
			st:    st,
			et:    et,
			intvl: 30 * 24 * time.Hour,
		}
	case ONEYEAR:
		st := time.Date(startTime.Year()+1, 1, 0, startTime.Hour(), startTime.Minute(),
			startTime.Second(), 0, time.UTC)

		return &Period{
			prd:   period,
			st:    st,
			et:    et,
			intvl: 365 * 24 * time.Hour,
		}
	default:
		return nil
	}
}

func (p *Period) GetMatchingTimestamps() []string {
	// Round the start time to hour.
	p.st = p.st.Round(60 * time.Minute)

	// Generate the periodic timestamps based on intervals
	var ptlist []string
	for t := p.st; t.Before(p.et); t = t.Add(p.intvl) {
		if p.prd == ONEMONTH || p.prd == ONEYEAR {
			t = lastDateOfMonth(t.Year(), t.Month(), t.Hour())
		}
		ptlist = append(ptlist, t.Format(SUPPORTEDFORMAT))
	}
	return ptlist
}

func lastDateOfMonth(year int, month time.Month, hour int) time.Time {
	return time.Date(year, month+1, 0, hour, 0, 0, 0, time.UTC)
}
