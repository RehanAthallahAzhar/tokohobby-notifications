package messaging

import "time"

// OrderStatusChangedEvent represents order status change
type OrderStatusChangedEvent struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	ItemCount   int       `json:"item_count"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderCreatedEvent represents new order creation
type OrderCreatedEvent struct {
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	TotalAmount   float64   `json:"total_amount"`
	ItemCount     int       `json:"item_count"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}

// OrderShippedEvent represents order shipment
type OrderShippedEvent struct {
	OrderID          string    `json:"order_id"`
	UserID           string    `json:"user_id"`
	TrackingNumber   string    `json:"tracking_number"`
	Courier          string    `json:"courier"`
	EstimatedArrival time.Time `json:"estimated_arrival"`
	ShippedAt        time.Time `json:"shipped_at"`
}

// OrderPaidEvent represents payment confirmation
type OrderPaidEvent struct {
	OrderID        string    `json:"order_id"`
	UserID         string    `json:"user_id"`
	PaidAmount     float64   `json:"paid_amount"`
	PaymentMethod  string    `json:"payment_method"`
	PaymentGateway string    `json:"payment_gateway"`
	TransactionID  string    `json:"transaction_id"`
	PaidAt         time.Time `json:"paid_at"`
}

// OrderDeliveredEvent represents successful delivery
type OrderDeliveredEvent struct {
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	ReceiverName  string    `json:"receiver_name"`
	DeliveryProof string    `json:"delivery_proof"`
	DeliveredAt   time.Time `json:"delivered_at"`
}

// OrderCancelledEvent represents order cancellation
type OrderCancelledEvent struct {
	OrderID         string    `json:"order_id"`
	UserID          string    `json:"user_id"`
	CancelledBy     string    `json:"cancelled_by"`
	CancelReason    string    `json:"cancel_reason"`
	RefundAmount    float64   `json:"refund_amount"`
	CancellationFee float64   `json:"cancellation_fee"`
	CancelledAt     time.Time `json:"cancelled_at"`
}

// OrderRefundedEvent represents refund processing
type OrderRefundedEvent struct {
	OrderID         string    `json:"order_id"`
	UserID          string    `json:"user_id"`
	RefundAmount    float64   `json:"refund_amount"`
	RefundMethod    string    `json:"refund_method"`
	RefundReference string    `json:"refund_reference"`
	ExpectedCredit  time.Time `json:"expected_credit"`
	RefundedAt      time.Time `json:"refunded_at"`
}
