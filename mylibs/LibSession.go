package mylibs

import (
	"log"
	"net"
)

type LibSessionHandler interface {
	exceptionCaught(session *LibSession, err error)
	sessionOpened(session *LibSession)
	messageReceived(session *LibSession, data interface{})
	sessionClosed(session *LibSession)
	sessionIdle(session *LibSession)
}

type LibSession struct {
	Id      int64
	Host    string
	Conn    net.Conn
	Handler LibSessionHandler
	buffer  *NetStream
	Code    *LibCodec
}

func NewLibSession(id int64, code *LibCodec, conn net.Conn, h LibSessionHandler) *LibSession {
	a := &LibSession{
		Id:      id,
		Conn:    conn,
		Handler: h,
		buffer:  NewNetStream(1024),
		Code:    code,
	}
	return a
}

func runConnReceive(c *LibSession) {
	buffer := make([]byte, 1024)
	for {
		n, err := c.Conn.Read(buffer)
		if err != nil {
			log.Printf("%s read error %s", c.Conn.RemoteAddr().String(), err.Error())
			c.Handler.exceptionCaught(c, err)
			return
		}
		if n == 0 {
			continue
		}

		c.buffer.Writes(buffer, n)

		for {
			if c.buffer.PrefixedDataAvailable() {
				ln, er := c.buffer.ReadNetUint16()
				if er != nil {
					break
				}
				msg := NewNetMsg()
				msg.Id = c.Id
				//log.Printf("==>%d size of buff %d num:%d r:%d,w:%d", ln, n, c.buffer.AvailableNum(), c.buffer.Rpos, c.buffer.Wpos)
				msg.Cmd, _ = c.buffer.ReadInt16()
				log.Printf("cmd : %d", msg.Cmd)
				dtln := int(ln) - 2
				if dtln > 0 {
					msg.Data = make([]byte, dtln)
					c.buffer.ReadBytes(msg.Data, 0, dtln)
				}
				c.buffer.ClearRead()
				c.Handler.messageReceived(c, msg)
			} else {
				break
			}
		}
	}
}

func (self *LibSession) Send(data *NetMsg) (err error) {
	rt := self.Code.Encode(data)
	ln := 2 + rt.Wpos
	stream := NewNetStream(ln)
	stream.WriteNetUint16(uint16(ln))
	stream.WriteBytes(data.Data, 0, len(data.Data))
	self.Conn.Write(stream.Data)
	return nil
}

func (self *LibSession) Run() (err error) {
	go runConnReceive(self)
	return nil
}

func (self *LibSession) Open() (err error) {
	self.Handler.sessionOpened(self)
	return nil
}

func (self *LibSession) Close() (err error) {
	self.Handler.sessionClosed(self)
	self.Conn.Close()
	return nil
}
