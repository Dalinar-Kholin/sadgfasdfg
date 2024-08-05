package hurtownie

import "net/http"

type HurtName int

const (
	Eurocash HurtName = 1 << iota
	Specjal
	Sot
	Tedi
)

type Token string

type IHurt interface {
	CheckToken(client *http.Client) bool
	RefreshToken(client *http.Client) bool
	TakeToken(login, password string, client *http.Client) bool // te 3 powinny być oddelegowane do obiektu tokena
	SearchProduct(Ean string, client *http.Client) (interface{}, error)
	SearchMany(list WishList, client *http.Client) ([]SearchManyProducts, error)
	AddToCart(list WishList, client *http.Client) bool // tutaj będą potrzebne konkretne instancje itemów z konkretnych hurotwni, czy można dać typy generyczne ?
	GetName() HurtName
}

type SearchManyProducts struct {
	Item interface{}
	Ean  string
}

type WishList struct {
	Items []Items `json:"Items"`
}
type Items struct {
	Ean      string   `json:"Ean"`
	Amount   int      `json:"Amount"`
	HurtName HurtName `json:"HurtName,omitempty"`
}
