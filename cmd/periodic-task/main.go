package main

import (
	"net/http"
	"os"
	"periodic-task/period"
	"periodic-task/server"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	defaultServerAddr    = "0.0.0.0:8080"
	defaultRWTimeout     = "15"
	defaultIdleTimeout   = "15"
	defaultServerTimeout = "15"
)

// Run - sets up our application
func Run() error {
	// Build a production logger
	logger, _ := zap.NewProduction()
	defer func() {
		err := logger.Sync() // flushes buffer, if any
		logger.Error(err.Error())
	}()
	log := logger.Sugar()

	log.Info("setting up periodic task")

	// Setup period service
	var ps period.Service
	ps = period.NewService(log)
	ps = period.NewLoggingService(log, ps)
	// ps = period.NewInstrumentingService(
	// 	kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
	// 		Namespace: "api",
	// 		Subsystem: "period_service",
	// 		Name:      "request_count",
	// 		Help:      "Number of requests received.",
	// 	}, fieldKeys),
	// 	kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	// 		Namespace: "api",
	// 		Subsystem: "period_service",
	// 		Name:      "request_latency_microseconds",
	// 		Help:      "Total duration of requests in microseconds.",
	// 	}, fieldKeys),
	// 	ts)

	srv := server.New(ps, log)

	// Get the timeouts from the enviroment variable
	rwTimeout, err := strconv.ParseInt(envString("RW_TIMEOUT", defaultRWTimeout), 10, 0)
	if err != nil {
		log.Error("failed to parse RW_TIMEOUT")
		return err
	}
	rwt := time.Duration(rwTimeout) * time.Second

	idleTimeout, err := strconv.ParseInt(envString("IDLE_TIMEOUT", defaultIdleTimeout), 10, 0)
	if err != nil {
		log.Error("failed to parse IDLE_TIMEOUT")
		return err
	}
	idlet := time.Duration(idleTimeout) * time.Second

	server := &http.Server{
		Addr: envString("SERVER_ADDR", defaultServerAddr),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: rwt,
		ReadTimeout:  rwt,
		IdleTimeout:  idlet,
		Handler:      srv,
	}

	if err := srv.Serve(server, envString("SERVER_TIMEOUT", defaultServerTimeout)); err != nil {
		log.Error("failed to gracefully serve periodic task")
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		zap.S().Error(err)
		zap.S().Panic("Error starting up periodic task")
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
