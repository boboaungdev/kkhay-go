package kkhay

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

var (
	// ErrInvalidSignature is returned when the webhook HMAC signature does not match.
	ErrInvalidSignature = errors.New("kkhay: invalid webhook signature")
)

// VerifyWebhookSignature verifies the HMAC-SHA256 signature sent with incoming IPN webhooks.
func VerifyWebhookSignature(payload []byte, signature string, ipnSecret string) bool {
	sig := strings.TrimSpace(signature)
	secret := strings.TrimSpace(ipnSecret)
	if sig == "" || secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(expectedMAC), []byte(sig)) == 1
}

// ParseWebhookEvent validates signature and parses the raw JSON body into WebhookEvent.
func ParseWebhookEvent(rawBody []byte, signature string, ipnSecret string) (*WebhookEvent, error) {
	if !VerifyWebhookSignature(rawBody, signature, ipnSecret) {
		return nil, ErrInvalidSignature
	}

	var event WebhookEvent
	if err := json.Unmarshal(rawBody, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

