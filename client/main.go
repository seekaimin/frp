package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var debug bool

var l sync.Mutex

func main() {

	run()
}
func run() {
	//读取配置 开始
	flag.BoolVar(&debug, "debug", false, "is debug")
	var path string
	flag.StringVar(&path, "path", "conf.json", "config path")
	flag.Parse()
	logFmtI("debug=%t	path=%s", debug, path)
	confs := LoadConfig(path)
	logI("配置文件:", confs)
	if 0 == len(confs) {
		//没有读到配置
		logI("没有读到配置信息")
		return
	}
	//读取配置 结束

	//启动所有数据通道
	for _, conf := range confs {
		channel := NewServer(conf)
		channel.Start()
	}

	//保持程序不退出
	for {
		time.Sleep(time.Second)
	}
}

func fun1(needLock bool) {
	if needLock {
		l.Lock()
		defer l.Unlock()
	}
	fmt.Println("fun1")
	time.Sleep(5 * time.Second)
}

func fun2(needLock bool) {
	if needLock {
		l.Lock()
		defer l.Unlock()
	}
	fmt.Println("fun2")
	time.Sleep(5 * time.Second)

	fun1(false)
}
