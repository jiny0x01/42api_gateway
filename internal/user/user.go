package user

import (
	"encoding/json"
	"github.com/jinykim0x80/42api_gateway/internal/api"
	"github.com/jinykim0x80/42api_gateway/internal/api/token"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
}
type Users []User

var users Users

func Get() Users {
	return users
}

func Set(u Users) {
	users = u
}

func IsValid(name string) bool {
	users := Get()
	if len(users) == 0 {
		return false
	}
	for _, user := range users {
		if user.Name == name {
			return true
		}
	}
	return false
}

func (u *User) GetValid(user []string, vu *[]string) error {
	for i := range user {
		if IsValid(user[i]) {
			*vu = append(*vu, user[i])
		}
	}
	return nil
}

func Load() error {
	const startID = 68848
	var err error
	// 고루틴으로 처리하면 1초당 API request 제한걸림
	builder := strings.Builder{}
	builder.WriteString(api.Endpoint)
	builder.WriteString("/campus/29/users?campus_id=29&sort=id&page[size]=100&page[number]=")
	base := builder.String()
	for pn := 1; ; pn++ {
		builder.Reset()
		builder.WriteString(base)
		builder.WriteString(strconv.Itoa(pn))
		req, err := http.NewRequest("GET", builder.String(), nil)
		if err != nil {
			log.Println("fail to request")
			return nil
		}
		req.Header = *token.GetHeader()

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			break
		}
		defer res.Body.Close()

		type UserInfo struct {
			ID    int    `json:"id"`
			Login string `json:"login"`
			URL   int    `json:"url"`
		}
		var userInfo []UserInfo
		decoder := json.NewDecoder(res.Body)
		decoder.Decode(&userInfo)
		if err != nil || len(userInfo) == 0 {
			break
		}
		for _, u := range userInfo {
			if u.ID >= startID {
				users = append(users, User{u.ID, u.Login})
			}
		}
		time.Sleep(time.Millisecond * 500) // 0.5sec
	}
	/*
		for _, u := range user {
			// 68848부터 유효한 사용자
			if u.ID >= 68848 {
				users.User = append(users.User, u)
			}
		}
	*/
	return err
}

func Upsert() {
	// insert user to DB
	//	u, err := GetAll()
}
