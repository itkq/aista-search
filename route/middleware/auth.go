package middleware

import (
	"aista-search/db"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.Query("token")

		if t == "" {
			authError(c)
			return
		}

		if db.IsValidToken(t) {
			c.Next()
		} else {
			authError(c)
		}
	}
}

func authError(c *gin.Context) {
	c.JSON(401, gin.H{"status": "bad", "msg": "token is invalid"})
	c.Abort()
}
