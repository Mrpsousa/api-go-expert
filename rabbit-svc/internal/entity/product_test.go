package entity_test

import (
	"svc/rabbitMq.com/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	p, err := entity.NewProduct("Product-1", 100)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, "Product-1", p.Name)
	assert.Equal(t, 100.00, p.Price)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	p, err := entity.NewProduct("", 10)
	assert.Nil(t, p)
	assert.Equal(t, entity.ErrNameIsRequired, err)

}

func TestProductWhenPriceIsRequired(t *testing.T) {
	p, err := entity.NewProduct("Product-1", 0)
	assert.Nil(t, p)
	assert.Equal(t, entity.ErrPriceIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	p, err := entity.NewProduct("Product-1", -10)
	assert.Nil(t, p)
	assert.Equal(t, entity.ErrInvalidPrice, err)
}

func TestProductValidate(t *testing.T) {
	p, err := entity.NewProduct("Product-1", 10)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Nil(t, p.Validate())
}
