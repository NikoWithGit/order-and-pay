package model

import "time"

const (
	CREATED   uint8 = 1
	COMPLITED uint8 = 2
)

type Order struct {
	Id       string           `json:"id"`
	Short    uint             `json:"short"`
	Products []ProductInOrder `json:"products"`
	Payments []Payment        `json:"payment"`
	Date     time.Time        `json:"date"`
	Status   string           `json:"status"`
}

type OrderReport struct {
	Id       string           `json:"id"`
	Total    float32          `json:"total"`
	Products []ProductInOrder `json:"products"`
	Date     time.Time        `json:"date"`
}

func (o *Order) GetOrderReport() *OrderReport {
	return &OrderReport{o.Id, o.getTotal(), o.Products, o.Date}
}

func (o *Order) getTotal() float32 {
	var result float32

	for i := range o.Payments {
		result += o.Payments[i].Total
		result -= o.Payments[i].Change
	}
	return result
}
