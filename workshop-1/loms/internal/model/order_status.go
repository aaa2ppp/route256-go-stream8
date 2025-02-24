package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type OrderStatus int

const (
	_ OrderStatus = iota
	OrderStatusNew
	OrderStatusAwaitingPayment
	OrderStatusFailed
	OrderStatusPayed
	OrderStatusCancelled
)

func ParseOrderStatus(s string) (OrderStatus, error) {
	switch s {
	case "new":
		return OrderStatusNew, nil
	case "awaiting payment":
		return OrderStatusAwaitingPayment, nil
	case "failed":
		return OrderStatusFailed, nil
	case "payed":
		return OrderStatusPayed, nil
	case "cancelled":
		return OrderStatusCancelled, nil
	}
	return 0, fmt.Errorf("ParseOrderStatus: unknown value %q", s)
}

func (os OrderStatus) String() string {
	switch os {
	case OrderStatusNew:
		return "new"
	case OrderStatusAwaitingPayment:
		return "awaiting payment"
	case OrderStatusFailed:
		return "failed"
	case OrderStatusPayed:
		return "payed"
	case OrderStatusCancelled:
		return "cancelled"
	default:
		return fmt.Sprintf("OrderStatus(%d)", os)
	}
}

// MarshalJSON implements json.Marshaler.
func (os OrderStatus) MarshalJSON() ([]byte, error) {
	return []byte(os.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (os *OrderStatus) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return fmt.Errorf("OrderSttus.UnmarshalJSON: %w", err)
	}
	v, err := ParseOrderStatus(s)
	if err != nil {
		return err
	}
	*os = v
	return nil
}

var _ json.Marshaler = OrderStatus(0)
var _ json.Unmarshaler = new(OrderStatus)
