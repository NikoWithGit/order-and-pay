package model

import "errors"

type Payment struct {
	Total   float32 `json:"total"`
	Change  float32 `json:"change"`
	OrderId string  `json:"-"`
}

func (p *Payment) Validate() error {
	if p.Total < 0 {
		return errors.New("PAYMENT VALIDATION ERROR: total must be positive")
	}
	if p.Change < 0 {
		return errors.New("PAYMENT VALIDATION ERROR: change must be positive")
	}
	return nil
}
