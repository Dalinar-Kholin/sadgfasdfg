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

func MakeNewUser(dbUser user.DataBaseUserObject) (userInstance *user.User) {
	var hurtTab []hurtownie.IHurt

	var wg sync.WaitGroup

	/*proxyURL, err := url.Parse("http://127.0.0.1:8000")
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return
	}*/

	// Step 2: Create a transport that uses the proxy
	/*transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}*/

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		//Transport: transport,
	}

	for _, creds := range dbUser.Creds { // do poprawy nie działa gdy użytkownik ma różną liczbę credsów względem hurtowni
		fmt.Printf("creds := %v\n", creds)
		instance, _ := factory.HurtFactory(creds.HurtName)
		hurtTab = append(hurtTab, instance)
		wg.Add(1)
		go func(hurt *hurtownie.IHurt, wg *sync.WaitGroup) {
			defer wg.Done()
			(*hurt).TakeToken(creds.Login, creds.Password, client)
		}(&instance, &wg)
	}

	go func() {
		wg.Wait()
	}()

	userInstance = &user.User{Client: client, Hurts: hurtTab, Creds: dbUser.Creds}
	fmt.Printf("userInstance := %v\n", userInstance)
	return
}

func (a AccountEndpoint) Login(c *gin.Context) {
	connection := DbConnect.Collection(UserCollection)

	request, err := io.ReadAll(c.Request.Body)
	if err != nil {
		println("poggerus")
		panic(err)
	}
	defer c.Request.Body.Close()
	var reqBody LoginBodyData
	err = json.Unmarshal(request, &reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
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
	userInstance := MakeNewUser(dataBaseResponse)

	newToken := make([]byte, 64)
	if _, err := rand.Read(newToken); err != nil {
		panic(err)
	}
	resToken := sha3.Sum256(newToken)
	shildedToken := hex.EncodeToString(resToken[:])
	//shieldedToken := base64.StdEncoding.EncodeToString(newToken) // token jest użytkowany tylko podczas sesji, więc nie ma potrzeby przechowywania go w bazie danych
	Users[shildedToken] = *userInstance
	c.SetCookie("accessToken", shildedToken, 0, "/", "127.0.0.1", true, false) // httpOnly musi być false, aby js mógł odczytać ciasteczko i dołączyć je do kończenia sesji
	c.JSON(200, gin.H{"result": "success"})
}

//
// makro
//
