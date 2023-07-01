package main

import (
	"database/sql"
	"order-and-pay/controller"
	"order-and-pay/env"
	repoimpl "order-and-pay/repo-impl"
	"order-and-pay/service"
	"order-and-pay/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	utils.InitLogger()
	utils.Logger.Sync()

	err := godotenv.Load(".env")

	if err != nil {
		utils.Logger.Panic(err.Error())
	}

	dsn := env.GetDbDsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer db.Close()

	orderRepoImpl := repoimpl.NewOrderRepoImpl(db)
	orderService := service.NewOrderService(orderRepoImpl)
	orderController := controller.NewOrderController(orderService)

	r := gin.Default()
	r.GET("/orders", orderController.GetAll)
	r.GET("/orders/:id", orderController.Get)
	r.POST("/orders/create", orderController.Create)
	r.PUT("/orders/:id/add-product", orderController.AddProduct)
	r.PUT("/orders/:id/add-payment", orderController.AddPayment)
	r.PUT("/orders/:id/finish", orderController.Finish)
	r.Run()
}
