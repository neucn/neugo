package neugo

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestHelperMatch(t *testing.T) {
	a := assert.New(t)
	type matchTestCase struct {
		Re        *regexp.Regexp
		Content   string
		ExpectErr bool
		ExpectLen int
	}

	singleTestCases := []*matchTestCase{
		{Re: regexp.MustCompile(`([a-z])`), Content: "9", ExpectErr: true, ExpectLen: 0},
		{Re: regexp.MustCompile(`([a-z])`), Content: "az", ExpectErr: false, ExpectLen: 1},
		{Re: regexp.MustCompile(`([a-z]+)`), Content: "az", ExpectErr: false, ExpectLen: 2},
	}
	for _, c := range singleTestCases {
		result, err := matchSingle(c.Re, c.Content)
		a.Equal(c.ExpectErr, err != nil)
		a.Equal(c.ExpectLen, len(result))
	}

	multipleTestCases := []*matchTestCase{
		{Re: regexp.MustCompile(`([a-z])`), Content: "9", ExpectErr: true, ExpectLen: 0},
		{Re: regexp.MustCompile(`([a-z])`), Content: "az", ExpectErr: false, ExpectLen: 2},
		{Re: regexp.MustCompile(`([a-z]+)`), Content: "az", ExpectErr: false, ExpectLen: 1},
	}
	for _, c := range multipleTestCases {
		result, err := matchMultiple(c.Re, c.Content)
		a.Equal(c.ExpectErr, err != nil)
		a.Equal(c.ExpectLen, len(result))
	}
}

func TestHelperExtractDomain(t *testing.T) {
	a := assert.New(t)
	type testCase struct{ URL, Expect string }
	testCases := []*testCase{
		{URL: "https://webvpn.neu.edu.cn/", Expect: "https://webvpn.neu.edu.cn"},
		{URL: "https://219-216-96-4.webvpn.neu.edu.cn/eams/homeExt.action", Expect: "https://219-216-96-4.webvpn.neu.edu.cn"},
		{URL: "https://webvpn.neu.edu.cn/https/77726476706e69737468656265737421e0f6528f693e6d45300d8db9d6562d/tpass/login?service=https%3A%2F%2Fportal-443.webvpn.neu.edu.cn%2Ftp_up%2F",
			Expect: "https://webvpn.neu.edu.cn"},
		{URL: "https://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fportal.neu.edu.cn%2Ftp_up%2F",
			Expect: "https://pass.neu.edu.cn"},
		{URL: "http://pass.neu.edu.cn/tpass/login?service=https%3A%2F%2Fportal.neu.edu.cn%2Ftp_up%2F", Expect: "http://pass.neu.edu.cn"},
		{URL: "//", Expect: "https://pass.neu.edu.cn"},
	}
	for _, c := range testCases {
		d := extractDomain(c.URL)
		a.Equal(c.Expect, d)
	}
}

func TestHelperCookie(t *testing.T) {
	a := assert.New(t)
	client := NewSession()
	cookie := &http.Cookie{
		Name:   "test",
		Value:  "yes",
		Path:   "/",
		Domain: "neu.test",
	}
	setCookie(client, cookie)
	u := &url.URL{
		Scheme: "https",
		Host:   "neu.test",
		Path:   "/",
	}
	_, err := getCookie(client.Jar.Cookies(u), "whatever")
	a.NotNil(err)
	v, err := getCookie(client.Jar.Cookies(u), cookie.Name)
	a.Nil(err)
	a.Equal(cookie.Value, v)
}

func TestEncryptWebVPNUrl(t *testing.T) {
	a := assert.New(t)

	testCases := map[string]string{
		"http://202.118.8.7:8991/F/29DK3KT4SV9VBRI548R8UD3MBIT991BXE4HLXENCFEGE54551T-22111?func=find-b-0": "https://webvpn.neu.edu.cn/http-8991/77726476706e69737468656265737421a2a713d27661301e2646de/F/29DK3KT4SV9VBRI548R8UD3MBIT991BXE4HLXENCFEGE54551T-22111?func=find-b-0",
		"http://219.216.96.4/eams/":  "https://webvpn.neu.edu.cn/http/77726476706e69737468656265737421a2a618d275613e1e275ec7f8/eams/",
		"https://portal.neu.edu.cn/": "https://webvpn.neu.edu.cn/https/77726476706e69737468656265737421e0f85388263c265e7b1dc7a99c406d369a/",
	}

	for origin, encrypted := range testCases {
		a.Equal(EncryptWebVPNUrl(origin), encrypted)
	}
}
