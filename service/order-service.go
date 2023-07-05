package service

import (
	"errors"
	"math"
	"order-and-pay/intrface"
	"order-and-pay/model"
	"time"
)

type OrderService struct {
	repo   intrface.OrderRepo
	logger intrface.Ilogger
}

func NewOrderService(or intrface.OrderRepo, l intrface.Ilogger) *OrderService {
	return &OrderService{or, l}
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
	if err := os.repo.Begin(); err != nil {
		return err
	}
	defer os.repo.Rollback()

	if p.Num == 0 {
		err := os.repo.DeleteProduct(p)
		return err
	}

	productId, err := os.repo.GetProductId(p)
	if err != nil {
		return err
	}
	if productId != -1 {
		err = os.repo.UpdateProductNumById(p.Num, uint(productId))
	} else {
		err = os.repo.AddProduct(p)
	}
	if err != nil {
		return err
	}
	err = os.repo.Commit()
	return err
}

func (os *OrderService) AddPayment(p *model.Payment) error {
	err := os.repo.AddPayment(p)
	return err
}

func (os *OrderService) Finish(orderId string) (bool, error, error) {

	if err := os.repo.Begin(); err != nil {
		return false, nil, err
	}
	defer os.repo.Rollback()

	status, err := os.repo.GetOrderStatus(orderId)
	if err != nil {
		return false, nil, err
	}
	if status == model.COMPLITED {
		return false, nil, nil
	}
	paymentGot, err := os.repo.GetPaymentsSumByOrderId(orderId)
	if err != nil {
		return false, nil, err
	}
	paymentNeed, err := os.repo.GetProductsPriceSumByOrderId(orderId)
	if err != nil {
		return false, nil, err
	}
	if !floatEq(paymentNeed, paymentGot) {
		return false, errors.New("WRONG TRANSACTION PAYMENTS"), nil
	}

	if err = os.repo.UpdateOrderStatusToComplete(orderId); err != nil {
		return false, nil, err
	}

	if err = os.repo.Commit(); err != nil {
		return false, nil, err
	}
	return true, nil, nil
}

func floatEq(f1, f2 float32) bool {
	return math.Abs(float64(f1-f2)) < 0.00001
}
