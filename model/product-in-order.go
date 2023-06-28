package model

type ProductInOrder struct {
	Uuid        string  `json:'uuid'`
	Num         uint    `json:'num'`
	PricePerOne float32 `json:'pricePerOne'`
	OrderId     string  `json:'order_id, omitempty'`
}
