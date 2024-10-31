package neugo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptWebVPNUrl(t *testing.T) {
	a := assert.New(t)

	testCases := map[string]string{
		"http://219.216.96.4/eams/":  "https://webvpn.neu.edu.cn/http/62304135386136393339346365373340e2b0fd71d8941093ab4e2527/eams/",
		"https://portal.neu.edu.cn/": "https://webvpn.neu.edu.cn/https/62304135386136393339346365373340a0eeb62b8bc908d3f70d257682109b84a2/",
		"//ipgw.neu.edu.cn":          "https://webvpn.neu.edu.cn/http/62304135386136393339346365373340b9f1a328c4cb43c8bc1d6f66c806db",
	}

	for origin, encrypted := range testCases {
		a.Equal(encrypted, EncryptURLToWebVPN(origin))
	}
}
