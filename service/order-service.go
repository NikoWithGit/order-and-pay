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
}

func NewOrderService(or intrface.OrderRepo) *OrderService {
	return &OrderService{or}
}

func (os *OrderService) Create() (string, uint) {
	return os.repo.Create()
}

func (os *OrderService) GetAll(from time.Time, to time.Time) []model.Order {
	return os.repo.GetAll(from, to)
}

func (os *OrderService) Get(orderId string) model.Order {
	return *os.repo.GetById(orderId)
}

func (os *OrderService) AddProduct(p *model.ProductInOrder) {
	if p.Num == 0 {
		os.repo.DeleteProduct(p)
		return
	}
	productId, isExists := os.repo.CheckAndGetProductId(p)
	if isExists {
		os.repo.UpdateProductNumById(p.Num, productId)
	} else {
		os.repo.AddProduct(p)
	}
}

func (os *OrderService) AddPayment(p *model.Payment) {
	os.repo.AddPayment(p)
}

func (os *OrderService) Finish(orderId string) error {
	paymentGot := os.repo.GetPaymentsSumByOrderId(orderId)
	paymentNeed := os.repo.GetProductsPriceSumByOrderId(orderId)
	if !floatEq(paymentNeed, paymentGot) {
		return errors.New("WRONG TRANSACTION PAYMENTS")
	}
	os.repo.UpdateOrderStatusToComplete(orderId)
	return nil
}

func floatEq(f1, f2 float32) bool {
	return math.Abs(float64(f1-f2)) < 0.00001
}
