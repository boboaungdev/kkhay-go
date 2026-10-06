package kkhay

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("kkhay_live_test_123")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	if client.baseURL != "https://api.kkhay.com" {
		t.Errorf("expected default baseURL https://api.kkhay.com, got %s", client.baseURL)
	}

	// Custom options
	custom, err := NewClient("kkhay_custom", WithBaseURL("http://localhost:3000"), WithTimeout(5*time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if custom.baseURL != "http://localhost:3000" {
		t.Errorf("expected custom URL, got %s", custom.baseURL)
	}

	// Empty key error
	_, err = NewClient("")
	if err == nil {
		t.Errorf("expected error for empty API key")
	}
}

func TestAPIErrorFormatting(t *testing.T) {
	err := &APIError{
		Status:  404,
		Message: "Invoice not found",
		Code:    "INVOICE_NOT_FOUND",
	}

	if err.Error() != "kkhay: API error HTTP 404 [INVOICE_NOT_FOUND]: Invoice not found" {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}

