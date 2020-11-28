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
