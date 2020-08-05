package main

import (
	//"encoding/binary"

	"net"
	"sync"
	tc "utils/common"
)

//Server server
type Server struct {
	//server addr
	listener net.Listener
	Addr     string
	channels map[int32]*Channel //存储所有已注册的客户端
	running  bool
	minPort  int32 //最小Channel监听端口
	maxPort  int32 //最大Channel监听端口
	l        sync.Mutex
}

//NewServer 创建server
func NewServer(debug bool, addr string) *Server {
	server := &Server{
		Addr:     addr,
		channels: make(map[int32]*Channel),
		minPort:  9102,
		maxPort:  9110,
		l:        sync.Mutex{},
	}
	return server
}

func (p *Server) lock() {
	p.l.Lock()
}
func (p *Server) unlock() {
	p.l.Unlock()
}

//Running running
func (p *Server) Running() bool {
	return p.running
}

//Start start
func (p *Server) Start() {
	var err error
	//开启监听服务
	p.listener, err = net.Listen("tcp", p.Addr)
	if err != nil {
		logFmtI("启动服务失败！ error=%s", err)
		return
	}
	defer p.Stop()
	p.running = true
	logFmtI("server:[%s]开启成功!", p.Addr)
	defer p.Stop() //关闭监听的端口
	for p.Running() {
		conn, err := p.listener.Accept() //客户端连接
		if err != nil {
			logFmtD("监听客户端失败！ error=%s", err)
			break
		}
		go p.doClient(conn)
	}
}

func (p *Server) addChannel(channel *Channel, lock bool) {
	if lock {
		p.lock()
		defer p.unlock()
	}
	//移除原有客户端
	p.removeChannel(channel.Key, false)
	//新增客户端
	p.channels[channel.Key] = channel
}
func (p *Server) removeChannel(key int32, lock bool) {
	if lock {
		p.lock()
		defer p.unlock()
	}
	c, ok := p.channels[key]
	if ok {
		c.Stop()
		delete(p.channels, key)
	}
}

func (p *Server) registDataConn(key int32, conn net.Conn, lock bool) {
	if lock {
		p.lock()
		defer p.unlock()
	}
	channel, ok := p.channels[key]
	if ok && channel.Running() {
		//已注册
		channel.AddRegister(conn)
	} else {
		//未注册
		CloseConn(conn)
	}
}

//客户端数据解析
func dataParse(conn net.Conn) (result bool, mtype uint16, xBuffer *tc.XBuffer) {
	result = false
	remoteAddr := conn.RemoteAddr().String()
	logFmtD("客户端:[%s]连接", remoteAddr)
	buffer := make([]byte, 255)
	size, err := conn.Read(buffer)
	if err != nil || size <= 0 {
		logD("接收客户端数据失败！	size=%d	err=%s", size, err)
		return
	}
	logD("接收数据[", size, "] 内容:", tc.Buffer2String(buffer[0:size]))
	// head(4) + op(2) + len(2) + data
	minSize := 4 + 2 + 2
	if size < minSize {
		//数据解析失败
		return
	}
	xBuffer = tc.NewXBuffer(buffer[0:size], true)
	mhead := xBuffer.GetUInt32()
	mtype = xBuffer.GetUInt16()
	mlength := int(xBuffer.GetUInt16())
	if xBuffer.GetIndex()+mlength > size {
		//数据结构长度不够
		//数据解析失败
		return
	}
	if mhead != HEAD {
		//消息头识别错误
		return
	}
	result = true
	return
}

// 处理客户端数据
func (p *Server) doClient(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	logFmtD("客户端:[%s]连接", remoteAddr)
	result, mtype, xBuffer := dataParse(conn)
	if false == result {
		//数据解析失败
		CloseConn(conn)
		return
	}
	if mtype == TYPEREGISTER {
		// 注册
		listenPort := xBuffer.GetInt32()
		if listenPort > p.maxPort || listenPort < p.minPort {
			// 不存在的监听端口
			CloseConn(conn)
			return
		}
		//移除已经注册的客户端
		p.removeChannel(listenPort, true)

		// 有新连接接入
		// 启动数据监听套接字
		f, channel := NewChannel(listenPort, conn)
		defer channel.Stop()
		if f {
			//注册成功
			temp := make([]byte, 4)
			conn.Write(temp)
		} else {
			//注册失败
			temp := make([]byte, 1)
			conn.Write(temp)
			return
		}
		logFmtD("注册成功:%d-%s", listenPort, remoteAddr)
		p.addChannel(channel, true)
		defer p.removeChannel(listenPort, true)
		channel.Start()
	} else if mtype == TYPECLIENTCONNECTION {
		// 注册
		port := xBuffer.GetInt32()
		p.registDataConn(port, conn, true)
		return
	} else {
		// 非法请求
		CloseConn(conn)
		return
	}
}

//Stop stop
func (p *Server) Stop() {
	if false == p.running {
		return
	}
	p.running = false
	if p.listener != nil {
		p.listener.Close()
	}
	logI("Stop")
}
