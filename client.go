package kkhay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client is the official Go HTTP client for K Khay Payment Gateway.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// ClientOption configures a Client instance.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(url, "/")
	}
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithTimeout sets a custom request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient constructs a new K Khay API client.
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	key := strings.TrimSpace(apiKey)
	if key == "" {
		return nil, errors.New("kkhay: API key cannot be empty")
	}

	c := &Client{
		apiKey:  key,
		baseURL: "https://api.kkhay.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Client) getURL(path string) string {
	cleanPath := path
	if !strings.HasPrefix(cleanPath, "/") {
		cleanPath = "/" + cleanPath
	}

	if strings.HasSuffix(c.baseURL, "/api") || strings.Contains(c.baseURL, "api.") {
		return c.baseURL + cleanPath
	}
	return c.baseURL + "/api" + cleanPath
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	fullURL := c.getURL(path)

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("kkhay: failed to encode request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("kkhay: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("User-Agent", "kkhay-go/1.0.0")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("kkhay: network error: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("kkhay: failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Message string      `json:"message"`
			Error   string      `json:"error"`
			Code    string      `json:"code"`
			Details interface{} `json:"details"`
		}
		_ = json.Unmarshal(respData, &errResp)

		errMsg := errResp.Message
		if errMsg == "" {
			errMsg = errResp.Error
		}
		if errMsg == "" {
			errMsg = string(respData)
		}

		return &APIError{
			Status:  resp.StatusCode,
			Message: errMsg,
			Code:    errResp.Code,
			Details: errResp.Details,
		}
	}

	if result != nil && len(respData) > 0 {
		if err := json.Unmarshal(respData, result); err != nil {
			return fmt.Errorf("kkhay: failed to decode JSON response: %w", err)
		}
	}

	return nil
}

// CreateInvoice creates a new sovereign crypto payment invoice.
func (c *Client) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	if req.PriceAmount <= 0 {
		return nil, errors.New("kkhay: PriceAmount must be greater than 0")
	}
	if req.PayNetwork == "" {
		return nil, errors.New("kkhay: PayNetwork is required (e.g. 'bsc', 'polygon')")
	}
	if req.PayToken == "" {
		return nil, errors.New("kkhay: PayToken is required (e.g. 'USDT', 'USDC')")
	}

	var res CreateInvoiceResponse
	if err := c.doRequest(ctx, http.MethodPost, "/v1/merchant/invoices", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetInvoice retrieves an existing invoice by UUID.
func (c *Client) GetInvoice(ctx context.Context, invoiceID string) (*GetInvoiceResponse, error) {
	id := strings.TrimSpace(invoiceID)
	if id == "" {
		return nil, errors.New("kkhay: invoiceID is required")
	}

	var res GetInvoiceResponse
	path := fmt.Sprintf("/v1/merchant/invoices/%s", url.PathEscape(id))
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ListInvoices queries and paginates merchant invoices.
func (c *Client) ListInvoices(ctx context.Context, query ListInvoicesQuery) (*ListInvoicesResponse, error) {
	params := url.Values{}
	if query.Page > 0 {
		params.Set("page", strconv.Itoa(query.Page))
	}
	if query.Limit > 0 {
		params.Set("limit", strconv.Itoa(query.Limit))
	}
	if query.Status != "" {
		params.Set("status", query.Status)
	}
	if query.Search != "" {
		params.Set("search", query.Search)
	}

	path := "/v1/merchant/invoices"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var res ListInvoicesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CheckHealth queries API health status.
func (c *Client) CheckHealth(ctx context.Context) (map[string]interface{}, error) {
	var res map[string]interface{}
	if err := c.doRequest(ctx, http.MethodGet, "/health", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

