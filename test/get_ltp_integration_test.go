package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eaxis/ltp-provider/internal/cache"
	"github.com/eaxis/ltp-provider/internal/kraken"
	"github.com/eaxis/ltp-provider/internal/services"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/eaxis/ltp-provider/internal/config"
	"github.com/eaxis/ltp-provider/internal/server"
)

func TestIntegrationGetLTP(t *testing.T) {
	cfg := &config.Config{
		KrakenHost: "https://api.kraken.com",
		HttpAddr:   ":8099",
	}

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

	go func() {
		err := srv.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("Server failed: %v", err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/ltp", cfg.HttpAddr))

	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	if len(body) == 0 {
		t.Fatalf("Expected response body, got empty")
	}

	var ltpResp server.LtpResponse

	if err := json.Unmarshal(body, &ltpResp); err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}

	if len(ltpResp.Entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(ltpResp.Entries))
	}

	err = srv.Shutdown(context.TODO())

	if err != nil {
		t.Fatalf("Server shutdown failed: %v", err)
	}
}
