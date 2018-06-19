## devfeel/connx
A simple tcp server/client go framework.

## Feature:
* 极其简单的使用代码
* 经典Server/Client模式
* 解决经典粘包问题
* Server支持主动关闭
* Client支持Send+Write\OnlySend模式
* 提供默认的Message协议，自包含版本、命令、数据，应用可直接使用

## Code:
#### server
```go
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
```
#### client
```go
func main(){
	client := connx.NewClient("127.0.0.1:7069", onConnHandler)
	go func(){
		for{
			err := client.Send(connx.RequestMessage("test client"))
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
```