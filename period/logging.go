package period

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type loggingService struct {
	logger *zap.SugaredLogger
	next   Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger *zap.SugaredLogger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) GetPTList(
	ctx context.Context, period string, t1, t2 time.Time,
) (ptlist []string, err error) {
	defer func(begin time.Time) {
		s.logger.Infow(
			"ptlist",
			zap.String("period", period),
			zap.Time("start point", t1),
			zap.Time("end point", t2),
			zap.Duration("took", time.Since(begin)),
			zap.Error(err),
		)
	}(time.Now())
	return s.next.GetPTList(ctx, period, t1, t2)
}

func (s *loggingService) Alive(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		s.logger.Infow(
			"alive",
			zap.Duration("took", time.Since(begin)),
			zap.Error(err),
		)
	}(time.Now())
	return s.next.Alive(ctx)
}
