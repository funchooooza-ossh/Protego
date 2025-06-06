package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginSuccessResponse struct {
	Msg string `json:"msg" example:"Welcome back, User"`
}
