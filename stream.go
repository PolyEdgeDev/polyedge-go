package polyedge

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	ErrFatalUnauthorized   = errors.New("polyedge: unauthorized (HTTP 401)")
	ErrFatalForbidden      = errors.New("polyedge: forbidden or quota exceeded (HTTP 403)")
	ErrFatalStreamNotFound = errors.New("polyedge: stream not found (HTTP 404)")
	ErrFatalBadRequest     = errors.New("polyedge: bad request (HTTP 400)")
	ErrWatchdogTimeout     = errors.New("polyedge: sse watchdog timed out with no data/heartbeat")
)

type StreamOptions struct {
	StreamBaseURL    string
	HeartbeatTimeout time.Duration
	InitialBackoff   time.Duration
	MaxBackoff       time.Duration
	HTTPClient       *http.Client
}

type StreamClient struct {
	streamID         string
	apiKey           string
	streamBaseURL    string
	heartbeatTimeout time.Duration
	initialBackoff   time.Duration
	maxBackoff       time.Duration
	httpClient       *http.Client
	lastEventID      atomic.Uint64
	onHeartbeat      func()
}

func NewStreamClient(streamID, apiKey string, opts ...StreamOptions) (*StreamClient, error) {
	if streamID == "" {
		return nil, errors.New("streamID cannot be empty")
	}
	if apiKey == "" {
		return nil, errors.New("apiKey cannot be empty")
	}

	sc := &StreamClient{
		streamID:         streamID,
		apiKey:           apiKey,
		streamBaseURL:    "https://stream.polyedge.dev",
		heartbeatTimeout: 45 * time.Second,
		initialBackoff:   1 * time.Second,
		maxBackoff:       30 * time.Second,
		httpClient: &http.Client{
			Timeout: 0, // Continuous streaming, no global request timeout
		},
	}

	if len(opts) > 0 {
		opt := opts[0]
		if opt.StreamBaseURL != "" {
			sc.streamBaseURL = strings.TrimRight(opt.StreamBaseURL, "/")
		}
		if opt.HeartbeatTimeout > 0 {
			sc.heartbeatTimeout = opt.HeartbeatTimeout
		}
		if opt.InitialBackoff > 0 {
			sc.initialBackoff = opt.InitialBackoff
		}
		if opt.MaxBackoff > 0 {
			sc.maxBackoff = opt.MaxBackoff
		}
		if opt.HTTPClient != nil {
			sc.httpClient = opt.HTTPClient
		}
	}

	return sc, nil
}

func (s *StreamClient) SetOnHeartbeat(fn func()) {
	s.onHeartbeat = fn
}

func (s *StreamClient) LastEventID() uint64 {
	return s.lastEventID.Load()
}

func (s *StreamClient) Subscribe(ctx context.Context) (<-chan *LiveTransaction, <-chan error) {
	txCh := make(chan *LiveTransaction, 1024)
	errCh := make(chan error, 16)

	go func() {
		defer close(txCh)
		defer close(errCh)

		backoff := s.initialBackoff

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			err := s.connectAndStream(ctx, txCh)
			if err == nil {
				return
			}

			// Check if fatal error
			if errors.Is(err, ErrFatalUnauthorized) ||
				errors.Is(err, ErrFatalForbidden) ||
				errors.Is(err, ErrFatalStreamNotFound) ||
				errors.Is(err, ErrFatalBadRequest) {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			// Non-fatal: notify error and retry with exponential backoff + jitter
			select {
			case errCh <- err:
			default:
			}

			jitter := float64(backoff) * (0.8 + rand.Float64()*0.4)
			sleepDur := time.Duration(jitter)

			select {
			case <-ctx.Done():
				return
			case <-time.After(sleepDur):
			}

			backoff = time.Duration(float64(backoff) * 2.0)
			if backoff > s.maxBackoff {
				backoff = s.maxBackoff
			}
		}
	}()

	return txCh, errCh
}

func (s *StreamClient) connectAndStream(ctx context.Context, txCh chan<- *LiveTransaction) error {
	reqURL, err := url.Parse(fmt.Sprintf("%s/streams/%s", s.streamBaseURL, s.streamID))
	if err != nil {
		return fmt.Errorf("invalid stream URL: %w", err)
	}

	lastID := s.lastEventID.Load()
	if lastID > 0 {
		q := reqURL.Query()
		q.Set("lastEventId", strconv.FormatUint(lastID, 10))
		reqURL.RawQuery = q.Encode()
	}

	reqCtx, reqCancel := context.WithCancel(ctx)
	defer reqCancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("X-PolyEdge-Key", s.apiKey)
	req.Header.Set("Cache-Control", "no-cache")
	if lastID > 0 {
		req.Header.Set("Last-Event-ID", strconv.FormatUint(lastID, 10))
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return ErrFatalUnauthorized
		case http.StatusForbidden:
			return ErrFatalForbidden
		case http.StatusNotFound:
			return ErrFatalStreamNotFound
		case http.StatusBadRequest:
			return ErrFatalBadRequest
		default:
			return fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
		}
	}

	// Setup watchdog timer
	watchdog := time.NewTimer(s.heartbeatTimeout)
	defer watchdog.Stop()

	// Watchdog timeout monitor goroutine
	watchdogErrCh := make(chan error, 1)
	go func() {
		select {
		case <-watchdog.C:
			watchdogErrCh <- ErrWatchdogTimeout
			reqCancel() // Cancel request to unblock reader
		case <-reqCtx.Done():
		}
	}()

	reader := bufio.NewReader(resp.Body)
	var curID uint64
	var curEvent string = "message"
	var curData bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			select {
			case wErr := <-watchdogErrCh:
				return wErr
			default:
				return err
			}
		}

		// Reset watchdog timer whenever any bytes/lines arrive
		if !watchdog.Stop() {
			select {
			case <-watchdog.C:
			default:
			}
		}
		watchdog.Reset(s.heartbeatTimeout)

		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			// End of SSE block
			if curID > 0 {
				s.lastEventID.Store(curID)
			}

			if curEvent == "tx" && curData.Len() > 0 {
				var tx LiveTransaction
				if err := json.Unmarshal(curData.Bytes(), &tx); err == nil {
					select {
					case txCh <- &tx:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			}

			curID = 0
			curEvent = "message"
			curData.Reset()
			continue
		}

		if bytes.HasPrefix(trimmed, []byte(":")) {
			// Heartbeat comment
			if s.onHeartbeat != nil {
				s.onHeartbeat()
			}
			continue
		} else if bytes.HasPrefix(trimmed, []byte("id:")) {
			idStr := strings.TrimSpace(string(trimmed[3:]))
			if parsed, pErr := strconv.ParseUint(idStr, 10, 64); pErr == nil {
				curID = parsed
			}
		} else if bytes.HasPrefix(trimmed, []byte("event:")) {
			curEvent = strings.TrimSpace(string(trimmed[6:]))
		} else if bytes.HasPrefix(trimmed, []byte("data:")) {
			payload := bytes.TrimPrefix(trimmed, []byte("data:"))
			curData.Write(bytes.TrimSpace(payload))
		}
	}
}
