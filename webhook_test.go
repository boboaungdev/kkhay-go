package kkhay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestWebhookSignatureVerification(t *testing.T) {
	secret := "whsec_test_secret_123"
	payload := []byte(`{"event":"payment.finished","invoice_id":"inv_123"}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSig := hex.EncodeToString(mac.Sum(nil))

	// Valid
	if !VerifyWebhookSignature(payload, validSig, secret) {
		t.Errorf("expected signature to be valid")
	}

	// Tampered
	tampered := []byte(`{"event":"payment.finished","invoice_id":"inv_999"}`)
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
}

func TestParseWebhookEvent(t *testing.T) {
	secret := "whsec_test_secret_123"
	payload := []byte(`{"event":"payment.finished","invoice_id":"inv_abc123","pay_amount":"50.00","pay_token":"USDT","status":"paid"}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSig := hex.EncodeToString(mac.Sum(nil))

	event, err := ParseWebhookEvent(payload, validSig, secret)
	if err != nil {
		t.Fatalf("unexpected error parsing webhook: %v", err)
	}

	if event.Event != "payment.finished" {
		t.Errorf("expected event payment.finished, got %s", event.Event)
	}
	if event.InvoiceID != "inv_abc123" {
		t.Errorf("expected invoice_id inv_abc123, got %s", event.InvoiceID)
	}
	if event.PayToken != "USDT" {
		t.Errorf("expected pay_token USDT, got %s", event.PayToken)
	}

	// Fails with bad sig
	_, err = ParseWebhookEvent(payload, "invalid_sig", secret)
	if err != ErrInvalidSignature {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}
}

