package server

import "order-and-pay/controller"

func (s *server) RegisterRoutes(oc *controller.OrderController) {
	s.Gin.GET("/orders", oc.GetAll)
	s.Gin.GET("/orders/:id", oc.Get)
	s.Gin.POST("/orders/create", oc.Create)
	s.Gin.PUT("/orders/:id/add-product", oc.AddProduct)
	s.Gin.PUT("/orders/:id/add-payment", oc.AddPayment)
	s.Gin.PUT("/orders/:id/finish", oc.Finish)
}
