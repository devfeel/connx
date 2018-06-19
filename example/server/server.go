package main

import (
	"github.com/devfeel/connx"
	"fmt"
	"encoding/gob"
)

type LoginInfo struct{
	UserName string
	Password string
	LoginFrom string
}

func init(){
	gob.Register(LoginInfo{})
}


func main(){
	server, err := connx.NewServer("127.0.0.1:7069", onConnHandler)
	if err != nil{
		fmt.Println("GetNewServer error", err)
		return
	}
	fmt.Println("GetNewServer begin listen")
	server.Start()
}


func onConnHandler(conn *connx.Connection) error{
	msg, err := conn.ParseMessage()
	fmt.Println(msg, err)
	return nil
}