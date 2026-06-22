// Command api is the flight-meta search service entrypoint.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flightmeta/internal/config"
	"flightmeta/internal/httpapi"
	"flightmeta/internal/search"
	"flightmeta/internal/sources"
	"flightmeta/internal/sources/mock"
	"flightmeta/internal/visa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	var srcs []sources.Adapter
	if cfg.EnableMock {
		srcs = append(srcs, mock.New())
	}
	// Real adapters (Kiwi, Travelpayouts, Amadeus) are registered in later
	// phases once credentials exist; they are read from cfg server-side only.
	if len(srcs) == 0 {
		log.Error("no data sources configured; set FM_ENABLE_MOCK=true or add a real source")
		os.Exit(1)
	}

	resolver, err := visa.Load()
	if err != nil {
		log.Error("failed to load transit-visa data", "err", err)
		os.Exit(1)
	}

	orch := search.New(log, cfg.SourceTimeout, resolver, srcs...)
	handler := httpapi.New(orch, log, cfg.CORSOrigin)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("listening", "addr", cfg.Addr, "mock", cfg.EnableMock)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info("stopped")
}
