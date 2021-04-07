package user

import (
	"encoding/json"
	"github.com/jinykim0x80/42api_gateway/internal/api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	URL   string `json:"url"`
}

type Users []string
type ValidUsers struct {
	users Users
}

func IsValidUser(user_id string) error {
	req, err := http.NewRequest("GET", api.Endpoint+"/users/"+user_id+"/campus_users", nil)
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header = *api.GetHeader()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	bytes, _ := ioutil.ReadAll(res.Body)
	// response로 만료되었다고 메시지가 오면?
	log.Printf("API Result: %s", string(bytes))
	return nil
}

func (users *Users) GetValidUsers(u Users, vu *ValidUsers) error {
	for _, user := range u {
		log.Printf("user: %s\n", user)
		if err := IsValidUser(user); err != nil {
			vu.users = append(vu.users, user)
			log.Printf("vu: %v\n", vu)
		}
	}
	return nil
}

func GetAll() {
	var err error
	// 고루틴으로 처리하면 1초당 API request 제한걸림
	builder := strings.Builder{}
	builder.WriteString(api.Endpoint)
	builder.WriteString("/campus/29/users?campus_id=29&range[id]68848,200000&page[size]=100&page[number]=")
	base := builder.String()
	for pn := 1; ; pn++ {
		builder.Reset()
		builder.WriteString(base)
		builder.WriteString(strconv.Itoa(pn))
		req, err := http.NewRequest("GET", builder.String(), nil)
		if err != nil {
			log.Println("fail to request")
			break
		}
		req.Header = *api.GetHeader()

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			break
		}
		defer res.Body.Close()

		var u []User
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&u)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 500) // 0.5sec
	}
	log.Println(err)
}

func Upsert() {
	// insert user to DB
}
