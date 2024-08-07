package eurocash

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"optimaHurt/hurtownie"
	"strconv"
	"sync"
)

type EurocashObject struct {
	Token string
}

func (e *EurocashObject) CheckToken(client *http.Client) bool {
	url := "https://ehurtapi.eurocash.pl/api/offer/getExtraBanner?placement=4"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false

	}
	req.Header.Set("Authorization", "Bearer "+e.Token)
	makeRequest(req)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func (e *EurocashObject) RefreshToken(client *http.Client) bool {
	return false
}

func (e *EurocashObject) GetName() hurtownie.HurtName {
	return hurtownie.Eurocash
}

func (e *EurocashObject) SearchMany(list hurtownie.WishList, client *http.Client) ([]hurtownie.SearchManyProducts, error) {
	ch := make(chan hurtownie.SearchManyProducts)
	var wg sync.WaitGroup
	for _, i := range list.Items {
		wg.Add(1)
		go func(wg *sync.WaitGroup, ch chan<- hurtownie.SearchManyProducts, ean string) {
			defer wg.Done()
			res, err := e.SearchProduct(i.Ean, client)
			if err != nil {
				ch <- hurtownie.SearchManyProducts{
					Item: nil,
					Ean:  ean,
				}
				return
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

func (e *EurocashObject) TakeToken(login, password string, client *http.Client) bool {

	csrf, location, veryfyer, CsrfCookie := takeLoginSiteAndCSRF(client)

	if csrf == "" || location == "" || veryfyer == "" || CsrfCookie == nil {
		return false
	}
	location, cookies := sendCredentials(client, csrf, location, login, password, CsrfCookie)
	if location == "" || cookies == nil {
		return false
	}
	cookies = append(cookies, CsrfCookie)
	code := takeCode(client, cookies, location)
	if code == "" {
		return false
	}
	accessToken := takeTokeRequest(client, code, veryfyer)
	e.Token = accessToken
	fmt.Printf("accessToken := %v\n", accessToken)
	return true
}

func (e *EurocashObject) SearchProduct(Ean string, client *http.Client) (interface{}, error) {
	request := EurocashRequest{
		GetLastOrderedCount: true,
		PaginationOptions: EurocashPaginationOptions{
			Skip:          0,
			Size:          35,
			GetTotalCount: true,
		},
		EurocashFilter: EuroCashFilter{
			Search:            Ean,
			PromotionsOnly:    false,
			Capacity:          []int{},
			MarketOffer:       false,
			Shelf:             false,
			FilterUsedCoupons: false,
			Contracts:         nil,
			LeafletId:         0,
		},
		Sort: 1,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(request)
	req, err := http.NewRequest("POST", "https://ehurtapi.eurocash.pl/api/offer/getOfferListWithPromotions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Błąd przy tworzeniu żądania:", err)
		return nil, errors.New("Błąd przy tworzeniu żądania")
	}
	fmt.Printf("req := %v\n", req)
	makeRequest(req)
	req.Header.Set("Authorization", "Bearer "+e.Token)
	req.Header.Set("Content-Length", strconv.Itoa(len(jsonData)))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Błąd przy wykonaniu żądania pogger:", err)
		return nil, errors.New("Błąd przy wykonaniu żądania")
	}
	var itemData EurocashResponse
	jsonReader := json.NewDecoder(resp.Body)
	err = jsonReader.Decode(&itemData)
	if err != nil {
		fmt.Println("Błąd przy dekodowaniu odpowiedzi:", err)
		return nil, errors.New("Błąd przy dekodowaniu odpowiedzi")

	}
	return itemData, nil
}

func (e *EurocashObject) AddToCart(list hurtownie.WishList, client *http.Client) bool {
	// ICOM: do zrobienia aby pobioerało dane z wishList

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Manually create a part for the file with custom Content-Type
	partHeaders := make(textproto.MIMEHeader)
	partHeaders.Set("Content-Disposition", `form-data; name="uploadFile"; filename="sot-wroclaw-20240802T1440.txt"`)
	partHeaders.Set("Content-Type", "text/plain")
	fileWriter, err := writer.CreatePart(partHeaders)
	if err != nil {
		fmt.Println(err)
		return false
	}

	data := ""
	prodForEurocash := make([]hurtownie.Items, 0)
	for _, i := range list.Items {
		if i.HurtName == hurtownie.Eurocash {
			data += fmt.Sprintf("%v;%v;1;a\n", i.Ean, i.Amount)
			prodForEurocash = append(prodForEurocash, i)
		}
	}
	if len(prodForEurocash) == 0 {
		return true
	}
	_, err = fileWriter.Write([]byte(data))
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Create a form field for FILE_EXTENSION
	err = writer.WriteField("FILE_EXTENSION", "Z_KCM")
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Close the writer to set the terminating boundary
	if writer.Close() != nil {
		return false
	}

	req, err := http.NewRequest("POST", "https://ehurtapi.eurocash.pl/api/order/importHistory", &requestBody)
	if err != nil {
		println("err := %v\n", err)
		return false
	}
	req.Header.Set("Authorization", "Berear "+e.Token)
	req.Header.Set("Host", "ehurtapi.eurocash.pl")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Origin", "https://eurocash.pl")
	req.Header.Set("Referer", "https://eurocash.pl/")
	req.Header.Set("Business-Unit", "ECT")
	req.Header.Set("Sec-Ch-Ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.60 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if resp.StatusCode != 200 {

		return false
	}
	// trzeba potwierdzić koszyk
	req, err = http.NewRequest("GET", "https://ehurtapi.eurocash.pl/api/order/rewriteImportedData/false", nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", "Berear "+e.Token)
	req.Header.Set("Host", "ehurtapi.eurocash.pl")
	req.Header.Set("Origin", "https://eurocash.pl")
	req.Header.Set("Referer", "https://eurocash.pl/")
	req.Header.Set("Business-Unit", "ECT")
	req.Header.Set("Sec-Ch-Ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.60 Safari/537.36")
	resp, err = client.Do(req)

	if err != nil {
		fmt.Println(err)
		return false
	}

	var response EurocashAddAToCartResponse
	jsonReader := json.NewDecoder(resp.Body)
	err = jsonReader.Decode(&response)

	if err != nil {
		fmt.Println(err)
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}

	return response.Success
}

type EurocashAddAToCartResponse struct {
	Success bool `json:"Success"`
}

type ResponseFromAddingToCart struct {
	Status bool `json:"Status"`
}
