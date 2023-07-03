package intrface

import (
	"order-and-pay/model"
	"time"
)

type OrderRepo interface {
	Create() (string, uint, error)

	GetAll(from, to time.Time) ([]model.Order, error)
	GetById(orderId string) (*model.Order, error)
	GetOrderStatus(orderId string) (uint8, error)
	GetProductId(p *model.ProductInOrder) (int, error)
	GetPaymentsByOrderId(orderId string) ([]model.Payment, error)
	GetProductsByOrderId(orderId string) ([]model.ProductInOrder, error)
	GetPaymentsSumByOrderId(orderId string) (float32, error)
	GetProductsPriceSumByOrderId(orderId string) (float32, error)

	UpdateProductNumById(num uint, id uint) error
	UpdateOrderStatusToComplete(orderId string) error

	AddPayment(p *model.Payment) error
	AddProduct(p *model.ProductInOrder) error

	DeleteProduct(p *model.ProductInOrder) error

	Begin() error
	Rollback()
	Commit() error
}
