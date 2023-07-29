package periodictask

import (
	"context"
	"errors"
	"periodic-task/pkg/period"
	"time"

	"go.uber.org/zap"
)

// errUnsupportedPeriod is used when the requested period is not supported
var errUnsupportedPeriod = errors.New("unsupported period")

// Service is the interface that provides period-task methods
type Service interface {
	GetPTList(
		ctx context.Context, period string, t1, t2 time.Time, tz *time.Location,
	) ([]string, error)
}

func (s *service) GetPTList(
	ctx context.Context, p string, t1, t2 time.Time, tz *time.Location,
) ([]string, error) {
	// Get a period object
	period := period.NewPeriod(p)
	if period == nil {
		s.l.Error(p, " is unsupported period")
		return nil, errUnsupportedPeriod
	}
	// Return the matching timestamps
	return period.GetMatchingTimestamps(t1, t2, tz), nil
}

type service struct {
	l *zap.SugaredLogger
}

// NewService creates a period service with necessary dependencies
func NewService(logger *zap.SugaredLogger) Service {
	return &service{
		l: logger,
	}
}
