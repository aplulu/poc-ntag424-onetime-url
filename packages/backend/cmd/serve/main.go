package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aplulu/poc-ngtag424-onetime-url/packages/backend/internal/config"
	"github.com/aplulu/poc-ngtag424-onetime-url/packages/backend/internal/infrastructure/http"
)

func main() {
	if err := config.LoadConf(); err != nil {
		panic(err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-quitCh
		log.Info("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := http.StopServer(shutdownCtx); err != nil {
			log.Error(fmt.Sprintf("failed to stop server: %+v", err))
			os.Exit(1)
			return
		}
	}()

	if err := http.StartServer(log); err != nil {
		log.Error(fmt.Sprintf("failed to start server: %+v", err))
		os.Exit(1)
	}
}
