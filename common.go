package neugo

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// NewSession 获取带有 cookiejar 的 http.Client，默认 Timeout 为6秒
func NewSession() *http.Client {
	jar, _ := cookiejar.New(nil)
	n := &http.Client{
		Timeout: 6 * time.Second,
		Jar:     jar,
	}
	return n
}

// Platform 需要操作的平台
type Platform = byte

const (
	CAS    Platform = iota // 一网通 pass.neu.edu.cn
	WebVPN                 // webvpn webvpn.neu.edu.cn
)

type config struct {
	Client                    *http.Client
	Username, Password, Token string
	UseToken                  bool
	Platform                  Platform
}

// Use 接收一个 *http.Client，提供登陆动作的链式调用。如果 client 没有 cookiejar 会自动加上一个空 cookiejar.
func Use(client *http.Client) AuthSelector {
	if client.Jar == nil {
		jar, _ := cookiejar.New(nil)
		client.Jar = jar
	}
	return &useCtx{config: config{Client: client}}
}

// AuthSelector 选择鉴权方式
type AuthSelector interface {
	// WithAuth 通过一网通账号密码登陆
	WithAuth(username, password string) ActionSelector
	// WithToken 通过 Token 登陆，Token 与平台有关
	WithToken(token string) ActionSelector
}

// ActionSelector 选择要执行的动作
type ActionSelector interface {
	// Login 登陆指定平台
	Login(platform Platform) error

	// DebugLogin 返回页面内容，用于调试
	DebugLogin(platform Platform) (string, error)
}

type useCtx struct {
	config
}

var _ AuthSelector = &useCtx{}
var _ ActionSelector = &useCtx{}

// 使用账号密码
func (c *useCtx) WithAuth(username, password string) ActionSelector {
	c.UseToken = false
	c.Username = username
	c.Password = password
	return c
}

// 使用Token
func (c *useCtx) WithToken(token string) ActionSelector {
	c.UseToken = true
	c.Token = token
	return c
}

// 登陆
func (c *useCtx) Login(platform Platform) error {
	c.Platform = platform
	_, err := login(c.config)
	return err
}

// private 仅调试用，可通过类型强制转换使用
func (c *useCtx) DebugLogin(platform Platform) (string, error) {
	c.Platform = platform
	return login(c.config)
}

// About 接收一个 *http.Client，提供查询相关信息的链式调用
func About(client *http.Client) QuerySelector {
	return &aboutCtx{Client: client}
}

// QuerySelector 选择要查询的内容
type QuerySelector interface {
	Token(platform Platform) (string, error)
}

type aboutCtx struct {
	Client *http.Client
}

var _ QuerySelector = &aboutCtx{}

func (c *aboutCtx) Token(platform Platform) (string, error) {
	return getToken(c.Client, platform)
}
