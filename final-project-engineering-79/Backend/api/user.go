package api

import (
	"fmt"
	"net/http"
	"usedbooks/backend/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	userService repository.Repository
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nohp     string `json:"nohp"`
	Password string `json:"password"`
}

func NewUserHandler(userRepository repository.Repository) *userHandler {
	return &userHandler{userRepository}
}

func (h *userHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.FetchUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})
		return
	}
	var usersResponse []UserResponse

	for _, v := range users {
		userResponse := convertToUserResponse(v)

		usersResponse = append(usersResponse, userResponse)
	}
	c.JSON(http.StatusOK, gin.H{
		"data": usersResponse,
	})
}

func (h *userHandler) PostUserRegist(c *gin.Context) {
	var userRequest repository.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error on field %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errorMessages,
		})
		return
	}
	res, err := h.userService.InsertUser(userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

func convertToUserResponse(h repository.User) UserResponse {
	return UserResponse{
		ID:       h.ID,
		Username: h.Username,
		Email:    h.Email,
		Nohp:     h.Nohp,
		Password: h.Password,
	}
}
