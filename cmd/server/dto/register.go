package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterSuccessResponse struct {
	Msg string `json:"msg" example:"user created"`
}
