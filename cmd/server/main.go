package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"

	"github.com/mrpsousa/api/configs"
	_ "github.com/mrpsousa/api/docs"
	"github.com/mrpsousa/api/internal/entity"
	"github.com/mrpsousa/api/internal/infra/database"
	"github.com/mrpsousa/api/internal/infra/webserver/handlers"
	"github.com/mrpsousa/api/pkg/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// println(config.DBDriver)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)
	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, config.TokenAuth, config.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger) // look the documentation to more diferent chi middleware
	r.Use(middleware.Recoverer)
	r.Use(middlewares.LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth)) // get the token and inject it into the context
		r.Use(jwtauth.Authenticator)              // validate of token
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}
