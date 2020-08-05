package commontools

import (
	"encoding/binary"
)

//buffer  处理对象
type XBuffer struct {
	Data        []byte
	Index       int
	IsBigEndian bool
}

func NewXBuffer(buffer []byte, isBigEndian bool) *XBuffer {
	xbuffer := &XBuffer{Data: buffer, Index: 0, IsBigEndian: isBigEndian}
	return xbuffer
}

func (p *XBuffer) GetLength() int {
	return len(p.Data)
}
func (p *XBuffer) GetIndex() int {
	return p.Index
}
func (p *XBuffer) AddIndex(size int) {
	p.Index += size
}
func (p *XBuffer) Reset() {
	p.Index = 0
}

func (p *XBuffer) GetBytes(length int) []byte {
	size := length
	start := p.Index
	end := start + size
	temp := p.Data[start:end]
	p.AddIndex(size)
	return temp
}
func (p *XBuffer) GetInt8() int8 {
	num := p.GetUInt8()
	return int8(num)
}
func (p *XBuffer) GetUInt8() uint8 {
	size := 1
	buffer := p.GetBytes(size)
	return buffer[0]
}
func (p *XBuffer) GetInt16() int16 {
	return int16(p.GetUInt16())
}
func (p *XBuffer) GetUInt16() uint16 {
	size := 2
	buffer := p.GetBytes(size)
	var num uint16
	if p.IsBigEndian {
		num = binary.BigEndian.Uint16(buffer)
	} else {
		num = binary.LittleEndian.Uint16(buffer)
	}
	return num
}
func (p *XBuffer) GetInt32() int32 {
	return int32(p.GetUInt32())
}
func (p *XBuffer) GetUInt32() uint32 {
	size := 4
	buffer := p.GetBytes(size)
	var num uint32
	if p.IsBigEndian {
		num = binary.BigEndian.Uint32(buffer)
	} else {
		num = binary.LittleEndian.Uint32(buffer)
	}
	return num
}
func (p *XBuffer) GetUInt64() uint64 {
	size := 8
	buffer := p.GetBytes(size)
	var num uint64
	if p.IsBigEndian {
		num = binary.BigEndian.Uint64(buffer)
	} else {
		num = binary.LittleEndian.Uint64(buffer)
	}
	return num
}
func (p *XBuffer) GetInt64() int64 {
	num := p.GetUInt64()
	return int64(num)
}

func (p *XBuffer) CopyUInt8(src uint8) {
	temp := make([]byte, 1)
	temp[0] = src //byte(src)
	p.CopyBytes(temp)
}
func (p *XBuffer) CopyInt8(src int8) {
	p.CopyUInt8(uint8(src))
}
func (p *XBuffer) CopyUInt16(src uint16) {
	size := 2
	temp := make([]byte, size)
	if p.IsBigEndian {
		binary.BigEndian.PutUint16(temp, src)
	} else {
		binary.LittleEndian.PutUint16(temp, src)
	}
	p.CopyBytes(temp)
}
func (p *XBuffer) CopyInt16(src int16) {
	p.CopyUInt16(uint16(src))
}
func (p *XBuffer) CopyUInt32(src uint32) {
	size := 4
	temp := make([]byte, size)
	if p.IsBigEndian {
		binary.BigEndian.PutUint32(temp, src)
	} else {
		binary.LittleEndian.PutUint32(temp, src)
	}
	p.CopyBytes(temp)
}
func (p *XBuffer) CopyInt32(src int32) {
	p.CopyUInt32(uint32(src))
}
func (p *XBuffer) CopyUInt64(src uint64) {
	size := 8
	temp := make([]byte, size)
	if p.IsBigEndian {
		binary.BigEndian.PutUint64(temp, src)
	} else {
		binary.LittleEndian.PutUint64(temp, src)
	}
	p.CopyBytes(temp)
}
func (p *XBuffer) CopyInt64(src int64) {
	p.CopyUInt64(uint64(src))
}
func (p *XBuffer) CopyBytes(src []byte) {
	size := len(src)
	start := p.Index
	temp := src
	for i := 0; i < size; i++ {
		p.Data[start+i] = temp[i]
	}
	p.AddIndex(size)
}
