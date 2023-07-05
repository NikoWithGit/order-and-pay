package main

import (
	"database/sql"
	"order-and-pay/controller"
	"order-and-pay/env"
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

	err = godotenv.Load(".env")

	if err != nil {
		zaplogger.Panic(err.Error())
	}

	dsn := env.GetDbDsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		zaplogger.Panic(err.Error())
	}
	defer db.Close()

	orderRepoImpl := repoimpl.NewOrderRepoImpl(db, zaplogger)
	orderService := service.NewOrderService(orderRepoImpl, zaplogger)
	orderController := controller.NewOrderController(orderService, zaplogger)

	s := server.NewServer()
	s.RegisterRoutes(orderController)
	if err = s.Start(); err != nil {
		zaplogger.Panic(err.Error())
	}
}
