package middleware

import (
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
)

func CheckToken(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")

	if token == "" {
		c.JSON(400, gin.H{
			"error": "where Token?",
		})
		return
	}
	var ok bool
	if _, ok = constAndVars.Users[token]; !ok {
		c.JSON(400, gin.H{
			"error": "where logowanie?",
		})
		return
	}
	c.Next()
} // globalna mapa mapująca TOKEN na usera
