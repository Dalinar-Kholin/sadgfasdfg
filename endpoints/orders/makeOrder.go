package orders

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/user"
	"sync"
)

type Order struct {
}

func (o *Order) MakeOrder(c *gin.Context) {
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
	fmt.Printf("list := %v\n", list)
	var wg sync.WaitGroup

	ch := make(chan Result)

	for _, hurt := range userInstance.Hurts {
		wg.Add(1)
		go func(hurt *hurtownie.IHurt, wg *sync.WaitGroup, ch chan<- Result) {
			defer wg.Done()
			ch <- Result{
				Name:   (*hurt).GetName(),
				Status: (*hurt).AddToCart(list, userInstance.Client)}
		}(&hurt, &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := make([]Result, len(userInstance.Hurts))
	i := 0

	for x := range ch {
		result[i] = x
		i++
	}
	c.JSON(200, result)
}

type Result struct {
	Status bool               `json:"status"`
	Name   hurtownie.HurtName `json:"hurtName"`
}
