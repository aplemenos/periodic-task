package period

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(
	counter metrics.Counter, latency metrics.Histogram, s Service,
) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		next:           s,
	}
}

func (s *instrumentingService) GetPTList(
	ctx context.Context, period string, t1, t2 time.Time,
) (ptlist []string, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ptlist").Add(1)
		s.requestLatency.With("method", "ptlist").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.GetPTList(ctx, period, t1, t2)
}

func (s *instrumentingService) Alive(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "alive").Add(1)
		s.requestLatency.With("method", "alive").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.Alive(ctx)
}
