package handlers

import (
	"net/http"

	"github.com/funchooooza-ossh/protego/cmd/server/dto"
	httpHelpers "github.com/funchooooza-ossh/protego/cmd/server/httpHelpers"
	"github.com/funchooooza-ossh/protego/internal/config"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/usecases"
	"github.com/gin-gonic/gin"
)

// AuthForward godoc
// @Summary ForwardAuth. Main endpoint to validate user permissions
// @Description Forward auth для любого зарегестрированного ресурса
// @Tags auth
// @Accept json
// @Produce json
// @Security AccessToken
// @Param X-Resource-Code header string true "Запрашиваемый ресурс"
// @Param X-Forwarded-Method header string true "Запрашиваемый метод"
// @Failure 404 {object} dto.ErrorResponse "resource not found"
// @Failure 422 {object} dto.ErrorResponse "invalid input"
// @Failure 401 {object} dto.ErrorResponse "unauthenticated"
// @Failure 403 {object} dto.ErrorResponse "access denied"
// @Failure 409 {object} dto.ErrorResponse "conflict"
// @Failure 429 {object} dto.ErrorResponse "too many requests"
// @Failure 401 {object} dto.ErrorResponse "token expired"
// @Failure 401 {object} dto.ErrorResponse "token invalid"
// @Failure 504 {object} dto.ErrorResponse "request timed out"
// @Failure 503 {object} dto.ErrorResponse "service unavailable"
// @Failure 500 {object} dto.ErrorResponse "internal error"
// @Failure 500 {object} dto.ErrorResponse "unexpected internal error"
// @Router /auth/forward [get]
func authHandler(c *gin.Context, u usecases.AuthUsecaseInterface, cfg *config.Config) {
	var req dto.AuthRequest
	if err := c.ShouldBindHeader(&req); err != nil { //проверяем наличие необходимых заголовков
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "no requested path in headers"})
		return
	}
	tokens, err := httpHelpers.ParseAuthCookies(c) // проверяем наличие токенов
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "no auth credentials provided"})
		return
	}

	action := httpHelpers.ExtractParams(req) // извлекаем из заголовков метод запроса и маппим его к внутренней логике
	ctx := c.Request.Context()

	allowed, newAccess, err := u.Execute(ctx, req.Resource, action, tokens.AccessToken, tokens.RefreshToken) // проверяем доступ пользователя к ресурсу
	// если доступа нет, мы все равно можем вернуть ему новый токен, хоть и не пропустить его дальше

	if err != nil {
		code, msg := e.ToHTTPResponse(err)
		c.JSON(code, dto.ErrorResponse{Error: msg})
		return
	}

	if newAccess != "" && newAccess != tokens.AccessToken {
		httpHelpers.SetCookie(
			c, "access_token", newAccess, int(cfg.AccessTtl.Seconds()),
		) //ставим новый токен в куку
	}

	if !allowed {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access forbidden"})
		return // 403 если доступа нет
	}

	c.Status(http.StatusOK) // иначе всегда 200
}

func MakeAuthHandler(u usecases.AuthUsecaseInterface, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHandler(c, u, cfg)
	}
}
