package server

import (
	"encoding/json"
	"net/http"
	"periodic-task/period"
	"periodic-task/period/timestamp"
	"time"

	"github.com/go-chi/chi"

	"go.uber.org/zap"
)

var (
	periodRequired      = "period required"
	startPointRequired  = "start point required"
	endPointRequired    = "end point required"
	startAftertEndPoint = "start point should be before end point"
	timezoneRequired    = "timezone required"
)

type periodHandler struct {
	s period.Service

	logger *zap.SugaredLogger
}

// router sets up all the routes for period service
func (h *periodHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.ptlist)

	return r
}

// ptlist retrieves the matching timestamps of a periodic task
func (h *periodHandler) ptlist(w http.ResponseWriter, r *http.Request) {
	// Get the period from url
	period := r.URL.Query().Get("period")
	if period == "" {
		h.logger.Error("no period found")
		httpError(w, http.StatusBadRequest, periodRequired)
		return
	}

	// Get the t1 from url
	st1 := r.URL.Query().Get("t1")
	if st1 == "" {
		h.logger.Error("no start point found")
		httpError(w, http.StatusBadRequest, startPointRequired)
		return
	}

	// Convert t1 as time
	t1, err := time.Parse(timestamp.SUPPORTEDFORMAT, st1)
	if err != nil {
		h.logger.Error("no supported format for " + st1)
		httpError(w, http.StatusBadRequest, errNoSupportedFormat(st1))
		return
	}

	// Get the t2 from url
	st2 := r.URL.Query().Get("t2")
	if st2 == "" {
		h.logger.Error("no end point found")
		httpError(w, http.StatusBadRequest, endPointRequired)
		return
	}

	// Convert t2 as time
	t2, err := time.Parse(timestamp.SUPPORTEDFORMAT, st2)
	if err != nil {
		h.logger.Error("no supported format for " + st2)
		httpError(w, http.StatusBadRequest, errNoSupportedFormat(st2))
		return
	}

	// t1 should be before t2
	if t1.After(t2) {
		h.logger.Error("t1 is after t2")
		httpError(w, http.StatusBadRequest, startAftertEndPoint)
		return
	}

	// Get the timezone from url
	tz := r.URL.Query().Get("tz")
	if tz == "" {
		h.logger.Error("no timezone found")
		httpError(w, http.StatusBadRequest, timezoneRequired)
		return
	}

	// Verify the requested timezone
	_, err = time.LoadLocation(tz)
	if err != nil {
		h.logger.Error(tz+" is invalid timezone")
		httpError(w, http.StatusBadRequest, errInvalidTimezone(tz))
		return
	}

	ptlist, err := h.s.GetPTList(r.Context(), period, t1, t2)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(ptlist); err != nil {
		h.logger.Error(err.Error())
		httpError(w, http.StatusInternalServerError, err.Error())
	}
}

// errNoSupportedFormat is used when an invocation point (timestamp) could not be parsed
func errNoSupportedFormat(t string) string {
	return t + " is not a supported format. A valid timestamp format is " +
		timestamp.SUPPORTEDFORMAT
}

// errInvalidTimezone is used when the timezone is invalid
func errInvalidTimezone(tz string) string {
	return tz + " is not a valid timezone."
}
