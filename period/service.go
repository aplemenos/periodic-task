package period

import (
	"context"
	"errors"
	"periodic-task/period/timestamp"
	"time"

	"go.uber.org/zap"
)

// errUnsupportedPeriod is used when the requested period is not supported
var errUnsupportedPeriod = errors.New("unsupported period")

// Service is the interface that provides period-task methods
type Service interface {
	GetPTList(ctx context.Context, period string, t1, t2 time.Time) ([]string, error)
}

func (s *service) GetPTList(
	ctx context.Context, period string, t1, t2 time.Time,
) ([]string, error) {
	// Get a timestamp object
	t := timestamp.NewTimestamp(period)
	if t == nil {
		s.logger.Error(period, " is unsupported period")
		return nil, errUnsupportedPeriod
	}
	// Get the matching timestamps
	ptl := t.GetMatchingTimestamps(t1, t2)

	return ptl, nil
}

type service struct {
	logger *zap.SugaredLogger
}

// NewService creates a period service with necessary dependencies
func NewService(logger *zap.SugaredLogger) Service {
	return &service{
		logger: logger,
	}
}
