package intrface

import (
	"order-and-pay/model"
	"time"
)

type OrderRepo interface {
	Create() (string, uint)

	GetAll(from, to time.Time) []model.Order
	GetById(orderId string) *model.Order
	GetProductId(p *model.ProductInOrder) int
	GetPaymentsByOrderId(orderId string) []model.Payment
	GetProductsByOrderId(orderId string) []model.ProductInOrder
	GetPaymentsSumByOrderId(orderId string) float32
	GetProductsPriceSumByOrderId(orderId string) float32

	UpdateProductNumById(num uint, id uint) *model.ProductInOrder
	UpdateOrderStatusToComplete(orderId string)

	AddPayment(p *model.Payment)
	AddProduct(p *model.ProductInOrder)

	DeleteProduct(p *model.ProductInOrder) *model.ProductInOrder

	IsExists(orderId string) bool
}
