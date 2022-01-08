package neugo

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	// generated from EncryptURLToWebVPN("https://pass.neu.edu.cn/tpass/login")
	webvpnLoginURL = "https://webvpn.neu.edu.cn/https/77726476706e69737468656265737421e0f6528f693e6d45300d8db9d6562d/tpass/login"
	casLoginURL    = "https://pass.neu.edu.cn/tpass/login"

	webvpnCookieDomain = ".webvpn.neu.edu.cn"
	casCookieDomain    = "pass.neu.edu.cn"
)

// Perform login request and returns the body of response.
//
// If login failed, an error will be returned with the body.
func login(c config) (string, error) {
	var resp *http.Response
	var err error
	var loginURL string
	var request *http.Request
	if c.Platform == CAS {
		loginURL = casLoginURL
	} else {
		loginURL = webvpnLoginURL
	}
	if c.UseToken {
		setToken(c.Client, c.Token, c.Platform)
		request = buildGetRequest(loginURL)
	} else {
		lt, err := getLT(c.Client, loginURL)
		if err != nil {
			return "", err
		}
		request = buildAuthRequest(c.Username, c.Password, lt, loginURL)
	}
	resp, err = c.Client.Do(request)
	if err != nil {
		return "", err
	}

	body := extractBody(resp)

	_, err = isLogged(body)

	return body, err
}

var (
	ltExp           = regexp.MustCompile(`name="lt" value="(.+?)"`)
	ErrorLTNotFound = errors.New("LT not found")
)

// Get LT by performing a pre-request.
func getLT(client *http.Client, requestURL string) (string, error) {
	req := buildGetRequest(requestURL)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body := extractBody(resp)

	lt := matchSingle(ltExp, body)
	if lt == "" {
		return "", ErrorLTNotFound
	}
	return lt, nil
}

// Build a POST *http.Request for login.
func buildAuthRequest(username, password, lt, reqURL string) (req *http.Request) {
	data := "rsa=" + username + password + lt +
		"&ul=" + strconv.Itoa(len(username)) +
		"&pl=" + strconv.Itoa(len(password)) +
		"&lt=" + lt +
		"&execution=e1s1" +
		"&_eventId=submit"

	req, _ = http.NewRequest("POST",
		reqURL,
		strings.NewReader(data))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", reqURL)
	return
}

// Build a GET *http.Request.
func buildGetRequest(reqURL string) (req *http.Request) {
	req, _ = http.NewRequest("GET", reqURL, nil)
	return
}

var (
	titleExp = regexp.MustCompile(`<title>(.+?)</title>`)
)

var (
	ErrorAccountBanned     = errors.New("account is banned")
	ErrorAccountNeedsReset = errors.New("account needs reset")
	ErrorAuthFailed        = errors.New("incorrect username or password or cookie")
)

// Check if logged in successfully.
func isLogged(body string) (bool, error) {
	// FIXME: use better verification methods instead of comparing titles

	title := matchSingle(titleExp, body)
	if title == "" {
		return true, nil
	}

	switch title {
	case "智慧东大--统一身份认证":
		return false, ErrorAuthFailed
	case "智慧东大":
		return false, ErrorAccountNeedsReset
	case "系统提示":
		return false, ErrorAccountBanned
	}
	return true, nil
}
