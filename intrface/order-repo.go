package intrface

import (
	"order-and-pay/model"
	"time"
)

type OrderRepo interface {
	Create() (string, uint, error)

	GetAll(from, to time.Time) ([]model.Order, error)
	GetById(orderId string) (*model.Order, error)
	GetOrderStatus(tx Idb, orderId string) (uint8, error)
	GetProductId(tx Idb, p *model.ProductInOrder) (int, error)
	GetPaymentsByOrderId(orderId string) ([]model.Payment, error)
	GetProductsByOrderId(orderId string) ([]model.ProductInOrder, error)
	GetPaymentsSumByOrderId(tx Idb, orderId string) (float32, error)
	GetProductsPriceSumByOrderId(tx Idb, orderId string) (float32, error)

	UpdateProductNumById(tx Idb, num uint, id uint) error
	UpdateOrderStatusToComplete(tx Idb, orderId string) error

	AddPayment(p *model.Payment) error
	AddProduct(Idb, *model.ProductInOrder) error

	DeleteProduct(p *model.ProductInOrder) error

	GetDb() Idb
}
