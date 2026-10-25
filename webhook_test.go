package kkhay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestWebhookSignatureVerification(t *testing.T) {
	secret := "whsec_test_secret_123"
	payload := []byte(`{"event":"payment.paid","invoice_id":"inv_123","gateway_id":"gw_456"}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSig := hex.EncodeToString(mac.Sum(nil))

	// Valid
	if !VerifyWebhookSignature(payload, validSig, secret) {
		t.Errorf("expected signature to be valid")
	}

	// Tampered payload
	tampered := []byte(`{"event":"payment.paid","invoice_id":"inv_999","gateway_id":"gw_456"}`)
	if VerifyWebhookSignature(tampered, validSig, secret) {
		t.Errorf("expected signature to fail on tampered payload")
	}

	// Wrong secret
	if VerifyWebhookSignature(payload, validSig, "wrong_secret") {
		t.Errorf("expected signature to fail on wrong secret")
	}

	// Empty sig
	if VerifyWebhookSignature(payload, "", secret) {
		t.Errorf("expected signature to fail on empty signature")
	}

	// Empty secret
	if VerifyWebhookSignature(payload, validSig, "") {
		t.Errorf("expected signature to fail on empty secret")
	}
}

func TestParseWebhookEvent(t *testing.T) {
	secret := "whsec_test_secret_123"
	payload := []byte(`{
		"event": "payment.paid",
		"invoice_id": "inv_0192a3b4c5d6",
		"gateway_id": "gw_9876543210",
		"invoice_type": "PAYMENT",
		"order_id": "ORDER-5544",
		"amount": 100.0,
		"currency": "USD",
		"price_amount": 100.0,
		"price_currency": "USD",
		"pay_amount": "100.00000000",
		"pay_token": "USDT",
		"pay_network": "bsc",
		"deposit_address": "0x1234567890123456789012345678901234567890",
		"status": "PAID",
		"timestamp": "2026-10-26T00:00:00Z"
	}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSig := hex.EncodeToString(mac.Sum(nil))

	event, err := ParseWebhookEvent(payload, validSig, secret)
	if err != nil {
		t.Fatalf("unexpected error parsing webhook: %v", err)
	}

	if event.Event != EventPaymentPaid {
		t.Errorf("expected event %s, got %s", EventPaymentPaid, event.Event)
	}
	if event.InvoiceID != "inv_0192a3b4c5d6" {
		t.Errorf("expected invoice_id inv_0192a3b4c5d6, got %s", event.InvoiceID)
	}
	if event.GatewayID != "gw_9876543210" {
		t.Errorf("expected gateway_id gw_9876543210, got %s", event.GatewayID)
	}
	if event.PayToken != "USDT" {
		t.Errorf("expected pay_token USDT, got %s", event.PayToken)
	}
	if event.PayNetwork != "bsc" {
		t.Errorf("expected pay_network bsc, got %s", event.PayNetwork)
	}
	if event.Status != StatusPaid {
		t.Errorf("expected status PAID, got %s", event.Status)
	}

	// Fails with bad sig
	_, err = ParseWebhookEvent(payload, "invalid_sig", secret)
	if err != ErrInvalidSignature {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}

	// Fails with invalid json
	invalidJSON := []byte(`{not-json`)
	mac2 := hmac.New(sha256.New, []byte(secret))
	mac2.Write(invalidJSON)
	sig2 := hex.EncodeToString(mac2.Sum(nil))
	_, err = ParseWebhookEvent(invalidJSON, sig2, secret)
	if err == nil {
		t.Errorf("expected JSON unmarshal error")
	}
}
