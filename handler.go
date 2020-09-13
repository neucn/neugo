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
	errorAccountBanned = errors.New("账号受限")
	errorWrongSetting  = errors.New("一网通设置有误")
	errorAuthFailed    = errors.New("账号密码错误或Token失效")
)

// 【账号登陆】根据title判断是否登陆成功，若不成功则结束并报错
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

// 生成登陆所要请求的URL
func genRequestURL(service, baseURL string) string {
	return baseURL + url.QueryEscape(service)
}

var (
	webVpnCookieUrl = &url.URL{
		Scheme: "https",
		Path:   "/",
		Host:   webvpnCookieDomain,
	}

	casCookieUrl = &url.URL{
		Scheme: "https",
		Path:   "/tpass/",
		Host:   casDomain,
	}
)

func setToken(client *http.Client, token string, platform Platform) {
	cookie := &http.Cookie{
		Value: token,
	}

	if platform == WebVPN {
		cookie.Domain = webvpnCookieDomain
		cookie.Name = "wengine_vpn_ticketwebvpn_neu_edu_cn"
		cookie.Path = "/"
	} else {
		cookie.Domain = casDomain
		cookie.Name = "CASTGC"
		cookie.Path = "/tpass/"
	}
	setCookie(client, cookie)
}

func getToken(client *http.Client, platform Platform) (string, error) {
	if platform == WebVPN {
		cookies := client.Jar.Cookies(webVpnCookieUrl)
		return getCookie(cookies, "wengine_vpn_ticketwebvpn_neu_edu_cn")
	}

	cookies := client.Jar.Cookies(casCookieUrl)
	return getCookie(cookies, "CASTGC")
}
