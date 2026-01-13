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
