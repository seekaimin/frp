package main

import (
	//"encoding/binary"
	"net"
	"time"
	tc "utils/common"
)

var (
	//服务端注册
	HEAD_SERVER = 0x010F
	//客户端注册
	HEAD_CLIENT = 0x010E

	//服务端操作区域
	//创建订阅/发送订阅数据
	SERVER_OP_CREATE_OR_SEND = 0x01
	//取消订阅
	SERVER_OP_CANCEL = 0x02
)

type MSServer struct {
	//server addr
	listener  net.Listener
	Addr      string
	channels  map[int32]*Channel //存储所有已注册的客户端
	isRunning bool
	minPort   int32 //最小Channel监听端口
	maxPort   int32 //最大Channel监听端口
}

func NewMSServer(addr string) *MSServer {
	server := &MSServer{
		Addr:     addr,
		channels: make(map[int32]*Channel),
		minPort:  9102,
		maxPort:  9110,
	}
	server.channels = make(map[int32]*Channel)
	return server
}
func (p *MSServer) Running() bool {
	return p.isRunning
}
func (p *MSServer) Start() {
	var err error
	p.listener, err = net.Listen("tcp", p.Addr) //使用协议是tcp，监听的地址是addr
	if err != nil {
		logI(err)
		return
	}
	logFmtI("server:[%s]开启成功!", p.Addr)
	p.isRunning = true
	defer p.Stop() //关闭监听的端口
	for p.Running() {
		var conn net.Conn
		var err error
		conn, err = p.listener.Accept() //用conn接收链接
		if err != nil {
			logD(err)
			time.Sleep(1 * time.Second)
			p.Stop()
			return
		}
		go p.doClient(conn)
	}
}
func (p *MSServer) addChannel(channel *Channel) {
	p.removeChannel(channel.Key)
	p.channels[channel.Key] = channel
}
func (p *MSServer) removeChannel(key int32) {
	_, ok := p.channels[key]
	if ok {
		delete(p.channels, key)
	}
}
func (p *MSServer) doClient(conn net.Conn) {
	closeConn := true
	remoteAddr := conn.RemoteAddr().String()
	defer func() {
		if closeConn {
			logFmtD("客户端:[%s]已经退出", remoteAddr)
			conn.Close()
		}
	}()
	var listenPort int32
	logFmtD("客户端:[%s]已经接入", remoteAddr)
	buffer := make([]byte, 255)
	size, err := conn.Read(buffer)
	if err != nil {
		logD(err)
	}
	logD(size, buffer[0:size])
	// head(4) + type(2) + len(2) + data
	minSize := 4 + 2 + 2
	if size < minSize {
		return
	}
	xBuffer := tc.NewXBuffer(buffer[0:size], true)
	mhead := xBuffer.GetUInt32()
	mtype := xBuffer.GetUInt16()
	mlength := int(xBuffer.GetUInt16())
	if xBuffer.GetIndex()+mlength > size {
		// 数据结构长度不够
		// 消息头错误
		return
	}
	if mhead != HEAD {
		//消息头识别错误
		return
	}
	if mtype == TYPE_REGISTER {
		// 注册
		listenPort = xBuffer.GetInt32()
		if listenPort > p.maxPort || listenPort < p.minPort {
			// 不存在的监听端口
			return
		}
		// 有新连接接入
		// 启动数据监听套接字
		temp := make([]byte, 4)
		conn.Write(temp)
		logFmtD("注册成功:%d-%s", listenPort, remoteAddr)
		_, ok := p.channels[listenPort]
		if ok {
			//关闭原有注册客户端
			p.channels[listenPort].Stop()
		}
		channel := NewChannel(listenPort, conn)
		p.channels[listenPort] = channel
		defer func() {
			channel.Stop()
			p.removeChannel(listenPort)
		}()
		channel.Start()
		closeConn = false
	} else if mtype == TYPE_CLIENT_CONNECTION {
		// 注册
		port := xBuffer.GetInt32()
		_, ok := p.channels[port]
		if false == ok {
			//尚未注册
			return
		}
		p.channels[port].AddRegister(conn)
		closeConn = false
	} else {
		// 非法请求
		return
	}
}

func (p *MSServer) Stop() {
	if false == p.isRunning {
		return
	}
	p.isRunning = false
	p.listener.Close()
	logI("Stop")
}
