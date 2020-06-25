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
		{Content: "<title>", Expect: errorNoTitleFound},
		{Content: "<title>智慧东大</title>", Expect: errorWrongSetting},
		{Content: "<title>系统提示</title>", Expect: errorAccountBanned},
		{Content: "<title>智慧东大--统一身份认证</title>", Expect: errorAuthFailed},
		{Content: "<title><||__&^$></title>", Expect: nil},
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
		{Service: "https://portal.neu.edu.cn/tp_up/", VPN: true, Expect: "https://pass-443.webvpn.neu.edu.cn/tpass/login?service=https%3A%2F%2Fportal.neu.edu.cn%2Ftp_up%2F"},
		{Service: "https://219-216-96-4.webvpn.neu.edu.cn/eams/homeExt.action", VPN: true,
			Expect: "https://pass-443.webvpn.neu.edu.cn/tpass/login?service=https%3A%2F%2F219-216-96-4.webvpn.neu.edu.cn%2Feams%2FhomeExt.action"},
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
	_, err := getToken(client, webvpnDomain)
	a.NotNil(err)
	token := "test-token"
	setToken(client, token, webvpnDomain)
	result, err := getToken(client, webvpnDomain)
	a.Nil(err)
	a.Equal(token, result)

	setToken(client, token, casDomain)
	result, err = getToken(client, casDomain)
	a.Nil(err)
	a.Equal(token, result)
}
