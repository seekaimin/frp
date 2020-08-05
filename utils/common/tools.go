package commontools

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//CalcMD5 计算MD5
//src []byte  需要计算的buff
func CalcMD5(src []byte) []byte {
	m := md5.New()
	m.Write(src)
	r := m.Sum(nil)
	return r
}

//CalcMD5String 计算MD5
//src string  字符串以UTF-8编码
func CalcMD5String(src string) string {
	d := []byte(src)
	r := CalcMD5(d)
	return hex.EncodeToString(r)
}

// HTTPPostJSON 发送 http json post请求
func HTTPPostJSON(url string, args []byte) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("url is null")
	}
	//生成client 参数为默认
	client := &http.Client{}
	//提交请求
	reqest, err := http.NewRequest("POST", url, bytes.NewBuffer(args))
	if err != nil {
		return nil, err
	}
	reqest.Header.Set("Content-Type", "application/json; charset=utf-8")
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}
	//将结果定位到标准输出 也可以直接打印出来 或者定位到其他地方进行相应的处理
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// HTTPPost 发送普通httppost请求
// url 指定地址
// args 参数
func HTTPPost(url string, args []byte) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("url is null")
	}
	//生成client 参数为默认
	client := &http.Client{}
	//提交请求
	reqest, err := http.NewRequest("POST", url, bytes.NewBuffer(args))
	if err != nil {
		return nil, err
	}
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}
	//将结果定位到标准输出 也可以直接打印出来 或者定位到其他地方进行相应的处理
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//TotalMinuteSecond 计算纳秒
func TotalMinuteSecond(small time.Time, big time.Time) int64 {
	return UnixMinuteSecond(big) - UnixMinuteSecond(small)
}

//TotalUnixNano 计算纳秒秒
func TotalUnixNano(small time.Time, big time.Time) int64 {
	return big.UnixNano() - small.UnixNano()
}

//TotalSecond 计算秒
func TotalSecond(small time.Time, big time.Time) int64 {
	return UnixSecond(big) - UnixSecond(small)
}

//UnixMinuteSecond 计算毫秒
func UnixMinuteSecond(p time.Time) int64 {
	return p.UnixNano() / 1000000
}

//UnixSecond 计算秒
func UnixSecond(p time.Time) int64 {
	return p.UnixNano() / 1000000000
}

//Combine 路径合并
func Combine(root string, paths ...string) string {
	result := ""
	sp := GetSystemSeparator()
	if EndWith(root, sp) {
		result = root
	} else if len(root) > 0 {
		result = root + sp
	}

	for _, v := range paths {
		if len(result) > 0 && false == EndWith(result, sp) {
			result = result + sp
		}
		t := v
		if StartWith(t, sp) {
			t = t[len(sp):]
		}
		result = result + t
	}
	return result
}

// GetSystemSeparator 获取系统路劲分隔符
func GetSystemSeparator() string {
	sp := "\\"
	if !os.IsPathSeparator('\\') {
		sp = "/"
	}
	return sp
}

// StartWith HasPrefix tests whether the string s begins with prefix.
func StartWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndWith HasSuffix tests whether the string s ends with suffix.
func EndWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

//Format 时间格式化 yyyy-MM-dd HH:mm:ss
func Format(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// MyLock 定义锁定接口
type MyLock struct {
	l sync.Mutex
}

// NewMyLock 创建自定义锁定对象
func NewMyLock() *MyLock {
	return &MyLock{l: sync.Mutex{}}
}

// Lock 同步锁 l 自定义同步锁 action 需要做的事情
func (p *MyLock) Lock(action func()) {
	p.l.Lock()
	defer p.l.Unlock()
	action()
}

//LockReturn 同步锁 l 自定义同步锁 action 需要做的事情  请使用ReturnT
func (p *MyLock) LockReturn(action func() interface{}) interface{} {
	p.l.Lock()
	defer p.l.Unlock()
	return action()
}

// IndexOfString 数组中筛选需要元素的位置 返回下标
//	buffer 队列
//	src 比较目标
func IndexOfString(buffer []string, src string) int {
	var result int
	result = -1
	EachBreakString(buffer, func(index int, value string) bool {
		if value == src {
			result = index
			return true
		}
		return false
	})
	return result
}

// IndexOfByte 数组中筛选需要元素的位置 返回下标
//	buffer 队列
//	src 比较目标
func IndexOfByte(buffer []byte, src byte) int {
	var result int
	result = -1
	EachBreakByte(buffer, func(index int, value byte) bool {
		if value == src {
			result = index
			return true
		}
		return false
	})
	return result
}

// EachByte 数组循环
//	buffer 队列
//	fun(k,v) res   k:index v:value
func EachByte(buffer []byte, fun func(int, byte)) {
	if fun == nil {
		return
	}
	for k, v := range buffer {
		fun(k, v)
	}
}

// EachString 数组循环
//	buffer 队列
//	fun(k,v) res   k:index v:value
func EachString(buffer []string, fun func(int, string)) {
	if fun == nil {
		return
	}
	for k, v := range buffer {
		fun(k, v)
	}
}

// EachBreakByte 数组循环
//	buffer 队列
//	fun(k,v) isbreak   k:index v:value  isbreak:是否退出
func EachBreakByte(buffer []byte, fun func(int, byte) bool) {
	if fun == nil {
		return
	}
	for k, v := range buffer {
		if fun(k, v) {
			break
		}
	}
}

// EachBreakString 数组循环
//	buffer 队列
//	fun(k,v) isbreak   k:index v:value  isbreak:是否退出
func EachBreakString(buffer []string, fun func(int, string) bool) {
	if fun == nil {
		return
	}
	for k, v := range buffer {
		if fun(k, v) {
			break
		}
	}
}

//SocketCopy socket流拷贝 不支持异常处理(包括超时)
// sender 发送者
// receiver 接收者
// cacheSize 缓存大小
func SocketCopy(sender net.Conn, receiver net.Conn, cacheSize int) {
	buffer := make([]byte, cacheSize)
	for sender != nil && receiver != nil {
		l, err := receiver.Read(buffer[0:])
		if err != nil || l <= 0 {
			//fmt.Println("copy1=", err)
			break
		}
		if sender != nil {
			l, err = sender.Write(buffer[0:l])
			if err != nil || l <= 0 {
				//fmt.Println("copy2=", err)
				break
			}
		}
	}
}

// Buffer2String 数组转十六进制字符串
func Buffer2String(buffer []byte) string {
	return fmt.Sprintf("%x", buffer)
}
