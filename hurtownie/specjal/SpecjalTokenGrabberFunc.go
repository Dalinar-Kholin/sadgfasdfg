package specjal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"optimaHurt/hurtownie"
	"strconv"
	"strings"
)

func firstRequest() ([]*http.Cookie, string, string, string) {
	res, err := http.Get("https://specjal.ehurtownia.pl/")
	if err != nil {
		return nil, "", "", ""
	}
	cookies := res.Cookies()
	tempBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", "", ""
	}
	body := string(tempBody)
	index := strings.Index(body, " action=\"")

	link := body[index+8 : index+strings.Index(body[index+9:], "\"")+9] // powinno być poprawnie

	index = strings.Index(link, "session_code=") + 13
	sessionCode := link[index : index+strings.Index(link[index:], "&")] //
	index = strings.Index(link, "tab_id=") + 7
	tabId := link[index:] //
	index = strings.Index(link, "execution") + 10
	execution := link[index : index+strings.Index(link[index:], "&")] //

	return cookies, sessionCode, execution, tabId
}

func secondRequest(client *http.Client, cookies []*http.Cookie, sessionCode, execution, tabId, username, password string) []*http.Cookie {

	body := "username=" + username +
		"&password=" + password +
		"&credentialId="

	req, err := http.NewRequest("POST",
		"https://sso.infinite.pl/auth/realms/InfiniteEH/login-actions/authenticate?"+
			"session_code="+sessionCode+
			"&execution="+execution+
			"&client_id=tema-seam"+
			"&tab_id="+tabId,
		bytes.NewBuffer([]byte(body))) // powinniśmy dodać redirecta i to jest w chuj ważne
	if err != nil {
		return nil
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Referer", "https://sso.infinite.pl/")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://sso.infinite.pl")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Content-Length", strconv.Itoa(len(body)))

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 302 {
		return nil

	}
	responseCookies := resp.Cookies()

	return responseCookies
}

func thirdRequest(client *http.Client, cookies []*http.Cookie) ([]*http.Cookie, string) {
	nonce := hurtownie.GenerateUUID()
	state := hurtownie.GenerateUUID()
	req, err := http.NewRequest("GET", "https://sso.infinite.pl/auth/realms/InfiniteEH/protocol/openid-connect/auth?client_id=ehurtownia-panel-frontend"+
		"&redirect_uri=https://nowaspecjal.ehurtownia.pl/"+
		"&state="+state+
		"&nonce="+nonce+
		"&response_mode=fragment"+
		"&response_type=code"+
		"&scope=openid", nil)
	if err != nil {
		return nil, ""
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Host", "sso.infinite.pl")
	req.Header.Add("Referer", "https://nowaspecjal.ehurtownia.pl/")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, ""
	}
	responseCookies := response.Cookies()

	code := ""
	location := response.Header.Get("Location")
	if location != "" {
		code = location[strings.Index(location, "code=")+5:]
	}
	return responseCookies, code

}

func tokenRequest(client *http.Client, cookies []*http.Cookie, code string) hurtownie.SotAndSpecjalTokenResponse {

	requestBody := "code=" + code +
		"&grant_type=authorization_code&client_id=ehurtownia-panel-frontend&redirect_uri=https%3A%2F%2Fnowaspecjal.ehurtownia.pl%2F"
	req, err := http.NewRequest("POST", "https://sso.infinite.pl/auth/realms/InfiniteEH/protocol/openid-connect/token", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return hurtownie.SotAndSpecjalTokenResponse{}
	}
	req.Header.Set("Host", "sso.infinite.pl")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(requestBody)))
	req.Header.Set("Referer", "https://nowaspecjal.ehurtownia.pl")
	req.Header.Set("Origin", "https://nowaspecjal.ehurtownia.pl")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Connection", "keep-alive")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return hurtownie.SotAndSpecjalTokenResponse{}
	}
	var specjalTokenResponse hurtownie.SotAndSpecjalTokenResponse
	responseReaderJson := json.NewDecoder(resp.Body)
	err = responseReaderJson.Decode(&specjalTokenResponse)
	if err != nil {
		return hurtownie.SotAndSpecjalTokenResponse{}
	}

	return specjalTokenResponse
}
