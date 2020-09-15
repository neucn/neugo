package neugo

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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

// 获取某服务在webvpn上的url，例如
// http://219.216.96.4/eams/homeExt.action =>
// https://webvpn.neu.edu.cn/http/77726476706e69737468656265737421a2a618d275613e1e275ec7f8/eams/homeExt.action
func EncryptWebVPNUrl(url string) string {
	// protocol
	var protocol, port string
	if strings.HasPrefix(url, "https://") {
		protocol = "https"
		url = url[8:]
	} else {
		protocol = "http"
		if strings.HasPrefix(url, "http://") {
			url = url[7:]
		} else if strings.HasPrefix(url, "//") {
			url = url[2:]
		}
	}

	segments := strings.Split(strings.Split(url, "?")[0], ":")

	if len(segments) > 1 {
		length := len(segments[0])
		port = strings.Split(segments[1], "/")[0]
		url = url[:length] + url[length+len(port)+1:]
	}

	index := strings.Index(url, "/")
	key := "wrdvpnisthebest!"

	if index == -1 {
		url = encrypt(url, key, key)
	} else {
		host := url[:index]
		path := url[index:]
		url = encrypt(host, key, key) + path
	}

	if len(port) > 0 {
		url = fmt.Sprintf("https://webvpn.neu.edu.cn/%s-%s/%s", protocol, port, url)
	} else {
		url = fmt.Sprintf("https://webvpn.neu.edu.cn/%s/%s", protocol, url)
	}

	return url
}

func encrypt(url string, key string, iv string) string {
	return hex.EncodeToString(
		encryptCFB(toBytes(url), toBytes(key), toBytes(iv)),
	)
}

func toBytes(text string) []byte {
	text = url.QueryEscape(text)
	length := len(text)
	result := make([]byte, 0, len(text))
	for i := 0; i < length; {
		c := text[i]
		i++
		if c == 37 {
			b, _ := hex.DecodeString(text[i : i+2])
			result = append(result, b[0])
			i += 2
		} else {
			result = append(result, c)
		}
	}
	return result
}

func encryptCFB(origData, key, iv []byte) []byte {
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(origData))
	_iv := iv[:aes.BlockSize]
	ciphertext = append(_iv, ciphertext...)
	stream := cipher.NewCFBEncrypter(block, _iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], origData)
	return ciphertext
}
