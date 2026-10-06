package kkhay

// CreateInvoiceRequest parameters for generating a new crypto invoice.
type CreateInvoiceRequest struct {
	PriceAmount    float64                `json:"priceAmount"`
	PriceCurrency  string                 `json:"priceCurrency,omitempty"`
	PayNetwork     string                 `json:"payNetwork"`
	PayToken       string                 `json:"payToken"`
	OrderID        string                 `json:"orderId,omitempty"`
	Title          string                 `json:"title,omitempty"`
	CustomerName   string                 `json:"customerName,omitempty"`
	CustomerEmail  string                 `json:"customerEmail,omitempty"`
	RedirectURL    string                 `json:"redirectUrl,omitempty"`
	CancelURL      string                 `json:"cancelUrl,omitempty"`
	IPNCallbackURL string                 `json:"ipnCallbackUrl,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// Invoice represents a K Khay invoice record.
type Invoice struct {
	ID             string  `json:"id"`
	MerchantID     string  `json:"merchantId"`
	OrderID        *string `json:"orderId"`
	Title          *string `json:"title"`
	PriceAmount    float64 `json:"priceAmount"`
	PriceCurrency  string  `json:"priceCurrency"`
	PayAmount      string  `json:"payAmount"`
	PayToken       string  `json:"payToken"`
	PayNetwork     string  `json:"payNetwork"`
	DepositAddress string  `json:"depositAddress"`
	Status         string  `json:"status"`
	FeeAmount      string  `json:"feeAmount"`
	NetAmount      string  `json:"netAmount"`
	SettlementMode *string `json:"settlementMode"`
	PaymentMethod  *string `json:"paymentMethod"`
	HostedURL      string  `json:"hostedUrl"`
	RedirectURL    *string `json:"redirectUrl"`
	CancelURL      *string `json:"cancelUrl"`
	ExpiresAt      string  `json:"expiresAt"`
	CreatedAt      string  `json:"createdAt"`
}

// PaymentRecord details on-chain transaction hashes and confirmations.
type PaymentRecord struct {
	ID              string  `json:"id"`
	TxHash          string  `json:"txHash"`
	AmountReceived  string  `json:"amountReceived"`
	Confirmations   int     `json:"confirmations"`
	Status          string  `json:"status"`
	ForwardedTxHash *string `json:"forwardedTxHash"`
	CreatedAt       string  `json:"createdAt"`
}

// CreateInvoiceResponse contains the created invoice and status.
type CreateInvoiceResponse struct {
	OK      bool    `json:"ok"`
	Invoice Invoice `json:"invoice"`
}

// GetInvoiceResponse contains the requested invoice and transaction records.
type GetInvoiceResponse struct {
	OK       bool            `json:"ok"`
	Invoice  Invoice         `json:"invoice"`
	Payments []PaymentRecord `json:"payments"`
}

// ListInvoicesQuery parameters for pagination and filtering.
type ListInvoicesQuery struct {
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Status string `json:"status,omitempty"`
	Search string `json:"search,omitempty"`
}

// ListInvoicesResponse paginated list of merchant invoices.
type ListInvoicesResponse struct {
	OK          bool      `json:"ok"`
	Items       []Invoice `json:"items"`
	TotalCount  int       `json:"totalCount"`
	TotalPages  int       `json:"totalPages"`
	CurrentPage int       `json:"currentPage"`
	Limit       int       `json:"limit"`
}

// WebhookEvent payload sent via IPN.
type WebhookEvent struct {
	Event          string  `json:"event"`
	InvoiceID      string  `json:"invoice_id"`
	OrderID        *string `json:"order_id"`
	PriceAmount    float64 `json:"price_amount"`
	PriceCurrency  string  `json:"price_currency"`
	PayAmount      string  `json:"pay_amount"`
	PayToken       string  `json:"pay_token"`
	PayNetwork     string  `json:"pay_network"`
	DepositAddress string  `json:"deposit_address"`
	TxHash         *string `json:"tx_hash"`
	Status         string  `json:"status"`
	Timestamp      string  `json:"timestamp"`
}

