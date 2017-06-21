package mylibs

type NetMsg struct {
	Id   int64
	Cmd  int16
	Data []byte
}

func NewNetMsg() *NetMsg {
	a := &NetMsg{
		Id:  0,
		Cmd: 0,
	}
	return a
}
