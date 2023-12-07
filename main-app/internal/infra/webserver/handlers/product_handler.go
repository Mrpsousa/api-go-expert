package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mrpsousa/api/internal/dto"
	"github.com/mrpsousa/api/internal/entity"
	rb "github.com/mrpsousa/api/internal/infra/rabbitmq"
)

const (
	exchange    = "mysqlEx"
	routing_key = "product"
)

type ProductHandler struct {
	RabbitMq *rb.RabbitMq
}

func NewProductHandler(rmq_adress string) *ProductHandler {
	conn, err := rb.NewRabbitMq(rmq_adress)
	if err != nil {
		log.Fatal(err)
	}
	return &ProductHandler{
		RabbitMq: conn,
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
func (h *ProductHandler) PublishProduct(w http.ResponseWriter, r *http.Request) {
	// type Person struct {
	// 	Name  string `json:"name"`
	// 	Age   int    `json:"age"`
	// 	City  string `json:"city"`
	// 	Email string `json:"email,omitempty"`
	// }

	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonString, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
		return
	}

	//TODO:_ this func must to return err
	h.RabbitMq.Publisher(exchange, routing_key, string(jsonString))
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
}
