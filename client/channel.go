package main

import (
	"net"
	"sync"
	"time"
	tc "utils/common"
	ahttp "utils/http"
)

//Channel 数据通道
type Channel struct {
	remote   net.Conn
	local    net.Conn
	disposed bool
	l        sync.Mutex
	conf     Configuration
	addr     string
}

// NewChannle 创建一个新的数据通道
func NewChannle(conf Configuration) *Channel {
	p := &Channel{
		disposed: false,
		l:        sync.Mutex{},
		conf:     conf,
	}
	return p
}

// Lock lock
func (p *Channel) Lock() {
	p.l.Lock()
}

// Unlock unlock
func (p *Channel) Unlock() {
	p.l.Unlock()
}

// Start 开始
func (p *Channel) Start(buffer []byte) {
	defer p.Stop()
	//解析数据
	flag := p.dataParsing(buffer)
	if false == flag {
		return
	}
	//创建连接
	flag = p.createConns()
	if false == flag {
		return
	}
	if PTTCP == p.conf.ProtocolType {
		//TCP
		p.tcpCopy()
	} else if PTHTTP == p.conf.ProtocolType {
		//HTTP
		p.httpCopy()
	} else {
		//非法请求
	}
}

func (p *Channel) dataParsing(buffer []byte) bool {
	xBuffer := tc.NewXBuffer(buffer, true)
	mhead := xBuffer.GetUInt32()
	if mhead != HEAD {
		return false
	}
	mtype := xBuffer.GetUInt16()
	if TYPEREGISTER != mtype {
		return false
	}
	len := int(xBuffer.GetUInt16())
	d := xBuffer.GetBytes(len)
	p.addr = string(d)
	logFmtI("数据端口=%d	客户端[%s]接入", p.conf.DataPort, p.addr)
	return true
}

func (p *Channel) createConns() bool {
	//创建远程连接
	var remote net.Conn
	var local net.Conn
	var err error
	var size int
	//向远程发送注册信息
	buffer := make([]byte, 12)
	xBuffer := tc.NewXBuffer(buffer, true)
	xBuffer.CopyUInt32(HEAD)
	xBuffer.CopyUInt16(TYPECLIENTCONNECTION)
	xBuffer.CopyUInt16(4)
	xBuffer.CopyInt32(int32(p.conf.DataPort))
	remote, err = net.DialTimeout("tcp", p.conf.Server, 5*time.Second)
	if err != nil {
		logFmtD("创建远程连接失败! 数据端口=%d	server=%s	err=%s", p.conf.DataPort, p.conf.Server, err)
		return false
	}
	p.remote = remote
	size, err = remote.Write(buffer[0:xBuffer.Index])
	if err != nil || size <= 0 {
		logFmtD("创建远程连接失败! 数据端口=%d	server=%s	err=%s	size=%d", p.conf.DataPort, p.conf.Server, err, size)
		return false
	}
	local, err = net.DialTimeout("tcp", p.conf.LocalServer, 5*time.Second)
	if err != nil {
		logFmtD("创建本地连接失败! 数据端口=%d	localserver=%s	error=%s", p.conf.DataPort, p.conf.LocalServer, err)
		return false
	}
	p.local = local
	return true
}

// Stop 停止
func (p *Channel) Stop() {
	if p.disposed {
		return
	}
	defer p.Unlock()
	p.Lock()
	logFmtI("数据端口=%d	客户端[%s]退出", p.conf.DataPort, p.addr)
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
func (p *Channel) tcpCopy() {
	go func() {
		defer p.Stop()
		tc.SocketCopy(p.local, p.remote, CACHESIZE)
		//logD("Channel send data size:", l)
	}()
	tc.SocketCopy(p.remote, p.local, CACHESIZE)
}

func (p *Channel) httpCopy() {
	http := ahttp.NewHTTPProxy(p.conf.LocalServer, p.remote, p.local, func(msg string) {
		logD(msg)
	})
	http.Start()
}
