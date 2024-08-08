package takePrices

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/user"
	"sync"
)

func (t *TakePrices) TakeMultiple(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
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

	var list hurtownie.WishList
	responseReaderJson := json.NewDecoder(c.Request.Body)
	err := responseReaderJson.Decode(&list)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "bad list",
		})
		return
	}

	var wg sync.WaitGroup
	ch := make(chan Result)

	for _, hurt := range userInstance.Hurts {
		wg.Add(1)
		go func(hurt *hurtownie.IHurt, wg *sync.WaitGroup, ch chan<- Result) {
			defer wg.Done()
			res, err := (*hurt).SearchMany(list, userInstance.Client)
			if err != nil {
				ch <- Result{
					HurtName: (*hurt).GetName(),
					Result:   nil,
				}
				return
			}
			ch <- Result{
				HurtName: (*hurt).GetName(),
				Result:   res}
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

type Result struct {
	HurtName hurtownie.HurtName
	Result   []hurtownie.SearchManyProducts
}
