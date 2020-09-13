package neugo

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestUse(t *testing.T) {
	a := assert.New(t)
	session := NewSession()
	s1 := Use(session).WithAuth("test", "test").On(CAS)
	ctx, ok := s1.(*useCtx)
	a.True(ok)
	a.Equal("test", ctx.Launcher.Username)
	a.Equal("test", ctx.Launcher.Password)
	a.Equal(casBaseURL, ctx.Launcher.BaseURL)
	a.Equal(casDomain, ctx.Launcher.Domain)

	emptyClient := &http.Client{}
	s2 := Use(emptyClient).WithToken("test").On(CAS)
	ctx, ok = s2.(*useCtx)
	a.True(ok)
	a.NotNil(ctx.Client.Jar)
	a.Equal("test", ctx.Launcher.Token)
	a.Equal(true, ctx.Launcher.UseToken)

	s3 := Use(session)
	s3.WithToken("test")
	s3.WithAuth("test", "test").On(WebVPN)
	ctx, ok = s3.(*useCtx)
	a.True(ok)
	a.Equal(false, ctx.Launcher.UseToken)
	a.Equal("test", ctx.Launcher.Username)
	a.Equal("test", ctx.Launcher.Password)
	a.Equal(webvpnBaseURL, ctx.Launcher.BaseURL)
	a.Equal(webvpnDomain, ctx.Launcher.Domain)
}

func TestAll(t *testing.T) {
	a := assert.New(t)
	session := NewSession()
	if !flag.Parsed() {
		flag.Parse()
	}

	argList := flag.Args()
	if len(argList) < 2 {
		a.Fail("没有指定测试账号和密码")
		return
	}
	username := argList[0]
	password := argList[1]

	err := Use(session).WithAuth(username, password).On(CAS).Login()
	if err != nil && strings.Contains(err.Error(), "timeout") {
		a.Nil(err)
	}
	err = Use(session).WithAuth(username, password).On(WebVPN).Login()
	if err != nil && strings.Contains(err.Error(), "timeout") {
		a.Nil(err)
	}

	token, err := About(session).Token(CAS)
	a.Nil(err)
	a.NotZero(len(token))
	session1 := NewSession()
	err = Use(session1).WithToken(token).On(CAS).Login()
	if err != nil && strings.Contains(err.Error(), "timeout") {
		a.Nil(err)
	}

	token, err = About(session).Token(WebVPN)
	a.Nil(err)
	a.NotZero(len(token))
	session2 := NewSession()
	err = Use(session2).WithToken(token).On(WebVPN).Login()
	if err != nil && strings.Contains(err.Error(), "timeout") {
		a.Nil(err)
	}
}
