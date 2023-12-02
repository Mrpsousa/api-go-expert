package main

import (
	"log"

	"svc/rabbitMq.com/internal/entity"
	"svc/rabbitMq.com/internal/infra/database"
	rb "svc/rabbitMq.com/internal/infra/rabbitmq"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	address    = "amqp://guest:guest@localhost"
	exchange   = "mysqlEx"
	routingKey = "product"
	queueName  = "mysql"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{})
	productDB := database.NewProduct(db)
	rabbitmq, err := rb.NewRabbitMq(address, productDB)
	if err != nil {
		log.Fatal(err)
	}

	rabbitmq.Consumer(exchange, queueName, routingKey)
}
