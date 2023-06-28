package controller

import (
	"net/http"
	"order-and-pay/intrface"
	"order-and-pay/model"
	"order-and-pay/service"
	"time"

	"github.com/gin-gonic/gin"
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
		Id    string `json:'id'`
		Short uint   `json:'short'`
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
	id := ctx.Param("id")
	order := oc.service.Get(id)
	ctx.JSON(http.StatusOK, &order)
}

func (oc *OrderController) AddProduct(ctx *gin.Context) {
	var product model.ProductInOrder
	id := ctx.Param("id")
	product.OrderId = id
	bindErr := ctx.BindJSON(&product)
	if bindErr != nil {
		panic(bindErr)
	}
	oc.service.AddProduct(&product)
}

func (oc *OrderController) AddPayment(ctx *gin.Context) {
	var payment model.Payment
	id := ctx.Param("id")
	payment.OrderId = id
	bindErr := ctx.BindJSON(&payment)
	if bindErr != nil {
		panic(bindErr)
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
