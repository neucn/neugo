package neugo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleToken(t *testing.T) {
	a := assert.New(t)
	client := NewSession()
	{
		result := getToken(client, WebVPN)
		a.Empty(result)
	}
	{
		token := "test-token"
		setToken(client, token, WebVPN)
		result := getToken(client, WebVPN)
		a.Equal(token, result)
	}

	{
		token := "test-token"
		setToken(client, token, CAS)
		result := getToken(client, CAS)
		a.Equal(token, result)
	}
}
