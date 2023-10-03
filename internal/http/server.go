package http

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	periodictask "periodic-task/pkg/periodic-task"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Server holds the dependencies for a HTTP server.
type Server struct {
	Period periodictask.Service

	Logger *zap.SugaredLogger

	router chi.Router
}

// New returns a new HTTP server.
func New(ps periodictask.Service, logger *zap.SugaredLogger) *Server {
	s := &Server{
		Period: ps,
		Logger: logger,
	}

	r := chi.NewRouter()

	r.Use(s.recovery)
	r.Use(s.accessControl)
	r.Use(s.jsonMiddleware)
	r.Use(s.timeoutMiddleware)
	r.Use(s.loggingMiddleware)

	r.Route("/api/v1", func(r chi.Router) {
		ph := periodictask.PeriodHandler{
			S: s.Period,
			L: s.Logger,
		}
		r.Mount("/ptlist", ph.Router())
	})

	r.Get("/alive", s.aliveCheck)

	s.router = r

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) jsonMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		h.ServeHTTP(w, r)
	})
}

func (s *Server) timeoutMiddleware(h http.Handler) http.Handler {
	timeout := os.Getenv("SERVER_TIMEOUT")
	serverTimeout, _ := strconv.ParseInt(timeout, 10, 0)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(serverTimeout)*time.Second)
		defer cancel()
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggingMiddleware is a handy middleware function that logs out incoming requests
func (s *Server) loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r) // serve the original request

		duration := time.Since(start)

		// Log request details
		s.Logger.Infow("logging",
			zap.String("url", uri),
			zap.String("method", method),
			zap.Duration("took", duration))
	})
}

// recovery is a wrapper which will try to recover from any panic error and report it
func (s *Server) recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {
				s.Logger.Error("Failed to recover the panic: ", err)

				w.WriteHeader(http.StatusInternalServerError)

				json.NewEncoder(w).Encode(map[string]string{
					"status": "error",
					"desc":   "There was an internal server error",
				})
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func (s *Server) accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// response object
type response struct {
	Message string `json:"message"`
}

func (s *Server) aliveCheck(w http.ResponseWriter, r *http.Request) {
	// It could be used on the future to check the DB or REDIS aliveness
	// Now, it returns always "I am alive"
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response{Message: "I am Alive!"}); err != nil {
		s.Logger.Error("Failed to send Alive: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Serve gracefully serves our newly set up handler function
func (s *Server) Serve(server *http.Server, timeout int64) error {
	go func() {
		if err := server.ListenAndServe(); err != nil {
			s.Logger.Error("Failed to run the server: ", err)
		}
	}()

	// Create a deadline to wait for
	s.Logger.Debug("the server timeout is ", timeout)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(timeout)*time.Second)
	defer cancel()

	// Shut downs gracefully the server
	if err := server.Shutdown(ctx); err != nil {
		s.Logger.Error("Failed to shut off the server: ", err)
		return err
	}

	s.Logger.Info("shutting down gracefully")
	return nil
}
