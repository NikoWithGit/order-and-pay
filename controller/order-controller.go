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
}

func NewOrderController(os *service.OrderService) *OrderController {
	return &OrderController{os}
}

func (oc *OrderController) Create(ctx *gin.Context) {
	id, short := oc.service.Create()
	type Response struct {
		Id    string `json:"id"`
		Short uint   `json:"short"`
	}
	ctx.JSON(http.StatusOK, &Response{id, short})
}

func (oc *OrderController) GetAll(ctx *gin.Context) {
	from, _ := time.Parse("2006-01-02", ctx.Query("from"))
	to, _ := time.Parse("2006-01-02", ctx.Query("to"))
	orders := oc.service.GetAll(from, to)
	ctx.JSON(http.StatusOK, &orders)
}

func (oc *OrderController) Get(ctx *gin.Context) {
	orderId := ctx.Param("id")
	_, err := uuid.Parse(orderId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong id")
		return
	}

	order := oc.service.Get(orderId)
	if order == nil {
		ctx.String(http.StatusBadRequest, "Order with id "+orderId+"doesn't exists")
		return
	}

	ctx.JSON(http.StatusOK, &order)
}

func (oc *OrderController) AddProduct(ctx *gin.Context) {
	orderId := ctx.Param("id")
	_, err := uuid.Parse(orderId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong id")
		return
	}

	var product model.ProductInOrder
	err = ctx.BindJSON(&product)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong product's json data")
		return
	}

	errValidate := product.Validate()
	if err != nil {
		ctx.String(http.StatusBadRequest, errValidate.Error())
		return
	}

	oc.service.AddProduct(&product)
}

func (oc *OrderController) AddPayment(ctx *gin.Context) {
	orderId := ctx.Param("id")
	_, err := uuid.Parse(orderId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong id")
		return
	}

	var payment model.Payment
	err = ctx.BindJSON(&payment)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Wrong product's json data")
		return
	}

	errValidate := payment.Validate()
	if err != nil {
		ctx.String(http.StatusBadRequest, errValidate.Error())
		return
	}

	oc.service.AddPayment(&payment)
}

func (oc *OrderController) Finish(ctx *gin.Context) {
	orderId := ctx.Param("id")
	err := oc.service.Finish(orderId)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
}
