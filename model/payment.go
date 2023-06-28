package model

type Payment struct {
	Total   float32 `json:'total'`
	Change  float32 `json:'change'`
	OrderId string  `json:'orderId, omitempty'`
}
