package main

import (
	"order-and-pay/controller"
	"order-and-pay/db"
	"order-and-pay/logger"
	"order-and-pay/producer"
	repoimpl "order-and-pay/repo-impl"
	"order-and-pay/server"
	"order-and-pay/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//Logger init
	zaplogger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}

	//Env load
	if err = godotenv.Load(".env"); err != nil {
		zaplogger.Panic(err.Error())
	}

	//DB connection
	db, err := db.NewSqlDb()
	if err != nil {
		zaplogger.Panic(err.Error())
	}
	defer db.Close()

	//Producer init
	producer, err := producer.NewKafkaProducer([]string{"order-and-pay-kafka-1:9092"}, zaplogger)
	if err != nil {
		zaplogger.Error(err.Error())
	}
	defer producer.Close()

	orderRepoImpl := repoimpl.NewOrderRepoImpl(db)
	orderService := service.NewOrderService(orderRepoImpl)
	orderController := controller.NewOrderController(orderService, producer, zaplogger)

	s := server.NewServer()
	s.RegisterRoutes(orderController)
	err = s.Start()
	if err != nil {
		zaplogger.Panic(err.Error())
	}
}
