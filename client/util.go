package main

import (
	"fmt"
	"os"
	"time"
)

var (
	//HEAD head
	HEAD = uint32(0x46474849)
	// TYPEREGISTER register
	TYPEREGISTER = uint16(0x0101)
	// TYPEHEART heart
	TYPEHEART = uint16(0x0102)
	// TYPECLIENTCONNECTION new connection
	TYPECLIENTCONNECTION = uint16(0x0103)
	// CACHESIZE cache size
	CACHESIZE = 1 * 1024

	// HEADSERVER 服务端注册
	HEADSERVER = 0x010F
	// HEADCLIENT 客户端注册
	HEADCLIENT = 0x010E

	//PTTCP tcp
	PTTCP = 0
	//PTHTTP http
	PTHTTP = 1
)

func logPrintln(a interface{}) {
	fmt.Fprintln(os.Stdout, time.Now().Format("2006-01-02 15:04:05"), a)
}

// info
func logI(a ...interface{}) {
	logPrintln(a)
}

// foramt
func logFmtI(format string, args ...interface{}) {
	a := fmt.Sprintf(format, args...)
	logPrintln(a)
}

// debug
func logD(a ...interface{}) {
	if debug {
		logPrintln(a)
	}
}

// debug
func logFmtD(format string, args ...interface{}) {
	if debug {
		a := fmt.Sprintf(format, args...)
		logPrintln(a)
	}
}
