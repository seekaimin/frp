package main

import (
	//"bufio"

	"net"
	"sync"

	//"os"
	"time"
	tc "utils/common"
)

//Server server
type Server struct {
	running  bool          //是否运行状态
	registed bool          //是否已经注册
	conn     net.Conn      //注册套接字
	conf     Configuration //配置信息
	l        sync.Mutex    //lock
}

// NewServer 创建一个数据通道
// conf 通道配置信息
func NewServer(conf Configuration) *Server {
	logI("配置:", conf)
	p := &Server{
		running:  true,
		conn:     nil,
		registed: false,
		conf:     conf,
		l:        sync.Mutex{},
	}
	return p
}

// Start 开启一个数据通道
func (p *Server) Start() {
	//服务端保持连接
	go p.checkConnection()
	//等待服务端的接入指令
	go p.waitNewChannel()
}

// Lock lock
func (p *Server) Lock() {
	p.l.Lock()
}

// Unlock unlock
func (p *Server) Unlock() {
	p.l.Unlock()
}

//与服务端的连接保持
func (p *Server) checkConnection() {

	buffer := make([]byte, 200)
	for p.running {
		//注册检核
		p.checkRegistState(buffer)
		//心跳发送
		p.checkHeart(buffer)
		if p.registed {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

//等待新连接接入
func (p *Server) waitNewChannel() {
	buffer := make([]byte, 200)
	for p.running {
		if false == p.registed {
			time.Sleep(1 * time.Second)
			continue
		}
		size, err := p.conn.Read(buffer)
		if err != nil || size <= 0 {
			p.disposeConnection(true)
			logFmtD("等待数据连接,接收数据出现异常!数据端口=%d	size=%d err=%s", p.conf.DataPort, size, err)
			time.Sleep(1 * time.Second)
			continue
		}
		//数据解析
		go p.startChannel(buffer[0:size])
	}
}
func (p *Server) startChannel(buffer []byte) {
	c := NewChannle(p.conf)
	c.Start(buffer)
}

func (p *Server) checkHeart(buffer []byte) {
	defer p.Unlock()
	p.Lock()
	if false == p.running {
		return
	}
	if false == p.registed {
		return
	}
	//心跳数据包封装
	// head(4) + op(2) + len(2) + data
	xBuffer := tc.NewXBuffer(buffer, true)
	xBuffer.CopyUInt32(HEAD)
	xBuffer.CopyUInt16(TYPEHEART)
	xBuffer.CopyUInt16(4)
	xBuffer.CopyInt32(int32(p.conf.DataPort))
	size, err := p.conn.Write(buffer[0:xBuffer.Index])
	if err != nil || size <= 0 {
		logFmtD("心跳发送失败! 数据端口=%d size=%d error=%s", p.conf.DataPort, size, err)
		p.disposeConnection(false)
		return
	}
}

//释放连接
func (p *Server) disposeConnection(lock bool) {
	if lock {
		defer p.Unlock()
		p.Lock()
	}
	if false == p.registed {
		return
	}
	p.registed = false
	if p.conn != nil {
		p.conn.Close()
	}
	p.conn = nil
}

//连接状态检核
func (p *Server) checkRegistState(buffer []byte) {
	defer p.Unlock()
	p.Lock()
	if false == p.running {
		return
	}
	if true == p.registed {
		return
	}

	//注册
	// head(4) + op(2) + len(2) + data
	xBuffer := tc.NewXBuffer(buffer, true)
	xBuffer.CopyUInt32(HEAD)
	xBuffer.CopyUInt16(TYPEREGISTER)
	xBuffer.CopyUInt16(4)
	xBuffer.CopyInt32(int32(p.conf.DataPort))
	//创建连接
	conn, err := net.Dial("tcp", p.conf.Server)
	if err != nil {
		logFmtD("注册失败! 数据端口=%d	error=%s", p.conf.DataPort, err)
		return
	}
	size, err := conn.Write(buffer[0:xBuffer.Index])
	if err != nil || size <= 0 {
		logFmtD("发送注册数据失败! 数据端口=%d size=%d error=%s", p.conf.DataPort, size, err)
		return
	}
	size, err = conn.Read(buffer[0:])
	if err != nil || size <= 0 {
		logFmtD("接收注册数据失败! 数据端口=%d size=%d error=%s", p.conf.DataPort, size, err)
		return
	}
	p.conn = conn
	p.registed = true
	logFmtI("注册成功! 数据端口=%d", p.conf.DataPort)
}
