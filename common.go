package neugo

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// NewSession 获取带有cookiejar的http.Client，默认Timeout为3秒
func NewSession() *http.Client {
	n := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func() func(req *http.Request, via []*http.Request) error {
			redirects := 0
			return func(req *http.Request, via []*http.Request) error {
				if redirects > 16 {
					return errors.New("stopped after 16 redirects")
				}
				redirects++
				return nil
			}
		}(),
	}
	jar, _ := cookiejar.New(nil)
	// 绑定session
	n.Jar = jar
	return n
}

// Platform 需要操作的平台
type Platform = byte

const (
	CAS    Platform = iota // 一网通 pass.neu.edu.cn
	WebVPN                 // webvpn webvpn.neu.edu.cn
)

// Use 接收一个 *http.Client，提供登陆动作的链式调用。如果 client 没有 cookiejar 会自动加上一个空 cookiejar.
func Use(client *http.Client) AuthSelector {
	if client.Jar == nil {
		jar, _ := cookiejar.New(nil)
		// 绑定session
		client.Jar = jar
	}
	return &useCtx{Client: client, Launcher: &launcher{}}
}

// AuthSelector 选择鉴权方式
type AuthSelector interface {
	WithAuth(username, password string) PlatformSelector
	WithToken(token string) PlatformSelector
}

// PlatformSelector 选择平台
type PlatformSelector interface {
	On(platform Platform) ActionSelector
}

// ActionSelector 选择要执行的动作
type ActionSelector interface {
	Login() error
	LoginService(url string) (string, error)
}

type useCtx struct {
	// 请求客户端
	Client *http.Client

	Launcher *launcher
}

var _ AuthSelector = &useCtx{}
var _ PlatformSelector = &useCtx{}
var _ ActionSelector = &useCtx{}

// 选择一网通平台或 Webvpn平台
func (c *useCtx) On(platform Platform) ActionSelector {
	c.Launcher.Platform = platform
	if platform == WebVPN {
		c.Launcher.Domain = webvpnDomain
		c.Launcher.BaseURL = webvpnBaseURL
	} else {
		c.Launcher.Domain = casDomain
		c.Launcher.BaseURL = casBaseURL
	}
	return c
}

// 使用账号密码
func (c *useCtx) WithAuth(username, password string) PlatformSelector {
	c.Launcher.UseToken = false
	c.Launcher.Username = username
	c.Launcher.Password = password
	return c
}

// 使用Token
func (c *useCtx) WithToken(token string) PlatformSelector {
	c.Launcher.UseToken = true
	c.Launcher.Token = token
	return c
}

// 登陆
func (c *useCtx) Login() error {
	_, err := c.LoginService(portalURL)
	return err
}

// 登陆指定服务，url需要是服务的完整地址，例如
// https://219.216.96.4/eams/homeExt.action
// 返回页面内容，如果登陆失败会返回error
func (c *useCtx) LoginService(url string) (string, error) {
	c.Launcher.ServiceURL = url
	return c.Launcher.Login(c.Client)
}

// TODO 查询
// About 接收一个 *http.Client，提供查询相关信息的链式调用
func About(client *http.Client) QuerySelector {
	return &aboutCtx{Client: client}
}

// QuerySelector 选择要查询的内容
type QuerySelector interface {
	Token(platform Platform) (string, error)
	Info(platform Platform) (*PersonalInfo, error)
}

type aboutCtx struct {
	Client *http.Client
}

var _ QuerySelector = &aboutCtx{}

func (c *aboutCtx) Token(platform Platform) (string, error) {
	return getToken(c.Client, platform)
}

func (c *aboutCtx) Info(platform Platform) (*PersonalInfo, error) {
	panic("not implemented yet")
}
