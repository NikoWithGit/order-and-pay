package service

import (
	"errors"
	"math"
	"order-and-pay/intrface"
	"order-and-pay/model"
	"time"
)

type OrderService struct {
	repo intrface.OrderRepo
	db   intrface.Idb
}

func NewOrderService(or intrface.OrderRepo) *OrderService {
	return &OrderService{or, or.GetDb()}
}

func (os *OrderService) Create() (string, uint, error) {
	return os.repo.Create()
}

func (os *OrderService) GetAll(from time.Time, to time.Time) ([]model.Order, error) {
	return os.repo.GetAll(from, to)
}

func (os *OrderService) Get(orderId string) (*model.Order, error) {
	return os.repo.GetById(orderId)
}

func (os *OrderService) AddProduct(p *model.ProductInOrder) error {
	if p.Num == 0 {
		err := os.repo.DeleteProduct(p)
		return err
	}

	tx, err := os.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	productId, err := os.repo.GetProductId(tx, p)
	if err != nil {
		return err
	}
	if productId != -1 {
		err = os.repo.UpdateProductNumById(tx, p.Num, uint(productId))
	} else {
		err = os.repo.AddProduct(tx, p)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (os *OrderService) AddPayment(p *model.Payment) error {
	return os.repo.AddPayment(p)
}

func (os *OrderService) Finish(orderId string) (bool, error, error) {
	tx, err := os.db.Begin()
	if err != nil {
		return false, nil, err
	}
	defer tx.Rollback()

	status, err := os.repo.GetOrderStatus(tx, orderId)
	if err != nil {
		return false, nil, err
	}
	if status == model.COMPLITED {
		return false, nil, nil
	}
	paymentGot, err := os.repo.GetPaymentsSumByOrderId(tx, orderId)
	if err != nil {
		return false, nil, err
	}
	paymentNeed, err := os.repo.GetProductsPriceSumByOrderId(tx, orderId)
	if err != nil {
		return false, nil, err
	}
	if !floatEq(paymentNeed, paymentGot) {
		badRequestError := errors.New("WRONG TRANSACTION PAYMENTS")
		return false, badRequestError, nil
	}

	if err = os.repo.UpdateOrderStatusToComplete(tx, orderId); err != nil {
		return false, nil, err
	}
	err = tx.Commit()
	if err != nil {
		return false, nil, err
	}
	return true, nil, nil
}

func floatEq(f1, f2 float32) bool {
	return math.Abs(float64(f1-f2)) < 0.00001
}
