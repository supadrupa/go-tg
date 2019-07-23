package tg

import (
	"encoding/json"
)

// SuccessfulPayment object contains basic information about a successful payment.
type SuccessfulPayment struct {
	// Three-letter ISO 4217 currency code.
	Currency string `json:"currency"`

	// Total price in the smallest units of the currency.
	// For example, for a price of US$ 1.45 pass amount = 145.
	// See the exp parameter in currencies.json, it shows the number
	// of digits past the decimal point for each currency (2 for the majority of currencies).
	TotalAmount int `json:"total_amount"`

	// Bot specified invoice payload.
	InvoicePayload string `json:"invoice_payload"`

	// Optional. Identifier of the shipping option chosen by the user
	ShippingOptionID string `json:"shipping_option_id,omitempty"`

	// Optional. Order info provided by the user
	OrderInfo json.RawMessage `json:"order_info,omitempty"`

	// Telegram payment identifier.
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`

	// Provider payment identifier.
	ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
}

// OrderInfo object represents information about an order.
type OrderInfo struct {
	// Optional. User name
	Name string `json:"name,omitempty"`

	// Optional. User's phone number
	PhoneNumber string `json:"phone_number,omitempty"`

	// Optional. User email.
	Email string `json:"email,omitempty"`

	// Optional. User shipping address
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// ShippingAddress object represents a shipping address.
type ShippingAddress struct {
	// ISO 3166-1 alpha-2 country code
	CountryCode string `json:"country_code"`

	// State, if applicable
	State string `json:"state"`

	// City
	City string `json:"city"`

	// First line for the address
	StreetLine1 string `json:"street_line_1"`

	// Second line for the address
	StreetLine2 string `json:"street_line_2"`

	// Address post code
	PostCode string `json:"post_code"`
}

// Invoice object contains basic information about an invoice.
type Invoice struct {
	// Product name
	Title string `json:"title"`

	// Product description
	Descrption string `json:"descrption"`

	// Unique bot deep-linking parameter that can be used to generate this invoice
	StartParameter string `json:"start_parameter"`

	// Three-letter ISO 4217 currency code.
	Currency string `json:"currency"`

	// Total price in the smallest units of the currency.
	TotalAmount int `json:"total_amount"`
}

// ShippingQueryID unique shipping query identifier.
type ShippingQueryID string

// ShippingQuery object contains information about an incoming shipping query.
type ShippingQuery struct {
	// Unique query identifier.
	ID ShippingQueryID `json:"id"`

	// User who sent the query.
	From User `json:"from"`

	// Bot specified invoice payload.
	InvoicePayload string `json:"invoice_payload"`

	// User specified shipping address.
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQueryID pre-checkout query identifier.
type PreCheckoutQueryID string

// PreCheckoutQuery object contains information about an incoming pre-checkout query.
type PreCheckoutQuery struct {
	// Unique query identifier.
	ID PreCheckoutQueryID `json:"id"`

	// User who sent the query.
	From User `json:"from"`

	// Three-letter ISO 4217 currency code.
	Currency string `json:"currency"`

	// Total price in the smallest units of the currency.
	TotalAmount int `json:"total_amount"`

	// Bot specified invoice payload.
	InvoicePayload string `json:"invoice_payload"`

	// Optional. Identifier of the shipping option chosen by the user.
	ShippingOptionID string `json:"shipping_option_id,omitempty"`

	// Optional. Order info provided by the user.
	OrderInfo *OrderInfo `json:"order_info,omitempty"`
}
