package coindesk

import (
	"context"
	"fmt"
	en "github.com/100bench/cryptocurrency_provider.git/internal/entities"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// структура клиента и сам клиент http, подумать про поля структуры + конструктор + методы

type ClientCoinDesk struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClientCoinDesk(httpClient *http.Client, baseURL, apiKey string) (*ClientCoinDesk, error) {
	if httpClient == nil {
		return nil, errors.Wrap(en.ErrNilDependency, "coindesk client")
	}
	if baseURL == "" {
		return nil, errors.Wrap(en.ErrEmptyBaseURL, "coindesk url is empty")
	}
	return &ClientCoinDesk{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}, nil
}

func (c *ClientCoinDesk) GetRates(ctx context.Context, currencies []string) ([]en.Rate, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}

	// query params
	q := u.Query()
	q.Set("fsyms", strings.Join(currencies, ",")) // BTC,ETH,SOL
	q.Set("tsyms", "USD")                         // USD
	q.Set("tryConversion", "true")
	u.RawQuery = q.Encode()

	// готовим запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// ключ лучше в заголовок
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Apikey "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "provider cryptocompare: do request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, errors.Wrapf(cases.ErrProviderUnavailable, "cryptocompare status=%d body=%s", resp.StatusCode, string(b))
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "provider cryptocompare: read body")
	}
	
	return nil, nil
}
