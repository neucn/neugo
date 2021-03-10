package neugo

import (
	"errors"
	"net/http"
	"net/url"
)

var (
	ErrorCookieNotFound = errors.New("no such entry found in the cookies")
)

// 设置cookie
func setCookie(client *http.Client, cookie *http.Cookie) {
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   cookie.Domain,
		Path:   cookie.Path,
	}, []*http.Cookie{cookie})
}

// 提取出cookie中指定名称的值
func getCookie(cookies []*http.Cookie, name string) (string, error) {
	for _, item := range cookies {
		if item.Name == name {
			return item.Value, nil
		}
	}
	return "", ErrorCookieNotFound
}

var (
	webvpnCookieUrl = &url.URL{
		Scheme: "https",
		Path:   "/",
		Host:   webvpnCookieDomain,
	}

	casCookieUrl = &url.URL{
		Scheme: "https",
		Path:   "/tpass/",
		Host:   casCookieDomain,
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
		cookie.Domain = casCookieDomain
		cookie.Name = "CASTGC"
		cookie.Path = "/tpass/"
	}
	setCookie(client, cookie)
}

func getToken(client *http.Client, platform Platform) (string, error) {
	if platform == WebVPN {
		cookies := client.Jar.Cookies(webvpnCookieUrl)
		return getCookie(cookies, "wengine_vpn_ticketwebvpn_neu_edu_cn")
	}

	cookies := client.Jar.Cookies(casCookieUrl)
	return getCookie(cookies, "CASTGC")
}
