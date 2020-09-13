package neugo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsLogged(t *testing.T) {
	a := assert.New(t)
	type testCase struct {
		Content string
		Expect  error
	}
	testCases := []*testCase{
		{Content: "<title>智慧东大--统一身份认证</title>", Expect: errorAuthFailed},
		{Content: "<title>智慧东大</title>", Expect: errorWrongSetting},
		{Content: "<title>系统提示</title>", Expect: errorAccountBanned},
		{Content: "<title><||__&^$></title>", Expect: nil},
		{Content: "<>", Expect: nil},
	}

	for _, c := range testCases {
		_, err := isLogged(c.Content)
		a.Equal(c.Expect, err)
	}
}

func TestGenRequestURL(t *testing.T) {
	a := assert.New(t)
	type testCase struct {
		Service string
		VPN     bool
		Expect  string
	}
	testCases := []*testCase{
		{Service: "https://portal.neu.edu.cn/tp_up/", VPN: false, Expect: "https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fportal.neu.edu.cn%2Ftp_up%2F"},
		{Service: "https://portal.neu.edu.cn/tp_up/", VPN: true, Expect: "https://webvpn.neu.edu.cn/https/77726476706e69737468656265737421e0f6528f693e6d45300d8db9d6562d/tpass/login?service=https%3A%2F%2Fportal.neu.edu.cn%2Ftp_up%2F"},
		{Service: "https://219-216-96-4.webvpn.neu.edu.cn/eams/homeExt.action", VPN: true,
			Expect: "https://webvpn.neu.edu.cn/https/77726476706e69737468656265737421e0f6528f693e6d45300d8db9d6562d/tpass/login?service=https%3A%2F%2F219-216-96-4.webvpn.neu.edu.cn%2Feams%2FhomeExt.action"},
		{Service: "http://219.216.96.4/eams/homeExt.action", VPN: false,
			Expect: "https://pass.neu.edu.cn/tpass/login?service=http%3A%2F%2F219.216.96.4%2Feams%2FhomeExt.action"},
	}
	for _, c := range testCases {
		var bu string
		if c.VPN {
			bu = webvpnBaseURL
		} else {
			bu = casBaseURL
		}
		result := genRequestURL(c.Service, bu)
		a.Equal(c.Expect, result)
	}
}

func TestHandlerToken(t *testing.T) {
	a := assert.New(t)
	client := NewSession()
	_, err := getToken(client, WebVPN)
	a.NotNil(err)
	token := "test-token"
	setToken(client, token, WebVPN)
	result, err := getToken(client, WebVPN)
	a.Nil(err)
	a.Equal(token, result)

	setToken(client, token, CAS)
	result, err = getToken(client, CAS)
	a.Nil(err)
	a.Equal(token, result)
}

func TestEncryptWebVPNUrl(t *testing.T) {
	a := assert.New(t)

	testCases := map[string]string{
		"http://202.118.8.7:8991/F/29DK3KT4SV9VBRI548R8UD3MBIT991BXE4HLXENCFEGE54551T-22111?func=find-b-0": "/http-8991/77726476706e69737468656265737421a2a713d27661301e2646de/F/29DK3KT4SV9VBRI548R8UD3MBIT991BXE4HLXENCFEGE54551T-22111?func=find-b-0",
		"http://219.216.96.4/eams/":  "/http/77726476706e69737468656265737421a2a618d275613e1e275ec7f8/eams/",
		"https://portal.neu.edu.cn/": "/https/77726476706e69737468656265737421e0f85388263c265e7b1dc7a99c406d369a/",
	}

	for origin, encrypted := range testCases {
		a.Equal(EncryptWebVPNUrl(origin), encrypted)
	}
}
