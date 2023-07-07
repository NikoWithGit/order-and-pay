package service

import (
	mock_intrface "order-and-pay/mock"
	"order-and-pay/model"
	"testing"

	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderRepo := mock_intrface.NewMockOrderRepo(ctrl)

	var id string = uuid.NewString()
	var short uint = 100000
	var err error = nil
	mockOrderRepo.EXPECT().Create().Return(id, short, err)

	service := NewOrderService(mockOrderRepo)
	resId, resShort, resErr := service.Create()

	if resId != id || resShort != short || resErr != err {
		t.Errorf("TEST ERROR: \n"+
			"expected id=%s, short=%d, err=%v; \n"+
			"result id=%s, short=%d, err=%v;",
			id, short, err,
			resId, resShort, resErr,
		)
	}
}

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrderRepo := mock_intrface.NewMockOrderRepo(ctrl)

	uid := uuid.NewString()
	var num uint = 7
	var price float32 = 6.8
	orderId := uuid.NewString()
	product1 := &model.ProductInOrder{uid, num, price, orderId}
	product2 := &model.ProductInOrder{uid, num, price, orderId}
	product3 := &model.ProductInOrder{uid, 0, price, orderId}

	gomock.InOrder(
		mockOrderRepo.EXPECT().Begin(),
		mockOrderRepo.EXPECT().GetProductId(product1).Return(-1, nil),
		mockOrderRepo.EXPECT().AddProduct(product1),
		mockOrderRepo.EXPECT().Commit(),
		mockOrderRepo.EXPECT().Rollback(),
	)

	gomock.InOrder(
		mockOrderRepo.EXPECT().Begin(),
		mockOrderRepo.EXPECT().GetProductId(product2).Return(19, nil),
		mockOrderRepo.EXPECT().UpdateProductNumById(product2.Num, uint(19)),
		mockOrderRepo.EXPECT().Commit(),
		mockOrderRepo.EXPECT().Rollback(),
	)

	gomock.InOrder(
		mockOrderRepo.EXPECT().DeleteProduct(product3),
	)

	service := NewOrderService(mockOrderRepo)
	service.AddProduct(product1)
	service.AddProduct(product2)
	service.AddProduct(product3)
}
