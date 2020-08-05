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

//ServerRun run
func run() {
	var addr string
	debug, addr = LoadConfig()
	logFmtI("debug=[%t]--addr=[%s]", debug, addr)
	server := NewServer(debug, addr)
	defer server.Stop()
	server.Start()
}

//LoadConfig 初始化
func LoadConfig() (bool, string) {
	var d bool
	var a string
	flag.BoolVar(&d, "debug", false, "is debug")
	flag.StringVar(&a, "addr", ":9101", "name")
	flag.Parse()
	return d, a
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
