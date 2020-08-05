package commontools

import (
	"encoding/binary"
	"net"
	"sync"
)

//buffer  处理对象
type SocketInputStream struct {
	//连接套接字
	conn net.Conn
	//缓存大小
	MaxLength int
	//每次读取数据长度 默认1k
	Length int
	//缓存 最大长度为 MaxLength
	buffer []uint8
	//标记位置结束位置
	start int
	//标记起始位置
	end int
	//数据读写方式
	IsBigEndian bool
	//lock
	l sync.Mutex
}

func NewSocketInputStream(maxlength int, bigendian bool) *SocketInputStream {
	stream := &SocketInputStream{
		MaxLength:   maxlength,
		start:       0,
		end:         0,
		Length:      1024,
		IsBigEndian: bigendian,
		l:           sync.Mutex{},
	}
	return stream
}
func (p *SocketInputStream) reset() {
	if p.start <= 0 {
		return
	}
	if p.end <= p.start {
		return
	}
	len := p.end - p.start
	copy(p.buffer[p.start:p.end], p.buffer[0:])
	p.start = 0
	p.end = len
}

func (p *SocketInputStream) size() int {
	return p.end - p.start
}

func (p *SocketInputStream) Size() int {
	defer func() {
		p.l.Unlock()
	}()
	p.l.Lock()
	return p.size()
}
func (p *SocketInputStream) read() error {
	defer func() {
		p.l.Unlock()
	}()
	p.l.Lock()
	p.reset()
	len := p.end + p.Length
	if len > p.MaxLength {
		len = p.MaxLength - p.end
	}
	size, err := p.conn.Read(p.buffer[p.start : len+p.start])
	if size > 0 {
		p.end += size
	}
	return err
}
func (p *SocketInputStream) readCkeck(len int) error {
	if p.size() <= len {
		return p.read()
	}
	return nil
}

func (p *SocketInputStream) readBytes(len int) ([]uint8, error) {
	var err error
	for err == nil {
		if p.size() >= len {
			break
		} else {
			err = p.readCkeck(len)
		}
	}
	if err == nil {
		var val []uint8
		val = p.buffer[p.start : p.start+len]
		return val, err
	}
	return nil, err
}
func (p *SocketInputStream) ReadBytes(len int) ([]uint8, error) {
	defer func() {
		p.l.Unlock()
	}()
	p.l.Lock()
	var err error
	for err == nil {
		if p.size() >= len {
			break
		} else {
			err = p.readCkeck(len)
		}
	}
	if err == nil {
		var val []byte
		val = p.buffer[p.start : p.start+len]
		return val, err
	}
	return nil, err
}

func (p *SocketInputStream) ReadUInt8() (uint8, error) {
	defer func() {
		p.l.Unlock()
	}()
	p.l.Lock()
	var val uint8
	var len int
	var err error
	len = 1
	for err == nil {
		p.readCkeck(len)
		val = p.buffer[p.start]
		p.start++
		break
	}
	return val, err
}

func (p *SocketInputStream) ReadInt8() (int8, error) {
	val, err := p.ReadUInt8()
	if err == nil {
		return int8(val), err
	} else {
		return 0, err
	}
}
func (p *SocketInputStream) ReadUInt16() (uint16, error) {
	defer func() {
		p.l.Unlock()
	}()
	p.l.Lock()
	len := 2
	buf, err := p.readBytes(len)
	var val uint16
	if err == nil {
		if p.IsBigEndian {
			val = binary.BigEndian.Uint16(buf)
		} else {
			val = binary.LittleEndian.Uint16(buf)
		}
	}
	return val, err
}
func (p *SocketInputStream) ReadInt16() (int16, error) {
	val, err := p.ReadUInt16()
	if err == nil {
		return int16(val), err
	} else {
		return 0, err
	}
}
func (p *SocketInputStream) ReadUInt32() (uint32, error) {
	defer func() {
		p.l.Unlock()
	}()
	p.l.Lock()
	len := 4
	buf, err := p.readBytes(len)
	var val uint32
	if err == nil {
		if p.IsBigEndian {
			val = binary.BigEndian.Uint32(buf)
		} else {
			val = binary.LittleEndian.Uint32(buf)
		}
	}
	return val, err
}
func (p *SocketInputStream) ReadInt32() (int32, error) {
	val, err := p.ReadUInt32()
	if err == nil {
		return int32(val), err
	} else {
		return 0, err
	}
}
func (p *SocketInputStream) ReadUInt64() (uint64, error) {
	defer func() {
		p.l.Unlock()
	}()
	p.l.Lock()
	len := 8
	buf, err := p.readBytes(len)
	var val uint64
	if err == nil {
		if p.IsBigEndian {
			val = binary.BigEndian.Uint64(buf)
		} else {
			val = binary.LittleEndian.Uint64(buf)
		}
	}
	return val, err
}
func (p *SocketInputStream) ReadInt64() (int64, error) {
	val, err := p.ReadUInt64()
	if err == nil {
		return int64(val), err
	} else {
		return 0, err
	}
}

func (p *SocketInputStream) Dispose() {
	p.conn.Close()
}
