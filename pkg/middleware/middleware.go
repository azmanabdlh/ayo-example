package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Provider interface {
	ValidateToken(ctx context.Context, token string) (int64, error)
}

type Authenticate interface {
	IsAccountValid(ctx context.Context, userID int64) bool
}

const RequestContextID = "trace_id"

func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(RequestContextID, requestID)

		ctx := context.WithValue(c.Request.Context(), RequestContextID, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Trace-ID", requestID)

		c.Next()
	}
}

func RequiredAuthentication(provider Provider, auth Authenticate) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value := strings.TrimSpace(
			ctx.GetHeader("Authorization"),
		)

		token := strings.TrimPrefix(
			value,
			"Bearer ",
		)

		if token == "" {
			ctx.JSON(401, gin.H{
				"message": "Required Auth Token",
			})

			ctx.Abort()
			return
		}

		userID, err := provider.ValidateToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.JSON(401, gin.H{
				"message": err.Error(),
			})

			ctx.Abort()
			return
		}

		if !auth.IsAccountValid(ctx, userID) {
			ctx.JSON(401, gin.H{
				"message": "invalid user account",
			})

			ctx.Abort()
		}

		ctx.Set("user_id", userID)

		ctx.Next()
	}
}
