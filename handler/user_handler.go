package handler

import (
	"github.com/gin-gonic/gin"
	"hrm/domain"
	"net/http"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(router *gin.Engine, userService domain.UserService) {
	handler := &UserHandler{userService: userService}
	router.POST("/signup", handler.SignUp)
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (userHandler *UserHandler) SignUp(context *gin.Context) {
	var req SignUpRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := userHandler.userService.SignUp(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
