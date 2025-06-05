package domain

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenClaims struct {
	UserID    string
	Role      Role
	Type      TokenType
	JTI       string
	IssuedAt  int64
	ExpiresAt int64
}
