package neugo

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	portalURL = "https://portal.neu.edu.cn/tp_up/"

	webvpnBaseURL = "https://pass-443.webvpn.neu.edu.cn/tpass/login?service="
	casBaseURL    = "https://pass.neu.edu.cn/tpass/login?service="

	webvpnDomain = "pass-443.webvpn.neu.edu.cn"
	casDomain    = "pass.neu.edu.cn"
)

type cas struct {
	Username, Password, Token string
	UseToken                  bool
	ServiceURL                string
	Domain, BaseURL           string
}

// 登陆一网通
func (c *cas) Login(client *http.Client) (string, error) {
	requestURL := genRequestURL(c.ServiceURL, c.BaseURL)
	var resp *http.Response
	var err error
	if c.UseToken {
		setToken(client, c.Token, c.Domain)
		resp, err = client.Get(requestURL)
	} else {
		lt, postURL, err := getArgs(client, requestURL)
		if err != nil {
			return "", err
		}

		request := buildAuthRequest(c.Username, c.Password, lt, postURL, requestURL)

		resp, err = client.Do(request)
		if err != nil {
			return "", err
		}
	}

	body := extractBody(resp)

	_, err = isLogged(body)
	if err != nil {
		return "", err
	}

	return body, nil
}

var (
	ltExp             = regexp.MustCompile(`name="lt" value="(.+?)"`)
	postURLExp        = regexp.MustCompile(`id="loginForm" action="(.+?)"`)
	errorArgsNotFound = errors.New("页面参数不全")
)

// 获取LT和PostURL
func getArgs(client *http.Client, requestURL string) (lt, postURL string, err error) {
	req, _ := http.NewRequest("GET", requestURL, nil)
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	body := extractBody(resp)

	lt, err = matchSingle(ltExp, body)
	if err != nil {
		return lt, postURL, errorArgsNotFound
	}
	postURL, err = matchSingle(postURLExp, body)
	if err != nil {
		return lt, postURL, errorArgsNotFound
	}
	return
}

// 构造登陆请求
func buildAuthRequest(username, password, lt, postURL, reqURL string) (req *http.Request) {
	data := "rsa=" + username + password + lt +
		"&ul=" + strconv.Itoa(len(username)) +
		"&pl=" + strconv.Itoa(len(password)) +
		"&lt=" + lt +
		"&execution=e1s1" +
		"&_eventId=submit"

	req, _ = http.NewRequest("POST",
		extractDomain(reqURL)+postURL,
		strings.NewReader(data))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", reqURL)
	return
}
