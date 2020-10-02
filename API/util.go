package API

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrInvalidLoginResponse = errors.New("Invalid login response from server")
)

func GetClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		},
	}
}

func Contains(array []string, str string) bool {
	for _, a := range array {
		if a == str {
			return true
		}
	}
	return false
}
