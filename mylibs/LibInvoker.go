package mylibs

import (
	"errors"
	"fmt"
)

type IHandler func(int16, *NetMsg) error

type LibInvoker struct {
	HandleMp map[int16]IHandler
}

type LibInvokerHandler struct {
	*LibInvoker
	LibSessionHandler
}

func NewLibInvoker() *LibInvoker {
	a := &LibInvoker{}
	a.HandleMp = make(map[int16]IHandler)
	return a
}

func (self *LibInvoker) Register(cmd int16, handler IHandler) error {
	_, ok := self.HandleMp[cmd]
	if ok {
		return errors.New(fmt.Sprintf("has cmd: %d", cmd))
	}
	self.HandleMp[cmd] = handler
	return nil
}

func (self *LibInvoker) UnRegister(cmd int16, handler IHandler) error {
	_, ok := self.HandleMp[cmd]
	if !ok {
		return errors.New(fmt.Sprintf("can't find cmd:%d", cmd))
	}
	delete(self.HandleMp, cmd)
	return nil
}

func (self *LibInvoker) Call(cmd int16, data *NetMsg) error {
	h, ok := self.HandleMp[cmd]
	if !ok {
		return errors.New(fmt.Sprintf("can't find cmd:%d", cmd))
	}

	return h(cmd, data)
}

func NewLibInvokerHandler() *LibInvokerHandler {
	a := &LibInvokerHandler{
		LibInvoker: NewLibInvoker(),
	}
	return a
}

func (self *LibInvokerHandler) exceptionCaught(session *LibSession, err error) {
}

func (self *LibInvokerHandler) sessionOpened(session *LibSession) {

}

func (self *LibInvokerHandler) messageReceived(session *LibSession, data interface{}) {
	ptr := data.(*NetMsg)
	self.Call(ptr.Cmd, ptr)
}

func (self *LibInvokerHandler) sessionClosed(session *LibSession) {

}

func (self *LibInvokerHandler) sessionIdle(session *LibSession) {

}
