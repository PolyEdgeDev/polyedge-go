// Package polyedge provides the official Go client for PolyEdge's ultra-low latency
// Polymarket mempool trade streaming and real-time on-chain analytics API.
//
// Engineered for Polymarket copy-trading bots, professional market makers, and prediction market builders.
// PolyEdge peers directly with 100+ Polygon validator P2P gossip nodes to deliver pending CLOB trade
// executions pre-block, coupled with institutional-grade on-chain trade analytics.
//
// # Key Features
//
//   - Pre-Block Polymarket Trade Streaming: Peered with 100+ Polygon nodes to capture pending Polymarket
//     trades pre-block with verified speed advantage.
//   - Smart Money & Trader Analytics: Track smart money, profitable trader leaderboards, and institutional-grade on-chain analytics.
//   - Full OpenAPI 3.1 Domain Coverage: 4 core service namespaces covering Streams, Analytics, Account, and Subscription.
//   - Ultra-Low Latency SSE Engine: Automatic reconnection, 45-second heartbeat watchdog, and Last-Event-ID buffer resumption.
//
// For more information, benchmarks, and documentation, visit https://polyedge.dev and https://polyedge.dev/docs.
package polyedge
