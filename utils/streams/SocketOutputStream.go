package stream

import (
	"errors"
	"net"
	"sync"
)

type SocketOutputStream struct {
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
	//写lock
	l         sync.Mutex
	AutoFlush bool        //自动输出缓存
	Error     func(error) //出现异常
}

func (p *SocketOutputStream) lock() {
	p.l.Lock()
}

func (p *SocketOutputStream) unlock() {
	p.l.Unlock()
}
func (p *SocketOutputStream) reset() {
	if p.start <= 0 {
		return
	}
	len := p.length()
	copy(p.buffer[p.start:p.end], p.buffer[0:])
	p.start = 0
	p.end = len
}

func (p *SocketOutputStream) length() uint64 {
	return p.end - p.start
}
func (p *SocketOutputStream) Length() uint64 {
	defer p.unlock()
	p.lock()
	return p.length()
}
func (p *SocketOutputStream) write(buf []uint8, start uint64, end uint64) bool {
	len := end - start
	for {
		if p.length()+len > p.MaxLength {
			//首先flush
			if false == p.flush() {
				return false
			}
		} else {
			copy(p.buffer[p.end:], buf[start:end])
			p.end = p.end + len
			break
		}
	}
	return true
}
func (p *SocketOutputStream) Write(buf []uint8, start uint64, end uint64) {
	defer p.unlock()
	p.lock()
	p.write(buf, start, end)
}

func (p *SocketOutputStream) flush() bool {
	for p.length() > 0 {
		size, err := p.conn.Write(p.buffer[p.start:p.end])
		if err != nil {
			p.Error(err)
			return false
		} else if size <= 0 {
			p.Error(errors.New(SOCKET_ERROR))
			return false
		} else {
			p.start = p.start + uint64(size)
		}
	}
	return true
}
func (p *SocketOutputStream) Flush() bool {
	defer p.unlock()
	p.lock()
	p.reset()
	return p.flush()
}
