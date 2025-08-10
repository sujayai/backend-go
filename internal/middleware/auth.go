package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuth(adminKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("x-admin-key")
		if key != adminKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, x-admin-key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
