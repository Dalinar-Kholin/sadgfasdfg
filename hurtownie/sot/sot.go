package sot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"optimaHurt/hurtownie"
	"strconv"
	"sync"
)

type Sot struct {
	Token   hurtownie.SotAndSpecjalTokenResponse
	Cookies []*http.Cookie
}

func (s *Sot) CheckToken(client *http.Client) bool {
	req, err := http.NewRequest("GET", "https://sot.ehurtownia.pl/eh-one-backend/rest/71/5/1503697/oferta?lang=EN"+
		"&offset=0"+
		"&limit=1"+
		"&sortAsc=nazwa",
		nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// Dodanie nagłówków do zapytania
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "sot.ehurtownia.pl")
	req.Header.Set("Referer", "https://sot.ehurtownia.pl/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-LOC", "6-6-8-5")
	req.Header.Set("Authorization", "Bearer "+s.Token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func (s *Sot) RefreshToken(client *http.Client) bool {
	return s.RefreshTokenFunc(client)
}

func (s *Sot) GetName() hurtownie.HurtName {
	return hurtownie.Sot
}

/*
tylko i wyłącznie effekty uboczne w postaci aktualizacji Token i RefreshToken
*/
func (s *Sot) RefreshTokenFunc(client *http.Client) bool {
	body := "grant_type=refresh_token" +
		"&refresh_token=" + s.Token.RefreshToken +
		"&client_id=ehurtownia-panel-frontend"

	req, err := http.NewRequest("POST", "https://sso.infinite.pl/auth/realms/InfiniteEH/protocol/openid-connect/token", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return false
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", strconv.Itoa(len(body)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "sso.infinite.pl")
	req.Header.Add("Origin", "https://sot.ehurtownia.pl")
	req.Header.Add("Referer", "https://sot.ehurtownia.pl/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")

	for _, i := range s.Cookies {
		req.AddCookie(i)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nfatal error := %v\n", err)
		return false
	}

	var sotToken hurtownie.SotAndSpecjalTokenResponse
	responseReader := json.NewDecoder(resp.Body)
	err = responseReader.Decode(&sotToken)
	if err != nil {
		return false
	}
	s.Token = sotToken
	return true
}

/*
efektem ubocznym funckji jest ustawienie parametrów Token, RefreshToken, SessionID w obiekcie pierwotnym
*/
func (s *Sot) TakeToken(login, password string, client *http.Client) bool {

	AuthSessionIDCookie, firstRequestCookie, sessionCode, tabId := firstRequestForToken(client)
	if tabId == "" || AuthSessionIDCookie == nil || firstRequestCookie == nil {
		return false
	}

	secondRequestCookies, code := secondRequestForToken(client, firstRequestCookie, sessionCode, tabId, login, password)

	if code == "" || secondRequestCookies == nil {
		return false
	}

	s.Cookies = append(secondRequestCookies, AuthSessionIDCookie)

	token := thirdRequestForToken(code, client, secondRequestCookies, AuthSessionIDCookie)

	s.Token = token
	if token.AccessToken == "" {
		return false
	}
	return true
}

func (s *Sot) SearchProduct(Ean string, client *http.Client) (interface{}, error) {
	req, err := http.NewRequest("GET", "https://sot.ehurtownia.pl/eh-one-backend/rest/71/5/1503697/oferta?lang=EN"+
		"&offset=0"+
		"&limit=20"+
		"&sortAsc=nazwa"+
		"&sugestia="+Ean+
		"&cechaWartosc="+Ean,
		nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Dodanie nagłówków do zapytania
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "sot.ehurtownia.pl")
	req.Header.Set("Referer", "https://sot.ehurtownia.pl/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-LOC", "6-6-8-5")
	req.Header.Set("Authorization", "Bearer "+s.Token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	var itemek hurtownie.SotAndSpecjalResponse
	err = json.Unmarshal(responseBody, &itemek)

	if err != nil {
		return nil, err
	}
	if itemek.CountPozycji == 0 {
		return nil, errors.New("Brak produktu w bazie SOT")
	}
	return itemek, nil
}

func (s *Sot) SearchMany(list hurtownie.WishList, client *http.Client) ([]hurtownie.SearchManyProducts, error) {
	ch := make(chan hurtownie.SearchManyProducts) // do tego chana wejdzie TYLKO I WYŁĄCZNIE albo -1 albo Item
	var wg sync.WaitGroup
	for _, i := range list.Items {
		wg.Add(1)
		go func(wg *sync.WaitGroup, ch chan<- hurtownie.SearchManyProducts, ean string) {
			defer wg.Done()
			res, err := s.SearchProduct(i.Ean, client)
			if err != nil {
				ch <- hurtownie.SearchManyProducts{Item: -1, Ean: ean}
				return
			}
			ch <- hurtownie.SearchManyProducts{Item: res, Ean: ean}
		}(&wg, ch, i.Ean)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	result := make([]hurtownie.SearchManyProducts, len(list.Items))
	i := 0
	for x := range ch {
		result[i] = x
		i++
	}
	fmt.Printf("result := %v\n", result)
	return result, nil
}

func (s *Sot) AddToCart(list hurtownie.WishList, client *http.Client) bool {

	setHeader := func(req *http.Request) {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Host", "sot.ehurtownia.pl")
		req.Header.Set("Referer", "https://sot.ehurtownia.pl/")
		req.Header.Set("Origin", "https://sot.ehurtownia.pl")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-LOC", "6-6-8-5")
		req.Header.Set("Authorization", "Bearer "+s.Token.AccessToken)
	}

	body := ""
	prodForSot := make([]hurtownie.Items, 0)
	for _, i := range list.Items {
		if i.HurtName == hurtownie.Sot {
			body += fmt.Sprintf("%v;%v;1;a\n", i.Ean, i.Amount)
			prodForSot = append(prodForSot, i)
		}
	}
	if len(prodForSot) == 0 {
		return true
	}
	basedBody := base64.StdEncoding.EncodeToString([]byte(body))

	req, err := http.NewRequest("POST", "https://sot.ehurtownia.pl/eh-one-backend/rest/71/5/1503697/KCF_KOSZ_HZ/0/0/sot-wroclaw-20240802T1440.csv/upload_v2?lang=EN", bytes.NewBuffer([]byte(basedBody)))
	if err != nil {
		return false
	}
	setHeader(req)
	// poprawić headery
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	var resJson SorResponse
	responseReader := json.NewDecoder(resp.Body)
	err = responseReader.Decode(&resJson)
	if err != nil {
		return false
	}
	req, err = http.NewRequest("GET", "https://sot.ehurtownia.pl/eh-one-backend/rest/71/329060/1503697/kosz-import-towary/"+
		strconv.Itoa(resJson.Codes[0])+
		"?lang=EN&offset=0&limit=6000&sortAsc=lp", nil)
	if err != nil {
		return false
	}
	setHeader(req)

	resp, err = client.Do(req)
	if err != nil {
		return false
	}
	var productRes []ProductResponse
	responseReader = json.NewDecoder(resp.Body)
	err = responseReader.Decode(&productRes)
	/*ICOM : WAŻNE W CHUJ, ZKAŁAKDAM ŻE KOLEJNOŚĆ PRODUKTÓW JEST TAKA JAK NA LIŚCIE
	mogę więc iterować po liście i odnościć się do ilości produktów*/

	newOrder := NewOrder{
		Id:       1,
		Origin:   "PLIK",
		Products: []Prod{},
	}

	for i, item := range prodForSot {
		newOrder.Products = append(newOrder.Products, Prod{
			Id:        nil,
			Amount:    item.Amount,
			Index:     productRes[i].Index,
			SetCode:   0,
			PromoCode: 0,
		})
	}

	itemizedBody, err := json.Marshal(newOrder)
	if err != nil {
		fmt.Printf("Error marshalling body := %v\n", err)
		return false
	}

	req, err = http.NewRequest("POST", "https://sot.ehurtownia.pl/eh-one-backend/rest/71/5/1503697/wstaw-import-kosz?lang=EN", bytes.NewBuffer(itemizedBody))
	setHeader(req)
	resp, err = client.Do(req)
	if err != nil {
		return false
	}

	return resp.StatusCode == 201
}

type SorResponse struct {
	Codes []int `json:"hdrKody"`
}

type ProductResponse struct {
	Index string `json:"indeks"`
}

type Prod struct {
	Id        interface{} `json:"id"` // aby można było nil-ować
	Amount    int         `json:"ilZam"`
	Index     string      `json:"indeks"`
	PromoCode int         `json:"kodPromocji"`
	SetCode   int         `json:"kodZestawu"`
}

type NewOrder struct {
	Id       int    `json:"id"`
	Origin   string `json:"pochodzenie"`
	Products []Prod `json:"towary"`
}
