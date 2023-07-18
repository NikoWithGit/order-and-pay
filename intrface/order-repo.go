package intrface

import (
	"order-and-pay/model"
	"time"
)

type OrderRepo interface {
	Create() (string, uint, error)

	GetAll(from, to time.Time) ([]model.Order, error)
	GetById(orderId string) (*model.Order, error)
	GetOrderStatus(tx Itx, orderId string) (uint8, error)
	GetProductId(tx Itx, p *model.ProductInOrder) (int, error)
	GetPaymentsByOrderId(orderId string) ([]model.Payment, error)
	GetProductsByOrderId(orderId string) ([]model.ProductInOrder, error)
	GetPaymentsSumByOrderId(tx Itx, orderId string) (float32, error)
	GetProductsPriceSumByOrderId(tx Itx, orderId string) (float32, error)

	UpdateProductNumById(tx Itx, num uint, id uint) error
	UpdateOrderStatusToComplete(tx Itx, orderId string) error

	AddPayment(p *model.Payment) error
	AddProduct(Itx, *model.ProductInOrder) error

	DeleteProduct(p *model.ProductInOrder) error

	GetDb() Idb
}
