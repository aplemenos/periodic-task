package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"periodic-task/period/timestamp"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock implementation of period.Service for testing purposes
type mockPeriodService struct {
	mock.Mock
}

func (mps *mockPeriodService) GetPTList(
	ctx context.Context, period string, t1, t2 time.Time,
) ([]string, error) {
	args := mps.Called(ctx, period, t1, t2)
	return args.Get(0).([]string), args.Error(1)
}

func (mps *mockPeriodService) Alive(ctx context.Context) error {
	args := mps.Called(ctx)
	return args.Error(0)
}

func TestPeriodHandler_PTList(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create a mock period service
	mockService := new(mockPeriodService)

	// Create the period handler with the mock logger and service
	ph := &periodHandler{
		s:      mockService,
		logger: logger.Sugar(),
	}

	// Create a router and add the handler function
	r := chi.NewRouter()
	r.Get("/", ph.ptlist)

	// Helper function to create a request and execute it on the router
	makeRequest := func(method, path string, body []byte) (*http.Response, error) {
		req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		return rr.Result(), nil
	}

	// Test cases
	t.Run("ValidRequest", func(t *testing.T) {
		// Set up the mock service expectations and return values
		t1, _ := time.Parse(timestamp.SUPPORTEDFORMAT, "20210729T000000Z")
		t2, _ := time.Parse(timestamp.SUPPORTEDFORMAT, "20210729T040000Z")
		mockService.On("GetPTList", mock.Anything, "1h", t1, t2).
			Return([]string{
				"20210729T000000Z",
				"20210729T010000Z",
				"20210729T020000Z",
				"20210729T030000Z",
			}, nil)

		// Make a request to the router with valid query parameters
		resp, err := makeRequest("GET",
			"/?period=1h&t1=20210729T000000Z&t2=20210729T040000Z&tz=Europe/Athens", nil)
		assert.NoError(t, err, "Expected no error")

		// Check the response status code
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK")

		// Parse the response body
		var ptlist []string
		err = json.NewDecoder(resp.Body).Decode(&ptlist)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response body
		expected := []string{
			"20210729T000000Z",
			"20210729T010000Z",
			"20210729T020000Z",
			"20210729T030000Z",
		}
		assert.Len(t, ptlist, len(expected), "Expected number of timestamps")
		for i := range ptlist {
			assert.Equal(t, expected[i], ptlist[i], "Expected timestamp match")
		}
	})

	t.Run("MissingPeriod", func(t *testing.T) {
		// Make a request with missing period query parameter
		resp, err := makeRequest("GET",
			"/?t1=20210729T000000Z&t2=20210729T040000Z&tz=Europe/Athens", nil)
		assert.NoError(t, err, "Expected no error")

		var errorMsg responseError
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status Bad Request")
		assert.Equal(t, periodRequired, errorMsg.Desc)
	})

	t.Run("MissingStartPoint", func(t *testing.T) {
		// Make a request with missing period query parameter
		resp, err := makeRequest("GET",
			"/?period=1&t2=20210729T040000Z&tz=Europe/Athens", nil)
		assert.NoError(t, err, "Expected no error")

		var errorMsg responseError
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status Bad Request")
		assert.Equal(t, startPointRequired, errorMsg.Desc)
	})

	t.Run("MissingEndPoint", func(t *testing.T) {
		// Make a request with missing period query parameter
		resp, err := makeRequest("GET",
			"/?period=1&t1=20210729T000000Z&tz=Europe/Athens", nil)
		assert.NoError(t, err, "Expected no error")

		var errorMsg responseError
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status Bad Request")
		assert.Equal(t, endPointRequired, errorMsg.Desc)
	})

	t.Run("UnsupportedStartPointFormat", func(t *testing.T) {
		// Make a request with missing period query parameter
		resp, err := makeRequest("GET",
			"/?period=1&t1=29072021T000000Z&t2=20210729T040000Z&tz=Europe/Athens", nil)
		assert.NoError(t, err, "Expected no error")

		var errorMsg responseError
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status Bad Request")
		assert.Equal(t, errNoSupportedFormat("29072021T000000Z"), errorMsg.Desc)
	})

	t.Run("UnsupportedEndPointFormat", func(t *testing.T) {
		// Make a request with missing period query parameter
		resp, err := makeRequest("GET",
			"/?period=1&t1=20210729T000000Z&t2=29072021T040000Z&tz=Europe/Athens", nil)
		assert.NoError(t, err, "Expected no error")

		var errorMsg responseError
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status Bad Request")
		assert.Equal(t, errNoSupportedFormat("29072021T040000Z"), errorMsg.Desc)
	})

	t.Run("StartAfterEndPoint", func(t *testing.T) {
		// Make a request with missing period query parameter
		resp, err := makeRequest("GET",
			"/?period=1&t2=20210729T000000Z&t1=20210729T040000Z&tz=Europe/Athens", nil)
		assert.NoError(t, err, "Expected no error")

		var errorMsg responseError
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status Bad Request")
		assert.Equal(t, startAftertEndPoint, errorMsg.Desc)
	})

	t.Run("UnsupportedTimezone", func(t *testing.T) {
		// Make a request with missing period query parameter
		resp, err := makeRequest("GET",
			"/?period=1&t1=20210729T000000Z&t2=20210729T040000Z&tz=DangerZone", nil)
		assert.NoError(t, err, "Expected no error")

		var errorMsg responseError
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status Bad Request")
		assert.Equal(t, errInvalidTimezone("DangerZone"), errorMsg.Desc)
	})

	t.Run("MissingTimezone", func(t *testing.T) {
		// Make a request with missing period query parameter
		resp, err := makeRequest("GET",
			"/?period=1&t1=20210729T000000Z&t2=20210729T040000Z", nil)
		assert.NoError(t, err, "Expected no error")

		var errorMsg responseError
		err = json.NewDecoder(resp.Body).Decode(&errorMsg)
		assert.NoError(t, err, "Expected no error while decoding JSON")

		// Check the response status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status Bad Request")
		assert.Equal(t, timezoneRequired, errorMsg.Desc)
	})
}
