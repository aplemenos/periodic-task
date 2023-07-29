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

// Period defines the interface of getting the matching timestamps
// following the rules of the strategy pattern.
type Period interface {
	GetMatchingTimestamps(t1, t2 time.Time) []string
}

func lastDateOfMonth(year int, month time.Month, hour int) time.Time {
	return time.Date(year, month+1, 0, hour, 0, 0, 0, time.UTC)
}

// NewPeriod returns the behavior of the matching timestamps at the runtime
// based on the requested period.
func NewPeriod(period string) Period {
	switch period {
	case ONEHOUR:
		return OneHourPeriod{}
	case ONEDAY:
		return OneDayPeriod{}
	case ONEMONTH:
		return OneMonthPeriod{}
	case ONEYEAR:
		return OneYearPeriod{}
	default:
		return nil
	}
}