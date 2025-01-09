package main

import (
	"context"
	"errors"
	"github.com/eaxis/ltp-provider/internal/cache"
	"github.com/eaxis/ltp-provider/internal/kraken"
	"github.com/eaxis/ltp-provider/internal/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eaxis/ltp-provider/internal/config"
	"github.com/eaxis/ltp-provider/internal/server"
)

func main() {
	cfg := config.Read()

	ltpCache := cache.NewLTPCache(1 * time.Minute) // can be configurable
	krakenClient := kraken.NewClient(cfg.KrakenHost)

	ltpService := services.NewLtpService(ltpCache, krakenClient)

	httpServer := server.NewHttpServer(ltpService)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/ltp", httpServer.GetLtp).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    cfg.HttpAddr,
		Handler: router,
	}

	stopped := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", cfg.HttpAddr)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Stopped")
}
