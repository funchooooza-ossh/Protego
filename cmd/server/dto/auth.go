package dto

type AuthRequest struct {
	Resource string `header:"X-Resource-Code" binding:"required"`
	Method   string `header:"X-Forwarded-Method" binding:"required"`
}

type AuthCookies struct {
	AccessToken  string
	RefreshToken string
}
