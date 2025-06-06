package handlers

import (
	"net/http"

	httpHelpers "github.com/funchooooza-ossh/protego/cmd/server/httpHelpers"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/usecases"
	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context, u usecases.LogoutUsecaseInterface) {
	access, err := httpHelpers.ParseCookie(c, "access_token")
	if err != nil {
		return
	}
	ctx := c.Request.Context()

	if err = u.Execute(ctx, access); err != nil {
		code, msg := e.ToHTTPResponse(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
