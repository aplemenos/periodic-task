package timestamp

import "time"

// Constants for all supported periods
const (
	ONEHOUR  = "1h"
	ONEDAY   = "1d"
	ONEMONTH = "1mo"
	ONEYEAR  = "1y"
)

const SUPPORTEDFORMAT = "20060102T150405Z"

// Timestamp defines the interface of getting the matching timestamps
// following the rules of the strategy pattern.
type Timestamp interface {
	GetMatchingTimestamps(t1, t2 time.Time) []string
}

func lastDateOfMonth(year int, month time.Month, hour int) time.Time {
	return time.Date(year, month+1, 0, hour, 0, 0, 0, time.UTC)
}

// NewTimestamp returns the behavior of the timestamp at the runtime
// based on the requested period.
func NewTimestamp(period string) Timestamp {
	switch period {
	case ONEHOUR:
		return OneHourTimestamp{}
	case ONEDAY:
		return OneDayTimestamp{}
	case ONEMONTH:
		return OneMonthTimestamp{}
	case ONEYEAR:
		return OneYearTimestamp{}
	default:
		return nil
	}
}
