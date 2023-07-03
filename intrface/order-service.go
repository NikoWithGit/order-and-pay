package intrface

import (
	"order-and-pay/model"
	"time"
)

type OrderService interface {
	Create() (string, uint)
	GetAll(from time.Time, to time.Time) []model.Order
	Get(orderId string) *model.Order
	AddProduct(p *model.ProductInOrder)
	AddPayment(p *model.Payment)
	Finish(orderId string) (bool, error)
}
