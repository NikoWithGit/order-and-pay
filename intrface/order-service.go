package intrface

import (
	"order-and-pay/model"
	"time"
)

type OrderService interface {
	Create() (string, uint, error)
	GetAll(from time.Time, to time.Time) ([]model.Order, error)
	Get(orderId string) (*model.Order, error)
	AddProduct(p *model.ProductInOrder) error
	AddPayment(p *model.Payment) error
	Finish(orderId string) (bool, error, error)
}
