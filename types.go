package polyedge

// Side represents the trade execution order side (BUY or SELL).
type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

// ActiveSessionInfo ActiveSessionInfo represents the ActiveSessionInfo data model.
type ActiveSessionInfo struct {
	// ApiKey: Masked API credential used to authenticate the stream connection
	ApiKey string `json:"api_key,omitempty"`
	// ClientIp: Remote IP address of the connected client
	ClientIp string `json:"client_ip,omitempty"`
	// ConnectedAt: ISO 8601 timestamp when the SSE connection was established
	ConnectedAt string `json:"connected_at,omitempty"`
	// DurationSeconds: Total active connection uptime in seconds
	DurationSeconds int64 `json:"duration_seconds,omitempty"`
	// PushedTxCount: Total number of live matched order events delivered over this connection
	PushedTxCount int64 `json:"pushed_tx_count,omitempty"`
	// StreamId: Unique persistent stream identifier (UUID)
	StreamId string `json:"stream_id,omitempty"`
	// UserAgent: HTTP User-Agent identifier of the client library or application
	UserAgent string `json:"user_agent,omitempty"`
}

// CreateAPIKeyRequest CreateAPIKeyRequest represents the CreateAPIKeyRequest data model.
type CreateAPIKeyRequest struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

// DeleteKeyResponse DeleteKeyResponse represents the DeleteKeyResponse data model.
type DeleteKeyResponse struct {
	ID string `json:"id,omitempty"`
}

// DeleteStreamResponse DeleteStreamResponse represents the DeleteStreamResponse data model.
type DeleteStreamResponse struct {
	ID string `json:"id,omitempty"`
}

// DepositItem DepositItem represents the DepositItem data model.
type DepositItem struct {
	Address      string `json:"address,omitempty"`
	DepositCount int64  `json:"deposit_count,omitempty"`
	// FirstDepositAt: Earliest deposit time in 60-day window (RFC3339)
	FirstDepositAt string `json:"first_deposit_at,omitempty"`
	FirstTradeAt   string `json:"first_trade_at,omitempty"`
	// LastDepositAt: Most recent deposit time (RFC3339)
	LastDepositAt string `json:"last_deposit_at,omitempty"`
	LastTradeAt   string `json:"last_trade_at,omitempty"`
	// TotalDeposit: 6 decimals micro-USD integer as string (e.g. "100000000")
	TotalDeposit string `json:"total_deposit,omitempty"`
}

// DepositQueryResult DepositQueryResult represents the DepositQueryResult data model.
type DepositQueryResult struct {
	Items     []DepositItem `json:"items,omitempty"`
	Total     int64         `json:"total,omitempty"`
	UpdatedAt string        `json:"updated_at,omitempty"`
}

// ErrorResponse ErrorResponse represents the ErrorResponse data model.
type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// HistoryMarket HistoryMarket represents the HistoryMarket data model.
type HistoryMarket struct {
	// Icon: Market thumbnail image URL
	Icon string `json:"icon,omitempty"`
	// ID: Unique numeric prediction market identifier
	ID int64 `json:"id,omitempty"`
	// Outcomes: Raw JSON array of outcome display names, e.g. ["Yes", "No"]
	Outcomes [2]string `json:"outcomes,omitempty"`
	// Question: Full title question describing the market condition
	Question string `json:"question,omitempty"`
	// ResolvedAt: Timestamp when market officially settled (RFC 3339 UTC), or null if active
	ResolvedAt string `json:"resolved_at,omitempty"`
	// Result: Resolution outcome (-2: unresolved, -1: void/invalid, 0-100: payout percentage for Outcome 0, e.g. 100 = Outcome 0 won, 0 = Outcome 1 won, 50 = 50/50 split)
	Result int64 `json:"result,omitempty"`
	// SeriesSlug: URL slug of the recurring sports league or tournament series
	SeriesSlug string `json:"series_slug,omitempty"`
	// Slug: URL-friendly slug identifying the market on Polymarket
	Slug string `json:"slug,omitempty"`
}

// HourlyStat HourlyStat represents the HourlyStat data model.
type HourlyStat struct {
	// CostBasis: Aggregate capital cost basis in micro-USD (6 decimals)
	CostBasis int64 `json:"cost_basis,omitempty"`
	// Fee: Total trading fees incurred in micro-USD (6 decimals)
	Fee int64 `json:"fee,omitempty"`
	// Hour: 1-hour UTC bucket timestamp (ISO 8601)
	Hour string `json:"hour,omitempty"`
	// MakerVolume: Maker trading volume in micro-USD (6 decimals)
	MakerVolume int64 `json:"maker_volume,omitempty"`
	// MarketCount: Number of distinct prediction markets active or settled in this hour
	MarketCount int64 `json:"market_count,omitempty"`
	// PnL: Realized net PnL in micro-USD (6 decimals)
	PnL int64 `json:"pnl,omitempty"`
	// TakerVolume: Taker trading volume in micro-USD (6 decimals)
	TakerVolume int64 `json:"taker_volume,omitempty"`
	// WinCount: Number of profitable markets settled in this hour
	WinCount int64 `json:"win_count,omitempty"`
}

// LeaderboardResponse LeaderboardResponse represents the LeaderboardResponse data model.
type LeaderboardResponse struct {
	// TotalCount: Total count of matching traders for pagination
	TotalCount int64 `json:"total_count,omitempty"`
	// Traders: Ranked array of trader performance summaries
	Traders []TraderSummary `json:"traders,omitempty"`
}

// LiveOrder LiveOrder represents the LiveOrder data model.
type LiveOrder struct {
	// Fee: Execution and relayer fee paid in micro-USDC
	Fee string `json:"fee,omitempty"`
	// Order: Underlying signed off-chain CLOB limit order metadata
	Order OrderInfo `json:"order,omitempty"`
	// Outcome: Market outcome label (e.g. "Yes", "No", "Up", "Down")
	Outcome string `json:"outcome,omitempty"`
	// Shares: Filled outcome shares amount represented as a 10^6 fixed-point decimal string
	Shares string `json:"shares,omitempty"`
	// Side: Order execution direction: "BUY" or "SELL"
	Side Side `json:"side,omitempty"`
	// TokenIdsIndex: 0-based index matching market.token_ids and market.outcomes
	TokenIdsIndex uint8 `json:"token_ids_index,omitempty"`
	// USDC: Filled USDC notional amount in micro-USDC (10^6 scale) decimal string
	USDC string `json:"usdc,omitempty"`
	// User: Trader account identity (wallet address and optional display pseudonym)
	User MonitorTrader `json:"user,omitempty"`
}

// LiveTransaction LiveTransaction represents the LiveTransaction data model.
type LiveTransaction struct {
	// Makers: List of matched maker limit orders filled in this transaction
	Makers []LiveOrder `json:"makers,omitempty"`
	// Market: Associated Polymarket prediction market condition and metadata
	Market Market `json:"market,omitempty"`
	// Taker: Taker order execution details that initiated the match against the CLOB
	Taker LiveOrder `json:"taker,omitempty"`
	// Timestamp: ISO 8601 UTC timestamp (millisecond precision) indicating when the pending transaction was detected in the mempool
	Timestamp string `json:"timestamp,omitempty"`
	// TxHash: Canonical Ethereum/Polygon on-chain transaction hash confirming the match execution
	TxHash string `json:"tx_hash,omitempty"`
}

// Market Market represents the Market data model.
type Market struct {
	// ConditionId: 32-byte hexadecimal condition ID from Gnosis Conditional Tokens
	ConditionId string `json:"condition_id,omitempty"`
	// EventSlug: URL slug of the parent event containing this market
	EventSlug string `json:"event_slug,omitempty"`
	// GroupItemTitle: Sub-category item title within a grouped multi-market event
	GroupItemTitle string `json:"group_item_title,omitempty"`
	// ID: Polymarket numeric prediction market identifier (e.g. 3688221)
	ID int64 `json:"id,omitempty"`
	// NegRisk: True if market operates under multi-outcome negative risk adapter
	NegRisk bool `json:"neg_risk,omitempty"`
	// Outcomes: Array of market outcome labels: [Outcome 1, Outcome 2] (e.g. ["Yes", "No"])
	Outcomes [2]string `json:"outcomes,omitempty"`
	// Question: Full question describing the prediction market condition
	Question string `json:"question,omitempty"`
	// SeriesSlug: URL slug of the recurring sports league or tournament series
	SeriesSlug string `json:"series_slug,omitempty"`
	// Slug: URL-friendly slug identifying the market on Polymarket
	Slug string `json:"slug,omitempty"`
	// SportsMarketType: Sports market classification (e.g. moneyline, spread, over/under)
	SportsMarketType string `json:"sports_market_type,omitempty"`
	// StartDate: ISO 8601 UTC timestamp when trading began
	StartDate string `json:"start_date,omitempty"`
	// TagsSlug: Canonical categorical tags for content filtering
	TagsSlug []string `json:"tags_slug,omitempty"`
	// TokenIds: ERC-1155 token IDs: [Token 1, Token 2]
	TokenIds [2]string `json:"token_ids,omitempty"`
}

// MarketDetailResponse MarketDetailResponse represents the MarketDetailResponse data model.
type MarketDetailResponse struct {
	// Market: Complete standardized market metadata matching exporter format
	Market MarketMetadata `json:"market,omitempty"`
	// TopEarners: Ranked list of top 20 profitable traders in this market
	TopEarners []MarketEarner `json:"top_earners,omitempty"`
	// TotalFees: Total protocol and liquidity fees collected in micro-USD (6 decimals)
	TotalFees int64 `json:"total_fees,omitempty"`
	// TotalPnlDistributed: Total positive realized PnL distributed to winners in micro-USD (6 decimals)
	TotalPnlDistributed int64 `json:"total_pnl_distributed,omitempty"`
	// TotalTraders: Total count of distinct wallet addresses that traded this market
	TotalTraders int64 `json:"total_traders,omitempty"`
	// TotalVolume: Cumulative trading volume across all participants in micro-USD (6 decimals)
	TotalVolume int64 `json:"total_volume,omitempty"`
	// TotalWinners: Total number of distinct traders with positive realized PnL
	TotalWinners int64 `json:"total_winners,omitempty"`
}

// MarketEarner MarketEarner represents the MarketEarner data model.
type MarketEarner struct {
	// Orders: Executed fill orders for this earner (included when include_orders=true)
	Orders []UserOrder `json:"orders,omitempty"`
	// PnL: Realized net profit and loss in micro-USD (6 decimals, e.g. 1000000 = 1 USD)
	PnL int64 `json:"pnl,omitempty"`
	// Profile: Public identity dossier of the earner
	Profile TraderIdentity `json:"profile,omitempty"`
	// Rank: Rank position among top earners in this market (1-based)
	Rank int64 `json:"rank,omitempty"`
	// Roi: Return on investment percentage (e.g. 25.5 for 25.5%)
	Roi float64 `json:"roi,omitempty"`
	// Volume: Total traded volume in micro-USD (6 decimals)
	Volume int64 `json:"volume,omitempty"`
}

// MarketMetadata MarketMetadata represents the MarketMetadata data model.
type MarketMetadata struct {
	// ConditionId: 0x-prefixed 66-character CTF condition identifier
	ConditionId string `json:"condition_id,omitempty"`
	// CreatedAt: Timestamp when market was indexed in database (RFC 3339 UTC)
	CreatedAt string `json:"created_at,omitempty"`
	// Description: Extended details and market resolution criteria
	Description string `json:"description,omitempty"`
	// EndDate: Estimated resolution date or scheduled closing time (RFC 3339 UTC)
	EndDate string `json:"end_date,omitempty"`
	// EventSlug: URL slug of the parent event containing this market
	EventSlug string `json:"event_slug,omitempty"`
	// FeeSchedule: Detailed fee schedule JSON specification
	FeeSchedule map[string]any `json:"fee_schedule,omitempty"`
	// FeeType: Fee model classification (e.g. "dynamic", "fixed")
	FeeType string `json:"fee_type,omitempty"`
	// FeesEnabled: Whether trading fees are activated on this market
	FeesEnabled bool `json:"fees_enabled,omitempty"`
	// GroupItemThreshold: Threshold for numerical interval outcomes if applicable
	GroupItemThreshold int64 `json:"group_item_threshold,omitempty"`
	// GroupItemTitle: Sub-category item title within a grouped multi-market event
	GroupItemTitle string `json:"group_item_title,omitempty"`
	// Icon: Market thumbnail image URL
	Icon string `json:"icon,omitempty"`
	// ID: Unique numeric prediction market identifier
	ID int64 `json:"id,omitempty"`
	// NegRisk: True if market operates under multi-outcome negative risk adapter
	NegRisk bool `json:"neg_risk,omitempty"`
	// Outcomes: Raw JSON array of outcome display names, e.g. ["Yes", "No"]
	Outcomes [2]string `json:"outcomes,omitempty"`
	// Question: Full title question describing the market condition
	Question string `json:"question,omitempty"`
	// ResolutionSource: Source URL or oracle specification for resolution
	ResolutionSource string `json:"resolution_source,omitempty"`
	// ResolvedAt: Timestamp when the market officially settled (RFC 3339 UTC), or null if active
	ResolvedAt string `json:"resolved_at,omitempty"`
	// Result: Resolution outcome (-2: unresolved, -1: void/invalid, 0-100: payout percentage for Outcome 0, e.g. 100 = Outcome 0 won, 0 = Outcome 1 won, 50 = 50/50 split)
	Result int64 `json:"result,omitempty"`
	// SeriesSlug: URL slug of the recurring sports league or tournament series
	SeriesSlug string `json:"series_slug,omitempty"`
	// Slug: URL-friendly slug identifying the market on Polymarket
	Slug string `json:"slug,omitempty"`
	// SportsMarketType: Sports market classification type if applicable
	SportsMarketType string `json:"sports_market_type,omitempty"`
	// StartDate: Timestamp when trading began (RFC 3339 UTC)
	StartDate string `json:"start_date,omitempty"`
	// TagsSlug: Raw JSON array of topic tag URL slugs
	TagsSlug []string `json:"tags_slug,omitempty"`
	// TokenIds: Array of 2 ERC-1155 token IDs corresponding to outcomes
	TokenIds [2]string `json:"token_ids,omitempty"`
}

// MarketSettlementDetail MarketSettlementDetail represents the MarketSettlementDetail data model.
type MarketSettlementDetail struct {
	// LastTradeAt: ISO 8601 UTC timestamp of the trader's last trade in this market
	LastTradeAt string `json:"last_trade_at,omitempty"`
	// Market: Concise standardized market metadata
	Market HistoryMarket `json:"market,omitempty"`
	// Orders: Detailed order fill events (present only when include_orders=true)
	Orders []UserOrder `json:"orders,omitempty"`
	// PnL: Realized net PnL in micro-USD (6 decimals)
	PnL int64 `json:"pnl,omitempty"`
	// Roi: Return on investment percentage in this market
	Roi float64 `json:"roi,omitempty"`
	// Volume: Total traded volume in this market in micro-USD (6 decimals)
	Volume int64 `json:"volume,omitempty"`
}

// MonitorTrader MonitorTrader represents the MonitorTrader data model.
type MonitorTrader struct {
	// Address: Canonical 42-character hexadecimal Ethereum/Polygon wallet address
	Address string `json:"address,omitempty"`
	// Name: Polymarket public username
	Name string `json:"name,omitempty"`
}

// OrderInfo OrderInfo represents the OrderInfo data model.
type OrderInfo struct {
	// OrderHash: Cryptographic EIP-712 order hash identifying the unique off-chain order
	OrderHash string `json:"order_hash,omitempty"`
	// Shares: Original clob order shares capacity in 10^6 scale
	Shares string `json:"shares,omitempty"`
	// SignatureType: Signature scheme encoding (0: Direct EOA, 1: Magic / Proxy Wallet, 2: Gnosis Safe, 3: Deposit Wallet (ERC-1271))
	SignatureType int64 `json:"signature_type,omitempty"`
	// Timestamp: Unix epoch seconds when the order was signed
	Timestamp string `json:"timestamp,omitempty"`
	// USDC: Original clob order notional value in micro-USDC (10^6 scale)
	USDC string `json:"usdc,omitempty"`
}

// QuoteSubscriptionRequest QuoteSubscriptionRequest represents the QuoteSubscriptionRequest data model.
type QuoteSubscriptionRequest struct {
	BillingCycle string `json:"billing_cycle,omitempty"`
	PromoCode    string `json:"promo_code,omitempty"`
	Tier         string `json:"tier,omitempty"`
}

// QuoteSubscriptionResponse QuoteSubscriptionResponse represents the QuoteSubscriptionResponse data model.
type QuoteSubscriptionResponse struct {
	AmountToPay         int64  `json:"amount_to_pay,omitempty"`
	BillingCycle        string `json:"billing_cycle,omitempty"`
	CanAfford           bool   `json:"can_afford,omitempty"`
	CurrentBalance      int64  `json:"current_balance,omitempty"`
	CurrentBillingCycle string `json:"current_billing_cycle,omitempty"`
	CurrentLicense      string `json:"current_license,omitempty"`
	CurrentTier         string `json:"current_tier,omitempty"`
	DiscountAmount      int64  `json:"discount_amount,omitempty"`
	DurationDays        int64  `json:"duration_days,omitempty"`
	IsActive            bool   `json:"is_active,omitempty"`
	License             string `json:"license,omitempty"`
	Message             string `json:"message,omitempty"`
	OriginalPrice       int64  `json:"original_price,omitempty"`
	ProratedCredit      int64  `json:"prorated_credit,omitempty"`
	Tier                string `json:"tier,omitempty"`
}

// SubscribeRequest SubscribeRequest represents the SubscribeRequest data model.
type SubscribeRequest struct {
	BillingCycle   string `json:"billing_cycle,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	PromoCode      string `json:"promo_code,omitempty"`
	Tier           string `json:"tier,omitempty"`
}

// SubscribeResponse SubscribeResponse represents the SubscribeResponse data model.
type SubscribeResponse struct {
	Subscription UserSubscription `json:"subscription,omitempty"`
	User         UserProfile      `json:"user,omitempty"`
}

// TagStat TagStat represents the TagStat data model.
type TagStat struct {
	// MakerVolume: Passive maker volume in micro-USD (6 decimals)
	MakerVolume int64 `json:"maker_volume,omitempty"`
	// MarketsCount: Number of markets traded in this tag
	MarketsCount int64 `json:"markets_count,omitempty"`
	// PnL: Realized net PnL in micro-USD (6 decimals)
	PnL int64 `json:"pnl,omitempty"`
	// Roi: Return on investment percentage for this tag
	Roi float64 `json:"roi,omitempty"`
	// Slug: Category tag URL slug (e.g. "crypto", "politics", "sports")
	Slug string `json:"slug,omitempty"`
	// TakerVolume: Aggressive taker volume in micro-USD (6 decimals)
	TakerVolume int64 `json:"taker_volume,omitempty"`
	// Volume: Aggregated trading volume in micro-USD (6 decimals)
	Volume int64 `json:"volume,omitempty"`
	// WinRate: Win rate percentage for this tag (0.0 to 100.0)
	WinRate float64 `json:"win_rate,omitempty"`
	// WinsCount: Number of profitable markets in this tag
	WinsCount int64 `json:"wins_count,omitempty"`
}

// Tier Tier represents the Tier data model.
type Tier struct {
	// AllowRestApi: Whether the tier is permitted to access REST API endpoints
	AllowRestApi bool `json:"allow_rest_api,omitempty"`
	// AllowStream: Whether the tier is permitted to access real-time stream endpoints
	AllowStream bool `json:"allow_stream,omitempty"`
	// AllowTrial: Whether this tier supports a trial subscription
	AllowTrial bool `json:"allow_trial,omitempty"`
	// CreatedAt: Record creation timestamp
	CreatedAt string `json:"created_at,omitempty"`
	// ID: Unique tier identifier (e.g. "free", "starter", "individual", "builder")
	ID string `json:"id,omitempty"`
	// Level: Relative tier rank for upgrade comparison (0=free, 1=starter, etc.)
	Level int64 `json:"level,omitempty"`
	// License: License type: "personal" or "builder"
	License string `json:"license,omitempty"`
	// MaxDailyRequests: Maximum API queries allowed per 24-hour UTC window (0 = unlimited)
	MaxDailyRequests int64 `json:"max_daily_requests,omitempty"`
	// MaxDailyStreamSeconds: Maximum aggregated streaming seconds allowed per 24-hour UTC window
	MaxDailyStreamSeconds int64 `json:"max_daily_stream_seconds,omitempty"`
	// MaxFilterAddrs: Maximum wallet addresses that can be filtered per stream
	MaxFilterAddrs int64 `json:"max_filter_addrs,omitempty"`
	// MaxFilterSeries: Maximum market series allowed across active streams (0 = disabled)
	MaxFilterSeries int64 `json:"max_filter_series,omitempty"`
	// MaxFilterTags: Maximum market tags allowed across active streams (0 = disabled)
	MaxFilterTags int64 `json:"max_filter_tags,omitempty"`
	// MaxKeys: Maximum concurrent API keys allowed per account
	MaxKeys int64 `json:"max_keys,omitempty"`
	// MaxNakedStreams: Maximum unfiltered streams allowed
	MaxNakedStreams int64 `json:"max_naked_streams,omitempty"`
	// MaxQps: Maximum API queries allowed per second
	MaxQps int64 `json:"max_qps,omitempty"`
	// MaxSessions: Maximum concurrent active SSE connections per account
	MaxSessions int64 `json:"max_sessions,omitempty"`
	// MaxStreams: Maximum persistent stream IDs the user can create
	MaxStreams int64 `json:"max_streams,omitempty"`
	// MonthlyCredits: Monthly credit quota allocated to this tier
	MonthlyCredits int64 `json:"monthly_credits,omitempty"`
	// Name: Display title of the tier
	Name string `json:"name,omitempty"`
	// PriceMonthly: 6 decimals, micro-USD (e.g. 99_000_000 = 99 USD)
	PriceMonthly int64 `json:"price_monthly,omitempty"`
	// PriceQuarterly: 6 decimals, micro-USD (e.g. 237_000_000 = 237 USD)
	PriceQuarterly int64 `json:"price_quarterly,omitempty"`
	// PriceTrial: 6 decimals, micro-USD for trial (e.g. 10_000_000 = 10 USD)
	PriceTrial int64 `json:"price_trial,omitempty"`
	// PriceYearly: 6 decimals, micro-USD (e.g. 499_000_000 = 499 USD)
	PriceYearly int64 `json:"price_yearly,omitempty"`
	// TrialDurationDays: Duration in days for trial subscription (e.g. 5)
	TrialDurationDays int64 `json:"trial_duration_days,omitempty"`
	// UpdatedAt: Last modification timestamp
	UpdatedAt string `json:"updated_at,omitempty"`
}

// TraderHistoryResponse TraderHistoryResponse represents the TraderHistoryResponse data model.
type TraderHistoryResponse struct {
	// History: Paginated list of market settlement details
	History []MarketSettlementDetail `json:"history,omitempty"`
	// TotalCount: Total count of matching market positions for pagination
	TotalCount int64 `json:"total_count,omitempty"`
}

// TraderHourlyStatsResponse TraderHourlyStatsResponse represents the TraderHourlyStatsResponse data model.
type TraderHourlyStatsResponse struct {
	// Stats: Chronological array of 1-hour performance metrics
	Stats []HourlyStat `json:"stats,omitempty"`
}

// TraderIdentity TraderIdentity represents the TraderIdentity data model.
type TraderIdentity struct {
	// Address: Checksummed 0x-prefixed EVM wallet address
	Address string `json:"address,omitempty"`
	// Name: Custom display name or Polymarket profile username
	Name string `json:"name,omitempty"`
	// ProfileCreatedAt: ISO 8601 UTC timestamp of Polymarket profile creation
	ProfileCreatedAt string `json:"profile_created_at,omitempty"`
	// ProfileImage: Public avatar image URL
	ProfileImage string `json:"profile_image,omitempty"`
	// XUsername: Linked X (formerly Twitter) social handle
	XUsername string `json:"x_username,omitempty"`
}

// TraderProfileResponse TraderProfileResponse represents the TraderProfileResponse data model.
type TraderProfileResponse struct {
	// LastTradeAt: ISO 8601 UTC timestamp of trader's most recent trade across the platform
	LastTradeAt string `json:"last_trade_at,omitempty"`
	// Profile: Public identity dossier of the trader
	Profile TraderIdentity `json:"profile,omitempty"`
	// PusdBalance: Trader's pUSD balance in micro-pUSD (6 decimals, e.g. 1000000 = 1 pUSD)
	PusdBalance int64 `json:"pusd_balance,omitempty"`
	// TakerTier: Polymarket Taker Rebate Program tier level (0-6: 0=Tier 0, 1=Bronze, 2=Silver, 3=Gold, 4=Platinum, 5=Diamond, 6=Obsidian)
	TakerTier int64 `json:"taker_tier,omitempty"`
	// TakerTierName: Descriptive taker rebate tier name ("Tier 0", "Bronze", "Silver", "Gold", "Platinum", "Diamond", "Obsidian", or empty string if unsynced)
	TakerTierName string `json:"taker_tier_name,omitempty"`
	// TopTags: Performance breakdown across top 5 categories in last 30 days
	TopTags []TagStat `json:"top_tags,omitempty"`
	// WeightedVolume: Rolling 30-day taker Weighted Volume (wV) in USD, calculated by Trade Size * (1 - Entry Price) * Category Weight * Bonuses
	WeightedVolume float64 `json:"weighted_volume,omitempty"`
}

// TraderSummary TraderSummary represents the TraderSummary data model.
type TraderSummary struct {
	// LastTradeAt: ISO 8601 UTC timestamp of trader's most recent trade across the platform
	LastTradeAt string `json:"last_trade_at,omitempty"`
	// MarketsCount: Number of distinct prediction markets traded in timeframe
	MarketsCount int64 `json:"markets_count,omitempty"`
	// Profile: Public identity dossier of the trader
	Profile TraderIdentity `json:"profile,omitempty"`
	// Rank: Leaderboard ranking position (1-based)
	Rank int64 `json:"rank,omitempty"`
	// Roi: Return on investment percentage over timeframe (e.g. 25.5 for 25.5%)
	Roi float64 `json:"roi,omitempty"`
	// TotalPnl: Realized net profit and loss in micro-USD (6 decimals, e.g. 1000000 = 1 USD)
	TotalPnl int64 `json:"total_pnl,omitempty"`
	// TotalVolume: Aggregated trading volume in micro-USD (6 decimals)
	TotalVolume int64 `json:"total_volume,omitempty"`
	// WinRate: Percentage of profitable closed positions (0.0 to 100.0)
	WinRate float64 `json:"win_rate,omitempty"`
	// WinsCount: Number of profitable settled markets in timeframe
	WinsCount int64 `json:"wins_count,omitempty"`
}

// UpdateSubscriptionRequest UpdateSubscriptionRequest represents the UpdateSubscriptionRequest data model.
type UpdateSubscriptionRequest struct {
	Addresses []string `json:"addresses,omitempty"`
	Series    []string `json:"series,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

// UserAPIKey UserAPIKey represents the UserAPIKey data model.
type UserAPIKey struct {
	CreatedAt string `json:"created_at,omitempty"`
	Key       string `json:"key,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    int64  `json:"status,omitempty"`
}

// UserActiveSessionsResponse UserActiveSessionsResponse represents the UserActiveSessionsResponse data model.
type UserActiveSessionsResponse struct {
	// Data: Array of active connected client session metadata
	Data []ActiveSessionInfo `json:"data,omitempty"`
	// TotalCount: Total concurrent active sessions currently running across all user streams
	TotalCount int64 `json:"total_count,omitempty"`
}

// UserCreateStreamRequest UserCreateStreamRequest represents the UserCreateStreamRequest data model.
type UserCreateStreamRequest struct {
	Addresses []string `json:"addresses,omitempty"`
	Name      string   `json:"name,omitempty"`
	Series    []string `json:"series,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

// UserKeysResponse UserKeysResponse represents the UserKeysResponse data model.
type UserKeysResponse struct {
	Data       []UserAPIKey `json:"data,omitempty"`
	TotalCount int64        `json:"total_count,omitempty"`
}

// UserLedgerEntry UserLedgerEntry represents the UserLedgerEntry data model.
type UserLedgerEntry struct {
	Amount       int64  `json:"amount,omitempty"`
	BalanceAfter int64  `json:"balance_after,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	ReferenceId  string `json:"reference_id,omitempty"`
	Type         string `json:"type,omitempty"`
}

// UserLedgerResponse UserLedgerResponse represents the UserLedgerResponse data model.
type UserLedgerResponse struct {
	Data       []UserLedgerEntry `json:"data,omitempty"`
	TotalCount int64             `json:"total_count,omitempty"`
}

// UserOrder UserOrder represents the UserOrder data model.
type UserOrder struct {
	// Fee: Protocol trading fee deducted in micro-USD (6 decimals)
	Fee int64 `json:"fee,omitempty"`
	// IsTaker: True if executed against resting book liquidity as a taker
	IsTaker bool `json:"is_taker,omitempty"`
	// Outcome: Prediction outcome label (e.g. "Yes", "No")
	Outcome string `json:"outcome,omitempty"`
	// Shares: Filled outcome token shares (6 decimals)
	Shares int64 `json:"shares,omitempty"`
	// Side: Order execution side: "BUY" or "SELL"
	Side Side `json:"side,omitempty"`
	// Time: ISO 8601 UTC timestamp of order fill execution
	Time string `json:"time,omitempty"`
	// USDC: Total collateral USDC transferred (6 decimals micro-USD)
	USDC int64 `json:"usdc,omitempty"`
}

// UserOrdersResponse UserOrdersResponse represents the UserOrdersResponse data model.
type UserOrdersResponse struct {
	// Orders: Chronological list of compacted order fill events
	Orders []UserOrder `json:"orders,omitempty"`
}

// UserProfile UserProfile represents the UserProfile data model.
type UserProfile struct {
	BillingCycle   string `json:"billing_cycle,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	DepositAddress string `json:"deposit_address,omitempty"`
	DepositBalance int64  `json:"deposit_balance,omitempty"`
	Email          string `json:"email,omitempty"`
	IsExpired      bool   `json:"is_expired,omitempty"`
	Status         int64  `json:"status,omitempty"`
	TelegramId     int64  `json:"telegram_id,omitempty"`
	Tier           string `json:"tier,omitempty"`
	TierExpiresAt  string `json:"tier_expires_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// UserSessionHistoryResponse UserSessionHistoryResponse represents the UserSessionHistoryResponse data model.
type UserSessionHistoryResponse struct {
	// Data: Array of historical stream session records
	Data []UserSessionRecord `json:"data,omitempty"`
	// TotalCount: Total count of historical session records
	TotalCount int64 `json:"total_count,omitempty"`
}

// UserSessionRecord UserSessionRecord represents the UserSessionRecord data model.
type UserSessionRecord struct {
	// ApiKey: Masked API credential used to authenticate the stream connection
	ApiKey string `json:"api_key,omitempty"`
	// ClientIp: Remote IP address of the connected client
	ClientIp string `json:"client_ip,omitempty"`
	// CloseReason: Diagnostic termination reason (e.g. client_closed, admin_force_disconnect, quota_exceeded)
	CloseReason string `json:"close_reason,omitempty"`
	// ConnectedAt: ISO 8601 timestamp when the SSE connection was established
	ConnectedAt string `json:"connected_at,omitempty"`
	// DisconnectedAt: ISO 8601 timestamp when the connection was closed, or null if currently active
	DisconnectedAt string `json:"disconnected_at,omitempty"`
	// DurationSeconds: Total connection duration in seconds
	DurationSeconds int64 `json:"duration_seconds,omitempty"`
	// PushedTxCount: Total number of live matched order events delivered over this connection
	PushedTxCount int64 `json:"pushed_tx_count,omitempty"`
	// StreamId: Target stream unique identifier (UUID)
	StreamId string `json:"stream_id,omitempty"`
	// UserAgent: HTTP User-Agent identifier of the client library or application
	UserAgent string `json:"user_agent,omitempty"`
}

// UserStreamResponse UserStreamResponse represents the UserStreamResponse data model.
type UserStreamResponse struct {
	ActiveSessions int64    `json:"active_sessions,omitempty"`
	Addresses      []string `json:"addresses,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	Enabled        bool     `json:"enabled,omitempty"`
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Series         []string `json:"series,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
}

// UserStreamsResponse UserStreamsResponse represents the UserStreamsResponse data model.
type UserStreamsResponse struct {
	Data       []UserStreamResponse `json:"data,omitempty"`
	TotalCount int64                `json:"total_count,omitempty"`
}

// UserSubscription UserSubscription represents the UserSubscription data model.
type UserSubscription struct {
	AmountPaid     int64  `json:"amount_paid,omitempty"`
	BillingCycle   string `json:"billing_cycle,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	DiscountAmount int64  `json:"discount_amount,omitempty"`
	DurationDays   int64  `json:"duration_days,omitempty"`
	ExpiresAt      string `json:"expires_at,omitempty"`
	OriginalPrice  int64  `json:"original_price,omitempty"`
	PlanTier       string `json:"plan_tier,omitempty"`
	PromoCode      string `json:"promo_code,omitempty"`
	ProratedCredit int64  `json:"prorated_credit,omitempty"`
	StartsAt       string `json:"starts_at,omitempty"`
	Status         string `json:"status,omitempty"`
}

// UserTelemetryResponse UserTelemetryResponse represents the UserTelemetryResponse data model.
type UserTelemetryResponse struct {
	Date          string         `json:"date,omitempty"`
	Endpoints     map[string]any `json:"endpoints,omitempty"`
	TotalRequests int64          `json:"total_requests,omitempty"`
}

// UserUpdateKeyRequest UserUpdateKeyRequest represents the UserUpdateKeyRequest data model.
type UserUpdateKeyRequest struct {
	Name string `json:"name,omitempty"`
	// Status: 1: Active, 0: Disabled
	Status int64 `json:"status,omitempty"`
}

// UserUpdateStreamMetadataRequest UserUpdateStreamMetadataRequest represents the UserUpdateStreamMetadataRequest data model.
type UserUpdateStreamMetadataRequest struct {
	Enabled bool   `json:"enabled,omitempty"`
	Name    string `json:"name,omitempty"`
}

// UserWithdrawRequest UserWithdrawRequest represents the UserWithdrawRequest data model.
type UserWithdrawRequest struct {
	Network   string `json:"network,omitempty"`
	ToAddress string `json:"to_address,omitempty"`
	Token     string `json:"token,omitempty"`
}

// UserWithdrawResponse UserWithdrawResponse represents the UserWithdrawResponse data model.
type UserWithdrawResponse struct {
	User       UserProfile `json:"user,omitempty"`
	Withdrawal Withdrawal  `json:"withdrawal,omitempty"`
}

// ValidatePromoCodeRequest ValidatePromoCodeRequest represents the ValidatePromoCodeRequest data model.
type ValidatePromoCodeRequest struct {
	BillingCycle string `json:"billing_cycle,omitempty"`
	Code         string `json:"code,omitempty"`
	Tier         string `json:"tier,omitempty"`
}

// ValidatePromoCodeResponse ValidatePromoCodeResponse represents the ValidatePromoCodeResponse data model.
type ValidatePromoCodeResponse struct {
	AmountToPay    int64  `json:"amount_to_pay,omitempty"`
	Code           string `json:"code,omitempty"`
	DiscountAmount int64  `json:"discount_amount,omitempty"`
	Message        string `json:"message,omitempty"`
	OriginalPrice  int64  `json:"original_price,omitempty"`
	ProratedCredit int64  `json:"prorated_credit,omitempty"`
	Valid          bool   `json:"valid,omitempty"`
}

// Withdrawal Withdrawal represents the Withdrawal data model.
type Withdrawal struct {
	// Amount: Full balance amount withdrawn in micro-USD
	Amount int64 `json:"amount,omitempty"`
	// CreatedAt: Withdrawal request creation timestamp
	CreatedAt string `json:"created_at,omitempty"`
	ID        int64  `json:"id,omitempty"`
	// LastDepositAt: User's most recent deposit timestamp
	LastDepositAt string `json:"last_deposit_at,omitempty"`
	// Network: "bsc" or "polygon"
	Network string `json:"network,omitempty"`
	// RejectReason: Reason recorded upon rejection
	RejectReason string `json:"reject_reason,omitempty"`
	// ReviewedAt: Administrative review timestamp
	ReviewedAt string `json:"reviewed_at,omitempty"`
	// Status: "pending", "completed", "rejected"
	Status string `json:"status,omitempty"`
	// ToAddress: Destination EVM wallet address
	ToAddress string `json:"to_address,omitempty"`
	// Token: "USDC" or "USDT"
	Token string `json:"token,omitempty"`
	// TxHash: Transaction hash on destination network
	TxHash string `json:"tx_hash,omitempty"`
	// UpdatedAt: Last modification timestamp
	UpdatedAt string `json:"updated_at,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	UserId    int64  `json:"user_id,omitempty"`
}
