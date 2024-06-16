package http

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/rs/cors"

	"github.com/aplulu/poc-ngtag424-onetime-url/packages/backend/internal/config"
	"github.com/aplulu/poc-ngtag424-onetime-url/packages/backend/internal/util"
)

var server *http.Server

// StartServer starts the server
func StartServer(log *slog.Logger) error {
	serverMux := http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	apiMux.HandleFunc("GET /validate/{x}", func(w http.ResponseWriter, r *http.Request) {
		x := r.PathValue("x")
		// xが0-F 36文字でない場合は400
		if len(x) != 36 || !isHex(x) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		uid, err := hex.DecodeString(x[:14])
		if err != nil {
			log.Error(fmt.Sprintf("failed to decode UID: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		counter, err := hex.DecodeString(x[14:20])
		if err != nil {
			log.Error(fmt.Sprintf("failed to decode counter: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cmac, err := hex.DecodeString(x[20:])
		if err != nil {
			log.Error(fmt.Sprintf("failed to decode CMAC: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		expectedCMAC, err := util.CalculateSDMShortMacAES(config.Key(), uid, counter)
		if err != nil {
			log.Error(fmt.Sprintf("failed to calculate CMAC: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !bytes.Equal(expectedCMAC, cmac) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	serverMux.Handle("/v1/", http.StripPrefix("/v1", apiMux))

	ch := cors.New(cors.Options{
		AllowedOrigins:   config.CORSAllowedOrigins(),
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Connect-Protocol-Version", "Content-Type", "Authorization", "Admin-Authorization"},
		MaxAge:           config.CORSMaxAge(),
	})

	server = &http.Server{
		Addr:    net.JoinHostPort(config.Listen(), config.Port()),
		Handler: ch.Handler(serverMux),
	}

	listenHost := config.Listen()
	if listenHost == "" {
		listenHost = "localhost"
	}
	log.Info(fmt.Sprintf("Server started at http://%s:%s", listenHost, config.Port()))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// StopServer stops the server
func StopServer(ctx context.Context) error {
	return server.Shutdown(ctx)
}
