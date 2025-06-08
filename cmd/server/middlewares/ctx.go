package middlewares

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderRequestID = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(HeaderRequestID)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("request_id", requestID)

		ctx := context.WithValue(c.Request.Context(), logger.CtxKeyRequestID{}, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set(HeaderRequestID, requestID)

		c.Next()
	}
}
