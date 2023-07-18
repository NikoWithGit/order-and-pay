package service

import (
	"errors"
	mock_intrface "order-and-pay/mock"
	"order-and-pay/model"
	"testing"

	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
)

type mocks struct {
	repo *mock_intrface.MockOrderRepo
	db   *mock_intrface.MockIdb
	tx   *mock_intrface.MockIdb
}

func TestAddProduct(t *testing.T) {
	productStandard := &model.ProductInOrder{
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
		prepareMock func(m *mocks, p *model.ProductInOrder)
		arg         *model.ProductInOrder
		isErr       bool
	}{
		{
			name: "test with successful adding new product (db hasn't got this kind of product)",
			prepareMock: func(m *mocks, p *model.ProductInOrder) {
				var err error = nil
				gomock.InOrder(
					m.db.EXPECT().Begin().Return(m.tx, err),
					m.repo.EXPECT().GetProductId(m.tx, p).Return(-1, err),
					m.repo.EXPECT().AddProduct(m.tx, p).Return(err),
					m.tx.EXPECT().Commit().Return(err),
					m.tx.EXPECT().Rollback(),
				)
			},
			arg:   productStandard,
			isErr: false,
		},
		{
			name: "test with successful num updating (db has got this kind of product)",
			prepareMock: func(m *mocks, p *model.ProductInOrder) {
				var err error = nil
				productInDbRecordId := 19
				gomock.InOrder(
					m.db.EXPECT().Begin().Return(m.tx, err),
					m.repo.EXPECT().GetProductId(m.tx, p).Return(productInDbRecordId, err),
					m.repo.EXPECT().UpdateProductNumById(m.tx, p.Num, uint(productInDbRecordId)),
					m.tx.EXPECT().Commit().Return(err),
					m.tx.EXPECT().Rollback(),
				)
			},
			arg:   productStandard,
			isErr: false,
		},
		{
			name: "test with successful product deleting (num = 0)",
			prepareMock: func(m *mocks, p *model.ProductInOrder) {
				var err error = nil
				gomock.InOrder(
					m.repo.EXPECT().DeleteProduct(p).Return(err),
				)
			},
			arg:   productWithZeroNum,
			isErr: false,
		},
		{
			name: "test with internal error on GetProductId()",
			prepareMock: func(m *mocks, p *model.ProductInOrder) {
				var err error = nil
				var internalEerr error = errors.New("Internal error")
				gomock.InOrder(
					m.db.EXPECT().Begin().Return(m.tx, err),
					m.repo.EXPECT().GetProductId(m.tx, p).Return(0, internalEerr),
					m.tx.EXPECT().Rollback(),
				)
			},
			arg:   productStandard,
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOrderRepo := mock_intrface.NewMockOrderRepo(ctrl)
			mockIdb := mock_intrface.NewMockIdb(ctrl)
			mockItx := mock_intrface.NewMockIdb(ctrl)

			mockOrderRepo.EXPECT().GetDb().Return(mockIdb)
			m := &mocks{mockOrderRepo, mockIdb, mockItx}
			tt.prepareMock(m, tt.arg)
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
		prepareMock func(m *mocks, orderId string)
		arg         string
		exp         *Expected
	}{
		{
			name: "test with successful status updating",
			prepareMock: func(m *mocks, orderId string) {
				status := model.CREATED
				var err error = nil
				paymentsSum := float32(25.7)
				productsPricesSum := paymentsSum
				gomock.InOrder(
					m.db.EXPECT().Begin().Return(m.tx, err),
					m.repo.EXPECT().GetOrderStatus(m.tx, orderId).Return(status, err),
					m.repo.EXPECT().GetPaymentsSumByOrderId(m.tx, orderId).Return(paymentsSum, err),
					m.repo.EXPECT().GetProductsPriceSumByOrderId(m.tx, orderId).Return(productsPricesSum, err),
					m.repo.EXPECT().UpdateOrderStatusToComplete(m.tx, orderId).Return(err),
					m.tx.EXPECT().Commit().Return(err),
					m.tx.EXPECT().Rollback(),
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
			prepareMock: func(m *mocks, orderId string) {
				status := model.COMPLITED
				var err error = nil
				gomock.InOrder(
					m.db.EXPECT().Begin().Return(m.tx, err),
					m.repo.EXPECT().GetOrderStatus(m.tx, orderId).Return(status, err),
					m.tx.EXPECT().Rollback(),
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
			prepareMock: func(m *mocks, orderId string) {
				status := model.CREATED
				var err error = nil
				paymentsSum := float32(25.7)
				productsPricesSum := float32(13.9)
				gomock.InOrder(
					m.db.EXPECT().Begin().Return(m.tx, err),
					m.repo.EXPECT().GetOrderStatus(m.tx, orderId).Return(status, err),
					m.repo.EXPECT().GetPaymentsSumByOrderId(m.tx, orderId).Return(paymentsSum, err),
					m.repo.EXPECT().GetProductsPriceSumByOrderId(m.tx, orderId).Return(productsPricesSum, err),
					m.tx.EXPECT().Rollback(),
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
			name: "test with internalError on GetOrderStatus()",
			prepareMock: func(m *mocks, orderId string) {
				var err error = nil
				internalErr := errors.New("Internal error")
				gomock.InOrder(
					m.db.EXPECT().Begin().Return(m.tx, err),
					m.repo.EXPECT().GetOrderStatus(m.tx, orderId).Return(uint8(0), internalErr),
					m.tx.EXPECT().Rollback(),
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
			mockIdb := mock_intrface.NewMockIdb(ctrl)
			mockItx := mock_intrface.NewMockIdb(ctrl)

			mockOrderRepo.EXPECT().GetDb().Return(mockIdb)
			m := &mocks{mockOrderRepo, mockIdb, mockItx}
			tt.prepareMock(m, tt.arg)
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
