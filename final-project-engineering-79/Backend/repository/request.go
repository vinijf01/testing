package repository

//input dari user

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Nohp     string `json:"nohp" binding:"required"`
	Password string `json:"password" binding:"required"`
}
