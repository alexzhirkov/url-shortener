package middleware

import "github.com/gin-gonic/gin"

func EnsureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Authorization")
		if authToken != "Bearer "+"XBp41KwWuhqJP2pD" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization required"})
			return
		}
		c.Next()
	}
}
