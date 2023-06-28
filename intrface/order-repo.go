package intrface

import (
	"order-and-pay/model"
	"time"
)

type OrderRepo interface {
	Create() (string, uint)
	UpdateProductNumById(num uint, id uint)
	CheckAndGetProductId(p *model.ProductInOrder) (uint, bool)
	DeleteProduct(p *model.ProductInOrder)
	GetPaymentsByOrderId(orderId string) []model.Payment
	GetProductsByOrderId(orderId string) []model.ProductInOrder
	UpdateOrderStatusToComplete(orderId string)
	GetProductsPriceSumByOrderId(orderId string) float32
	GetPaymentsSumByOrderId(orderId string) float32
	GetAll(frim, to time.Time) []model.Order
	GetById(orderId string) *model.Order
	AddPayment(p *model.Payment)
	AddProduct(p *model.ProductInOrder)
}
