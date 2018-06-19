package connx

import (
	"sync"
	"net"
	"fmt"
)

const(
	defaultTcpPort    = ":7069"
)

type Server struct {
	listener *net.TCPListener
	handler OnConnHandle
	connectMap map[int64]*Connection
	locker     *sync.RWMutex
	stopChan 	chan struct{}
}

// handleConn loop handle new message from conn
func (s *Server) handleConn(conn *Connection) {
	s.addConnection(conn)
	defer s.removeConnection(conn)
	for {
		errRead := conn.readMessage()
		if errRead == nil {
			if s.handler != nil{
				errHandler := s.handler(conn)
				if errHandler != nil{
					connLogger.Error(fmt.Sprintf("Server.handler error ClientID:%v err:%v", conn.ConnIndex, errHandler))
				}
			}
		} else {
			//discard current request
			connLogger.Error(fmt.Sprintf("Server.handleConn error ClientID:%v err:%v", conn.ConnIndex, errRead))
			break
		}
	}
}

// SetOnConnHandle set handler on new conn receive
func (s *Server) SetOnConnHandle(handler OnConnHandle){
	s.handler = handler
}

// AddConnection add new connection
func (s *Server) addConnection(conn *Connection) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.connectMap[conn.ConnIndex] = conn
}

// RemoveConnection remove connection
func (s *Server) removeConnection(conn *Connection) {
	conn.Close()
	s.locker.Lock()
	defer s.locker.Unlock()
	delete(s.connectMap, conn.ConnIndex)
}

// Start start loop handler conn
func (s *Server) Start(){
	for {
		select {
		case <-s.stopChan:
			connLogger.Debug(fmt.Sprint("Server get stop signal, so stop server loop handle conn"))
			break
		default:
			conn, err := s.listener.Accept()
			if err != nil{
				connLogger.Error(fmt.Sprint("Server accept listener error ", err))
			}else{
				s.handleConn(NewConnction(conn))
			}
		}
	}
}

// Stop send stop signal
func (s *Server) Stop(){
	s.stopChan <- struct{}{}
}

// GetNewServer get new server with tcp port
func NewServer(tcpPort string, handler OnConnHandle) (*Server, error) {
	if tcpPort == ""{
		tcpPort = defaultTcpPort
	}

	var s *Server
	s = &Server{
		handler:handler,
		locker:     new(sync.RWMutex),
		connectMap: make(map[int64]*Connection),
		stopChan:make(chan struct{}),
	}

	connLogger.SetEnabledLog(true)

	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpPort)
	if err != nil {
		return s, err
	}
	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	return s, err
}

