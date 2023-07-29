package periodictask

import (
	"encoding/json"
	"net/http"
	"periodic-task/pkg/period"
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

type PeriodHandler struct {
	S Service

	L *zap.SugaredLogger
}

// Router sets up all the routes for period service
func (h *PeriodHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.ptlist)

	return r
}

// ptlist retrieves the matching timestamps of a periodic task
func (h *PeriodHandler) ptlist(w http.ResponseWriter, r *http.Request) {
	// Get the period from url
	p := r.URL.Query().Get("period")
	if p == "" {
		h.L.Error("no period found")
		httpError(w, http.StatusBadRequest, periodRequired)
		return
	}

	// Get the timezone from url
	tz := r.URL.Query().Get("tz")
	if tz == "" {
		h.L.Error("no timezone found")
		httpError(w, http.StatusBadRequest, timezoneRequired)
		return
	}

	// Verify the requested timezone
	tzone, err := time.LoadLocation(tz)
	if err != nil {
		h.L.Error(tz + " is invalid timezone")
		httpError(w, http.StatusBadRequest, errInvalidTimezone(tz))
		return
	}

	// Get the t1 from url
	st1 := r.URL.Query().Get("t1")
	if st1 == "" {
		h.L.Error("no start point found")
		httpError(w, http.StatusBadRequest, startPointRequired)
		return
	}

	// Convert t1 as time
	t1, err := time.ParseInLocation(period.SUPPORTEDFORMAT, st1, tzone)
	if err != nil {
		h.L.Error("no supported format for " + st1)
		httpError(w, http.StatusBadRequest, errNoSupportedFormat(st1))
		return
	}

	// Get the t2 from url
	st2 := r.URL.Query().Get("t2")
	if st2 == "" {
		h.L.Error("no end point found")
		httpError(w, http.StatusBadRequest, endPointRequired)
		return
	}

	// Convert t2 as time
	t2, err := time.ParseInLocation(period.SUPPORTEDFORMAT, st2, tzone)
	if err != nil {
		h.L.Error("no supported format for " + st2)
		httpError(w, http.StatusBadRequest, errNoSupportedFormat(st2))
		return
	}

	// t1 should be before t2
	if t1.After(t2) {
		h.L.Error("t1 is after t2")
		httpError(w, http.StatusBadRequest, startAftertEndPoint)
		return
	}

	ptlist, err := h.S.GetPTList(r.Context(), p, t1, t2)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(ptlist); err != nil {
		h.L.Error(err.Error())
		httpError(w, http.StatusInternalServerError, err.Error())
	}
}

// errNoSupportedFormat is used when an invocation point (timestamp) could not be parsed
func errNoSupportedFormat(t string) string {
	return t + " is not a supported format. A valid timestamp format is " +
		period.SUPPORTEDFORMAT
}

// errInvalidTimezone is used when the timezone is invalid
func errInvalidTimezone(tz string) string {
	return tz + " is not a valid timezone."
}

type responseError struct {
	Status string `json:"status"`
	Desc   string `json:"desc"`
}

func httpError(w http.ResponseWriter, status int, desc string) {
	errorResponse := responseError{
		Status: "error",
		Desc:   desc,
	}

	writeResponse(w, status, errorResponse)
}

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
