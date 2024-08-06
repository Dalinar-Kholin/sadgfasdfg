package sot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"optimaHurt/hurtownie"
	"strconv"
	"strings"
)

func firstRequestForToken(client *http.Client) (authCookies *http.Cookie, firstRequestCookie []*http.Cookie, sessionCode, tabID string) {
	state := hurtownie.GenerateUUID()
	nonce := hurtownie.GenerateUUID()

	req, err := http.NewRequest("GET",
		"https://sso.infinite.pl/auth/realms/InfiniteEH/protocol/openid-connect/auth?client_id=ehurtownia-panel-frontend&redirect_uri=https%3A%2F%2Fsot.ehurtownia.pl%2F&"+
			"state="+
			state+
			"&nonce="+
			nonce+
			"&response_mode=fragment&response_type=code&scope=openid", nil)
	// o dziwo state i nonce są generowane po stronie klienta, a następnie podpisywane przez serwer

	if err != nil {
		fmt.Printf("Błąd przy tworzeniu żądania: %v", err)
		return
	}

	req.Header.Set("Host", "sso.infinite.pl")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://sot.ehurtownia.pl/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "cross-site")

	// Utwórz transport z ustawionym proxy

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Błąd wykonania żądania do Sot: %v", err)
		return
	}
	firstRequestCookie = resp.Cookies()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Błąd odczytu danych:", err)
		return
	}
	sessionCodePlace := strings.Index(string(responseBody), "session_code=") + 13
	sessionCode = string(responseBody[sessionCodePlace : sessionCodePlace+43]) // 43 to długość session_code
	tabIdPlace := strings.Index(string(responseBody), "tab_id=") + 7
	tabID = string(responseBody[tabIdPlace : tabIdPlace+11]) // 11 to długość ID
	authCookies = resp.Cookies()[0]                          // powinno wziąć Auth ale kurwa głowy za to nie dam
	return
}

func secondRequestForToken(client *http.Client,
	firstRequestCookie []*http.Cookie,
	sessionCode, tabId, login, password string) (secondRequestCookie []*http.Cookie, code string) {

	urlCookie := http.Cookie{
		Name:     "url",
		Value:    "https://sot.ehurtownia.pl/",
		Domain:   "sso.infinite.pl",
		HttpOnly: false,
		Path:     "/",
	}
	loginPayload := "username=" + login + "&password=" + password + "&login=Log+in"

	req, err := http.NewRequest("POST",
		"https://sso.infinite.pl/auth/realms/InfiniteEH/login-actions/authenticate?"+
			"session_code="+sessionCode+
			"&execution=ce7d4a91-ac8f-4977-afa0-c28229423e07"+
			"&client_id=ehurtownia-panel-frontend"+
			"&tab_id="+tabId, bytes.NewBuffer([]byte(loginPayload)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.AddCookie(firstRequestCookie[0])
	req.AddCookie(firstRequestCookie[1])
	req.AddCookie(&urlCookie)
	req.Header.Add("Host", "sso.infinite.pl")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Referer", "https://sso.infinite.pl/")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(loginPayload))) // Update this if the payload changes
	req.Header.Add("Origin", "https://sso.infinite.pl")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")

	response, err := client.Do(req)

	if err != nil {
		fmt.Printf("\nfatal error := %v\n", err)
		return
	}
	secondRequestCookie = response.Cookies()

	location := response.Header.Get("Location")

	if location == "" {
		fmt.Println("Location is empty")
		return
	}

	code = location[strings.Index(location, "code=")+5:]
	return
}

func thirdRequestForToken(code string, client *http.Client, secondRequestCookies []*http.Cookie, AuthSessionIDCookie *http.Cookie) (token hurtownie.SotAndSpecjalTokenResponse) {
	body := "code=" + code +
		"&grant_type=authorization_code&" +
		"client_id=ehurtownia-panel-frontend&" +
		"redirect_uri=https%3A%2F%2Fsot.ehurtownia.pl%2F"

	req, err := http.NewRequest("POST", "https://sso.infinite.pl/auth/realms/InfiniteEH/protocol/openid-connect/token", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return hurtownie.SotAndSpecjalTokenResponse{}
	}

	req.Header.Set("Host", "sso.infinite.pl")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	req.Header.Set("Referer", "https://sot.ehurtownia.pl/")
	req.Header.Set("Origin", "https://sot.ehurtownia.pl")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Connection", "keep-alive")

	for _, i := range secondRequestCookies {
		req.AddCookie(i)
	}
	req.AddCookie(AuthSessionIDCookie)

	resp, err := client.Do(req)

	if err != nil {
		return hurtownie.SotAndSpecjalTokenResponse{}
	}

	var sotToken hurtownie.SotAndSpecjalTokenResponse
	responseReaderJson := json.NewDecoder(resp.Body)
	err = responseReaderJson.Decode(&sotToken)
	if err != nil {
		return hurtownie.SotAndSpecjalTokenResponse{}
	}

	return sotToken
}
