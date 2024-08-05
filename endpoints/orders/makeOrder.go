package orders

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"sync"
)

type Order struct {
}

func (o *Order) MakeOrder(c *gin.Context) {
	cookie, err := c.Cookie("accessToken")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "where Token?",
		})
		return
	} // po testach do zmiany
	userInstance := constAndVars.Users[cookie]

	var list hurtownie.WishList
	responseReaderJson := json.NewDecoder(c.Request.Body)
	err = responseReaderJson.Decode(&list)
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
