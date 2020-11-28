package neugo

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

// 获取某服务在 WebVPN 上的 url，若 url 中协议为空则默认为 http
//
// 例如 http://219.216.96.4/eams/homeExt.action =>
// https://webvpn.neu.edu.cn/http/77726476706e69737468656265737421a2a618d275613e1e275ec7f8/eams/homeExt.action
func EncryptToWebVPN(url string) string {
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
		return fmt.Sprintf("https://webvpn.neu.edu.cn/%s-%s/%s", protocol, port, url)
	}

	return fmt.Sprintf("https://webvpn.neu.edu.cn/%s/%s", protocol, url)
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
