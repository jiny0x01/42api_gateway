package main

import (
	"encoding/json"
	"github.com/jinykim0x80/42api_gateway/internal"
	"github.com/jinykim0x80/42api_gateway/internal/api/token"
	"github.com/jinykim0x80/42api_gateway/internal/user"

	//	"go.mongodb.org/mongo-driver/mongo"
	//	"context"
	//	"go.mongodb.org/mongo-driver/mongo/options"
	//	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net"
	"net/rpc"
)

func Server() {
	rpc.Register(new(user.Users))
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	log.Println("RCP Listen to Ready")
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()

		go rpc.ServeConn(conn)
	}
}

func Init() {
	log.Println("Init API Server")
	t := token.Get()
	if err := t.Refresh(); err != nil {
		return
	}
	log.Println("Token Set Done")
	token.SetHeader()
	log.Println("Header Set Done")

	var users user.Users
	if err := util.ReadJSON("user.json", &users); err == nil {
		log.Println(users)
		return
	}
	if len(users.User) == 0 {
		log.Println("Getting user list")
		u, err := user.GetAll()
		if err != nil {
			log.Println(err)
			return
		}
		bytes, _ := json.Marshal(u)
		err = util.WriteJSON("user.json", bytes)
		if err != nil {
			log.Println(err)
			return
		}
		users = u
	}
	log.Println("WriteResult")
	log.Println(users)

	/*
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		collection := client.Database("test").Collection("users")
		insertManyResult, err := collection.InsertMany(context.TODO(), u)
		if err != nil {
			log.Println(err)
			return
		}
	*/

}

func main() {
	Init()
	//	Server()
}
