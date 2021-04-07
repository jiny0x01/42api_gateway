package api

import (
	"github.com/jinykim0x80/42api_gateway/internal/api/token"
	"net/http"
)

const (
	Endpoint = "https://api.intra.42.fr/v2"
)

var hdr http.Header

func GetHeader() *http.Header {
	return &hdr
}

func SetHeader() {
	t := token.Get()
	hdr = http.Header{}
	hdr.Add("Authorization", t.TokenType+" "+t.AccessToken)
}

/*
type OAuth struct {
	Client_id     string `json:"client_id"`
	Redirect_uri  string `json:"redirect_uri"`
	Scope         string `json:"scope"`
	State         string `json:"state"`
	Response_type string `json:"response_type"`
}
*/
