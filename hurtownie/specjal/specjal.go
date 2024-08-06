package specjal

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

type Specjal struct {
	Token hurtownie.SotAndSpecjalTokenResponse
}

func (s *Specjal) CheckToken(client *http.Client) bool {
	req, err := http.NewRequest("GET", "https://nowaspecjal.ehurtownia.pl/eh-one-backend/rest/2004/18/9562/oferta?"+
		"lang=PL"+
		"&offset=0"+
		"&limit=1"+
		"&sortAsc=nazwa",
		nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// Set headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "nowaspecjal.ehurtownia.pl")
	req.Header.Set("Referer", "https://nowaspecjal.ehurtownia.pl/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "no-cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-LOC", "2004-4-4-18")
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

func (s *Specjal) RefreshToken(client *http.Client) bool {
	body := "grant_type=refresh_token" +
		"&refresh_token=" + s.Token.RefreshToken +
		"&client_id=ehurtownia-panel-frontend"

	req, err := http.NewRequest("POST", "https://sso.infinite.pl/auth/realms/InfiniteEH/protocol/openid-connect/token", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return false
	}

	req.Header.Set("Host", "sso.infinite.pl")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	req.Header.Set("Referer", "https://nowaspecjal.ehurtownia.pl")
	req.Header.Set("Origin", "https://nowaspecjal.ehurtownia.pl")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Connection", "keep-alive")
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
	if sotToken.AccessToken == "" {
		return false
	}
	return true
}

func (s *Specjal) GetName() hurtownie.HurtName {
	return hurtownie.Specjal
}

// ICOM: do zaimplementowania tego potrzebujemy około 5 requestów, będzie brzydko, ale będzie dziłać, w przeciwieństwie do SOTu nie możemy wygenerwoać sobie swoich nonców i statów, tylko musimy je pobrać z servera
func (s *Specjal) TakeToken(login, password string, client *http.Client) bool {
	firstRequestCookies, sessionCode, execiution, tabId := firstRequest()

	if sessionCode == "" || execiution == "" || tabId == "" || firstRequestCookies == nil {
		return false
	}

	var AuthCookie *http.Cookie
	for _, i2 := range firstRequestCookies {
		if i2.Name == "AUTH_SESSION_ID" {
			AuthCookie = i2
		}
	}
	secondRequestCookies := secondRequest(client, firstRequestCookies, sessionCode, execiution, tabId, login, password)
	if secondRequestCookies == nil {
		return false
	}

	accomodatedCookies := make([]*http.Cookie, 0)
	accomodatedCookies = append(accomodatedCookies, firstRequestCookies...)
	accomodatedCookies = append(accomodatedCookies, secondRequestCookies...)

	thirdResponseCookies, code := thirdRequest(client, accomodatedCookies)
	if thirdResponseCookies == nil {
		return false
	}
	cookiesForTokenRequest := append(thirdResponseCookies, AuthCookie)
	s.Token = tokenRequest(client, cookiesForTokenRequest, code)
	return true
}

func (s *Specjal) SearchProduct(Ean string, client *http.Client) (interface{}, error) {
	req, err := http.NewRequest("GET", "https://nowaspecjal.ehurtownia.pl/eh-one-backend/rest/2004/18/9562/oferta?"+
		"lang=PL"+
		"&offset=0"+
		"&limit=20"+
		"&sortAsc=nazwa"+
		"&sugestia="+Ean+
		"&cechaWartosc="+Ean, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "nowaspecjal.ehurtownia.pl")
	req.Header.Set("Referer", "https://nowaspecjal.ehurtownia.pl/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "no-cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-LOC", "2004-4-4-18")
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

func (s *Specjal) SearchMany(list hurtownie.WishList, client *http.Client) ([]hurtownie.SearchManyProducts, error) {
	ch := make(chan hurtownie.SearchManyProducts)
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
	return result, nil
}

func (s *Specjal) AddToCart(list hurtownie.WishList, client *http.Client) bool {

	setHeader := func(req *http.Request) {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Host", "sot.ehurtownia.pl")
		req.Header.Set("Referer", "https://nowaspecjal.ehurtownia.pl/")
		req.Header.Set("Origin", "https://nowaspecjal.ehurtownia.pl")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-LOC", "2004-4-4-18")
		req.Header.Set("Authorization", "Bearer "+s.Token.AccessToken)
	}

	body := ""
	prodForSpecjal := make([]hurtownie.Items, 0)
	for _, i := range list.Items {
		if i.HurtName == hurtownie.Specjal {
			body += fmt.Sprintf("%v;%v;1;a\n", i.Ean, i.Amount)
			prodForSpecjal = append(prodForSpecjal, i)
		}
	}
	if len(prodForSpecjal) == 0 {
		return true
	}
	basedBody := base64.StdEncoding.EncodeToString([]byte(body))

	req, err := http.NewRequest("POST", "https://nowaspecjal.ehurtownia.pl/eh-one-backend/rest/2006/18/9562/KCFMIN_KOSZ_TEMA/0/0/klahsdlkjbasdkfjbasd.txt/upload_v2?lang=EN", bytes.NewBuffer([]byte(basedBody)))
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
	req, err = http.NewRequest("GET", "https://nowaspecjal.ehurtownia.pl/eh-one-backend/rest/2006/4948/9562/kosz-import-towary/"+
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

	for i, item := range prodForSpecjal {
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

	req, err = http.NewRequest("POST", "https://nowaspecjal.ehurtownia.pl/eh-one-backend/rest/2006/18/9562/wstaw-import-kosz?lang=EN", bytes.NewBuffer(itemizedBody))
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
