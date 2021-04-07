package token

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	Endpoint = "https://api.intra.42.fr/oauth/token"
)

type Credential struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AccessInfo struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	CreatedAt   int    `json:"created_at"`
}

var t Token

func OpenFile(file string) (*AccessInfo, error) {
	data, err := os.Open(file)
	if err != nil {
		log.Printf("Fail to Open json file: %v", err)
		return nil, err
	}
	var access_info AccessInfo
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &access_info)
	return &access_info, nil
}

func Get() *Token {
	return &t
}

func (t *Token) Verify(file string) error {
	access, err := OpenFile(file)
	if err != nil {
		return err
	}

	cred := Credential{"client_credentials", access.ClientID, access.ClientSecret}
	mcred, _ := json.Marshal(cred)
	res, err := http.Post(Endpoint, "application/json", bytes.NewBuffer(mcred))
	if err != nil && res.StatusCode != http.StatusOK {
		log.Printf("couldn't create HTTP request: %v", err)
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&t)
	if err != nil {
		log.Printf("Fail to Decode token\n")
		return err
	}
	return nil
}

func (t *Token) Refresh() error {
	if err := t.Verify("api_access.json"); err != nil {
		log.Fatalln("Fail to getting token. Check api_access")
		return err
	}
	log.Printf("token: %v\n", t)
	return nil
}
