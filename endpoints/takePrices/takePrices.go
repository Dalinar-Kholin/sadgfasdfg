package takePrices

import (
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/user"
	"sync"
)

type TakePrices struct {
}

func (t *TakePrices) TakePrice(c *gin.Context) {
	cookie, err := c.Cookie("accessToken")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "where Token?",
		})
		return
	} // po testach do zmiany
	userInstance := constAndVars.Users[cookie]
	ean := c.Query("ean")
	var wg sync.WaitGroup
	ch := make(chan interface{})
	for _, hurt := range userInstance.Hurts {
		wg.Add(1)
		go func(hurt *hurtownie.IHurt, wg *sync.WaitGroup, ch chan<- interface{}) {
			defer wg.Done()
			res, err := (*hurt).SearchProduct(ean, userInstance.Client)
			if err != nil && err.Error() == "tokenError" {
				var creds user.UserCreds
				for _, i := range userInstance.Creds {
					if i.HurtName == (*hurt).GetName() {
						creds = i
					}
				}
				(*hurt).TakeToken(creds.Login, creds.Password, userInstance.Client) // tutaj zamienić na refreshToken, nie ma potrzeby ponownie pobierać tokena
				res, err = (*hurt).SearchProduct(ean, userInstance.Client)
			}
			if err != nil {
				ch <- SearchResult{Ean: ean, Result: nil, HurtName: (*hurt).GetName()}
				return
			}
			ch <- SearchResult{Ean: ean, Result: res, HurtName: (*hurt).GetName()}
		}(&hurt, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	result := make([]interface{}, len(userInstance.Hurts))
	i := 0

	for x := range ch {
		result[i] = x
		i++
	}
	c.JSON(200, result)

}

type SearchResult struct {
	Ean      string             `json:"ean"`
	HurtName hurtownie.HurtName `json:"hurtName"`
	Result   interface{}        `json:"result"`
}
