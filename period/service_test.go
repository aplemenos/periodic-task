package period

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestService_GetPTList(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	service := NewService(logger.Sugar())

	t1, _ := time.Parse("20060102T150405Z", "20210729T000000Z")
	t2, _ := time.Parse("20060102T150405Z", "20210729T050000Z")

	t.Run("SupportedPeriod", func(t *testing.T) {
		period := "1h"
		expected := []string{
			"20210729T000000Z",
			"20210729T010000Z",
			"20210729T020000Z",
			"20210729T030000Z",
			"20210729T040000Z",
		}

		result, err := service.GetPTList(context.Background(), period, t1, t2)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}

		if len(result) != len(expected) {
			t.Fatalf("Expected %d timestamps, but got %d", len(expected), len(result))
		}

		for i := range result {
			if result[i] != expected[i] {
				t.Errorf("Expected %s, but got %s", expected[i], result[i])
			}
		}
	})

	t.Run("UnsupportedPeriod", func(t *testing.T) {
		period := "invalid"
		_, err := service.GetPTList(context.Background(), period, t1, t2)

		if err == nil {
			t.Error("Expected error for unsupported period, but got no error")
		} else if err != errUnsupportedPeriod {
			t.Errorf("Expected unsupported period error, but got: %v", err)
		}
	})
}
