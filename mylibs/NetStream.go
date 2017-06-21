package mylibs

import (
	"encoding/binary"
	"errors"
)

var (
	ERROR_LESS_LEN = errors.New("less length")
	ERR_LEN_0      = errors.New("len == 0")
)

const (
	NILL_RET = 0
)

// 提供读取/写入网络整形长度的方式，还有普通的整形处理

type NetStream struct {
	Data     []byte
	Rpos     int
	Wpos     int
	Capacity int
}

func NewNetStream(capacity int) *NetStream {
	a := &NetStream{
		Rpos:     0,
		Wpos:     0,
		Capacity: capacity,
		Data:     make([]byte, capacity),
	}
	return a
}

func (self *NetStream) CanRead(n int) bool {
	return self.Wpos-self.Rpos >= n
}

func (self *NetStream) CanWrite(n int) bool {
	return self.Capacity-self.Wpos >= n
}

func (self *NetStream) ReadNetUint16() (uint16, error) {
	if !self.CanRead(2) {
		return NILL_RET, ERROR_LESS_LEN
	}

	v := binary.BigEndian.Uint16(self.Data[self.Rpos : self.Rpos+2])
	self.Rpos += 2
	return v, nil
}

func (self *NetStream) WriteNetUint16(p uint16) error {
	if !self.CanWrite(2) {
		return ERROR_LESS_LEN
	}

	mn := make([]byte, 2)
	binary.BigEndian.PutUint16(mn, p)
	self.WriteBytes(mn, 0, 2)
	return nil
}

func (self *NetStream) GetUInt16() (uint16, error) {
	if !self.CanRead(2) {
		return NILL_RET, ERROR_LESS_LEN
	}

	v := binary.BigEndian.Uint16(self.Data[self.Rpos : self.Rpos+2])
	return v, nil
}

// read
func (self *NetStream) ReadUint8() (uint8, error) {
	if !self.CanRead(1) {
		return NILL_RET, ERROR_LESS_LEN
	}

	v := uint8(self.Data[self.Rpos])
	self.Rpos += 1
	return v, nil
}
func (self *NetStream) ReadUint16() (uint16, error) {
	if !self.CanRead(2) {
		return NILL_RET, ERROR_LESS_LEN
	}

	v := uint16(uint16(self.Data[self.Rpos]) | (uint16(self.Data[self.Rpos+1]) << 8))
	self.Rpos += 2
	return v, nil
}
func (self *NetStream) ReadUint32() (uint32, error) {
	if !self.CanRead(4) {
		return NILL_RET, ERROR_LESS_LEN
	}

	v := uint32((uint32(self.Data[self.Rpos]) | (uint32(self.Data[self.Rpos+1]) << 8) |
		(uint32(self.Data[self.Rpos+2]) << 16) | (uint32(self.Data[self.Rpos+3]) << 24)))
	self.Rpos += 4
	return v, nil
}
func (self *NetStream) ReadUint64() (uint64, error) {
	if !self.CanRead(8) {
		return NILL_RET, ERROR_LESS_LEN
	}

	v := uint64((uint64(self.Data[self.Rpos]) | (uint64(self.Data[self.Rpos+1]) << 8) |
		(uint64(self.Data[self.Rpos+2]) << 16) | (uint64(self.Data[self.Rpos+3]) << 24) |
		(uint64(self.Data[self.Rpos+4]) << 32) | (uint64(self.Data[self.Rpos+5]) << 40) |
		(uint64(self.Data[self.Rpos+6]) << 48) | (uint64(self.Data[self.Rpos+7]) << 56)))
	self.Rpos += 8
	return v, nil
}

func (self *NetStream) ReadInt8() (int8, error) {
	v, err := self.ReadUint8()
	return int8(v), err
}
func (self *NetStream) ReadInt16() (int16, error) {
	v, err := self.ReadUint16()
	return int16(v), err
}
func (self *NetStream) ReadInt32() (int32, error) {
	v, err := self.ReadUint32()
	return int32(v), err
}
func (self *NetStream) ReadInt64() (int64, error) {
	v, err := self.ReadUint64()
	return int64(v), err
}

func (self *NetStream) ReadBytes(p []byte, offset int, ln int) error {
	if self.CanRead(ln) {
		return ERROR_LESS_LEN
	}
	copy(p[offset:offset+ln], self.Data[self.Rpos:self.Rpos+ln])
	self.Rpos += ln
	return nil
}

func (self *NetStream) Reads(p []byte, ln int) error {
	return self.ReadBytes(p, 0, ln)
}

//===========================================================================
// write
func (self *NetStream) WriteUint8(p uint8) error {
	if !self.CanWrite(1) {
		return ERROR_LESS_LEN
	}

	self.Data[self.Wpos] = byte(p)
	self.Wpos += 1
	return nil
}

func (self *NetStream) WriteUint16(p uint16) error {
	if !self.CanWrite(2) {
		return ERROR_LESS_LEN
	}

	self.Data[self.Wpos] = byte(p)
	self.Data[self.Wpos+1] = byte(p >> 8)
	self.Wpos += 2
	return nil
}

func (self *NetStream) WriteUint32(p uint32) error {
	if !self.CanWrite(4) {
		return ERROR_LESS_LEN
	}

	self.Data[self.Wpos] = byte(p)
	self.Data[self.Wpos+1] = byte(p >> 8)
	self.Data[self.Wpos+2] = byte(p >> 16)
	self.Data[self.Wpos+3] = byte(p >> 24)
	self.Wpos += 4
	return nil
}

func (self *NetStream) WriteUint64(p uint64) error {
	if !self.CanWrite(8) {
		return ERROR_LESS_LEN
	}

	self.Data[self.Wpos] = byte(p)
	self.Data[self.Wpos+1] = byte(p >> 8)
	self.Data[self.Wpos+2] = byte(p >> 16)
	self.Data[self.Wpos+3] = byte(p >> 24)
	self.Data[self.Wpos+5] = byte(p >> 32)
	self.Data[self.Wpos+6] = byte(p >> 40)
	self.Data[self.Wpos+7] = byte(p >> 48)
	self.Data[self.Wpos+8] = byte(p >> 56)
	self.Wpos += 8
	return nil
}

func (self *NetStream) WriteInt8(p int8) error {
	return self.WriteUint8(uint8(p))
}

func (self *NetStream) WriteInt16(p int16) error {
	return self.WriteUint16(uint16(p))
}

func (self *NetStream) WriteInt32(p int32) error {
	return self.WriteUint32(uint32(p))
}

func (self *NetStream) WriteInt64(p int64) error {
	return self.WriteUint64(uint64(p))
}

func (self *NetStream) WriteBytes(p []byte, offset int, ln int) error {
	if ln == 0 {
		return ERR_LEN_0
	}
	if !self.CanWrite(ln) {
		return ERROR_LESS_LEN
	}

	copy(self.Data[self.Wpos:], p[offset:offset+ln])
	self.Wpos += ln
	return nil
}

func (self *NetStream) Writes(p []byte, ln int) error {
	return self.WriteBytes(p, 0, ln)
}

func (self *NetStream) ClearRead() {
	if self.Rpos > 0 {
		if self.Rpos == self.Wpos {
			self.Rpos = 0
			self.Wpos = 0
		} else {
			copy(self.Data[0:], self.Data[self.Rpos:self.Wpos])
			self.Wpos = self.Wpos - self.Rpos
			self.Rpos = 0
		}
	}
}

func (self *NetStream) AvailableNum() int {
	return self.Wpos - self.Rpos
}

func (self *NetStream) PrefixedDataAvailable() bool {
	if self.AvailableNum() < 2 {
		return false
	}

	b1, err := self.ReadNetUint16()

	if err == nil {
		if self.AvailableNum() >= int(b1)+2 {
			return true
		}
	}
	return false
}
