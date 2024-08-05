package eurocash

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func takeLoginSiteAndCSRF(client *http.Client) (csrf, location, verifyer string, responseCookie *http.Cookie) {
	state, err := generateNonce(60)
	if err != nil {
		panic(err)
	}
	verifyer, _ = generateNonce(60)
	challenge := generateCodeChallenge(verifyer)
	rUrl := "https://logowanie.eurocash.pl/connect/authorize?response_type=code" +
		"&client_id=EplOnline" +
		"&state=" + state + "semicolonhttps%253A%252F%252Feurocash.pl%252Fang%252Fdashboard" +
		"&redirect_uri=https%3A%2F%2Feurocash.pl%2Fang%2Fdashboard" +
		"&scope=openid%20profile%20offline_access%20IdentityServerApi%20MarketplaceApi%20EplApi%20AgreementsApi%20CartSyncApi" +
		"&code_challenge=" + challenge +
		"&code_challenge_method=S256" +
		"&nonce=" + state +
		"&alt="
	req, err := http.NewRequest("GET", rUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	addHeadersToRequest := func(req *http.Request) {
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Referer", "https://eurocash.pl/")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Sec-Fetch-Dest", "iframe")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "same-site")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("TE", "trailers")
	}
	addHeadersToRequest(req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	// Handle response here as needed
	location = resp.Header.Get("Location")

	req, err = http.NewRequest("GET", location, nil)
	if err != nil {
		panic(err)
	}
	addHeadersToRequest(req)
	resp, err = client.Do(req)

	body, _ := io.ReadAll(resp.Body)
	index := strings.Index(string(body), "Cf")
	csrf = string(body[index : index+155])
	if err != nil {
		panic(err)
	}
	responseCookie = resp.Cookies()[0]

	return
}

func sendCredentials(client *http.Client, scrf, location, login, password string, cookies *http.Cookie) (resLocation string, resCookies []*http.Cookie) {

	var jsonStr = "ReturnUrl=" + location[strings.Index(location, "=")+1:] +
		"&Username=" + login +
		"&Password=" + password +
		"&recaptcha=" +
		"&RememberLogin=False" +
		"&CSRF-TOKEN=" + scrf
	req, err := http.NewRequest("POST", location, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(jsonStr)))
	req.Header.Set("origin", "null")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("TE", "trailers")
	req.Header.Set("upgrade-insecure-requests", "1")

	req.AddCookie(cookies)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	resCookies = resp.Cookies()
	resLocation = resp.Header.Get("Location")
	defer resp.Body.Close()
	return
}

func takeCode(client *http.Client, cookies []*http.Cookie, location string) (code string) {
	req, err := http.NewRequest("GET", "https://logowanie.eurocash.pl"+location, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("TE", "trailers")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	newLocation := resp.Header.Get("Location")
	index := strings.Index(newLocation, "=") + 1
	if index == 0 {
		return ""
	}
	code = newLocation[index : index+66]
	return
}

func takeTokeRequest(client *http.Client, code, veryfyer string) (accessToken string) {
	var jsonStr = "grant_type=authorization_code" +
		"&code=" + code +
		"&redirect_uri=https://eurocash.pl/ang/dashboard" +
		"&code_verifier=" + veryfyer +
		"&client_id=EplOnline&alt="
	req, err := http.NewRequest("POST", "https://logowanie.eurocash.pl/connect/token", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("TE", "trailers")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(jsonStr)))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyByte, _ := io.ReadAll(resp.Body)
	body := string(bodyByte)
	index := strings.Index(body, "\"access_token\":\"") + 16
	accessToken = body[index : index+66]
	return
}

func generateCodeChallenge(codeVerifier string) string {
	// Oblicz SHA-256
	hash := sha256.New()
	hash.Write([]byte(codeVerifier))
	hashed := hash.Sum(nil)

	// Zakoduj w base64url
	base64URL := base64.RawURLEncoding.EncodeToString(hashed)

	return base64URL
}

func generateNonce(targetLength int) (string, error) {
	// Obliczamy długość potrzebnych bajtów, aby po zakodowaniu uzyskać odpowiednią długość
	length := targetLength * 3 / 4
	randomBytes := make([]byte, length)

	// Generujemy losowe bajty
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Kodujemy w Base64 URL-safe i przycinamy do pożądanej długości
	nonce := base64.URLEncoding.EncodeToString(randomBytes)
	if len(nonce) > targetLength {
		nonce = nonce[:targetLength]
	}

	return nonce, nil
}

func makeRequest(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://eurocash.pl/")
	req.Header.Set("Business-Unit", "ECT")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Ch-Ua", `"Not-A.Brand";v="99", "Chromium";v="124"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.6367.60 Safari/537.36")
}

func decodeGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var uncompressedData bytes.Buffer
	_, err = io.Copy(&uncompressedData, reader)
	if err != nil {
		return nil, err
	}
	return uncompressedData.Bytes(), nil
}
