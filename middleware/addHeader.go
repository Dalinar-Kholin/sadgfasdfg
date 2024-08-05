package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

func AddHeaders(c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Session-Id")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	c.Next()
	return
	c.Writer.Header().Set("Cache-Control", "no-store")
	c.Writer.Header().Set("Date", time.Now().String())
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-Frame-Options", "DENY")
	c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
	c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	c.Writer.Header().Set("Accept-Encoding", "application/json")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}
