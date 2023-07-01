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
