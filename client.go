package mirror

import (
	"net/http"
	"errors"
)

//type Transport struct {
//	// Source supplies the token to add to outgoing requests'
//	// Authorization headers.
//	Source TokenSource
//
//	// Base is the base RoundTripper used to make HTTP requests.
//	// If nil, http.DefaultTransport is used.
//	Base http.RoundTripper
//
//	mu     sync.Mutex                      // guards modReq
//	modReq map[*http.Request]*http.Request // original -> modified
//}
//

func NewCLient(token string) *http.Client {
	return &http.Client{
		Transport: &Transport{
			token: token,
			base:  http.DefaultTransport,
		},
	}
}

type Transport struct {
	token string
	base  http.RoundTripper
}

// RoundTrip authorizes and authenticates the request with an
// access token. If no token exists or token is expired,
// tries to refresh/fetch a new token.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(t.token) == 0 {
		return nil, errors.New("Transport's Token is empty")
	}
	req.Header.Set("Authorization", "Bearer"+" "+t.token)
	return t.base.RoundTrip(req)
}
