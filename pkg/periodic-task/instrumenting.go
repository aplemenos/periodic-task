package periodictask

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	c metrics.Counter
	l metrics.Histogram
	n Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(
	counter metrics.Counter, latency metrics.Histogram, s Service,
) Service {
	return &instrumentingService{
		c: counter,
		l: latency,
		n: s, // Next service
	}
}

func (s *instrumentingService) GetPTList(
	ctx context.Context, period string, t1, t2 time.Time,
) (ptlist []string, err error) {
	defer func(begin time.Time) {
		s.c.With("method", "ptlist").Add(1)
		s.l.With("method", "ptlist").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.n.GetPTList(ctx, period, t1, t2)
}
