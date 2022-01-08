package neugo

import (
	"net/http"
	"net/url"
)

// INTERNAL. Add cookie into the client. Note that the scheme is forced to be `https`.
func setCookie(client *http.Client, cookie *http.Cookie) {
	client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   cookie.Domain,
		Path:   cookie.Path,
	}, []*http.Cookie{cookie})
}

// INTERNAL. Get the cookie of specific name from the client.
func getCookie(cookies []*http.Cookie, name string) string {
	for _, item := range cookies {
		if item.Name == name {
			return item.Value
		}
	}
	return ""
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

// Set platform-dependent token.
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

// Get token from the client. An empty string will be returned if not exists.
func getToken(client *http.Client, platform Platform) string {
	if platform == WebVPN {
		cookies := client.Jar.Cookies(webvpnCookieUrl)
		return getCookie(cookies, "wengine_vpn_ticketwebvpn_neu_edu_cn")
	}

	cookies := client.Jar.Cookies(casCookieUrl)
	return getCookie(cookies, "CASTGC")
}
