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

// getConnection get connection with ConnIndex
func (s *Server) getConnection(index int64) (*Connection, bool){
	conn, isExists := s.connectMap[index]
	return conn, isExists
}

// addConnection add new connection
func (s *Server) addConnection(conn *Connection) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.connectMap[conn.ConnIndex] = conn
	connLogger.Debug(fmt.Sprintf("Server addConnection %v", conn.RemoteAddr()))

}

// removeConnection remove connection
func (s *Server) removeConnection(conn *Connection) {
	s.locker.Lock()
	defer s.locker.Unlock()
	conn.Close()
	delete(s.connectMap, conn.ConnIndex)
	connLogger.Debug(fmt.Sprintf("Server removeConnection %v", conn.RemoteAddr()))
}

// GetConnCount get connection count on current Server
func (s *Server) GetConnectionCount() int{
	s.locker.RLock()
	defer s.locker.RUnlock()
	return len(s.connectMap)
}

// GetConnMap get connection map on current Server
func (s *Server) GetConnectionMap() map[int64]*Connection{
	return s.connectMap
}

// AddConnection add new connection
func (s *Server) AddConnection(conn *Connection){
	s.addConnection(conn)
}

// RemoveConnection remove connection with ConnIndex
func (s *Server) RemoveConnection(connIndex int64){
	s.locker.RLock()
	conn, isExists := s.getConnection(connIndex)
	s.locker.RUnlock()
	if !isExists{
		return
	}
	s.removeConnection(conn)
}

// Start start loop handler conn
func (s *Server) Start(){
	for {
		select {
		case <-s.stopChan:
			connLogger.Info(fmt.Sprint("Server get stop signal, so stop server loop handle conn"))
			break
		default:
			conn, err := s.listener.Accept()
			if err != nil{
				connLogger.Error(fmt.Sprint("Server accept listener error ", err))
			}else{
				connLogger.Debug(fmt.Sprintf("Server received new connection from %v", conn.RemoteAddr()))
				go s.handleConn(NewConnction(conn))
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
	connLogger.Debug(fmt.Sprintf("NewServer %v", tcpAddr))
	return s, err
}

