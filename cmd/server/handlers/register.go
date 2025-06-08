package handlers

import (
	"net/http"

	"github.com/funchooooza-ossh/protego/cmd/server/dto"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register new user
// @Desription Регистрация нового пользователя с базовыми правами.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Тело запроса"
// @Success 201 {object} dto.RegisterSuccessResponse "user created"
// @Failure 409 {object} dto.ErrorResponse "resource already exists"
// @Failure 422 {object} dto.ErrorResponse "invalid input"
// @Failure 409 {object} dto.ErrorResponse "conflict"
// @Failure 429 {object} dto.ErrorResponse "too many requests"
// @Failure 504 {object} dto.ErrorResponse "request timed out"
// @Failure 503 {object} dto.ErrorResponse "service unavailable"
// @Failure 500 {object} dto.ErrorResponse "internal error"
// @Failure 500 {object} dto.ErrorResponse "unexpected internal error"
// @Router /register [post]
func registerHandler(c *gin.Context, u contracts.RegisterUsecaseInterface) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid data"})
		return
	}

	ctx := c.Request.Context()
	_, err := u.Execute(ctx, req.Email, req.Password)

	if err != nil {
		code, msg := e.ToHTTPResponse(err)
		c.JSON(code, dto.ErrorResponse{Error: msg})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterSuccessResponse{Msg: "user created"})
}

func MakeRegisterHandler(u contracts.RegisterUsecaseInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		registerHandler(c, u)
	}
}
