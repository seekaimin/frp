package main

import (
	"fmt"
	"net"
	"sync"
	"time"
	tc "utils/common"
)

//Channel Channel
type Channel struct {
	dataServer net.Listener //数据监听客户端
	conn       net.Conn     //注册客户端连接
	Key        int32
	running    bool
	dataConn   chan net.Conn
	l          sync.Mutex
}

//Running running
func (p *Channel) Running() bool {
	return p.running
}

//NewChannel new channel
func NewChannel(key int32, conn net.Conn) (result bool, channel *Channel) {
	result = false
	channel = &Channel{
		conn:     conn,
		Key:      key,
		running:  true,
		l:        sync.Mutex{},
		dataConn: make(chan net.Conn),
	}
	flag := channel.createDataServer()
	if flag {
		//创建数据连接失败
		result = true
		return
	}
	return
}
func (p *Channel) lock() {
	p.l.Lock()
}
func (p *Channel) unlock() {
	p.l.Unlock()
}

//AddRegister regist
func (p *Channel) AddRegister(conn net.Conn) {
	if false == p.Running() {
		return
	}
	logD("注册客户端[", conn.RemoteAddr(), "]订阅[", p.Key, "]成功!")
	p.dataConn <- conn
}

//Start start
func (p *Channel) Start() {
	defer p.Stop()
	go p.dataAccept() //接收数据
	buffer := make([]byte, 255)
	//心跳检测
	for p.Running() {
		size, err := p.conn.Read(buffer)
		if err != nil || size <= 0 {
			logFmtD("心跳异常！ 端口=%d	size=%d	error=%s", p.Key, size, err)
			break
		}
		logD("客户端", p.Key, "数据", tc.Buffer2String(buffer[0:size]))
	}
}

//Stop stop
func (p *Channel) Stop() {
	if false == p.running {
		return
	}
	p.running = false

	logD("key:[", p.Key, "] is stoping")
	if p.dataServer != nil {
		p.dataServer.Close()
	}
	//关闭所有通道
	flag := true
	for flag {
		select {
		case c, closed := <-p.dataConn:
			{
				if closed {
					flag = false
				} else if c != nil {
					c.Close()
				} else {
					flag = false
				}
				break
			}
		case <-time.After(time.Second):
			{
				flag = false
				break
			}
		}
	}
	close(p.dataConn)
	if p.conn != nil {
		p.conn.Close()
	}
	logD("key:[", p.Key, "] is stoped")
}

func (p *Channel) createDataServer() bool {
	addr := fmt.Sprintf(":%d", p.Key)
	var err error
	p.dataServer, err = net.Listen("tcp", addr) //创建数据连接
	if err != nil {
		logFmtD("创建数据服务失败！ addr=%s	error=%s", addr, err)
		return false
	}
	logFmtD("创建数据服务成功！ addr=%s", addr)
	return true
}
func (p *Channel) dataAccept() {
	maxSize := 255
	buffer := make([]byte, maxSize)
	xBuffer := tc.NewXBuffer(buffer[0:maxSize], true)
	for p.Running() {
		//用conn接收链接
		var conn net.Conn
		var err error
		conn, err = p.dataServer.Accept() //用conn接收链接
		if err != nil {
			logD(err)
			return
		}
		addr := conn.RemoteAddr().String()
		data := []byte(addr)
		//size := 4 + 2 + 2 + len(addrData)
		xBuffer.Reset()
		xBuffer.CopyUInt32(HEAD)
		xBuffer.CopyUInt16(TYPEREGISTER)
		xBuffer.CopyInt16(int16(len(data)))
		xBuffer.CopyBytes(data)
		p.conn.Write(buffer[0:xBuffer.GetIndex()])
		go p.copyStream(conn)
	}
}

func (p *Channel) copyStream(remote net.Conn) {
	addr := remote.RemoteAddr().String()
	var local net.Conn
	var ok bool
	select {
	case local, ok = <-p.dataConn:
		break
	case <-time.After(20 * time.Second):
		ok = false
		break
	}
	datachannel := NewDataChannle(p.Key, addr, remote, local)
	defer datachannel.Stop()
	if false == ok {
		return
	}
	datachannel.Start()
}
