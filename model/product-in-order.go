package model

import (
	"errors"

	"github.com/google/uuid"
)

type ProductInOrder struct {
	Uuid        string  `json:"uuid"`
	Num         uint    `json:"num"`
	PricePerOne float32 `json:"pricePerOne"`
	OrderId     string  `json:"-"`
}

func (p *ProductInOrder) Validate() error {
	if _, err := uuid.Parse(p.Uuid); err != nil {
		return errors.New("PRODUCT VALIDATION ERROR: wrong product id")
	}
	if p.PricePerOne < 0 {
		return errors.New("PRODUCT VALIDATION ERROR: pricePerOne must be positive")
	}
	return nil
}
