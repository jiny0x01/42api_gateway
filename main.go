package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	//	"golang.org/x/oauth2"
	//	"golang.org/x/oauth2/google"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	Authorize_EndPoint = "https://api.intra.42.fr/oauth/authorize"
	Token_Endpoint     = "https://api.intra.42.fr/oauth/token"
	CallBackURL        = "/callback"
)

type API_Access struct {
	Client_id     string `json:"client_id"`
	Redirect_uri  string `json:"redirect_uri"`
	Scope         string `json:"scope"`
	State         string `json:"state"`
	Response_type string `json:"response_type"`
}

/*
func GetOAuthConf() *oauth2.Config {
	OAuthConf := &oauth2.Config{
		ClientID:     UID,
		ClientSecret: SECRET,
		Endpoint:     google.Endpoint, // 42 endpoint로 바꿔야함
		RedirectURL:  CallBackURL,
		Scopes:       []string{"public"},
	}
	return OAuthConf
}
*/

/*
func GetLoginURL(state string) string {
	return OAuthConf.AuthCodeURL(state)
}
*/

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, _ := template.ParseFiles(name)
	tmpl.Execute(w, data)
}

func RandToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello This is MainPage\n")
}
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello This is AuthPage\n")
	//OAuthConf := GetOAuthConf()
	//log.Printf("RedirectURL: %s\n", OAuthConf.AuthCodeURL(RandToken()))
}

type Credential struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	CreatedAt   int    `json:"created_at"`
}

type AccessKey struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func OpenJsonFile(file string) (AccessKey, error) {
	data, err := os.Open(file)
	if err != nil {
		log.Printf("Fail to Open json file: %v", err)
		return AccessKey{}, err
	}
	var access_key AccessKey
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &access_key)
	return access_key, nil
}

func GetToken() {
	access_key, err := OpenJsonFile("api_access.json")
	if err != nil {
		return
	}
	credential := Credential{"client_credentials", access_key.ClientId, access_key.ClientSecret}
	//
	marshal_credential, _ := json.Marshal(credential)
	res, err := http.Post("https://api.intra.42.fr/oauth/token", "application/json", bytes.NewBuffer(marshal_credential))
	if err != nil {
		log.Printf("couldn't create HTTP request: %v", err)
		return
	}
	defer res.Body.Close()

	var token Token
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&token)
	log.Println(token)
	if err != nil {
		log.Printf("Fail to Decode token\n")
		return
	}
}

func main() {
	GetToken()
	http.HandleFunc("/", MainHandler)
	http.ListenAndServe(":8000", nil)
}
