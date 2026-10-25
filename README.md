<div align="center">

# `kkhay-go`

**Official Go SDK for K Khay — The Non-Custodial Crypto Payment Gateway**

[![Go Reference](https://pkg.go.dev/badge/github.com/boboaungdev/kkhay-go.svg)](https://pkg.go.dev/github.com/boboaungdev/kkhay-go)
[![CI](https://github.com/boboaungdev/kkhay-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/boboaungdev/kkhay-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/boboaungdev/kkhay-go)](https://goreportcard.com/report/github.com/boboaungdev/kkhay-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)

Accept cryptocurrency payments across **190+ blockchain networks** and **500+ tokens** (EVM, Bitcoin, Solana, TON, Tron, and XRPL) with **zero custody** and automatic direct-to-wallet settlement. Built with **zero external dependencies** using the Go standard library.

[Gateway Dashboard](https://gateway.kkhay.com) • [API Keys](https://gateway.kkhay.com/dashboard/keys) • [Website](https://kkhay.com)

</div>

---

## ⚡ Key Highlights

- 🛡️ **100% Non-Custodial**: Payments settle directly to your destination blockchain address in seconds. K Khay never holds your private keys or customer funds.
- 🌐 **190+ Networks Supported**: 
  - **EVM Networks**: Ethereum, BNB Smart Chain, Polygon, Arbitrum, Base, Optimism, Avalanche, and 180+ more.
  - **UTXO Networks**: Bitcoin (`btc`), Dogecoin (`doge`), Litecoin (`ltc`), Bitcoin Cash (`bch`).
  - **High-Throughput Chains**: Solana (`sol`), TON (`ton`), Tron (`trx`), XRP Ledger (`xrp`).
- 🪙 **500+ Cryptocurrencies**: Native coins and tokens (USDT, USDC, BTC, ETH, SOL, TON, BNB, DOGE, XRP, TRX, LTC, and more).
- 🔗 **Branded Hosted Checkouts**: Automatic payment links with dynamic QR codes, live blockchain polling, and wallet deep linking.
- 🔔 **Cryptographic Webhooks**: HMAC-SHA256 signed Instant Payment Notifications (`x-kkhay-signature`).
- 🪶 **Zero Dependencies**: Pure Go standard library implementation (`net/http`, `crypto/hmac`, `crypto/sha256`, `encoding/json`).

---

## 📦 Installation

```bash
go get github.com/boboaungdev/kkhay-go@v1.0.4
```

*(Requires Go 1.20 or newer)*

---

## 🚀 Quickstart

### 1. Initialize the Client

Generate your API key in the [K Khay Gateway Dashboard](https://gateway.kkhay.com/dashboard/keys).

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/boboaungdev/kkhay-go"
)

func main() {
    client, err := kkhay.NewClient(os.Getenv("KKHAY_API_KEY")) // starts with "kkhay_live_"
    if err != nil {
        log.Fatalf("Failed to initialize client: %v", err)
    }

    ctx := context.Background()

    // 2. Create a crypto payment invoice
    resp, err := client.CreateInvoice(ctx, kkhay.CreateInvoiceRequest{
        Amount:        49.99,
        Currency:      "USD",
        PayNetwork:    "bsc",       // "bsc", "polygon", "base", "arbitrum", "ethereum", "solana", "bitcoin", "ton"
        PayToken:      "USDT",      // "USDT", "USDC", "BTC", "ETH", "SOL", "TON"
        OrderID:       "ORDER-9842",
        Title:         "Annual Pro Plan",
        CustomerName:  "Alice Smith",
        CustomerEmail: "alice@example.com",
        RedirectURL:   "https://myshop.com/orders/success",
        CancelURL:     "https://myshop.com/cart",
    })
    if err != nil {
        log.Fatalf("Failed to create invoice: %v", err)
    }

    fmt.Printf("Invoice ID: %s\n", resp.Invoice.ID)
    fmt.Printf("Gateway ID: %s\n", resp.Invoice.GatewayID)
    fmt.Printf("Deposit Address: %s\n", resp.Invoice.DepositAddress)
    fmt.Printf("Hosted Checkout URL: %s\n", resp.Invoice.HostedURL)
}
```

---

### 3. Query Invoice & Transaction Confirmations

```go
inv, err := client.GetInvoice(ctx, "inv_0192a3b4c5d6")
if err != nil {
    log.Fatalf("Failed to get invoice: %v", err)
}

fmt.Printf("Status: %s\n", inv.Invoice.Status) // "WAITING" | "PAID" | "EXPIRED" | "CANCELED"

for _, p := range inv.Payments {
    fmt.Printf("Payment received: %s (Confirmations: %d, Status: %s)\n",
        p.AmountReceived, p.Confirmations, p.Status)
    fmt.Printf("Blockchain Tx: %s\n", p.TxHash)
}
```

---

### 4. List & Filter Invoices

```go
res, err := client.ListInvoices(ctx, kkhay.ListInvoicesQuery{
    Page:   1,
    Limit:  20,
    Status: kkhay.StatusPaid,
})
if err != nil {
    log.Fatalf("Failed to list invoices: %v", err)
}

fmt.Printf("Total paid invoices: %d\n", res.TotalCount)
for _, item := range res.Items {
    fmt.Printf("- %s: %f %s (%s)\n", item.ID, item.Amount, item.Currency, item.Status)
}
```

---

### 5. Verify Incoming Webhooks (IPN)

When a customer pays and the blockchain confirms the transaction, K Khay sends a cryptographically signed `POST` webhook to your server with the `x-kkhay-signature` header.

```go
package main

import (
    "io"
    "log"
    "net/http"
    "os"

    "github.com/boboaungdev/kkhay-go"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Cannot read body", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    signature := r.Header.Get("x-kkhay-signature")
    secret := os.Getenv("KKHAY_IPN_SECRET")

    event, err := kkhay.ParseWebhookEvent(body, signature, secret)
    if err != nil {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    if event.Event == kkhay.EventPaymentPaid || event.Event == kkhay.EventPaymentFinished {
        log.Printf("Payment confirmed for invoice %s (Gateway: %s)!", event.InvoiceID, event.GatewayID)
        log.Printf("Received: %s %s on %s", event.PayAmount, event.PayToken, event.PayNetwork)
        if event.TxHash != nil {
            log.Printf("Transaction Hash: %s", *event.TxHash)
        }
        // Fulfill your customer's order in your database
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"ok": true}`))
}

func main() {
    http.HandleFunc("/webhooks/kkhay", webhookHandler)
    log.Println("Server listening on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 🌐 Supported Networks & Families

| Network Family | Popular Chains | Supported Tokens |
| :--- | :--- | :--- |
| **EVM (180+ Chains)** | Ethereum, BNB Smart Chain, Polygon, Arbitrum, Base, Optimism, Avalanche | USDT, USDC, ETH, BNB, POL, AVAX, DAI, and all ERC-20s |
| **UTXO** | Bitcoin, Dogecoin, Litecoin, Bitcoin Cash | BTC, DOGE, LTC, BCH |
| **Solana** | Solana Mainnet | SOL, USDT, USDC |
| **TON** | The Open Network | TON, USDT |
| **Tron (TVM)** | Tron Mainnet | TRX, USDT |
| **XRPL** | XRP Ledger | XRP |

---

## 🛡️ Error Handling

The SDK returns structured `*kkhay.APIError` instances with HTTP status codes and error messages:

```go
inv, err := client.GetInvoice(ctx, "inv_nonexistent")
if err != nil {
    if apiErr, ok := err.(*kkhay.APIError); ok {
        fmt.Printf("API Error HTTP %d [%s]: %s\n", apiErr.Status, apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

---

## 📄 License

MIT © [K Khay](https://kkhay.com)
