package neugo

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestGetArgs(t *testing.T) {
	a := assert.New(t)
	handler := http.NewServeMux()
	handler.HandleFunc("/success", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<form id="loginForm" action="/tpass/login?service=https%3A%2F%2Fportal.neu.edu.cn%2Ftp_up%2F" method="post">
				    <input type="hidden" id="rsa" name="rsa"/>
			        <input type="hidden" id="ul" name="ul"/>
			        <input type="hidden" id="pl" name="pl"/>
			        <input type="hidden" id="lt" name="lt" value="LT-324784-5WKhfINLQf4HWzozfafzSnEguyQ6Ox-tpass" />
			        <input type="hidden" name="execution" value="e3s1" />
			        <input type="hidden" name="_eventId" value="submit" />`))
	})
	handler.HandleFunc("/fail", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<form id="loginForm" action="/tpass/login?service=https%3A%2F%2Fportal.neu.edu.cn%2Ftp_up%2F" method="post">
				    <input type="hidden" id="rsa" name="rsa"/>
			        <input type="hidden" id="ul" name="ul"/>
			        <input type="hidden" id="pl" name="pl"/>
			        <input type="hidden" id="lt" name="lt"/>
			        <input type="hidden" name="execution" value="e3s1" />
			        <input type="hidden" name="_eventId" value="submit" />`))
	})
	srv := httptest.NewServer(handler)
	client := NewSession()
	lt, err := getLT(client, srv.URL+"/success")
	a.Nil(err)
	a.Equal("LT-324784-5WKhfINLQf4HWzozfafzSnEguyQ6Ox-tpass", lt)

	_, err = getLT(client, srv.URL+"/fail")
	a.NotNil(err)

	// meaningless
	_, err = getLT(client, "")
	a.NotNil(err)
}

func TestBuildAuthRequest(t *testing.T) {
	a := assert.New(t)
	r := buildAuthRequest("test", "test", "test",
		"https://pass.neu.edu.cn/tpass/login?service=http%3A%2F%2F219.216.96.4%2Feams%2FhomeExt.action")
	res, _ := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()
	a.Equal("rsa=testtesttest&ul=4&pl=4&lt=test&execution=e1s1&_eventId=submit",
		string(res))
	a.Equal("https://pass.neu.edu.cn/tpass/login?service=http%3A%2F%2F219.216.96.4%2Feams%2FhomeExt.action",
		r.URL.String())
}
