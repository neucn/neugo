package neugo

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	webvpnLoginURL = "https://webvpn.neu.edu.cn/https/77726476706e69737468656265737421e0f6528f693e6d45300d8db9d6562d/tpass/login"
	casLoginURL    = "https://pass.neu.edu.cn/tpass/login"

	webvpnCookieDomain = ".webvpn.neu.edu.cn"
	casCookieDomain    = "pass.neu.edu.cn"
)

// 登陆一网通
func login(c config) (string, error) {
	var resp *http.Response
	var err error
	var loginURL string
	if c.Platform == CAS {
		loginURL = casLoginURL
	} else {
		loginURL = webvpnLoginURL
	}

	if c.UseToken {
		setToken(c.Client, c.Token, c.Platform)
		resp, err = c.Client.Get(loginURL)
	} else {
		lt, err := getLT(c.Client, loginURL)
		if err != nil {
			return "", err
		}
		request := buildAuthRequest(c.Username, c.Password, lt, loginURL)

		resp, err = c.Client.Do(request)
	}
	if err != nil {
		return "", err
	}

	body := extractBody(resp)

	_, err = isLogged(body)

	return body, err
}

var (
	ltExp             = regexp.MustCompile(`name="lt" value="(.+?)"`)
	postURLExp        = regexp.MustCompile(`id="loginForm" action="(.+?)"`)
	errorArgsNotFound = errors.New("页面参数不全")
)

// 获取 LT
func getLT(client *http.Client, requestURL string) (lt string, err error) {
	req, _ := http.NewRequest("GET", requestURL, nil)
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	body := extractBody(resp)

	lt, err = matchSingle(ltExp, body)
	if err != nil {
		return lt, errorArgsNotFound
	}
	return
}

// 构造登陆请求
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

var (
	titleExp = regexp.MustCompile(`<title>(.+?)</title>`)
)

var (
	errorAccountBanned = errors.New("账号受限")
	errorWrongSetting  = errors.New("一网通设置有误")
	errorAuthFailed    = errors.New("账号密码错误或Token失效")
)

// 根据title判断是否登陆成功，若不成功则结束并报错
func isLogged(body string) (bool, error) {
	title, err := matchSingle(titleExp, body)
	if err != nil {
		return true, nil
	}

	switch title {
	case "智慧东大--统一身份认证":
		return false, errorAuthFailed
	case "智慧东大":
		return false, errorWrongSetting
	case "系统提示":
		return false, errorAccountBanned
	}
	return true, nil
}
