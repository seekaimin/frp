package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var MainServer *MSServer
var debug bool
var addr string

func main() {
	start()
}

// start
func start() {
	debug, addr = DoConf()
	logFmtI("debug=[%t]--addr=[%s]", debug, addr)
	MainServer := NewMSServer(addr)
	MainServer.Start()
}

// test
func test() {
	ok := flag.Bool("debug", false, "is ok")
	id := flag.Int("id", 0, "id")
	port := flag.String("port", ":8080", "http listen port")
	var name string
	flag.StringVar(&name, "name", "123", "name")

	flag.Parse()

	fmt.Println("ok:", *ok)
	fmt.Println("id:", *id)
	fmt.Println("port:", *port)
	fmt.Println("name:", name)
}

// 初始化
func DoConf() (bool, string) {
	var d bool
	var a string
	flag.BoolVar(&d, "debug", false, "is debug")
	flag.StringVar(&a, "addr", ":9101", "name")
	flag.Parse()
	return d, a
}

// info
func logI(a ...interface{}) {
	fmt.Fprintln(os.Stdout, time.Now(), a)
}

// foramt
func logFmtI(format string, args ...interface{}) {
	temp := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stdout, time.Now(), temp)
}

// debug
func logD(a ...interface{}) {
	if debug {
		fmt.Fprintln(os.Stdout, time.Now(), a)
	}
}

// debug
func logFmtD(format string, args ...interface{}) {
	if debug {
		temp := fmt.Sprintf(format, args...)
		fmt.Fprintln(os.Stdout, time.Now(), temp)
	}
}
