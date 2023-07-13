package main

import (
	"order-and-pay/controller"
	"order-and-pay/db"
	"order-and-pay/logger"
	repoimpl "order-and-pay/repo-impl"
	"order-and-pay/server"
	"order-and-pay/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	zaplogger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}

	if err = godotenv.Load(".env"); err != nil {
		zaplogger.Panic(err.Error())
	}

	db, err := db.NewSqlDb()
	if err != nil {
		zaplogger.Panic(err.Error())
	}
	defer db.Close()

	orderRepoImpl := repoimpl.NewOrderRepoImpl(db)
	orderService := service.NewOrderService(orderRepoImpl)
	orderController := controller.NewOrderController(orderService, zaplogger)

	s := server.NewServer()
	s.RegisterRoutes(orderController)
	err = s.Start()
	if err != nil {
		zaplogger.Panic(err.Error())
	}
}
