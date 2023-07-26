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
	Alive(ctx context.Context) error
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
	// Return the matching timestamps
	return t.GetMatchingTimestamps(t1, t2), nil
}

func (s *service) Alive(ctx context.Context) error {
	// TODO: Verify the DB aliveness
	return nil
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
