package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	ApiEndpoint   = "https://api.intra.42.fr/v2"
	TokenEndpoint = "https://api.intra.42.fr/oauth/token"
	CallBackURL   = "/callback"
)

type API_Access struct {
	Client_id     string `json:"client_id"`
	Redirect_uri  string `json:"redirect_uri"`
	Scope         string `json:"scope"`
	State         string `json:"state"`
	Response_type string `json:"response_type"`
}

var token Token

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

func OpenJsonFile(file string) (*AccessKey, error) {
	data, err := os.Open(file)
	if err != nil {
		log.Printf("Fail to Open json file: %v", err)
		return nil, err
	}
	var access_key AccessKey
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &access_key)
	return &access_key, nil
}

func InitToken(access_key_file string) error {
	access_key, err := OpenJsonFile(access_key_file)
	if err != nil {
		return err
	}
	credential := Credential{"client_credentials", access_key.ClientId, access_key.ClientSecret}
	//
	marshal_credential, _ := json.Marshal(credential)
	res, err := http.Post(TokenEndpoint, "application/json", bytes.NewBuffer(marshal_credential))
	if err != nil && res.StatusCode != http.StatusOK {
		log.Printf("couldn't create HTTP request: %v", err)
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&token)
	log.Println(token)
	if err != nil {
		log.Printf("Fail to Decode token\n")
		return err
	}
	return nil
}

func IsValidUser(user_id string) error {
	req, err := http.NewRequest("GET", ApiEndpoint+"/users/"+user_id+"/campus_users", nil)
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("Authorization", token.TokenType+" "+token.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	bytes, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("API Result: %s", string(bytes))
	fmt.Printf("API Result: %s", string(bytes))
	return nil
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	if err := InitToken("api_access.json"); err != nil {
		log.Fatalln("Fail to getting token. Check api_access")
	}
	IsValidUser("jinykim")
	http.HandleFunc("/", MainHandler)
	http.HandleFunc("/api", ApiHandler)
	http.ListenAndServe(":8000", nil)
}
