package database_test

import (
	"fmt"
	"math/rand"
	"testing"

	"svc/rabbitMq.com/internal/entity"
	"svc/rabbitMq.com/internal/infra/database"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func returnDBInstance() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateProduct(t *testing.T) {
	db, err := returnDBInstance()
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("Product - 1", 10.00)
	assert.NoError(t, err)
	ProductDB := database.NewProduct(db)
	err = ProductDB.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)

}

func TestFindAllProducts(t *testing.T) {
	db, err := returnDBInstance()
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product - %v", i), rand.Float64()*100)
		assert.NoError(t, err)
		db.Create(product)
	}
	productDB := database.NewProduct(db)
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product - 1", products[0].Name)
	assert.Equal(t, "Product - 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product - 11", products[0].Name)
	assert.Equal(t, "Product - 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product - 21", products[0].Name)
	assert.Equal(t, "Product - 23", products[2].Name)

}

func TestProductByID(t *testing.T) {
	db, err := returnDBInstance()
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("Product - 1", 10.00)
	assert.NoError(t, err)
	db.Create(product)
	productDB := database.NewProduct(db)
	product, err = productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product - 1", product.Name)
}

func TestUpdateProduct(t *testing.T) {
	db, err := returnDBInstance()
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("Product - 1", 10.00)
	assert.NoError(t, err)
	db.Create(product)
	productDB := database.NewProduct(db)
	product.Name = "Product - 2"
	err = productDB.Update(product)
	assert.NoError(t, err)
	product, err = productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product - 2", product.Name)
}

func TestDeleteProduct(t *testing.T) {
	db, err := returnDBInstance()
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("Product - 1", 10.00)
	assert.NoError(t, err)
	db.Create(product)
	productDB := database.NewProduct(db)
	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)
	_, err = productDB.FindByID(product.ID.String())
	assert.Error(t, err)

}
