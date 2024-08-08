package account

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/sha3"
	"io"
	"net/http"
	. "optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	"optimaHurt/hurtownie/factory"
	"optimaHurt/user"
	"sync"
)

type LoginBodyData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func MakeNewUser(dbUser user.DataBaseUserObject) (userInstance *user.User, availableHurts hurtownie.HurtName, resultsLogin []ChannelResponse) {
	var hurtTab []hurtownie.IHurt

	var wg sync.WaitGroup

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	ch := make(chan ChannelResponse)
	var credsEnum hurtownie.HurtName = 0
	for _, creds := range dbUser.Creds {
		instance, _ := factory.HurtFactory(creds.HurtName)
		hurtTab = append(hurtTab, instance)
		credsEnum += creds.HurtName
		wg.Add(1)
		go func(hurt *hurtownie.IHurt, wg *sync.WaitGroup, ch chan<- ChannelResponse) {
			defer wg.Done()
			res := (*hurt).TakeToken(creds.Login, creds.Password, client)
			name := (*hurt).GetName()
			if res {
				ch <- ChannelResponse{
					Hurt:    name,
					Success: true,
				}
			} else {
				ch <- ChannelResponse{
					Hurt:    name,
					Success: false,
				}
			}
		}(&instance, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	res := make([]ChannelResponse, len(hurtTab))
	i := 0
	for x := range ch {
		if x.Success {
			availableHurts += x.Hurt
		}
		res[i] = x
		i++
	}
	userInstance = &user.User{Client: client, Hurts: hurtTab, Creds: dbUser.Creds}
	resultsLogin = res
	fmt.Printf("userInstance := %v\n", userInstance)
	return
}

type ChannelResponse struct {
	Hurt    hurtownie.HurtName `json:"hurt"`
	Success bool               `json:"success"`
}

func (a AccountEndpoint) Login(c *gin.Context) {
	connection := DbConnect.Collection(UserCollection)

	request, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant read body"})
		return
	}
	defer c.Request.Body.Close()
	var reqBody LoginBodyData
	err = json.Unmarshal(request, &reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	var dataBaseResponse user.DataBaseUserObject

	err = connection.FindOne(ContextBackground, bson.M{
		"username": reqBody.Username,
		"password": reqBody.Password,
	}).Decode(&dataBaseResponse)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "bad credentials"})
		return
	}
	fmt.Printf("dataBaseResponse := %v\n", dataBaseResponse)
	userInstance, loggedHurts, loginLog := MakeNewUser(dataBaseResponse)
	if loggedHurts == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cant logged to selected Hurts"})
		return
	}
	newToken := make([]byte, 64)
	if _, err := rand.Read(newToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant generate token"})
		return
	}
	resToken := sha3.Sum256(newToken)
	shieldedToken := hex.EncodeToString(resToken[:])
	//shieldedToken := base64.StdEncoding.EncodeToString(newToken) // token jest użytkowany tylko podczas sesji, więc nie ma potrzeby przechowywania go w bazie danych
	Users[shieldedToken] = *userInstance
	c.JSON(200, gin.H{
		"result":         loginLog,
		"token":          shieldedToken,
		"availableHurts": loggedHurts},
	)
}

//
// makro
//
