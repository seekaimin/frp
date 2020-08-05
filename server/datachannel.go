package main

import (
	"net"
	tc "utils/common"
)

//DataChannel 数据通道
type DataChannel struct {
	remote   net.Conn
	local    net.Conn
	disposed bool
	addr     string
	Key      int32
}

// NewDataChannle 创建一个新的数据通道
func NewDataChannle(key int32, addr string, remote net.Conn, local net.Conn) *DataChannel {
	p := &DataChannel{
		disposed: false,
		remote:   remote,
		local:    local,
		addr:     addr,
		Key:      key,
	}
	return p
}

// Start 开始
func (p *DataChannel) Start() {
	defer p.Stop()
	//TCP
	p.tcpCopy()
}

// Stop 停止
func (p *DataChannel) Stop() {
	if p.disposed {
		return
	}
	logD("注册客户端[", p.addr, "]  退订[", p.Key, "]成功!")
	p.disposed = true
	if p.remote != nil {
		p.remote.Close()
		p.remote = nil
	}
	if p.local != nil {
		p.local.Close()
		p.local = nil
	}
}
func (p *DataChannel) tcpCopy() {
	go func() {
		defer p.Stop()
		tc.SocketCopy(p.local, p.remote, CACHESIZE)
		//logD("DataChannel send data size:", l)
	}()
	tc.SocketCopy(p.remote, p.local, CACHESIZE)
}
