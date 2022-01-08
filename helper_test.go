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
		ExpectLen int
	}

	singleTestCases := []*matchTestCase{
		{Re: regexp.MustCompile(`([a-z])`), Content: "9", ExpectLen: 0},
		{Re: regexp.MustCompile(`([a-z])`), Content: "az", ExpectLen: 1},
		{Re: regexp.MustCompile(`([a-z]+)`), Content: "az", ExpectLen: 2},
	}
	for _, c := range singleTestCases {
		result := matchSingle(c.Re, c.Content)
		a.Equal(c.ExpectLen, len(result))
	}

	multipleTestCases := []*matchTestCase{
		{Re: regexp.MustCompile(`([a-z])`), Content: "9", ExpectLen: 0},
		{Re: regexp.MustCompile(`([a-z])`), Content: "az", ExpectLen: 2},
		{Re: regexp.MustCompile(`([a-z]+)`), Content: "az", ExpectLen: 1},
	}
	for _, c := range multipleTestCases {
		result := matchMultiple(c.Re, c.Content)
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
	c := getCookie(client.Jar.Cookies(u), "whatever")
	a.Empty(c)
	c = getCookie(client.Jar.Cookies(u), cookie.Name)
	a.Equal(cookie.Value, c)
}
