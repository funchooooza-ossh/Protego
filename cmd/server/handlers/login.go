package handlers

import (
	"fmt"
	"net/http"

	"github.com/funchooooza-ossh/protego/cmd/server/dto"
	httpHelpers "github.com/funchooooza-ossh/protego/cmd/server/httpHelpers"
	"github.com/funchooooza-ossh/protego/internal/config"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/usecases"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary Авторизация
// @Description Авторизация, возвращает куки с токенами
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Тело запроса"
// @Success 200 {object} dto.LoginSuccessResponse "welcome back, user"
// @Failure 404 {object} dto.ErrorResponse "resource not found"
// @Failure 422 {object} dto.ErrorResponse "invalid input"
// @Failure 409 {object} dto.ErrorResponse "conflict"
// @Failure 429 {object} dto.ErrorResponse "too many requests"
// @Failure 504 {object} dto.ErrorResponse "request timed out"
// @Failure 503 {object} dto.ErrorResponse "service unavailable"
// @Failure 500 {object} dto.ErrorResponse "internal error"
// @Failure 500 {object} dto.ErrorResponse "unexpected internal error"
// @Router /login [post]
func loginHandler(c *gin.Context, u usecases.LoginUsecaseInterface, cfg *config.Config) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	ctx := c.Request.Context()
	access, refresh, err := u.Execute(ctx, req.Email, req.Password)
	if err != nil {
		code, msg := e.ToHTTPResponse(err)
		c.JSON(code, dto.ErrorResponse{Error: msg})
		return
	}

	tokens := &dto.AuthCookies{
		AccessToken:  access,
		RefreshToken: refresh,
	}
	httpHelpers.SetAuthCookies(c, tokens, cfg)
	c.JSON(http.StatusOK, dto.LoginSuccessResponse{Msg: fmt.Sprintf("welcome back, %s", req.Email)})
}

func MakeLoginHandler(u usecases.LoginUsecaseInterface, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		loginHandler(c, u, cfg)
	}
}
