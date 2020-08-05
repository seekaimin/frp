package stream

import (
	//"encoding/binary"
	//"errors"
	"encoding/binary"
	"net"
	"sync"
)

const (
	SOCKET_ERROR = "SOCKET_ERROR"
)

type SocketInputStream struct {
	//连接套接字
	conn net.Conn
	//缓存大小
	MaxLength uint64
	//缓存 最大长度为 MaxLength
	buffer []uint8
	//标记位置结束位置
	start uint64
	//标记起始位置
	end uint64
	//数据读写方式
	IsBigEndian bool
	//读lock
	l     sync.Mutex
	Error func(error) //出现异常
}

func (p *SocketInputStream) lock() {
	p.l.Lock()
}
func (p *SocketInputStream) unlock() {
	p.l.Unlock()
}
func (p *SocketInputStream) reset() {
	if p.start <= 0 {
		return
	}
	len := p.length()
	copy(p.buffer[p.start:p.end], p.buffer[0:])
	p.start = 0
	p.end = len
}

func (p *SocketInputStream) length() uint64 {
	return p.end - p.start
}
func (p *SocketInputStream) Length() uint64 {
	defer p.unlock()
	p.lock()
	return p.length()
}

func (p *SocketInputStream) read() bool {
	p.reset()
	size, err := p.conn.Read(p.buffer[p.start:])
	if err != nil {
		p.Error(err)
	} else if size <= 0 {
		p.Error(err)
		return false
	} else {
		p.end += uint64(size)
	}
	return true
}
func (p *SocketInputStream) readCkeck(len uint64) {
	if p.length() <= len {
		p.read()
	}
}

func (p *SocketInputStream) readBytes(len uint64) []uint8 {

	for {
		if p.length() >= len {
			break
		} else {
			p.readCkeck(len)
		}
	}
	var val []uint8
	val = p.buffer[p.start : p.start+len]
	p.start = p.start + len
	return val
}
func (p *SocketInputStream) Read() []uint8 {
	defer p.unlock()
	p.lock()
	if p.length() == 0 {
		p.read()
	}
	return p.readBytes(p.length())
}
func (p *SocketInputStream) ReadBytes(len uint64) []uint8 {
	defer p.unlock()
	p.lock()
	return p.readBytes(len)
}

func (p *SocketInputStream) ReadUInt8() uint8 {
	defer p.unlock()
	p.lock()
	var val uint8
	var len uint64
	len = 1
	p.readCkeck(len)
	val = p.buffer[p.start]
	p.start++
	return val
}

func (p *SocketInputStream) ReadInt8() int8 {
	val := p.ReadUInt8()
	return int8(val)
}
func (p *SocketInputStream) ReadUInt16() uint16 {
	defer p.unlock()
	p.lock()
	var len uint64
	len = 2
	buf := p.readBytes(len)
	var val uint16
	if p.IsBigEndian {
		val = binary.BigEndian.Uint16(buf)
	} else {
		val = binary.LittleEndian.Uint16(buf)
	}
	return val
}
func (p *SocketInputStream) ReadInt16() int16 {
	val := p.ReadUInt16()
	return int16(val)
}
func (p *SocketInputStream) ReadUInt32() uint32 {
	defer p.unlock()
	p.lock()
	var len uint64
	len = 4
	buf := p.readBytes(len)
	var val uint32
	if p.IsBigEndian {
		val = binary.BigEndian.Uint32(buf)
	} else {
		val = binary.LittleEndian.Uint32(buf)
	}
	return val
}
func (p *SocketInputStream) ReadInt32() int32 {
	val := p.ReadUInt32()
	return int32(val)
}
func (p *SocketInputStream) ReadUInt64() uint64 {
	defer p.unlock()
	p.lock()
	var len uint64
	len = 8
	buf := p.readBytes(len)
	var val uint64
	if p.IsBigEndian {
		val = binary.BigEndian.Uint64(buf)
	} else {
		val = binary.LittleEndian.Uint64(buf)
	}
	return val
}
func (p *SocketInputStream) ReadInt64() int64 {
	val := p.ReadUInt64()
	return int64(val)
}
