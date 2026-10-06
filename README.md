# K Khay Go SDK 🚀

Official Go SDK for the **[K Khay Sovereign Crypto Payment Gateway](https://kkhay.com)**.

Accept non-custodial and custodial crypto payments (USDT, USDC, BNB, ETH on BSC, Polygon, Arbitrum, Base, Ethereum) with **zero external dependencies** using the Go standard library.

---

## 📦 Installation

```bash
go get github.com/boboaungdev/kkhay-go
```

*(Requires Go 1.20 or newer)*

---

## ⚡ Quick Start

### 1. Initialize Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/boboaungdev/kkhay-go"
)

func main() {
    client, err := kkhay.NewClient("kkhay_live_your_api_key_here")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create an invoice
    resp, err := client.CreateInvoice(ctx, kkhay.CreateInvoiceRequest{
        PriceAmount:   49.99,
        PriceCurrency: "USD",
        PayNetwork:    "bsc",
        PayToken:      "USDT",
        OrderID:       "ORD-8821",
        Title:         "Pro Tier Subscription",
        RedirectURL:   "https://myshop.com/success",
    })
    if err != nil {
        log.Fatalf("Invoice error: %v", err)
    }

    fmt.Printf("Invoice ID: %s\n", resp.Invoice.ID)
    fmt.Printf("Checkout URL: %s\n", resp.Invoice.HostedURL)
}
```

### 2. Query Invoice & Transaction Confirmations

```go
inv, err := client.GetInvoice(ctx, "inv_9f81a7b2")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", inv.Invoice.Status)
for _, p := range inv.Payments {
    fmt.Printf("Received %s on tx %s (confirms: %d)\n", p.AmountReceived, p.TxHash, p.Confirmations)
}
```

---

## 🔐 Webhook / IPN Verification (`net/http`)

```go
package main

import (
    "io"
    "log"
    "net/http"
    "os"

    "github.com/boboaungdev/kkhay-go"
)

func handleKkhayWebhook(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    sig := r.Header.Get("x-kkhay-signature")
    secret := os.Getenv("KKHAY_IPN_SECRET")

    event, err := kkhay.ParseWebhookEvent(body, sig, secret)
    if err != nil {
        http.Error(w, "Unauthorized signature", http.StatusUnauthorized)
        return
    }

    if event.Event == "payment.finished" {
        log.Printf("Order %s was paid! Tx: %v", *event.OrderID, event.TxHash)
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"ok": true}`))
}
```

---

## 📄 License

MIT © [K Khay](https://kkhay.com)

