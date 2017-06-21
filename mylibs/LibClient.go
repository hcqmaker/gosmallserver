package mylibs

import (
	"log"
	"net"
)

type ClientHandler func(*LibClient, int16, *NetMsg) error

/*

package main

import (
	"fmt"
	"mylibs"
)

func HandlerTest(ptr *mylibs.LibClient, cmd int16, msg *mylibs.NetMsg) error {
	return nil
}

func main() {
	c := mylibs.NewLibClient(HandlerTest)
	err := c.Connect("127.0.0.1:9001")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Run()
	for {
		var msg string
		fmt.Print("input data:")
		fmt.Scanf("%s", &msg)

		dt := mylibs.NewNetMsg()
		dt.Cmd = 10
		dt.Data = []byte(msg)
		c.Send(dt)
	}
}

*/

type LibClient struct {
	Host      string
	Conn      net.Conn
	Handle    ClientHandler
	buffer    *NetStream
	Code      *LibCodec
	msgChan   chan *NetMsg
	timeHChan chan int
}

func NewLibClient(h ClientHandler) *LibClient {
	a := &LibClient{
		Handle:    h,
		buffer:    NewNetStream(1024),
		Code:      NewLibCodec(),
		msgChan:   make(chan *NetMsg, 1024),
		timeHChan: make(chan int, 2),
	}
	return a
}

func runClientReceive(c *LibClient) {
	buffer := make([]byte, 1024)
	for {
		n, err := c.Conn.Read(buffer)
		if err != nil {
			log.Printf("%s read error %s", c.Conn.RemoteAddr().String(), err.Error())
			return
		}

		c.buffer.WriteBytes(buffer, 0, n)
		for {
			if c.buffer.PrefixedDataAvailable() {
				ln, er := c.buffer.ReadNetUint16()
				if er != nil {
					break
				}
				msg := NewNetMsg()
				msg.Cmd, _ = c.buffer.ReadInt16()
				dtln := int(ln) - 2
				if dtln > 0 {
					msg.Data = make([]byte, dtln)
					c.buffer.ReadBytes(msg.Data, 0, dtln)
				}
				c.buffer.ClearRead()
				c.msgChan <- msg

			} else {
				break
			}
		}

	}
}

func runClientHandler(c *LibClient) {
	for {
		select {
		case ch := <-c.msgChan:
			c.Handle(c, ch.Cmd, ch)
		case <-c.timeHChan:
			log.Fatal("leave msg handle")
			return
		}
	}
}

func (self *LibClient) Connect(host string) (err error) {
	self.Host = host
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	self.Conn = conn
	return nil
}

func (self *LibClient) Send(data *NetMsg) (err error) {
	rt := self.Code.Encode(data)
	ln := 2 + rt.Wpos
	stream := NewNetStream(ln)
	stream.WriteNetUint16(uint16(ln))
	stream.WriteBytes(data.Data, 0, rt.Wpos)
	self.Conn.Write(stream.Data)
	return nil
}

func (self *LibClient) Run() (err error) {
	go runClientReceive(self)
	go runClientHandler(self)
	return nil
}

func (self *LibClient) Close() (err error) {
	self.Conn.Close()
	self.timeHChan <- 1
	return nil
}
