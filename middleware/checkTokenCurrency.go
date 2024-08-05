package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/user"
	"sync"
)

func CheckTokenCurrency(c *gin.Context) {
	token, err := c.Cookie("accessToken")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "where Token?",
		})
		return
	}
	var ok bool
	var userInstance user.User
	if userInstance, ok = constAndVars.Users[token]; !ok {
		c.JSON(400, gin.H{
			"error": "where logowanie?",
		})
	}

	var wg sync.WaitGroup

	for _, hurt := range userInstance.Hurts {
		wg.Add(1)
		go func(wg *sync.WaitGroup, hurt hurtownie.IHurt) {
			defer wg.Done()
			if !hurt.CheckToken(userInstance.Client) {
				fmt.Printf("check Token się nie udał := %v\n", hurt.GetName())
				if !hurt.RefreshToken(userInstance.Client) {
					fmt.Printf("refresh tokena hurt := %v\n", hurt.GetName())
					userCred := userInstance.TakeHurtCreds(hurt.GetName())
					hurt.TakeToken(userCred.Login, userCred.Password, userInstance.Client)
				}
			}
		}(&wg, hurt)
	}
	wg.Wait()
	c.Next()
}
