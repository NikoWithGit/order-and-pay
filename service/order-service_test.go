package service

import (
	"errors"
	mock_intrface "order-and-pay/mock"
	"order-and-pay/model"
	"testing"

	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
)

func TestAddProduct(t *testing.T) {
	productStandart := &model.ProductInOrder{
		Uuid:        uuid.NewString(),
		Num:         7,
		PricePerOne: 6.8,
		OrderId:     uuid.NewString(),
	}
	productWithZeroNum := &model.ProductInOrder{
		Uuid:        uuid.NewString(),
		Num:         0,
		PricePerOne: 6.8,
		OrderId:     uuid.NewString(),
	}

	tests := []struct {
		name        string
		prepareMock func(m *mock_intrface.MockOrderRepo, p *model.ProductInOrder)
		arg         *model.ProductInOrder
		isErr       bool
	}{
		{
			name: "test with successful adding new product (db hasn't got this kind of product)",
			prepareMock: func(m *mock_intrface.MockOrderRepo, p *model.ProductInOrder) {
				var err error = nil
				gomock.InOrder(
					m.EXPECT().Begin().Return(err),
					m.EXPECT().GetProductId(p).Return(-1, err),
					m.EXPECT().AddProduct(p).Return(err),
					m.EXPECT().Commit().Return(err),
					m.EXPECT().Rollback(),
				)
			},
			arg:   productStandart,
			isErr: false,
		},
		{
			name: "test with successful num updating (db has got this kind of product)",
			prepareMock: func(m *mock_intrface.MockOrderRepo, p *model.ProductInOrder) {
				var err error = nil
				productInDbRecordId := 19
				gomock.InOrder(
					m.EXPECT().Begin().Return(err),
					m.EXPECT().GetProductId(p).Return(productInDbRecordId, err),
					m.EXPECT().UpdateProductNumById(p.Num, uint(productInDbRecordId)),
					m.EXPECT().Commit(),
					m.EXPECT().Rollback(),
				)
			},
			arg:   productStandart,
			isErr: false,
		},
		{
			name: "test with successful product deleting (num = 0)",
			prepareMock: func(m *mock_intrface.MockOrderRepo, p *model.ProductInOrder) {
				var err error = nil
				gomock.InOrder(
					m.EXPECT().DeleteProduct(p).Return(err),
				)
			},
			arg:   productWithZeroNum,
			isErr: false,
		},
		{
			name: "test with internal error on Begin()",
			prepareMock: func(m *mock_intrface.MockOrderRepo, p *model.ProductInOrder) {
				var err error = errors.New("Internal error")
				gomock.InOrder(
					m.EXPECT().Begin().Return(err),
				)
			},
			arg:   productStandart,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOrderRepo := mock_intrface.NewMockOrderRepo(ctrl)
			tt.prepareMock(mockOrderRepo, tt.arg)
			service := NewOrderService(mockOrderRepo)
			err := service.AddProduct(tt.arg)
			isErrorNotNil := err == nil
			if (err != nil) != tt.isErr {
				t.Errorf("TEST ERROR: \n"+
					"Got: isErr = %t \n"+
					"Expected: isErr = %t \n",
					isErrorNotNil, tt.isErr,
				)
			}
		})
	}
}

func TestFinish(t *testing.T) {

	type Expected struct {
		isUpdated         bool
		isBadRequestError bool
		isErr             bool
	}

	tests := []struct {
		name        string
		prepareMock func(m *mock_intrface.MockOrderRepo, orderId string)
		arg         string
		exp         *Expected
	}{
		{
			name: "test with successful status updating",
			prepareMock: func(m *mock_intrface.MockOrderRepo, orderId string) {
				status := model.CREATED
				var err error = nil
				paymentsSum := float32(25.7)
				productsPricesSum := paymentsSum
				gomock.InOrder(
					m.EXPECT().Begin().Return(err),
					m.EXPECT().GetOrderStatus(orderId).Return(status, err),
					m.EXPECT().GetPaymentsSumByOrderId(orderId).Return(paymentsSum, err),
					m.EXPECT().GetProductsPriceSumByOrderId(orderId).Return(productsPricesSum, err),
					m.EXPECT().UpdateOrderStatusToComplete(orderId).Return(err),
					m.EXPECT().Commit(),
					m.EXPECT().Rollback(),
				)
			},
			arg: uuid.NewString(),
			exp: &Expected{
				isUpdated:         true,
				isBadRequestError: false,
				isErr:             false,
			},
		},
		{
			name: "test without status updating, coz its already updated",
			prepareMock: func(m *mock_intrface.MockOrderRepo, orderId string) {
				status := model.COMPLITED
				var err error = nil
				gomock.InOrder(
					m.EXPECT().Begin().Return(nil),
					m.EXPECT().GetOrderStatus(orderId).Return(status, err),
					m.EXPECT().Rollback(),
				)
			},
			arg: uuid.NewString(),
			exp: &Expected{
				isUpdated:         false,
				isBadRequestError: false,
				isErr:             false,
			},
		},
		{
			name: "test with badRequestError",
			prepareMock: func(m *mock_intrface.MockOrderRepo, orderId string) {
				status := model.CREATED
				var err error = nil
				paymentsSum := float32(25.7)
				productsPricesSum := float32(13.9)
				gomock.InOrder(
					m.EXPECT().Begin().Return(err),
					m.EXPECT().GetOrderStatus(orderId).Return(status, err),
					m.EXPECT().GetPaymentsSumByOrderId(orderId).Return(paymentsSum, err),
					m.EXPECT().GetProductsPriceSumByOrderId(orderId).Return(productsPricesSum, err),
					m.EXPECT().Rollback(),
				)
			},
			arg: uuid.NewString(),
			exp: &Expected{
				isUpdated:         false,
				isBadRequestError: true,
				isErr:             false,
			},
		},
		{
			name: "test with internalError on Begin()",
			prepareMock: func(m *mock_intrface.MockOrderRepo, orderId string) {
				err := errors.New("Internal error")
				gomock.InOrder(
					m.EXPECT().Begin().Return(err),
				)
			},
			arg: uuid.NewString(),
			exp: &Expected{
				isUpdated:         false,
				isBadRequestError: false,
				isErr:             true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOrderRepo := mock_intrface.NewMockOrderRepo(ctrl)
			tt.prepareMock(mockOrderRepo, tt.arg)
			service := NewOrderService(mockOrderRepo)
			isUpdated, badRequestError, err := service.Finish(tt.arg)
			isBadRequestError := (badRequestError != nil)
			isErr := (err != nil)
			if isUpdated != tt.exp.isUpdated || isBadRequestError != tt.exp.isBadRequestError || isErr != tt.exp.isErr {
				t.Errorf("TEST ERROR: \n"+
					"Got: isUpdated = %t, isBadRequestError = %t, isErr =  %t \n"+
					"Expected:  isUpdated = %t, isBadRequestError = %t, isErr =  %t \n",
					isUpdated, isBadRequestError, isErr,
					tt.exp.isUpdated, tt.exp.isBadRequestError, tt.exp.isErr,
				)
			}
		})
	}
}

/*func TestCreate(t *testing.T) {
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
}*/
