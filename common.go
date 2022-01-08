package neugo

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// NewSession returns a *http.Client with an empty cookie jar and timeout of 6s.
func NewSession() *http.Client {
	jar, _ := cookiejar.New(nil)
	n := &http.Client{
		Timeout: 6 * time.Second,
		Jar:     jar,
	}
	return n
}

// Platform is the platform providing authentication service.
type Platform = byte

const (
	// CAS pass.neu.edu.cn
	CAS Platform = iota
	// WebVPN webvpn.neu.edu.cn
	WebVPN
)

type config struct {
	Client                    *http.Client
	Username, Password, Token string
	UseToken                  bool
	Platform                  Platform
}

// Use receives a *http.Client, and will add an empty cookie jar if the client doesn't have.
func Use(client *http.Client) AuthSelector {
	if client.Jar == nil {
		jar, _ := cookiejar.New(nil)
		client.Jar = jar
	}
	return &useCtx{config: config{Client: client}}
}

// AuthSelector determines the authentication type
type AuthSelector interface {
	// WithAuth through username and password on CAS
	WithAuth(username, password string) ActionSelector
	// WithToken through platform-dependent token
	WithToken(token string) ActionSelector
}

// ActionSelector determines the action to be performed
type ActionSelector interface {
	// Login tries to log in to specific platform.
	Login(platform Platform) error

	// DebugLogin does the same as Login except returns response text.
	DebugLogin(platform Platform) (string, error)
}

type useCtx struct {
	config
}

var _ AuthSelector = &useCtx{}
var _ ActionSelector = &useCtx{}

func (c *useCtx) WithAuth(username, password string) ActionSelector {
	c.UseToken = false
	c.Username = username
	c.Password = password
	return c
}

func (c *useCtx) WithToken(token string) ActionSelector {
	c.UseToken = true
	c.Token = token
	return c
}

func (c *useCtx) Login(platform Platform) error {
	c.Platform = platform
	_, err := login(c.config)
	return err
}

func (c *useCtx) DebugLogin(platform Platform) (string, error) {
	c.Platform = platform
	return login(c.config)
}

// About receives a *http.Client so can query some info about the session.
func About(client *http.Client) QuerySelector {
	return &aboutCtx{Client: client}
}

// QuerySelector determines the info to be queried.
type QuerySelector interface {
	// Token returns the platform-dependent token.
	//
	// If no token exists in the session, returns an empty string.
	Token(platform Platform) string
}

type aboutCtx struct {
	Client *http.Client
}

var _ QuerySelector = &aboutCtx{}

func (c *aboutCtx) Token(platform Platform) string {
	return getToken(c.Client, platform)
}
