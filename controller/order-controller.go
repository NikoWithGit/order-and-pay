package controller

import (
	"net/http"
	"order-and-pay/intrface"
	"order-and-pay/model"
	"order-and-pay/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderController struct {
	service intrface.OrderService
	logger  intrface.Ilogger
}

func NewOrderController(os *service.OrderService, l intrface.Ilogger) *OrderController {
	return &OrderController{os, l}
}

func (oc *OrderController) Create(ctx *gin.Context) {
	id, short, err := oc.service.Create()
	if err != nil {
		oc.logger.Error(err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	type Response struct {
		Id    string `json:"id"`
		Short uint   `json:"short"`
	}
	ctx.JSON(http.StatusOK, &Response{id, short})
}

func (oc *OrderController) GetAll(ctx *gin.Context) {
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	if fromStr == "" || toStr == "" {
		ctx.String(http.StatusBadRequest, "'from' and 'to' parameters are required")
		return
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong 'from' value")
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong 'to' value")
		return
	}
	orders, err := oc.service.GetAll(from, to)
	if err != nil {
		oc.logger.Error(err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &orders)
}

func (oc *OrderController) Get(ctx *gin.Context) {
	orderId := ctx.Param("id")
	_, err := uuid.Parse(orderId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong order_id value")
		return
	}

	order, err := oc.service.Get(orderId)
	if err != nil {
		oc.logger.Error(err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if order == nil {
		ctx.String(http.StatusBadRequest, "Order with id "+orderId+" doesn't exists")
		return
	}

	ctx.JSON(http.StatusOK, &order)
}

func (oc *OrderController) AddProduct(ctx *gin.Context) {
	orderId := ctx.Param("id")
	_, err := uuid.Parse(orderId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong order_id value")
		return
	}

	//Order existence check?

	var product model.ProductInOrder
	err = ctx.BindJSON(&product)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong product's json data")
		return
	}

	errValidate := product.Validate()
	if errValidate != nil {
		ctx.String(http.StatusBadRequest, errValidate.Error())
		return
	}
	product.OrderId = orderId
	err = oc.service.AddProduct(&product)
	if err != nil {
		oc.logger.Error(err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, "Success!")
}

func (oc *OrderController) AddPayment(ctx *gin.Context) {
	orderId := ctx.Param("id")
	_, err := uuid.Parse(orderId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong order_id value")
		return
	}

	var payment model.Payment
	err = ctx.BindJSON(&payment)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong product's json data")
		return
	}

	//Order existence check?

	errValidate := payment.Validate()
	if errValidate != nil {
		ctx.String(http.StatusBadRequest, errValidate.Error())
		return
	}
	payment.OrderId = orderId
	err = oc.service.AddPayment(&payment)
	if err != nil {
		oc.logger.Error(err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, "Success!")
}

func (oc *OrderController) Finish(ctx *gin.Context) {
	orderId := ctx.Param("id")
	_, err := uuid.Parse(orderId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong order_id value")
		return
	}

	res, badRequestErr, err := oc.service.Finish(orderId)
	if err != nil {
		oc.logger.Error(err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	if badRequestErr != nil {
		ctx.String(http.StatusBadRequest, badRequestErr.Error())
		return
	}
	if !res {
		ctx.String(http.StatusOK, "Transaction has already been completed")
		return
	}
	ctx.String(http.StatusOK, "Transaction has been successfuly completed!")
}
