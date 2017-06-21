package mylibs

import (
	"errors"
	"log"
	"net"
)

/**

package main

import (
	"log"
	"mylibs"
	"time"
)

func HandlerTest1(cmd int16, msg *mylibs.NetMsg) error {
	log.Printf("cmd:%d,data:%s", cmd, string(msg.Data))
	return nil

}

func HandlerTest2(cmd int16, msg *mylibs.NetMsg) error {
	log.Printf("cmd:%d,data:%s", cmd, string(msg.Data))
	return nil
}

func main() {
	invoker := mylibs.NewLibInvokerHandler()
	invoker.Register(int16(10), HandlerTest1)
	invoker.Register(int16(11), HandlerTest2)

	code := mylibs.NewLibCodec()
	s := mylibs.NewLibServer(":9001", code, invoker)
	s.Run()

	time.Sleep(3)
}

*/

type LibServer struct {
	LibSessionHandler
	ConnId   int64
	Host     string
	ConnMap  map[int64]*LibSession
	Code     *LibCodec
	listener net.Listener
	Handler  LibSessionHandler
}

func NewLibServer(host string, code *LibCodec, handler LibSessionHandler) *LibServer {
	a := &LibServer{
		ConnId:  1,
		Host:    host,
		ConnMap: make(map[int64]*LibSession),
		Code:    code,
		Handler: handler,
	}
	return a
}

func (self *LibServer) Run() (err error) {
	listener, err := net.Listen("tcp", self.Host)

	if err != nil {
		log.Println("Listen fail:", err)
		return err
	}

	self.listener = listener
	defer listener.Close()

	log.Println("-----------------Waiting for clients-xx--------------------")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		id := self.ConnId
		log.Println("new conn id:", id)
		self.ConnId = self.ConnId + 1
		c := NewLibSession(id, self.Code, conn, self.Handler)
		c.Open()
		self.ConnMap[id] = c
		c.Run()
	}
	return nil
}

func (self *LibServer) Close() (err error) {
	self.listener.Close()
	return nil
}

func (self *LibServer) FindConnById(id int64) (cn *LibSession, err error) {
	if c, ok := self.ConnMap[id]; ok {
		return c, nil
	}
	return nil, errors.New("can't find key")
}

//=============================
func (self *LibServer) exceptionCaught(session *LibSession, err error) {
	log.Println("error id:", session.Id, err)
	delete(self.ConnMap, session.Id)
	self.exceptionCaught(session, err)
	//session.Close()
}

func (self *LibServer) sessionOpened(session *LibSession) {
	self.sessionOpened(session)
}

func (self *LibServer) messageReceived(session *LibSession, data interface{}) {
	self.messageReceived(session, data)
}

func (self *LibServer) sessionClosed(session *LibSession) {
	self.sessionClosed(session)
}

func (self *LibServer) sessionIdle(session *LibSession) {
	self.sessionIdle(session)
}
