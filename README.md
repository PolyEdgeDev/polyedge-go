# PolyEdge Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/PolyEdgeDev/polyedge-go.svg)](https://pkg.go.dev/github.com/PolyEdgeDev/polyedge-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Documentation](https://img.shields.io/badge/docs-polyedge.dev-cyan)](https://polyedge.dev/docs)
[![Benchmark](https://img.shields.io/badge/benchmark-100%2B_nodes-green)](https://github.com/PolyEdgeDev/polyedge-stream-benchmark)

Official Go client library for **[PolyEdge](https://polyedge.dev)** — **Ultra-Low-Latency Polymarket Mempool Trade Streaming & Real-Time On-Chain Analytics API.**

> Engineered for Polymarket copy-trading bots, professional traders, and prediction market builders. Capture pending trades pre-block and track smart money with institutional-grade on-chain analytics.

---

## Key Features

- **4 First-Class Namespaces**: Aligned with OpenAPI 3.1:
  - `client.Streams`: Real-time SSE trade streaming, sessions & filter management.
  - `client.Analytics`: Polymarket trade attribution, trader leaderboards, market settlement audit.
  - `client.Account`: User profile, API keys, request telemetry, and ledger.
  - `client.Subscription`: Pricing tiers catalog, upgrade quotes, and promo codes.
- **Pure Standard Library**: Zero external runtime dependencies.
- **High-Performance SSE Engine**:
  - **45s Keep-Alive Watchdog**: Native timer watchdog detects silent dead TCP half-open connections.
  - **Atomic `Last-Event-ID` Resumption**: Preserves event progress across reconnects.
  - **Jittered Exponential Backoff**: Prevents server connection storms.
  - **Go Channels**: Delivers events via `<-chan *LiveTransaction` and errors via `<-chan error`.
- **Pre-Block Polymarket Trade Streaming**: Peered with 100+ Polygon nodes to capture pending Polymarket trades pre-block. The World's Fastest — [verify yourself](https://github.com/PolyEdgeDev/polyedge-stream-benchmark).

---

## Installation

```bash
go get github.com/PolyEdgeDev/polyedge-go
```

---

## Quick Start

### 1. Initialize Client

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PolyEdgeDev/polyedge-go"
)

func main() {
	apiKey := os.Getenv("POLYEDGE_API_KEY")
	client, err := polyedge.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	ctx := context.Background()

	// 1. Query Leaderboard
	leaderboard, err := client.Analytics.GetLeaderboard(ctx, 5, 0)
	if err != nil {
		log.Fatalf("GetLeaderboard error: %v", err)
	}
	fmt.Printf("Total active traders: %d\n", leaderboard.TotalCount)
	for i, t := range leaderboard.Traders {
		pnlUsd := float64(t.TotalPnl) / 1e6
		fmt.Printf("#%d %s (%s) - PnL: $%.2f\n", i+1, t.Profile.Address, t.Profile.Name, pnlUsd)
	}

	// 2. Real-Time Mempool SSE Trade Stream
	streamClient, err := client.Streams.Connect("your_stream_id")
	if err != nil {
		log.Fatalf("Connect error: %v", err)
	}

	txCh, errCh := streamClient.Subscribe(ctx)
	for {
		select {
		case tx, ok := <-txCh:
			if !ok {
				return
			}
			fmt.Printf("⚡ [Live Trade] Tx: %s...\n", tx.TxHash[:18])
			fmt.Printf("   Market: \"%s\"\n", tx.Market.Question)
			fmt.Printf("   Taker: %s (%s) | Shares: %s\n", tx.Taker.Outcome, tx.Taker.Side, tx.Taker.Shares)
			fmt.Printf("   Matched Makers: %d orders\n", len(tx.Makers))
		case sErr := <-errCh:
			if sErr != nil {
				log.Printf("Stream warning/reconnect: %v\n", sErr)
			}
		}
	}
}
```

---

## Complete API Service Reference (28 Methods)

### 1. Streams Service (`client.Streams`)

| Method | HTTP | Description |
| :--- | :--- | :--- |
| `Connect(streamID, opts...)` | `GET /streams/{id}` (SSE) | Connect to ultra-low-latency real-time SSE mempool trade stream |
| `List(ctx)` | `GET /streams` | List all configured data streams |
| `Create(ctx, req)` | `POST /streams` | Provision a new targeted filter stream (tags, addresses, series) |
| `Get(ctx, id)` | `GET /streams/{id}` | Get stream metadata and configuration by ID |
| `UpdateMetadata(ctx, id, req)` | `PUT /streams/{id}/meta` | Update stream name and description |
| `UpdateSubscription(ctx, id, req)`| `PUT /streams/{id}/subscription` | Hot-reload market filters (addresses, tags, series) on active stream |
| `Delete(ctx, id)` | `DELETE /streams/{id}` | Delete stream by ID |
| `GetActiveSessions(ctx)` | `GET /sessions` | List all active live SSE connections across streams |
| `GetSessionHistory(ctx, limit, offset)`| `GET /sessions/history` | Query historical SSE connection logs and durations |

### 2. Analytics Service (`client.Analytics`)

| Method | HTTP | Description |
| :--- | :--- | :--- |
| `GetLeaderboard(ctx, limit, offset)` | `GET /v2/analytics/leaderboard` | Top profitable traders ranked by realized PnL, volume, and ROI |
| `GetTrader(ctx, address)` | `GET /v2/traders/{address}` | Comprehensive trader intelligence profile, taker tier, and top tags |
| `GetTraderHourlyStats(ctx, address)` | `GET /v2/traders/{address}/hourly_stats` | 24-hour hourly trading PnL and volume breakdown |
| `GetTraderMarkets(ctx, addr, limit, offset)` | `GET /v2/traders/{address}/markets` | Historical market positions and settled outcomes by trader |
| `GetTraderMarketOrders(ctx, addr, id, ...)` | `GET /v2/traders/{address}/markets/{id}/orders` | Order fill details for a specific trader in a specific market |
| `GetMarket(ctx, id)` | `GET /v2/markets/{id}` | Prediction market detail, top 20 earners, and volume attribution |
| `GetDeposits(ctx, limit, offset)` | `GET /v2/analytics/deposits` | 60-day aggregated trader deposit summaries (>= 100 pUSD) |

### 3. Account Service (`client.Account`)

| Method | HTTP | Description |
| :--- | :--- | :--- |
| `GetProfile(ctx)` | `GET /v2/user/me` | Current user profile, active tier, and deposit balance |
| `ListKeys(ctx)` | `GET /v2/user/keys` | List all active API keys and statuses |
| `CreateKey(ctx, req)` | `POST /v2/user/keys` | Generate a new API key with custom name |
| `UpdateKey(ctx, key, req)` | `PATCH /v2/user/keys/{key}` | Enable, disable, or rename an API key |
| `DeleteKey(ctx, key)` | `DELETE /v2/user/keys/{key}` | Revoke and delete an API key |
| `GetLedger(ctx, limit, offset)` | `GET /v2/user/ledger` | Balance ledger transaction records |
| `GetTelemetry(ctx)` | `GET /v2/user/telemetry` | Daily request count and per-endpoint usage telemetry |
| `Withdraw(ctx, req)` | `POST /v2/user/withdraw` | Request balance withdrawal |

### 4. Subscription Service (`client.Subscription`)

| Method | HTTP | Description |
| :--- | :--- | :--- |
| `ListTiers(ctx)` | `GET /v2/tiers` | Public tier catalog, feature allowances, and pricing specifications |
| `GetQuote(ctx, req)` | `POST /v2/user/subscribe/quote` | Calculate quote for plan upgrade or billing cycle change |
| `Subscribe(ctx, req)` | `POST /v2/user/subscribe` | Purchase or upgrade tier subscription |
| `ValidatePromoCode(ctx, req)` | `POST /v2/user/promo-codes/validate` | Validate promotional discount code |

---

## License

[MIT](LICENSE) © 2026 [PolyEdge Labs](https://polyedge.dev)
