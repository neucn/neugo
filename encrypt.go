package neugo

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

// EncryptURLToWebVPN encrypts a service url so that it can be accessed via WebVPN.
//
// If the protocol is missed in the provided url, it will be set to http.
//
// Example:
//   "http://219.216.96.4/eams/homeExt.action"
//   will be encrypted into
//   "https://webvpn.neu.edu.cn/http/77726476706e69737468656265737421a2a618d275613e1e275ec7f8/eams/homeExt.action"
func EncryptURLToWebVPN(url string) string {
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

// Re-implement from Javascript code.
//
// Source:
//	function (text, key, iv) {
//	  const textLength = text.length
//	  text = textRightAppend(text, 'utf8')
//	  const keyBytes = utf8.toBytes(key)
//	  const ivBytes = utf8.toBytes(iv)
//	  const textBytes = utf8.toBytes(text)
//	  const aesCfb = new AesCfb(keyBytes, ivBytes, 16)
//	  const encryptBytes = aesCfb.encrypt(textBytes)
//	  return hex.fromBytes(ivBytes) + hex.fromBytes(encryptBytes).slice(0, textLength * 2)
//	}
//	function decrypt(text, key) {
//	 const textLength = (text.length - 32) / 2
//	 text = textRightAppend(text, 'hex')
//	 const keyBytes = utf8.toBytes(key)
//	   const ivBytes = hex.toBytes(text.slice(0, 32))
//	   const textBytes = hex.toBytes(text.slice(32))
//	   const aesCfb = new AesCfb(keyBytes, ivBytes, 16)
//	   const decryptBytes = aesCfb.decrypt(textBytes)
//	   return utf8.fromBytes(decryptBytes).slice(0, textLength)
//	}
//	function encryptUrl(protocol, url) {
//	   let port = ""
//	   let segments = ""
//	   if (url.substring(0, 7) === "http://") {
//		   url = url.substr(7)
//	   } else if (url.substring(0, 8) === "https://") {
//		   url = url.substr(8)
//	   }
//	   let v6 = ""
//	   const match = /\[[0-9a-fA-F:]+?\]/.exec(url)
//	   if (match) {
//		   v6 = match[0]
//		   url = url.slice(match[0].length)
//	   }
//	   segments = url.split("?")[0].split(":")
//	   if (segments.length > 1) {
//		   port = segments[1].split("/")[0]
//		   url = url.substr(0, segments[0].length) + url.substr(segments[0].length + port.length + 1)
//	   }
//	   if (protocol != "connection") {
//		   const i = url.indexOf('/')
//		   if (i === -1) {
//			   if (v6 !== "") {
//				   url = v6
//			   }
//			   url = encrypt(url, "wrdvpnisthebest!", "wrdvpnisthebest!")
//		   } else {
//			   const host = url.slice(0, i)
//			   const path = url.slice(i)
//			   if (v6 !== "") {
//				   host = v6
//			   }
//			   url = encrypt(host, "wrdvpnisthebest!", "wrdvpnisthebest!") + path
//		   }
//	   }
//	   if (port !== "") {
//		   url = "/" + protocol + "-" + port + "/" + url
//	   } else {
//		   url = "/" + protocol + "/" + url
//	   }
//	   return url
//	}
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
