package kkhay

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("kkhay_live_test_123")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("expected default baseURL %s, got %s", DefaultBaseURL, client.baseURL)
	}

	// Custom options
	custom, err := NewClient("kkhay_custom", WithBaseURL("http://localhost:3000"), WithTimeout(5*time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if custom.baseURL != "http://localhost:3000" {
		t.Errorf("expected custom URL, got %s", custom.baseURL)
	}

	// Custom HTTP client
	customHTTP := &http.Client{Timeout: 10 * time.Second}
	customWithClient, err := NewClient("kkhay_custom", WithHTTPClient(customHTTP))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customWithClient.httpClient != customHTTP {
		t.Errorf("expected custom HTTP client to be assigned")
	}

	// Empty key error
	_, err = NewClient("")
	if err == nil {
		t.Errorf("expected error for empty API key")
	}

	// Whitespace only key error
	_, err = NewClient("   ")
	if err == nil {
		t.Errorf("expected error for whitespace API key")
	}
}

func TestGetURL(t *testing.T) {
	tests := []struct {
		baseURL  string
		path     string
		expected string
	}{
		{
			baseURL:  "https://api.kkhay.com",
			path:     "/v1/gateway/invoices",
			expected: "https://api.kkhay.com/v1/gateway/invoices",
		},
		{
			baseURL:  "https://kkhay.com",
			path:     "/v1/gateway/invoices",
			expected: "https://kkhay.com/api/v1/gateway/invoices",
		},
		{
			baseURL:  "https://kkhay.com/api",
			path:     "v1/gateway/invoices",
			expected: "https://kkhay.com/api/v1/gateway/invoices",
		},
		{
			baseURL:  "http://localhost:3000",
			path:     "/v1/gateway/invoices",
			expected: "http://localhost:3000/api/v1/gateway/invoices",
		},
	}

	for _, tt := range tests {
		client, err := NewClient("test_key", WithBaseURL(tt.baseURL))
		if err != nil {
			t.Fatalf("unexpected client init error: %v", err)
		}
		actual := client.getURL(tt.path)
		if actual != tt.expected {
			t.Errorf("for base %s and path %s, expected %s, got %s", tt.baseURL, tt.path, tt.expected, actual)
		}
	}
}

func TestAPIErrorFormatting(t *testing.T) {
	errWithCode := &APIError{
		Status:  404,
		Message: "Invoice not found",
		Code:    "INVOICE_NOT_FOUND",
	}

	if errWithCode.Error() != "kkhay: API error HTTP 404 [INVOICE_NOT_FOUND]: Invoice not found" {
		t.Errorf("unexpected error string: %s", errWithCode.Error())
	}

	errWithoutCode := &APIError{
		Status:  500,
		Message: "Internal server error",
	}

	if errWithoutCode.Error() != "kkhay: API error HTTP 500: Internal server error" {
		t.Errorf("unexpected error string: %s", errWithoutCode.Error())
	}
}

func TestCreateInvoiceValidation(t *testing.T) {
	client, _ := NewClient("kkhay_live_test")
	ctx := context.Background()

	// Missing amount
	_, err := client.CreateInvoice(ctx, CreateInvoiceRequest{
		PayNetwork: "bsc",
		PayToken:   "USDT",
	})
	if err == nil {
		t.Errorf("expected error for missing amount")
	}

	// Missing network
	_, err = client.CreateInvoice(ctx, CreateInvoiceRequest{
		Amount:   50.0,
		PayToken: "USDT",
	})
	if err == nil {
		t.Errorf("expected error for missing PayNetwork")
	}

	// Missing token
	_, err = client.CreateInvoice(ctx, CreateInvoiceRequest{
		Amount:     50.0,
		PayNetwork: "bsc",
	})
	if err == nil {
		t.Errorf("expected error for missing PayToken")
	}
}

func TestCreateInvoiceSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/v1/gateway/invoices") {
			t.Errorf("expected path ending in /v1/gateway/invoices, got %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "kkhay_live_test_key" {
			t.Errorf("expected x-api-key header")
		}
		if r.Header.Get("User-Agent") != "kkhay-go/1.0.4" {
			t.Errorf("expected User-Agent kkhay-go/1.0.4, got %s", r.Header.Get("User-Agent"))
		}

		var req CreateInvoiceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if req.Amount != 49.99 {
			t.Errorf("expected amount 49.99, got %f", req.Amount)
		}
		if req.PayNetwork != "bsc" {
			t.Errorf("expected payNetwork bsc, got %s", req.PayNetwork)
		}
		if req.PayToken != "USDT" {
			t.Errorf("expected payToken USDT, got %s", req.PayToken)
		}

		resp := CreateInvoiceResponse{
			OK: true,
			Invoice: Invoice{
				ID:             "inv_0192a3b4c5d6",
				GatewayID:      "gw_1234567890",
				Amount:         49.99,
				Currency:       "USD",
				PriceAmount:    49.99,
				PriceCurrency:  "USD",
				PayNetwork:     &req.PayNetwork,
				PayToken:       &req.PayToken,
				DepositAddress: "0x1234567890123456789012345678901234567890",
				Status:         StatusWaiting,
				HostedURL:      "https://kkhay.com/invoice/inv_0192a3b4c5d6",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient("kkhay_live_test_key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.CreateInvoice(context.Background(), CreateInvoiceRequest{
		Amount:     49.99,
		PayNetwork: "BSC",
		PayToken:   "usdt",
	})
	if err != nil {
		t.Fatalf("unexpected error creating invoice: %v", err)
	}

	if !resp.OK {
		t.Errorf("expected OK true")
	}
	if resp.Invoice.ID != "inv_0192a3b4c5d6" {
		t.Errorf("expected invoice ID inv_0192a3b4c5d6, got %s", resp.Invoice.ID)
	}
	if resp.Invoice.GatewayID != "gw_1234567890" {
		t.Errorf("expected gateway ID gw_1234567890, got %s", resp.Invoice.GatewayID)
	}
	if resp.Invoice.Status != StatusWaiting {
		t.Errorf("expected status WAITING, got %s", resp.Invoice.Status)
	}
	if resp.Invoice.HostedURL != "https://kkhay.com/invoice/inv_0192a3b4c5d6" {
		t.Errorf("expected hosted URL https://kkhay.com/invoice/inv_0192a3b4c5d6, got %s", resp.Invoice.HostedURL)
	}
}

func TestGetInvoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if strings.HasSuffix(r.URL.Path, "/v1/gateway/invoices/inv_not_found") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":   "Not Found",
				"message": "Invoice not found",
				"code":    "NOT_FOUND",
			})
			return
		}

		if !strings.HasSuffix(r.URL.Path, "/v1/gateway/invoices/inv_valid") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := GetInvoiceResponse{
			OK: true,
			Invoice: Invoice{
				ID:        "inv_valid",
				GatewayID: "gw_123",
				Status:    StatusPaid,
			},
			Payments: []PaymentRecord{
				{
					ID:             "pmt_001",
					InvoiceID:      "inv_valid",
					TxHash:         "0xabcdef1234567890",
					AmountReceived: "49.99",
					Confirmations:  12,
					Status:         PaymentStatusSettled,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient("test_key", WithBaseURL(server.URL))
	ctx := context.Background()

	// Empty ID error
	_, err := client.GetInvoice(ctx, "")
	if err == nil {
		t.Errorf("expected error for empty invoiceID")
	}

	// Success
	res, err := client.GetInvoice(ctx, "inv_valid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Invoice.Status != StatusPaid {
		t.Errorf("expected status PAID, got %s", res.Invoice.Status)
	}
	if len(res.Payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(res.Payments))
	}
	if res.Payments[0].TxHash != "0xabcdef1234567890" {
		t.Errorf("unexpected txHash: %s", res.Payments[0].TxHash)
	}

	// 404
	_, err = client.GetInvoice(ctx, "inv_not_found")
	if err == nil {
		t.Fatalf("expected error for not found invoice")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Status != 404 {
		t.Errorf("expected status 404, got %d", apiErr.Status)
	}
}

func TestListInvoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/v1/gateway/invoices") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		q := r.URL.Query()
		if q.Get("page") != "2" || q.Get("limit") != "10" || q.Get("status") != "PAID" || q.Get("type") != "PAYMENT" || q.Get("search") != "order" {
			t.Errorf("unexpected query parameters: %s", r.URL.RawQuery)
		}

		resp := ListInvoicesResponse{
			OK:          true,
			TotalCount:  1,
			TotalPages:  1,
			CurrentPage: 2,
			Limit:       10,
			Items: []Invoice{
				{
					ID:        "inv_item_1",
					GatewayID: "gw_123",
					Status:    StatusPaid,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient("test_key", WithBaseURL(server.URL))
	res, err := client.ListInvoices(context.Background(), ListInvoicesQuery{
		Page:   2,
		Limit:  10,
		Status: StatusPaid,
		Type:   TypePayment,
		Search: "order",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotalCount != 1 || len(res.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(res.Items))
	}
}

func TestCheckHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/health") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
		})
	}))
	defer server.Close()

	client, _ := NewClient("test_key", WithBaseURL(server.URL))
	res, err := client.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["status"] != "ok" {
		t.Errorf("expected status ok, got %v", res["status"])
	}
}
