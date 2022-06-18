package api

import (
	"net/http"
	"time"
	"usedbooks/backend/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type authHandler struct {
	userRepo repository.Repository
}

func NewAuthHandler(usersRepository repository.Repository) *authHandler {
	return &authHandler{usersRepository}
}

type User struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Email string `json:"email`
	Token string `json:"token"`
}

type AuthErrorResponse struct {
	Error string `json:"error"`
}

//jwt key fo signature
var jwtKey = []byte("secret")

type Claims struct {
	Email string
	Role  string
	jwt.StandardClaims
}

func (h *authHandler) LoginUser(c *gin.Context) {

	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.userRepo.LoginUser(user.Email, user.Password)
	c.Header("Content-Type", "application/json")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Email: *res,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	//encode claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})
	c.JSON(http.StatusOK, gin.H{
		"data": LoginResponse{
			Email: *res,
			Token: tokenString,
		},
	})
}
