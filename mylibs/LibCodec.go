package mylibs

type LibCodec struct {
}

func NewLibCodec() *LibCodec {
	a := &LibCodec{}
	return a
}

func (self *LibCodec) Encode(data *NetMsg) *NetStream {
	ln := len(data.Data)
	rt := NewNetStream(2 + ln)
	rt.WriteInt16(data.Cmd)
	if ln > 0 {
		rt.WriteBytes(data.Data, 0, len(data.Data))
	}

	return rt
}

func (self *LibCodec) Decode(data NetStream) *NetMsg {
	rt := NewNetMsg()
	rt.Cmd, _ = data.ReadInt16()
	ln := data.AvailableNum()
	if ln > 0 {
		rt.Data = make([]byte, ln)
		data.ReadBytes(rt.Data, 0, ln)
	}

	return rt
}
