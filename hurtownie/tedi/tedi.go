package tedi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"optimaHurt/hurtownie"
	"strconv"
	"strings"
	"sync"
)

// ICOM: aby dodać cos do koszyka wystarczy posiadać Ean i ile produktów chcemy dodać --  proste w chuj
type Tedi struct {
	Token hurtownie.SotAndSpecjalTokenResponse
}

func (t *Tedi) CheckToken(client *http.Client) bool {
	req, err := http.NewRequest("GET", "https://tedi-ws.ampli-solutions.com/product-attributes/?getAll=true&search=atdyhmdfghmdghfkmdhgi", nil)
	if err != nil {
		return false
	}
	req.Header.Add("Origin", "https://tedi.kd-24.pl")
	req.Header.Add("Referer", "https://tedi.kd-24.pl")
	req.Header.Add("Accept-Language", "PL")
	req.Header.Add("Amper_app_name", "B2B")
	req.Header.Add("Authorization", "Bearer "+t.Token.AccessToken)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

func (t *Tedi) RefreshToken(client *http.Client) bool {
	body := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		t.Token.RefreshToken,
	}
	jsoned, err := json.Marshal(body)
	if err != nil {
		return false
	}
	req, err := http.NewRequest("POST", "https://tedi-ws.ampli-solutions.com/auth/token-refresh/", bytes.NewBuffer(jsoned))
	if err != nil {
		return false
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsoned)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Origin", "https://tedi.kd-24.pl")
	req.Header.Add("Referer", "https://tedi.kd-24.pl/")
	req.Header.Add("Priority", "u=1, i")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nfatal error := %v\n", err)
		return false
	}

	var sotToken hurtownie.SotAndSpecjalTokenResponse
	responseReader := json.NewDecoder(resp.Body)
	err = responseReader.Decode(&sotToken)
	if err != nil {
		bdy, _ := io.ReadAll(resp.Body)
		panic("błąd przy parsowaniu response od tedi\nerror := " + err.Error() + "\n" + string(bdy))
		return false
	}
	t.Token = sotToken
	fmt.Printf("token := %v", t.Token)
	return true
} // do naprawy, nie działa

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (t *Tedi) TakeToken(login, password string, client *http.Client) bool {
	url := "https://tedi-ws.ampli-solutions.com/auth/?session_id=" + RandomString(128)
	body := Creds{
		login,
		password,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return false
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false
	}

	var tediTokenResponse hurtownie.SotAndSpecjalTokenResponse
	responseReaderJson := json.NewDecoder(resp.Body)
	err = responseReaderJson.Decode(&tediTokenResponse)
	if err != nil {
		return false
	}
	t.Token = tediTokenResponse
	return true
}

func (t *Tedi) SearchProduct(Ean string, client *http.Client) (interface{}, error) {
	req, err := http.NewRequest("GET", "https://tedi-ws.ampli-solutions.com/product-search/?limit=12&search="+Ean+"&offset=0", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "PL")
	req.Header.Add("Amper_app_name", "B2B")
	req.Header.Add("Origin", "https://tedi.kd-24.pl")
	req.Header.Add("Referer", "https://tedi.kd-24.pl")
	req.Header.Add("Authorization", "Bearer "+t.Token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var serverResponse ProductResponse
	responseReaderJson := json.NewDecoder(resp.Body)
	err = responseReaderJson.Decode(&serverResponse)
	if err != nil {
		return nil, err
	}
	return serverResponse, nil
}

func (t *Tedi) SearchMany(list hurtownie.WishList, client *http.Client) ([]hurtownie.SearchManyProducts, error) {
	ch := make(chan hurtownie.SearchManyProducts)
	var wg sync.WaitGroup
	for _, i := range list.Items {
		wg.Add(1)
		go func(wg *sync.WaitGroup, ch chan<- hurtownie.SearchManyProducts, ean string) {
			defer wg.Done()
			res, err := t.SearchProduct(i.Ean, client)
			if err != nil {
				ch <- hurtownie.SearchManyProducts{
					Item: -1,
					Ean:  ean,
				}
			}
			ch <- hurtownie.SearchManyProducts{
				Item: res,
				Ean:  ean,
			}
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

func (t *Tedi) AddToCart(list hurtownie.WishList, client *http.Client) bool {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	data := ""
	prodForTedi := make([]hurtownie.Items, 0)
	// wyciąganie danych z wishListy na KCfirme, gdize nie liczy się ani 3 ani 4 kolumna
	for _, i := range list.Items {
		if i.HurtName == hurtownie.Tedi {
			data += fmt.Sprintf("%v;%v;1;a\n", i.Ean, i.Amount)
			prodForTedi = append(prodForTedi, i)
		}
	}
	if len(prodForTedi) == 0 {
		return true
	}
	// Add the file content to the multipart form data
	part, err := writer.CreateFormFile("order_file", "order.txt")

	if err != nil {
		fmt.Println(err)
		return false
	}

	if _, err = io.Copy(part, strings.NewReader(data)); err != nil {
		fmt.Println(err)
		return false
	}

	if writer.WriteField("order_format", "kcfirma") != nil {
		fmt.Println(err)
		return false
	}

	if writer.Close() != nil {
		fmt.Println(err)
		return false
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://tedi-ws.ampli-solutions.com/create-order-from-file/", &body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Set the headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "PL")
	req.Header.Add("Amper_app_name", "B2B")
	req.Header.Add("Origin", "https://tedi.kd-24.pl")
	req.Header.Add("Referer", "https://tedi.kd-24.pl")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "pl")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "no-cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("AMPER_APP_NAME", "B2B")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+t.Token.AccessToken)
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	// Handle the response
	fmt.Println("Status:", resp.Status)
	return true
}

func (t *Tedi) GetName() hurtownie.HurtName {
	return hurtownie.Tedi
}
