package connx

import (
	"net"
	"sync"
	"sync/atomic"
	"bytes"
	"encoding/gob"
	"fmt"
	"errors"
)

const (
	HeadLenght = 20
)

var (
	connctionPool sync.Pool
	connectIndex uint64
	connectCreateCount uint64
)


func init() {
	connctionPool = sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&connectCreateCount, 1)
			return &Connection{lock: new(sync.RWMutex)}
		},
	}
}

type OnConnHandle func(conn *Connection) error


type Connection struct {
	ConnIndex  int64
	lock       *sync.RWMutex
	conn       net.Conn
	Head       *HeadInfo
	headBuf    []byte
	Body 	   []byte
}


// readHead read head message with HeadLenght
func (c *Connection) readHead() error {
	err := c.readSize(HeadLenght, &c.headBuf)
	if err != nil {
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Head = &HeadInfo{}
	c.Head.FromBytes(c.headBuf)
	return nil
}

// readBody read body message with head datalen
func (c *Connection) readBody(size int64) error{
	err := c.readSize(size, &c.Body)
	return err
}

// readMessage read head and body from conn, check head flag is match
func (c *Connection) readMessage() error{
	err := c.readHead()
	if err != nil {
		return err
	}
	if c.Head ==nil{
		err :=errors.New("read HeadInfo is nil")
		return err
	}
	if c.Head.Flag != HeadFlag{
		err :=errors.New(fmt.Sprintf("check head-flag not match readFlag:%v mustFlag:%v", c.Head.Flag, HeadFlag))
		return err
	}
	return c.readBody(int64(c.Head.DataLen))
}

// readSize read message with size
func (c *Connection) readSize(size int64, buf *[]byte) error {
	*buf = make([]byte, 0)
	var err error
	leftSize := size
	for {

		bufinner := make([]byte, leftSize)
		var n int
		n, err = c.conn.Read(bufinner)
		leftSize -= int64(n)
		if err == nil {
			*buf = slice_merge(*buf, bufinner)
			if leftSize <= 0 {
				//read end
				break
			}
		} else {
			break
		}
	}
	return err
}

func (c *Connection) Write(p []byte) (int, error) {
	return c.conn.Write(p)
}

func (c *Connection) WriteMerge(head []byte, body []byte) (int, error) {
	return c.conn.Write(slice_merge(head, body))
}

// GobEncode default encoding with gob mode
func (c *Connection) GobEncode(data interface{}) ([]byte, error){
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}else{
		return encoded.Bytes(), nil
	}
}

// GobDecode default decoding with gob mode
func (c *Connection) GobDecode(data []byte, dst interface{}) error{
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(dst)
	return err
}

// ParseMessage parse body to message ptr
func (c *Connection) ParseMessage()(*Message, error){
	msg := new(Message)
	err := c.GobDecode(c.Body, msg)
	return msg, err
}

// SendMessage send message to remote addr
func (c *Connection) SendMessage(msg *Message) error {
	head := DefaultHead()
	message, err := c.GobEncode(msg)
	if err != nil{
		return err
	}
	head.DataLen = uint64(len(message))
	_, err= c.WriteMerge(head.GetBytes(), message)
	return err
}


// Close close conn and put to pool
func (c *Connection) Close() {
	c.conn.Close()
	putConnction(c)
}

func (c *Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

// NewConnction create new conn with net.Conn
func NewConnction(con net.Conn) *Connection {
	atomic.AddUint64(&connectIndex, 1)

	c := connctionPool.Get()
	if h, ok := c.(*Connection); ok {
		h.conn = con
		return h
	}
	return &Connection{lock: new(sync.RWMutex), conn: con}
}


func putConnction(c *Connection) {
	c.ConnIndex = 0
	c.conn = nil
	c.Head = nil
	c.headBuf = nil
	c.Body = nil
	connctionPool.Put(c)
}

func slice_merge(slice1, slice2 []byte) (c []byte) {
	c = append(slice1, slice2...)
	return
}

