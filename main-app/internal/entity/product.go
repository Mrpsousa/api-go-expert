package entity

import (
	"errors"
	"time"

	"github.com/mrpsousa/api/pkg/entity"
)

var (
	ErrIDIsRequired    = errors.New("id_is_required")
	ErrInvalID         = errors.New("invalid_id")
	ErrNameIsRequired  = errors.New("name_is_required")
	ErrPriceIsRequired = errors.New("price_is_required")
	ErrInvalidPrice    = errors.New("invalid_price")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}
	err := product.Validate()
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalID
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0 {
		return ErrPriceIsRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
