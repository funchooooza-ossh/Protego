package httpHelpers

import (
	"fmt"
	"net/http"

	"github.com/funchooooza-ossh/protego/cmd/server/dto"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/gin-gonic/gin"
)

func ParseAuthCookies(c *gin.Context) (*dto.AuthCookies, error) {
	var accessToken, refreshToken string

	if val, err := c.Cookie("access_token"); err == nil {
		accessToken = val
	}

	if val, err := c.Cookie("refresh_token"); err == nil {
		refreshToken = val
	}

	// если обе пусты — ошибка
	if accessToken == "" && refreshToken == "" { // для реализации логики доступа нам необходим хотя бы один токен
		return nil, fmt.Errorf("no valid auth tokens found")
	}

	return &dto.AuthCookies{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ParseCookie(c *gin.Context, name string) (string, error) {
	value, err := c.Cookie(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request params"})
		return "", err
	}
	return value, nil
}

func SetAuthCookies(c *gin.Context, tokens *dto.AuthCookies, cfg *config.Config) {
	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		int(cfg.AccessTtl.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		int(cfg.RefreshTtl.Seconds()),
		"/",
		"",
		true,
		true,
	)
}

func SetCookie(c *gin.Context, name string, value string, ttl int) {
	c.SetCookie(
		name,
		value,
		ttl,
		"/",
		"",
		true,
		true,
	)
}
