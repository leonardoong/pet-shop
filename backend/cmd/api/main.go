// @title           PetShop API
// @version         1.0
// @description     REST API for PetShop e-commerce platform
// @termsOfService  http://swagger.io/terms/

// @contact.name   PetShop Dev
// @contact.email  dev@petshop.com

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the access token.

package main

import (
	_ "petshop/docs"
	"petshop/internal/server"
	"petshop/pkg/config"
	"petshop/pkg/database"
	jwtpkg "petshop/pkg/jwt"
	customvalidator "petshop/pkg/validator"
)

func main() {
	cfg := config.Load()
	customvalidator.RegisterCustomValidators()

	db := database.NewPostgres(&cfg.DB)
	_ = database.NewRedis(&cfg.Redis) // available for future use (cart caching, rate limiting)

	jwtManager := jwtpkg.NewManager(&cfg.JWT)

	srv := server.New(cfg, db, jwtManager)
	srv.Run()
}
