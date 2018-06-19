package connx

import (
	"net"
	"fmt"
)

type Client struct {
	remoteAddr string
	conn *Connection
	handler OnConnHandle
}


// handleConn loop handle new message from conn
func (c *Client) handleConn() {
	defer func(){
		c.conn.Close()
		c.conn = nil
	}()
	for {
		errRead := c.conn.readMessage()
		if errRead == nil {
			if c.handler != nil{
				errHandler := c.handler(c.conn)
				if errHandler != nil{
					connLogger.Error(fmt.Sprintf("Client.handler error ClientID:%v err:%v", c.conn.ConnIndex, errHandler))
				}
			}
		} else {
			//discard current request
			connLogger.Error(fmt.Sprintf("Client.handleConn error ClientID:%v err:%v", c.conn.ConnIndex, errRead))
			break
		}
	}
}

// SetOnConnHandle set handler on new conn receive
func (c *Client) SetOnConnHandle(handler OnConnHandle){
	c.handler = handler
}

func (c *Client) Write(head *HeadInfo, body []byte) error {
	if err:=c.Dial();err!= nil{
		return err
	}
	head.DataLen = uint64(len(body))
	_, err := c.conn.WriteMerge(head.GetBytes(), body)
	return err
}

// Dial dial remote addr
func (c *Client) Dial() error{
	if c.conn == nil{
		conn, err := net.Dial("tcp", c.remoteAddr)
		if err != nil{
			return err
		}
		c.conn = NewConnction(conn)
		//if handler is nil, no start loop handle receive message
		if c.handler != nil {
			go c.handleConn()
		}
	}
	return nil
}

// Send send data to remote addr
func (c *Client) Send(msg *Message) error {
	if err:=c.Dial();err!= nil{
		return err
	}
	return c.conn.SendMessage(msg)
}

// Close close connection
func (c *Client) Close(){
	c.conn.Close()
}

// NewClient get new client which can send and read message with remoteAddr and OnConnHandle
func NewClient(remoteAddr string, handler OnConnHandle) *Client {
	c := &Client{
		remoteAddr:remoteAddr,
		handler:handler,
	}
	return c
}

// NewRequestClient get new client which only send message with remoteAddr
func NewRequestOnlyClient(remoteAddr string) *Client {
	c := &Client{
		remoteAddr:remoteAddr,
	}
	return c
}