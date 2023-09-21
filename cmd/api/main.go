package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/romankravchuk/eldorado/internal/pkg/sl"
	"github.com/romankravchuk/eldorado/internal/server/http/api"
	"github.com/romankravchuk/eldorado/internal/server/http/handlers"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))
}

func main() {
	r := mux.NewRouter()
	r.NotFoundHandler = api.MakeHTTPHandler(handlers.Handle404)

	r.HandleFunc("/health", api.MakeHTTPHandler(handlers.HandleHealthCheck))

	srv := http.Server{
		Handler: r,
		Addr:    ":5555",
	}

	go func() {
		slog.Info("server starting", slog.String("port", "5555"))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("failed to serve server", sl.Err(err))
			os.Exit(1)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("failed to gracefully shutdown server", sl.Err(err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	slog.Info("server successfully sutdown")
	os.Exit(0)
}
