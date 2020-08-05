package stream

import (
	"net"
	"sync"
)

//buffer  处理对象
type SocketStream struct {
	//连接套接字
	conn     net.Conn
	Reader   *SocketInputStream
	Writer   *SocketOutputStream
	disposed bool //资源释放
}

func NewSocketStream(conn net.Conn, maxlength int, bigendian bool) *SocketStream {
	reader := &SocketInputStream{
		conn:        conn,
		MaxLength:   uint64(maxlength),
		start:       0,
		end:         0,
		IsBigEndian: bigendian,
		l:           sync.Mutex{},
		buffer:      make([]uint8, maxlength),
	}
	writer := &SocketOutputStream{
		conn:        conn,
		MaxLength:   uint64(maxlength),
		start:       0,
		end:         0,
		IsBigEndian: bigendian,
		l:           sync.Mutex{},
		AutoFlush:   false,
		buffer:      make([]uint8, maxlength),
	}
	stream := &SocketStream{
		conn:   conn,
		Reader: reader,
		Writer: writer,
	}
	return stream
}
func (p *SocketStream) IsDispose() bool {
	return p.disposed
}
func (p *SocketStream) Dispose() {
	if p.disposed {
		return
	}
	p.disposed = true
	p.dispose()
}

func (p *SocketStream) dispose() {
	if p.conn != nil {
		p.conn.Close()
	}
}
