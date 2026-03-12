// @title Arena Ban API
// @version 1.0
// @description This is the API documentation for Arena Ban Backend.
// @host localhost:3000
// @BasePath /api/v1
package main

import (
	"arena-ban/config"
	server "arena-ban/internal/delivery/http"
	"arena-ban/internal/delivery/http/handler"
	"arena-ban/internal/repository"
	"arena-ban/internal/usecase"
	util "arena-ban/pkg"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger := config.InitLogger()
	db := config.ConnectDB()

	redis, err := config.SetupRedis()
	if err != nil {
		panic(err)
	}
	config.RunMigrations(db)
	util.InitJwt()

	util.NewSMTP()

	userRepo := repository.NewUserRepository(db, redis)

	authUsecase := usecase.NewAuthUsecase(*&userRepo, logger)

	authHandler := handler.NewAuthHandler(authUsecase, logger)

	server.SetupServer(authHandler, logger)

}
