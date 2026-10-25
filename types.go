package kkhay

// Invoice status constants matching PostgreSQL GatewayInvoiceStatus enum.
const (
	StatusWaiting  = "WAITING"
	StatusPaid     = "PAID"
	StatusExpired  = "EXPIRED"
	StatusCanceled = "CANCELED"
)

// Invoice type constants matching PostgreSQL GatewayInvoiceType enum.
const (
	TypePayment      = "PAYMENT"
	TypeDonation     = "DONATION"
	TypeSubscription = "SUBSCRIPTION"
)

// Settlement mode constants matching PostgreSQL GatewaySettlementMode enum.
const (
	SettlementModeNonCustodial = "NON_CUSTODIAL"
	SettlementModeCustodial    = "CUSTODIAL"
)

// Payment method constants matching PostgreSQL GatewayPaymentMethod enum.
const (
	PaymentMethodDeposit = "DEPOSIT"
	PaymentMethodWallet  = "WALLET"
)

// Payment status constants matching PostgreSQL GatewayPaymentStatus enum.
const (
	PaymentStatusSettling = "SETTLING"
	PaymentStatusSettled  = "SETTLED"
	PaymentStatusFailed   = "FAILED"
)

// Webhook event constants sent via Instant Payment Notifications (IPN).
const (
	EventPaymentPaid       = "payment.paid"
	EventPaymentFinished   = "payment.finished"
	EventPaymentConfirming = "payment.confirming"
	EventPaymentFailed     = "payment.failed"
	EventInvoiceCanceled   = "invoice.canceled"
)

// CreateInvoiceRequest parameters for generating a new crypto payment invoice.
type CreateInvoiceRequest struct {
	Amount          float64                `json:"amount,omitempty"`
	Currency        string                 `json:"currency,omitempty"`
	PriceAmount     float64                `json:"priceAmount,omitempty"`
	PriceCurrency   string                 `json:"priceCurrency,omitempty"`
	PayNetwork      string                 `json:"payNetwork"`
	PayToken        string                 `json:"payToken"`
	SettlementAsset string                 `json:"settlementAsset,omitempty"`
	Type            string                 `json:"type,omitempty"`
	OrderID         string                 `json:"orderId,omitempty"`
	Title           string                 `json:"title,omitempty"`
	CustomerName    string                 `json:"customerName,omitempty"`
	CustomerEmail   string                 `json:"customerEmail,omitempty"`
	RedirectURL     string                 `json:"redirectUrl,omitempty"`
	CancelURL       string                 `json:"cancelUrl,omitempty"`
	IPNCallbackURL  string                 `json:"ipnCallbackUrl,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// Invoice represents a K Khay sovereign payment invoice record.
type Invoice struct {
	ID              string                 `json:"id"`
	GatewayID       string                 `json:"gatewayId"`
	Type            string                 `json:"type"`
	OrderID         *string                `json:"orderId"`
	Title           *string                `json:"title"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	PriceAmount     float64                `json:"priceAmount"`
	PriceCurrency   string                 `json:"priceCurrency"`
	SettlementAsset *string                `json:"settlementAsset"`
	PayAmount       *string                `json:"payAmount"`
	PayToken        *string                `json:"payToken"`
	PayNetwork      *string                `json:"payNetwork"`
	LockedRate      *string                `json:"lockedRate"`
	DepositAddress  string                 `json:"depositAddress"`
	DerivationIndex *int                   `json:"derivationIndex,omitempty"`
	PayoutAddress   string                 `json:"payoutAddress"`
	Status          string                 `json:"status"`
	FeeAmount       string                 `json:"feeAmount"`
	NetAmount       string                 `json:"netAmount"`
	CustomerName    *string                `json:"customerName"`
	CustomerEmail   *string                `json:"customerEmail"`
	RedirectURL     *string                `json:"redirectUrl"`
	CancelURL       *string                `json:"cancelUrl"`
	IPNCallbackURL  *string                `json:"ipnCallbackUrl"`
	PaymentMethod   *string                `json:"paymentMethod"`
	SettlementMode  *string                `json:"settlementMode"`
	HostedURL       string                 `json:"hostedUrl"`
	Metadata        map[string]interface{} `json:"metadata"`
	ExpiresAt       string                 `json:"expiresAt"`
	CreatedAt       string                 `json:"createdAt"`
	UpdatedAt       string                 `json:"updatedAt"`
}

// PaymentRecord details on-chain transaction hashes, confirmations, and settlement status.
type PaymentRecord struct {
	ID              string  `json:"id"`
	InvoiceID       string  `json:"invoiceId"`
	TxHash          string  `json:"txHash"`
	FromAddress     *string `json:"fromAddress"`
	AmountReceived  string  `json:"amountReceived"`
	Confirmations   int     `json:"confirmations"`
	Status          string  `json:"status"`
	ForwardedTxHash *string `json:"forwardedTxHash"`
	ForwardedAt     *string `json:"forwardedAt"`
	Token           *string `json:"token,omitempty"`
	Network         *string `json:"network,omitempty"`
	OrderID         *string `json:"orderId,omitempty"`
	Title           *string `json:"title,omitempty"`
	CreatedAt       string  `json:"createdAt"`
}

// CreateInvoiceResponse contains the created invoice and hosted checkout URL.
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
	Type   string `json:"type,omitempty"`
	Search string `json:"search,omitempty"`
}

// ListInvoicesResponse paginated list of gateway invoices.
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
	GatewayID      string  `json:"gateway_id,omitempty"`
	InvoiceType    string  `json:"invoice_type,omitempty"`
	OrderID        *string `json:"order_id"`
	Amount         float64 `json:"amount,omitempty"`
	Currency       string  `json:"currency,omitempty"`
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
