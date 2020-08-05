package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	commontools "utils/common"
	loghelper "utils/log"
	servicehelper "utils/service"
	"utils/set"
	stream "utils/streams"
)

// T t
type T struct {
}

// Start start
func (p *T) Start() {}

// Stop stop
func (p *T) Stop() {}

// Validate validate
func (p *T) Validate() bool {
	return false
}
func sum(s []int, sl time.Duration, c chan int) {
	fmt.Println("calc ", c)
	sum := 0
	for _, v := range s {
		sum += v
	}
	time.Sleep(sl * time.Second)
	c <- sum // send sum to c
}
func bytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
func main() {
	ch1 := make(chan int)
	go func() {
		s := 0
		count := 10
		for i := s; i < count+s; i++ {
			ch1 <- i
		}
	}()
	go func() {
		s := 10
		count := 10
		for i := s; i < count+s; i++ {
			ch1 <- i
		}
	}()

	go func() {
		s := 20
		count := 10
		for i := s; i < count+s; i++ {
			ch1 <- i
		}
	}()
	running := true
	for running {
		select {
		case c := <-ch1:
			{
				fmt.Println("ch1 = ", c)
				break
			}
		case <-time.After(1 * time.Second):
			{
				running = false
				fmt.Println("time out ")
				break
			}
		}
	}
	fmt.Println("exit")
}

func fff() {
	server, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("server:开启成功!")
	defer func() {
		if server != nil {
			server.Close()
		}
	}()
	for true {
		var conn net.Conn
		var err error
		conn, err = server.Accept() //用conn接收链接
		if err != nil {
			fmt.Println(err)
			break
		}
		go func(c net.Conn) {
			var s *stream.SocketStream
			s = stream.NewSocketStream(c, 1024, true)
			s.Writer.Error = func(e error) {
				fmt.Println(e)
			}
			s.Reader.Error = func(e error) {
				fmt.Println(e)
			}
			defer func() {
				s.Dispose()
			}()
			i := 0
			for i < 10 {
				v := s.Reader.Read()
				fmt.Println(v)
				var l uint64
				l = uint64(len(v))
				s.Writer.Write(v, 0, l)
				s.Writer.Flush()
				i = i + 1
			}
		}(conn)
	}
}
func paratest() {
	server := ":1001"
	count := 1
	for i := 1; i < len(os.Args); i = i + 2 {
		var k = os.Args[i]
		var v = os.Args[i+1]
		fmt.Println("k", k, "v", v)
		if k == "-s" {
			server = v
		} else if k == "-c" {
			count, _ = strconv.Atoi(v)
		}
	}
	fmt.Println(server, count)
}

// round 到最近的2的倍数
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
func ss() {
	//runtime.GOMAXPROCS(2)
	path := "/mp/mp3.html?123&456&name=范德萨"
	i := strings.Index(path, "?")
	uri := path
	//var args map[string]string
	//args = make(map[string]string)
	if i > 0 {
		uri = path[0:i]
		temp := path[i+1:]
		fmt.Println("temp", temp)
		i := strings.Index(temp, "&")
		items := strings.Split(temp, "&")
		for _, v := range items {
			kvs := strings.Split(v, "=")
			if len(kvs) > 1 {
				aa := kvs[0]
				vv := kvs[1]
				//key := strings.Trim(kvs[0])
				//val := strings.Trim(kvs[1])
				//args[key] = val
				fmt.Println("key", aa, "vv", vv)
			} else if len(kvs) > 0 {
				fmt.Println("key", "id", "vv", v)
			}
		}
		fmt.Println("i", i)
	}
	fmt.Println("uir", uri)
}

func split(s string, c string, cutc bool) (string, string) {
	i := strings.Index(s, c)
	if i < 0 {
		return s, ""
	}
	if cutc {
		return s[:i], s[i+len(c):]
	}
	return s[:i], s[i:]
}

// AA aa
type AA struct {
	ID   int
	Name string
}

func settest() {
	lst := set.New()
	a := &AA{ID: 1, Name: "a1"}
	aaa(lst, a)
	a = &AA{ID: 1, Name: "a2"}
	aaa(lst, a)
	a = &AA{ID: 1, Name: "a3"}
	aaa(lst, a)
	a = &AA{ID: 1, Name: "a4"}
	aaa(lst, a)
	prt(lst)
	lst.First()
	prt(lst)
	lst.First()
	prt(lst)
	lst.Last()
	prt(lst)
	a = &AA{ID: 1, Name: "a5"}
	aaa(lst, a)
	prt(lst)

	lst.First()
	lst.First()
	lst.First()
	lst.First()
	lst.First()
	prt(lst)
	a = &AA{ID: 1, Name: "a6"}
	aaa(lst, a)
	prt(lst)
}
func aaa(set *set.Set, a *AA) {
	//a.ID = set.Push(a)

	fmt.Println("add id", a.ID)
}
func prt(lst *set.Set) {
	fmt.Println("-----------------start--------------------")
	lst.Each(func(k int, it interface{}) {
		v, ok := it.(*AA)
		if ok {
			fmt.Println("item:", v.ID, v.Name)
		} else {
			fmt.Println("error")
		}
	})
	fmt.Println("-----------------end--------------------", lst.Count())
}

func logTest() {
	var l loghelper.DxLogger
	var s servicehelper.IService
	s = &T{}

	fmt.Println(l, s.Validate())

	ll := loghelper.New("", loghelper.DEBUG)
	ll.Debugf("ff")

	s1 := "100"

	//v, _ := strconv.Atoi(s1)
	v := commontools.ToUInt32(s1)
	fmt.Println(v)
}
func test() {
	for {
		sendToDevice([]byte("我是中文测试文本123abc"))
		time.Sleep(5 * time.Second)
	}

}

//发送数据到设备
func sendToDevice(buffer []byte) ([]byte, error) {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP("192.168.58.86"),
		Port: 7070,
	})
	t := time.Now()
	socket.SetDeadline(t.Add(time.Duration(5 * time.Second)))
	socket.SetWriteDeadline(t.Add(time.Duration(10 * time.Second)))
	socket.SetReadDeadline(t.Add(time.Duration(10 * time.Second)))
	if err != nil {
		return nil, err
	}
	defer socket.Close()
	_, err = socket.Write(buffer)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
