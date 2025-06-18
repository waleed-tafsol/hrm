package main

import (
	"github.com/gin-gonic/gin"
	"hrm/config"
	"hrm/handler"
	"hrm/repository"
	"hrm/usecase"
)

func main() {
	db := config.ConnectDB()
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	userService := usecase.NewUserService(userRepo)
	handler.NewUserHandler(r, userService)

	r.Run(":8080")
}
