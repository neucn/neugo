package neugo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleToken(t *testing.T) {
	a := assert.New(t)
	client := NewSession()
	_, err := getToken(client, WebVPN)
	a.NotNil(err)
	token := "test-token"
	setToken(client, token, WebVPN)
	result, err := getToken(client, WebVPN)
	a.Nil(err)
	a.Equal(token, result)

	setToken(client, token, CAS)
	result, err = getToken(client, CAS)
	a.Nil(err)
	a.Equal(token, result)
}
