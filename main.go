package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"optimaHurt/constAndVars"
	"optimaHurt/endpoints/account"
	"optimaHurt/endpoints/orders"
	"optimaHurt/endpoints/takePrices"
	"optimaHurt/middleware"
	"os"
)

func connectToDB() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("CONNECTION_STRING")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the serer
	client, err := mongo.Connect(constAndVars.ContextBackground, opts)
	if err != nil {
		panic(err)
	}
	constAndVars.DbConnect = client.Database(constAndVars.DbName)
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(constAndVars.ContextBackground, bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}

// potem automatycznie dołączam usera do requesta
func main() {

	connection := connectToDB()
	defer func() {
		connection.Disconnect(constAndVars.ContextBackground)
	}()

	r := gin.Default()
	r.Use(middleware.AddHeaders)
	accountEnd := account.AccountEndpoint{}
	order := orders.Order{}

	r.Static("/assets", "./frontend/dist/assets")

	// Obsługa głównego pliku index.html
	r.StaticFile("/", "./frontend/dist/index.html")

	// Obsługa aplikacji typu SPA - przekierowanie wszystkich nieznalezionych ścieżek do index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	prices := takePrices.TakePrices{}
	// Dodanie trasy API
	api := r.Group("/api")
	{
		api.POST("/checkCookie", func(c *gin.Context) {
			Response := func(c *gin.Context, response bool, status int) {
				c.JSON(status, gin.H{
					"response": response,
				})
			}
			token := c.Request.Header.Get("Authorization")
			fmt.Printf("token := %v\n", token)
			if token == "" {
				Response(c, false, http.StatusUnauthorized)
				return
			}
			fmt.Printf("cookie := %v\nmapa := %v\n", token, constAndVars.Users)
			_, ok := constAndVars.Users[token]
			if !ok {
				Response(c, false, http.StatusUnauthorized)
				return
			}
			Response(c, true, http.StatusOK) // ciasteczko jest prawidłowe
		})
		api.POST("/exit", func(c *gin.Context) {
			cookie, _ := c.Request.Cookie("accessToken")

			delete(constAndVars.Users, cookie.Value)
			fmt.Printf("deleted cookie %v\n", cookie.Value)
		})
		api.GET("/addUser", account.AddUser)
		api.GET("/takePrice", middleware.CheckToken, middleware.CheckTokenCurrency, prices.TakePrice)
		api.POST("/takePrices", middleware.CheckToken, middleware.CheckTokenCurrency, prices.TakeMultiple) // get nei może mieć body, więc robimy post
		api.POST("/makeOrder", middleware.CheckToken, middleware.CheckTokenCurrency, order.MakeOrder)
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.POST("/login", accountEnd.Login)
	}

	//r.Run(":8080")
	r.Run("127.0.0.1:" + os.Getenv("PORT"))
	//r.RunTLS("0.0.0.0:"+"443", "./cert.crt", "./key.key")
	return
}

// w makro alholole nie działają
// brak - zamiast brak na składzie
// dodać tedi - około 1k nice
