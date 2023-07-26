package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"periodic-task/period"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	Period period.Service

	Logger *zap.SugaredLogger

	router chi.Router
}

// New returns a new HTTP server.
func New(ps period.Service, logger *zap.SugaredLogger) *Server {
	s := &Server{
		Period: ps,
		Logger: logger,
	}

	r := chi.NewRouter()

	r.Use(accessControl)
	r.Use(jsonMiddleware)
	r.Use(timeoutMiddleware)

	r.Route("/api/v1", func(r chi.Router) {
		ph := periodHandler{s.Period, s.Logger}
		r.Mount("/ptlist", ph.router())
	})

	r.Get("/alive", s.aliveCheck)
	//r.Method("GET", "/metrics", promhttp.Handler())

	s.router = r

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func jsonMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		h.ServeHTTP(w, r)
	})
}

func timeoutMiddleware(h http.Handler) http.Handler {
	timeout := os.Getenv("SERVER_TIMEOUT")
	serverTimeout, _ := strconv.ParseInt(timeout, 10, 0)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(serverTimeout)*time.Second)
		defer cancel()
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		h.ServeHTTP(w, r)
	})
}

// response object
type response struct {
	Message string `json:"message"`
}

func (s *Server) aliveCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.Period.Alive(r.Context()); err != nil {
		s.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response{Message: "I am Alive!"}); err != nil {
		s.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Serve gracefully serves our newly set up handler function
func (s *Server) Serve(server *http.Server, stimeout string) error {
	go func() {
		if err := server.ListenAndServe(); err != nil {
			s.Logger.Error(err)
		}
	}()

	// Create a deadline to wait for
	serverTimeout, err := strconv.ParseInt(stimeout, 10, 0)
	if err != nil {
		s.Logger.Error(err)
		return err
	}
	s.Logger.Debug("the server timeout is ", serverTimeout)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	timeout := time.Duration(serverTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Shut downs gracefully the server
	if err := server.Shutdown(ctx); err != nil {
		s.Logger.Error(err)
		return err
	}

	s.Logger.Info("shutting down gracefully")
	return nil
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
