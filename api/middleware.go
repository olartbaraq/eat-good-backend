package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized request",
			})
			ctx.Abort()
			return
		}

		tokenSplit := strings.Split(token, " ")

		if len(tokenSplit) != 2 && strings.ToLower(tokenSplit[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token format",
			})
			ctx.Abort()
			return
		}

		userId, role, err := tokenManager.VerifyToken(tokenSplit[1])

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error":  err.Error(),
				"status": "failed to verify token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("id", userId)
		ctx.Set("role", role)

	}
}
