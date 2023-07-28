package periodictask

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type loggingService struct {
	l *zap.SugaredLogger
	n Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger *zap.SugaredLogger, s Service) Service {
	return &loggingService{
		l: logger,
		n: s, // Next service
	}
}

func (s *loggingService) GetPTList(
	ctx context.Context, period string, t1, t2 time.Time,
) (ptlist []string, err error) {
	defer func(begin time.Time) {
		s.l.Infow(
			"ptlist",
			zap.String("period", period),
			zap.Time("start point", t1),
			zap.Time("end point", t2),
			zap.Duration("took", time.Since(begin)),
			zap.Error(err),
		)
	}(time.Now())
	return s.n.GetPTList(ctx, period, t1, t2)
}
