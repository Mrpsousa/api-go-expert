package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/mrpsousa/api/configs"
	_ "github.com/mrpsousa/api/docs"
	"github.com/mrpsousa/api/internal/entity"
	"github.com/mrpsousa/api/internal/infra/database"
	rb "github.com/mrpsousa/api/internal/infra/rabbitmq"
	"github.com/mrpsousa/api/internal/infra/webserver/handlers"
	"github.com/mrpsousa/api/pkg/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rmq_adress = "amqp://guest:guest@localhost"
)

// @title           Go Expert API Example
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Marcos Rogerio
// @contact.url    http://www.example.com.br
// @contact.email  urameshi.uba@gmail.com

// @license.name   Full Cycle License
// @license.url    http://www.fullcycle.com.br

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func PrepareAmqp() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func main() {
	conn, ch, err := PrepareAmqp()
	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()
	defer ch.Close()

	config := configs.NewConfig()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rabbit := rb.NewRabbit(ch)
	db.AutoMigrate(&entity.User{})
	productHandler := handlers.NewProductHandler(rabbit)
	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, config.TokenAuth, config.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger) // look the documentation to more diferent chi middleware
	r.Use(middleware.Recoverer)
	r.Use(middlewares.LogRequest)

	r.Route("/products", func(r chi.Router) {
		// 	r.Use(jwtauth.Verifier(config.TokenAuth)) // get the token and inject it into the context
		// 	r.Use(jwtauth.Authenticator)              // validate of token
		r.Post("/", productHandler.PublishProduct)

	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	// r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://137.184.246.121/docs/doc.json")))
	r.Get("/ping", handlers.Healthz)
	http.ListenAndServe(":8000", r)
}
