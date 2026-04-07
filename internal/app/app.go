package app

import (
	"fmt"
	docs "website-dm/api/openapi"
	"website-dm/internal/config"
	"website-dm/internal/handler"
	"website-dm/internal/middleware"
	"website-dm/internal/repository"
	"website-dm/internal/service"
	"website-dm/internal/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	db, err := storage.Open(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repository.New(db)
	svc := service.New(repo, cfg.JWT.Secret)
	h := handler.New(svc)

	r := gin.Default()
	r.Use(middleware.CORS())

	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Title = "Pet & Plant Store API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "REST API для магазину товарів для тварин та рослин."
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.Port)
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/auth/register", h.Register)
	r.POST("/api/v1/auth/login", h.Login)
	r.GET("/api/v1/categories", h.Categories)
	r.GET("/api/v1/products", h.Products)
	r.GET("/api/v1/products/:id", h.ProductByID)

	authorized := r.Group("/api/v1")
	authorized.Use(middleware.Auth(svc))
	{
		authorized.POST("/cart/items", h.AddToCart)
		authorized.GET("/cart", h.Cart)
	}

	return r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
