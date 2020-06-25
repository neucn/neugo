package neugo

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var (
	errorNoMatched = errors.New("no matched")
)

// 使用正则在文本中匹配出所有结果，如果没有匹配到则返回error
func matchMultiple(re *regexp.Regexp, content string) ([][]string, error) {
	matched := re.FindAllStringSubmatch(content, -1)
	if len(matched) < 1 {
		return nil, errorNoMatched
	}
	return matched, nil
}

// 使用正则在文本中匹配出一个结果并取出，如果没有匹配到则返回error
func matchSingle(re *regexp.Regexp, content string) (string, error) {
	matched := re.FindAllStringSubmatch(content, -1)
	if len(matched) < 1 {
		return "", errorNoMatched
	}
	return matched[0][1], nil
}

// 读取响应体，并关闭resp.Body
func extractBody(resp *http.Response) string {
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(res)
}

var domainExp = regexp.MustCompile(`(https?://.+?)/`)

// 从URL中提取出域名
func extractDomain(url string) string {
	domain, err := matchSingle(domainExp, url)
	if err != nil {
		return "https://pass.neu.edu.cn"
	}
	return domain
}

var (
	errorCookieNotFound = errors.New("cookie中没有此项")
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
	return "", errorCookieNotFound
}
