package polyedge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// StreamsService handles communication with Streams related methods of the PolyEdge API.
type StreamsService struct {
	client *Client
}

// Connect initializes a real-time SSE order stream client.
func (s *StreamsService) Connect(streamID string, opts ...StreamOptions) (*StreamClient, error) {
	return s.client.Stream(streamID, opts...)
}

// List: List User Streams.
// Returns custom and managed real-time trade streams configured by the user.
func (s *StreamsService) List(ctx context.Context) (*UserStreamsResponse, error) {
	endpoint := "/streams"
	targetURL := fmt.Sprintf("%s%s", s.client.streamBaseURL, endpoint)
	var res UserStreamsResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// Create: Create Filtered Order Stream.
// Provisions a new custom real-time stream with optional address, tag, or series filters.
func (s *StreamsService) Create(ctx context.Context, req *UserCreateStreamRequest) (*UserStreamResponse, error) {
	endpoint := "/streams"
	targetURL := fmt.Sprintf("%s%s", s.client.streamBaseURL, endpoint)
	var res UserStreamResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "POST", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// Get: Get Stream Metadata.
// Fetches metadata, active sessions count, and filter configuration for a specific stream.
func (s *StreamsService) Get(ctx context.Context, id string) (*UserStreamResponse, error) {
	endpoint := fmt.Sprintf("/streams/%s/meta", id)
	targetURL := fmt.Sprintf("%s%s", s.client.streamBaseURL, endpoint)
	var res UserStreamResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// Delete: Delete Stream.
// Permanently removes a stream and gracefully disconnects all connected listener sessions.
func (s *StreamsService) Delete(ctx context.Context, id string) (*DeleteStreamResponse, error) {
	endpoint := fmt.Sprintf("/streams/%s", id)
	targetURL := fmt.Sprintf("%s%s", s.client.streamBaseURL, endpoint)
	var res DeleteStreamResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "DELETE", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// UpdateMetadata: Update Stream Metadata.
// Updates operational attributes of the stream, such as nickname or enabled/disabled status.
func (s *StreamsService) UpdateMetadata(ctx context.Context, id string, req *UserUpdateStreamMetadataRequest) (*UserStreamResponse, error) {
	endpoint := fmt.Sprintf("/streams/%s/meta", id)
	targetURL := fmt.Sprintf("%s%s", s.client.streamBaseURL, endpoint)
	var res UserStreamResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "PUT", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// UpdateSubscription: Update Stream Filter Subscription.
// Dynamically mutates monitored wallet addresses, market tags, or series slugs without disconnecting SSE clients.
func (s *StreamsService) UpdateSubscription(ctx context.Context, id string, req *UpdateSubscriptionRequest) (*UserStreamResponse, error) {
	endpoint := fmt.Sprintf("/streams/%s/subscription", id)
	targetURL := fmt.Sprintf("%s%s", s.client.streamBaseURL, endpoint)
	var res UserStreamResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "PUT", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetActiveSessions: List All User Active Sessions.
// Returns all currently connected active SSE listeners for this account.
func (s *StreamsService) GetActiveSessions(ctx context.Context) (*UserActiveSessionsResponse, error) {
	endpoint := "/sessions"
	targetURL := fmt.Sprintf("%s%s", s.client.streamBaseURL, endpoint)
	var res UserActiveSessionsResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetSessionHistory: List User Stream Session History.
// Returns historical SSE connection logs including duration and pushed transactions count.
func (s *StreamsService) GetSessionHistory(ctx context.Context, limit int, offset int) (*UserSessionHistoryResponse, error) {
	endpoint := "/sessions/history"
	targetURL := fmt.Sprintf("%s%s", s.client.streamBaseURL, endpoint)
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := parsedURL.Query()
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	parsedURL.RawQuery = q.Encode()
	targetURL = parsedURL.String()
	var res UserSessionHistoryResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// AnalyticsService handles communication with Analytics related methods of the PolyEdge API.
type AnalyticsService struct {
	client *Client
}

// GetDeposits: 60-Day Trader Deposit Analytics.
// Returns 60-day aggregated deposit summary and individual transaction breakdown for top traders.
func (s *AnalyticsService) GetDeposits(ctx context.Context, limit int, offset int) (*DepositQueryResult, error) {
	endpoint := "/v2/analytics/deposits"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := parsedURL.Query()
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	parsedURL.RawQuery = q.Encode()
	targetURL = parsedURL.String()
	var res DepositQueryResult
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetLeaderboard: Top Traders PnL Leaderboard.
// Queries ranked traders across timeframes with PnL, volume, win rate, and performance metrics.
func (s *AnalyticsService) GetLeaderboard(ctx context.Context, limit int, offset int) (*LeaderboardResponse, error) {
	endpoint := "/v2/analytics/leaderboard"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := parsedURL.Query()
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	parsedURL.RawQuery = q.Encode()
	targetURL = parsedURL.String()
	var res LeaderboardResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetMarket: Get Prediction Market Detail.
// Fetches standardized market metadata matching on-chain condition IDs and NegRisk parameters.
func (s *AnalyticsService) GetMarket(ctx context.Context, id string) (*MarketDetailResponse, error) {
	endpoint := fmt.Sprintf("/v2/markets/%s", id)
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res MarketDetailResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetTrader: Get Trader Intelligence Profile.
// Fetches trader identity dossier, pUSD balance, total PnL, win rates, and ranking metrics.
func (s *AnalyticsService) GetTrader(ctx context.Context, address string) (*TraderProfileResponse, error) {
	endpoint := fmt.Sprintf("/v2/traders/%s", address)
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res TraderProfileResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetTraderHourlyStats: Get Hourly PnL Equity Curve.
// Provides 1-hour bucketed historical equity curves and trade metrics for a specific trader.
func (s *AnalyticsService) GetTraderHourlyStats(ctx context.Context, address string) (*TraderHourlyStatsResponse, error) {
	endpoint := fmt.Sprintf("/v2/traders/%s/hourly_stats", address)
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res TraderHourlyStatsResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetTraderMarkets: Get Trader Market Participation History.
// Queries resolved and active prediction markets traded by the specified wallet address.
func (s *AnalyticsService) GetTraderMarkets(ctx context.Context, address string, limit int, offset int) (*TraderHistoryResponse, error) {
	endpoint := fmt.Sprintf("/v2/traders/%s/markets", address)
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := parsedURL.Query()
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	parsedURL.RawQuery = q.Encode()
	targetURL = parsedURL.String()
	var res TraderHistoryResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetTraderMarketOrders: Get Trader Orders in Market.
// Retrieves granular fill events, side, price, and token outcomes for a trader in a given market.
func (s *AnalyticsService) GetTraderMarketOrders(ctx context.Context, address string, id string, limit int, offset int) (*UserOrdersResponse, error) {
	endpoint := fmt.Sprintf("/v2/traders/%s/markets/%s/orders", address, id)
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := parsedURL.Query()
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	parsedURL.RawQuery = q.Encode()
	targetURL = parsedURL.String()
	var res UserOrdersResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// AccountService handles communication with Account related methods of the PolyEdge API.
type AccountService struct {
	client *Client
}

// GetProfile: Get Current Profile.
// Retrieves the authenticated user profile, tier subscription status, and available USDC balance.
func (s *AccountService) GetProfile(ctx context.Context) (*UserProfile, error) {
	endpoint := "/v2/user/me"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res UserProfile
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// ListKeys: List API Keys.
// Lists all active and revoked API keys issued to the authenticated account.
func (s *AccountService) ListKeys(ctx context.Context) (*UserKeysResponse, error) {
	endpoint := "/v2/user/keys"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res UserKeysResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// CreateKey: Create API Key.
// Generates a new authenticated API key token with optional memo label.
func (s *AccountService) CreateKey(ctx context.Context, req *CreateAPIKeyRequest) (*UserAPIKey, error) {
	endpoint := "/v2/user/keys"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res UserAPIKey
	var reqErr error
	reqErr = s.client.doRequest(ctx, "POST", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// UpdateKey: Update API Key.
// Modifies label memo or operational status of an existing API key.
func (s *AccountService) UpdateKey(ctx context.Context, key string, req *UserUpdateKeyRequest) (*UserAPIKey, error) {
	endpoint := fmt.Sprintf("/v2/user/keys/%s", key)
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res UserAPIKey
	var reqErr error
	reqErr = s.client.doRequest(ctx, "PATCH", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// DeleteKey: Revoke API Key.
// Permanently revokes and deactivates an API key token.
func (s *AccountService) DeleteKey(ctx context.Context, key string) (*DeleteKeyResponse, error) {
	endpoint := fmt.Sprintf("/v2/user/keys/%s", key)
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res DeleteKeyResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "DELETE", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetLedger: List User Balance Ledger.
// Lists chronological ledger entries (subscription billing charges, deposits, withdrawals).
func (s *AccountService) GetLedger(ctx context.Context, limit int, offset int) (*UserLedgerResponse, error) {
	endpoint := "/v2/user/ledger"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	q := parsedURL.Query()
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	parsedURL.RawQuery = q.Encode()
	targetURL = parsedURL.String()
	var res UserLedgerResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// GetTelemetry: Get Account Telemetry & Limits.
// Provides live usage statistics, quota limits, and remaining streaming bandwidth.
func (s *AccountService) GetTelemetry(ctx context.Context) (*UserTelemetryResponse, error) {
	endpoint := "/v2/user/telemetry"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res UserTelemetryResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// Withdraw: Request Balance Withdrawal.
// Submits a request to withdraw unspent USDC balance to a designated Polygon address.
func (s *AccountService) Withdraw(ctx context.Context, req *UserWithdrawRequest) (*UserWithdrawResponse, error) {
	endpoint := "/v2/user/withdraw"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res UserWithdrawResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "POST", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// SubscriptionService handles communication with Subscription related methods of the PolyEdge API.
type SubscriptionService struct {
	client *Client
}

// ListTiers: List Public Tiers Catalog.
// Returns plan tiers (Free, Starter, Pro, Growth, Whale) with limits and pricing details.
func (s *SubscriptionService) ListTiers(ctx context.Context) ([]Tier, error) {
	endpoint := "/v2/tiers"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res []Tier
	var reqErr error
	reqErr = s.client.doRequest(ctx, "GET", targetURL, nil, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return res, nil
}

// GetQuote: Get Subscription Quote.
// Calculates prorated billing amounts for subscribing to or upgrading a tier plan.
func (s *SubscriptionService) GetQuote(ctx context.Context, req *QuoteSubscriptionRequest) (*QuoteSubscriptionResponse, error) {
	endpoint := "/v2/user/subscribe/quote"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res QuoteSubscriptionResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "POST", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// Subscribe: Purchase / Upgrade Tier Subscription.
// Executes purchase, renewal, or upgrade of an active subscription tier plan.
func (s *SubscriptionService) Subscribe(ctx context.Context, req *SubscribeRequest) (*SubscribeResponse, error) {
	endpoint := "/v2/user/subscribe"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res SubscribeResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "POST", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// ValidatePromoCode: Validate Promo Code.
// Verifies validity, discount percentage, and applicable plans for a promotion coupon code.
func (s *SubscriptionService) ValidatePromoCode(ctx context.Context, req *ValidatePromoCodeRequest) (*ValidatePromoCodeResponse, error) {
	endpoint := "/v2/user/promo-codes/validate"
	targetURL := fmt.Sprintf("%s%s", s.client.apiBaseURL, endpoint)
	var res ValidatePromoCodeResponse
	var reqErr error
	reqErr = s.client.doRequest(ctx, "POST", targetURL, req, &res)
	if reqErr != nil {
		return nil, reqErr
	}
	return &res, nil
}

// Client manages communication with the PolyEdge API.
type Client struct {
	apiKey        string
	apiBaseURL    string
	streamBaseURL string
	httpClient    *http.Client

	// Streams handles real-time SSE streams and filters
	Streams *StreamsService
	// Analytics handles trader dossiers, leaderboard, and markets
	Analytics *AnalyticsService
	// Account handles profile, API keys, ledger, and telemetry
	Account *AccountService
	// Subscription handles plan tiers, quotes, and subscriptions
	Subscription *SubscriptionService
}

type ClientOptions struct {
	APIBaseURL    string
	StreamBaseURL string
	HTTPClient    *http.Client
}

func NewClient(apiKey string, opts ...ClientOptions) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey cannot be empty")
	}
	c := &Client{
		apiKey:        apiKey,
		apiBaseURL:    "https://api.polyedge.dev",
		streamBaseURL: "https://stream.polyedge.dev",
		httpClient:    &http.Client{Timeout: 15 * time.Second},
	}
	if len(opts) > 0 {
		if opts[0].APIBaseURL != "" {
			c.apiBaseURL = strings.TrimRight(opts[0].APIBaseURL, "/")
		}
		if opts[0].StreamBaseURL != "" {
			c.streamBaseURL = strings.TrimRight(opts[0].StreamBaseURL, "/")
		}
		if opts[0].HTTPClient != nil {
			c.httpClient = opts[0].HTTPClient
		}
	}
	c.Streams = &StreamsService{client: c}
	c.Analytics = &AnalyticsService{client: c}
	c.Account = &AccountService{client: c}
	c.Subscription = &SubscriptionService{client: c}
	return c, nil
}

func (c *Client) doRequest(ctx context.Context, method, targetURL string, reqBody, respData any) error {
	var bodyReader *bytes.Reader
	if reqBody != nil {
		buf, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(buf)
	} else {
		bodyReader = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, targetURL, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-PolyEdge-Key", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("polyedge API error [%d]", resp.StatusCode)
	}
	if respData != nil {
		return json.NewDecoder(resp.Body).Decode(respData)
	}
	return nil
}

// Stream establishes real-time order stream directly.
func (c *Client) Stream(streamID string, opts ...StreamOptions) (*StreamClient, error) {
	var merged StreamOptions
	if len(opts) > 0 {
		merged = opts[0]
	}
	if merged.StreamBaseURL == "" {
		merged.StreamBaseURL = c.streamBaseURL
	}
	return NewStreamClient(streamID, c.apiKey, merged)
}
