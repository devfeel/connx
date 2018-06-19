package main

import (
	"github.com/devfeel/connx"
	"fmt"
	"encoding/gob"
	"time"
)

type LoginInfo struct{
	UserName string
	Password string
	LoginFrom string
}

func main(){
	gob.Register(LoginInfo{})
	client := connx.NewClient("127.0.0.1:7069", onConnHandler)
	//client := connx.NewRequestOnlyClient("127.0.0.1:7020")

	login:= new(LoginInfo)
	login.UserName = "user"
	login.Password = "111111"
	login.LoginFrom = "test client"

	go func(){
		for{
			err := client.Send(connx.RequestMessage(login))
			if err != nil{
				fmt.Println("Send login message failed", err)
			}else{
				fmt.Println("Send login message success")
			}
			time.Sleep(time.Second*10)
		}
	}()

	for{
		select{}
	}
}

func onConnHandler(conn *connx.Connection) error{
	msg, err := conn.ParseMessage()
	fmt.Println(msg, err)
	return nil
}

