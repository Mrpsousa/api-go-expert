package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mrpsousa/api/internal/dto"
	"github.com/mrpsousa/api/internal/entity"
	"github.com/mrpsousa/api/internal/infra/rabbitmq"
	lg "github.com/mrpsousa/api/pkg/log"
)

const (
	exchange    = "mysqlEx"
	routing_key = "product"
)

type ProductHandler struct {
	Rabbit rabbitmq.RabbitCHInterface
}

func NewProductHandler(rabbit rabbitmq.RabbitCHInterface) *ProductHandler {
	return &ProductHandler{
		Rabbit: rabbit,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Create Product godoc
// @Summary      Create product
// @Description  Create products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateProductInput  true  "product request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /products [post]
// @Security ApiKeyAuth
func (ph *ProductHandler) PublishProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productResp, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productJsonString, err := json.Marshal(productResp)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = ph.Rabbit.Publisher("mysqlEx", "product", productJsonString)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}

	lg.LogMessage("PublishProduct", "AJUSTAR LOG COM CALMA", false)

	failOnError(err, "Failed to publish a message")

	w.WriteHeader(http.StatusOK)
}
