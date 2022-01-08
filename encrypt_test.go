package neugo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptWebVPNUrl(t *testing.T) {
	a := assert.New(t)

	testCases := map[string]string{
		"http://202.118.8.7:8991/F/29DK3KT4SV9VBRI548R8UD3MBIT991BXE4HLXENCFEGE54551T-22111?func=find-b-0": "https://webvpn.neu.edu.cn/http-8991/77726476706e69737468656265737421a2a713d27661301e2646de/F/29DK3KT4SV9VBRI548R8UD3MBIT991BXE4HLXENCFEGE54551T-22111?func=find-b-0",
		"http://219.216.96.4/eams/":  "https://webvpn.neu.edu.cn/http/77726476706e69737468656265737421a2a618d275613e1e275ec7f8/eams/",
		"https://portal.neu.edu.cn/": "https://webvpn.neu.edu.cn/https/77726476706e69737468656265737421e0f85388263c265e7b1dc7a99c406d369a/",
		"//ipgw.neu.edu.cn":          "https://webvpn.neu.edu.cn/http/77726476706e69737468656265737421f9e7468b693e6d45300d8db9d6562d",
	}

	for origin, encrypted := range testCases {
		a.Equal(encrypted, EncryptURLToWebVPN(origin))
	}
}
