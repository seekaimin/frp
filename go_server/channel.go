package main

import (
	"fmt"
	"io"
	"net"
	"time"
	tc "utils/common"
)

type Channel struct {
	listener  net.Listener
	conn      net.Conn
	Key       int32
	isRunning bool
	l         *tc.MyLock
	regConn   chan net.Conn
}

func (p *Channel) Running() bool {
	return p.isRunning
}
func NewChannel(key int32, conn net.Conn) *Channel {
	var channel *Channel
	channel = &Channel{
		conn:      conn,
		Key:       key,
		isRunning: false,
		l:         tc.NewMyLock(),
		regConn:   make(chan net.Conn),
	}
	return channel
}
func (p *Channel) validate() bool {
	addr := fmt.Sprintf(":%d", p.Key)
	var err error
	p.listener, err = net.Listen("tcp", addr) //使用协议是tcp，监听的地址是addr
	if err != nil {
		logD(err)
		return false
	}
	return true
}

func (p *Channel) AddRegister(conn net.Conn) {
	if false == p.Running() {
		return
	}
	logD("注册客户端[", conn.RemoteAddr(), "]订阅[", p.Key, "]成功!")
	p.regConn <- conn
}
func (p *Channel) addRemote(conn net.Conn) {
	if false == p.Running() {
		return
	}
	logD("远程客户端[", conn.RemoteAddr(), "]订阅[", p.Key, "]成功!")
	go p.copyStream(conn)
}

func (p *Channel) Start() {
	if p.Running() {
		return
	}
	if false == p.validate() {
		return
	}
	p.isRunning = true
	go p.accept()
	buffer := make([]byte, 255)
	for p.Running() {
		size, err := p.conn.Read(buffer)
		if err != nil {
			logD(err)
			break
		}
		logD("客户端", p.Key, "发送数据", buffer[0:size])
	}
}

func (p *Channel) Stop() {
	if p == nil {
		return
	}
	if false == p.isRunning {
		return
	}

	logD("key:[", p.Key, "] is stoping")
	p.isRunning = false
	p.conn.Close()
	p.listener.Close()
	logD("key:[", p.Key, "] is stoped")
}

func (p *Channel) accept() {
	defer p.Stop()
	maxSize := 255
	buffer := make([]byte, maxSize)
	xBuffer := tc.NewXBuffer(buffer[0:maxSize], true)
	for p.Running() {
		// tools.println("等待客户端接入:%s", this.listenPort);
		//用conn接收链接
		var conn net.Conn
		var err error
		conn, err = p.listener.Accept() //用conn接收链接
		if err != nil {
			logD(err)
			return
		}
		addr := conn.RemoteAddr().String()
		data := []byte(addr)
		//size := 4 + 2 + 2 + len(addrData)
		xBuffer.Reset()
		xBuffer.CopyUInt32(HEAD)
		xBuffer.CopyUInt16(TYPE_REGISTER)
		xBuffer.CopyInt16(int16(len(data)))
		xBuffer.CopyBytes(data)
		p.conn.Write(buffer[0:xBuffer.GetIndex()])
		p.addRemote(conn)
	}
}

func (p *Channel) copyStream(remote net.Conn) {
	waitTime := int64(2000)
	var conn net.Conn
	var ok bool
	go func() {
		var open bool
		conn, open = <-p.regConn
		if !open {
			logD("channel closed!")
		}
		ok = true
	}()

	defer func() {
		if conn != nil {
			conn.Close()
		}
		remote.Close()
	}()
	startTime := time.Now()
	timeSpan := int64(0)
	for false == ok {
		if false == p.isRunning {
			break
		}
		if timeSpan < waitTime {
			timeSpan = tc.TotalUnixNano(startTime, time.Now())
			time.Sleep(100 * time.Millisecond)
		} else {
			p.regConn <- nil
		}
	}
	if conn == nil {
		return
	}
	logD(conn.RemoteAddr())
	//io.CopyBuffer(remote, conn, buffer)
	/*
		size, err := conn.Read(buffer)
		if err == nil {
			fmt.Println(err)
		}
		remote.Write(buffer[0:size])
		time.Sleep(10 * time.Second)
	*/
	go func() {
		defer func() {
			conn.Close()
			remote.Close()
		}()
		buffer := make([]byte, CACHE_SIZE)
		l, _ := io.CopyBuffer(conn, remote, buffer)
		logD("server send data size:", l)
	}()
	buffer := make([]byte, CACHE_SIZE)
	l, _ := io.CopyBuffer(remote, conn, buffer)
	logD("client send data size:", l)
}
