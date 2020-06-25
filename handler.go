package neugo

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

var (
	titleExp = regexp.MustCompile(`<title>(.+?)</title>`)
)

var (
	errorNoTitleFound  = errors.New("无法识别页面")
	errorAccountBanned = errors.New("账号受限")
	errorWrongSetting  = errors.New("一网通设置有误")
	errorAuthFailed    = errors.New("账号密码错误或Token失效")
)

// 【账号登陆】根据title判断是否登陆成功，若不成功则结束并报错
func isLogged(body string) (bool, error) {
	title, err := matchSingle(titleExp, body)
	if err != nil {
		return false, errorNoTitleFound
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

// 生成登陆所要请求的URL
func genRequestURL(service, baseURL string) string {
	return baseURL + url.QueryEscape(service)
}

func setToken(client *http.Client, token, domain string) {
	cookie := &http.Cookie{
		Name:   "CASTGC",
		Value:  token,
		Path:   "/tpass/",
		Domain: domain,
	}
	setCookie(client, cookie)
}

func getToken(client *http.Client, domain string) (string, error) {
	var cookies []*http.Cookie
	dm := &url.URL{
		Scheme: "https",
		Path:   "/tpass/",
		Host:   domain,
	}

	cookies = client.Jar.Cookies(dm)

	return getCookie(cookies, "CASTGC")
}
