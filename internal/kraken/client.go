package kraken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	host string
}

var pairMap = map[string]string{
	"BTC/USD": "XXBTZUSD",
	"BTC/EUR": "XXBTZEUR",
	"BTC/CHF": "XBTCHF",
}

type krakenPair struct {
	LastTradeClosed []string `json:"c"`
}

type krakenResponse struct {
	Error  []string              `json:"error"`
	Result map[string]krakenPair `json:"result"`
}

func NewClient(host string) *Client {
	return &Client{
		host: host,
	}
}

func (s *Client) GetLastTradedPrice(ctx context.Context, pair string) (float64, error) {
	const op = "kraken.Client.GetLastTradedPrice"

	mappedPair, ok := pairMap[pair]

	if !ok {
		return 0, fmt.Errorf("(%s) this pair cannot be mapped: %s", op, pair)
	}

	url := fmt.Sprintf("%s/0/public/Ticker", s.host)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return 0, fmt.Errorf("(%s) cannot build the request: %w", op, err)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("(%s) incorrect status code: %d", op, resp.StatusCode)
	}

	var kResp krakenResponse

	if err := json.NewDecoder(resp.Body).Decode(&kResp); err != nil {
		return 0, fmt.Errorf("(%s) response cannot be decoded: %w", op, err)
	}

	if len(kResp.Error) > 0 {
		return 0, fmt.Errorf("(%s) response has errors: %s", op, strings.Join(kResp.Error, ", "))
	}

	pairData, ok := kResp.Result[mappedPair]

	if !ok {
		return 0, fmt.Errorf("(%s) no results for pair: %s", op, pair)
	}

	if len(pairData.LastTradeClosed) != 2 {
		return 0, fmt.Errorf("(%s) incorrect length of the LTC value: %s", op, pair)
	}

	var lastTrade float64
	_, err = fmt.Sscanf(pairData.LastTradeClosed[0], "%f", &lastTrade)

	if err != nil {
		return 0, fmt.Errorf("(%s) LTC value cannot be casted: %w", op, err)
	}

	return lastTrade, nil
}
