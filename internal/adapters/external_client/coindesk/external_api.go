package coindesk

import (
	"context"
	"encoding/json"
	"fmt"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// структура клиента и сам клиент http, подумать про поля структуры + конструктор + методы

type ClientCoinDesk struct {
	httpClient http.Client
	baseURL    string
}

func NewClientCoinDesk(baseURL string) (*ClientCoinDesk, error) {
	if baseURL == "" {
		return nil, errors.Wrap(en.ErrNilDependency, "coindesk url is empty")
	}
	client := http.Client{}
	return &ClientCoinDesk{
		httpClient: client,
		baseURL:    baseURL,
	}, nil
}

func (c *ClientCoinDesk) GetRatesFromClient(ctx context.Context, currencies []string) ([]en.Rate, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}

	// query params
	q := u.Query()
	q.Set("fsyms", strings.Join(currencies, ",")) // BTC,ETH,SOL
	q.Set("tsyms", "USD")
	u.RawQuery = q.Encode()

	// готовим запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "provider cryptocompare: do request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, errors.Wrapf(err, "cryptocompare status=%d body=%s", resp.StatusCode, string(b)) // ?
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "provider cryptocompare: read body")
	}

	var raw map[string]map[string]float64
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, errors.Wrap(err, "cryptocompare: decode json")
	}

	rates := make([]en.Rate, 0, len(raw))
	ts := time.Now().UTC()

	for k, v := range raw {
		for _, vv := range v {
			rate, err := en.NewRate(k, vv, ts)
			if err != nil {
				return nil, fmt.Errorf("cryptocompare: new rate: %w", err)
			}
			rates = append(rates, *rate)
		}
	}
	return rates, nil
}
