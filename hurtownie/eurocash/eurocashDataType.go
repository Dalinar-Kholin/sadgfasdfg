package eurocash

type EurocashResponse struct {
	Success bool `json:"Success"`
	Data    Data `json:"Data"`
}

type Data struct {
	Items      []EurocashItem `json:"Items"`
	TotalCount int            `json:"TotalCount"`
}

type EurocashItem struct {
	CenaBudzet      float64 `json:"CenaBudzet"`      // licznenie optymalności
	SposobPakowania float64 `json:"SposobPakowania"` // dodawanie do koszyka
}

type EurocashPaginationOptions struct {
	Skip          int  `json:"Skip"`
	Size          int  `json:"Size"`
	GetTotalCount bool `json:"GetTotalCount"`
}

type EuroCashFilter struct {
	Search            string      `json:"Search"`
	PromotionsOnly    bool        `json:"PromotionsOnly"`
	Capacity          []int       `json:"Capacity"`
	MarketOffer       bool        `json:"MarketOffer"`
	Shelf             bool        `json:"Shelf"`
	FilterUsedCoupons bool        `json:"FilterUsedCoupons"`
	Contracts         interface{} `json:"Contracts"`
	LeafletId         int         `json:"LeafletId"`
}

type EurocashRequest struct {
	GetLastOrderedCount bool                      `json:"GetLastOrderedCount"`
	PaginationOptions   EurocashPaginationOptions `json:"PaginationOptions"`
	EurocashFilter      EuroCashFilter            `json:"Filter"`
	Sort                int                       `json:"Sort"`
}

type AddToCartEurocash struct {
	GetShoppingBasketValue bool    `json:"GetShoppingBasketValue"`
	TowarId                string  `json:"TowarId"`
	Typ                    string  `json:"Typ"`
	Cena                   float64 `json:"Cena"`
	CenaBudzet             float64 `json:"CenaBudzet"`
	Wartosc                string  `json:"Wartosc"`
	WartoscPlus            string  `json:"WartoscPlus"`
	PromocjaBudzetowaId    string  `json:"PromocjaBudzetowaId"`
	Vat                    int     `json:"Vat"`
	Bonifikata             float64 `json:"Bonifikata"`
	PromocjaId             *string `json:"PromocjaId"`
}

/* przykładowa kwerenda dodawania do koszyka
Bonifikata: 45.06
Cena : 1.89
CenaBudzet: 1.89
GetShoppingBasketValue: true
PromocjaBudzetowaId: ""
PromocjaId: "0000600365"
PromocjaPakietId: "0000600366"
TowarId: "0000158765"
Typ: "IloscZam"
Vat: 23
Wartosc: "280"
WartoscPlus: "0"
*/
