package polyedge

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoSDK_StreamReceiveAndHeartbeat(t *testing.T) {
	var heartbeatEmitted atomic.Bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected flusher")
		}

		// 1. Send heartbeat
		fmt.Fprintf(w, ": heartbeat\n\n")
		flusher.Flush()

		// 2. Send transaction event
		txPayload := `{"tx_hash":"0xgotest","timestamp":"2026-10-09T00:00:00Z","market":{"id":1,"outcomes":["Yes","No"],"token_ids":["1","2"]}}`
		fmt.Fprintf(w, "id: 8888\nevent: tx\ndata: %s\n\n", txPayload)
		flusher.Flush()

		time.Sleep(50 * time.Millisecond)
	}))
	defer server.Close()

	client, err := NewStreamClient("test-stream", "test_key", StreamOptions{
		StreamBaseURL:    server.URL,
		HeartbeatTimeout: 2 * time.Second,
		InitialBackoff:   50 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client.SetOnHeartbeat(func() {
		heartbeatEmitted.Store(true)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	txCh, errCh := client.Subscribe(ctx)

	select {
	case tx, ok := <-txCh:
		if !ok {
			t.Fatal("txCh closed prematurely")
		}
		if tx.TxHash != "0xgotest" {
			t.Fatalf("unexpected tx hash: %s", tx.TxHash)
		}
	case err := <-errCh:
		t.Fatalf("unexpected error: %v", err)
	case <-ctx.Done():
		t.Fatal("timed out waiting for tx")
	}

	if client.LastEventID() != 8888 {
		t.Fatalf("expected lastEventID=8888, got %d", client.LastEventID())
	}
	if !heartbeatEmitted.Load() {
		t.Fatal("expected heartbeat to be emitted")
	}
}

func TestGoSDK_Fatal401FastFail(t *testing.T) {
	var connCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connCount.Add(1)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error":"Unauthorized"}`)
	}))
	defer server.Close()

	client, err := NewStreamClient("fatal-stream", "bad_key", StreamOptions{
		StreamBaseURL:  server.URL,
		InitialBackoff: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, errCh := client.Subscribe(ctx)

	select {
	case err := <-errCh:
		if err != ErrFatalUnauthorized {
			t.Fatalf("expected ErrFatalUnauthorized, got %v", err)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for fatal error")
	}

	time.Sleep(50 * time.Millisecond)
	if connCount.Load() != 1 {
		t.Fatalf("expected exactly 1 connection attempt, got %d", connCount.Load())
	}
}

func TestGoSDK_WatchdogTimeoutAndResume(t *testing.T) {
	var connCount atomic.Int32
	var lastEventIDs []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := connCount.Add(1)
		lastEventIDs = append(lastEventIDs, r.Header.Get("Last-Event-ID"))

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)

		flusher := w.(http.Flusher)

		if count == 1 {
			// First connection: send message id 1234, then stall to trigger watchdog
			fmt.Fprintf(w, "id: 1234\nevent: tx\ndata: {\"tx_hash\":\"0xfirst\"}\n\n")
			flusher.Flush()
			// Stall
			time.Sleep(1 * time.Second)
		} else {
			// Second connection: send heartbeat and tx
			fmt.Fprintf(w, ": heartbeat\n\n")
			fmt.Fprintf(w, "id: 1235\nevent: tx\ndata: {\"tx_hash\":\"0xsecond\"}\n\n")
			flusher.Flush()
		}
	}))
	defer server.Close()

	client, err := NewStreamClient("watchdog-stream", "test_key", StreamOptions{
		StreamBaseURL:    server.URL,
		HeartbeatTimeout: 200 * time.Millisecond, // 200ms watchdog
		InitialBackoff:   30 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	txCh, _ := client.Subscribe(ctx)

	var receivedTxs []string
	for len(receivedTxs) < 2 {
		select {
		case tx, ok := <-txCh:
			if !ok {
				t.Fatal("txCh closed unexpectedly")
			}
			receivedTxs = append(receivedTxs, tx.TxHash)
		case <-ctx.Done():
			t.Fatalf("timed out, received %d txs: %v", len(receivedTxs), receivedTxs)
		}
	}

	if len(receivedTxs) != 2 || receivedTxs[0] != "0xfirst" || receivedTxs[1] != "0xsecond" {
		t.Fatalf("unexpected txs: %v", receivedTxs)
	}

	if connCount.Load() < 2 {
		t.Fatalf("expected at least 2 connections, got %d", connCount.Load())
	}

	if lastEventIDs[1] != "1234" {
		t.Fatalf("expected second connection to have Last-Event-ID=1234, got %s", lastEventIDs[1])
	}
}
