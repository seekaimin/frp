package ahttp

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var (
	//METHODGET get
	METHODGET = "GET"
	//METHODPOST post
	METHODPOST = "POST"
	//MMETHODCONNECTION connection
	MMETHODCONNECTION = "Connection"
)

//HTTPContextProxy httpcontextproxy
type HTTPContextProxy struct {
	disposed   bool         //资源释放标志
	client     net.Conn     //本地请求对象
	server     net.Conn     //远程请求对象
	serverHost string       //需要转发的服务地址
	infoHandle func(string) //信息输出句柄
}

//NewHTTPProxy 创建HTTP代理对象
// client 请求端
// server 服务端
// handle 消息输出句柄
func NewHTTPProxy(serverHost string, client net.Conn, server net.Conn, handle func(string)) *HTTPContextProxy {
	v := &HTTPContextProxy{
		server:     server,
		client:     client,
		infoHandle: handle,
		disposed:   false,
		serverHost: serverHost,
	}
	return v
}

//Start start
func (p *HTTPContextProxy) Start() {
	//获取 request
	request := p.httpGetReadRequest()
	if request == nil {
		return
	}
	p.httpSendRequest(request)
	response := p.httpReceiveResponse(request)
	if response == nil {

	} else {
		p.httpSendResponse(response)
	}
}

//httpSendError do error
func httpSendError(p *HTTPContextProxy) {

}

//Stop stop
func (p *HTTPContextProxy) Stop() {
	if p.disposed {
		return
	}
	p.disposed = true
}

func (p *HTTPContextProxy) fi(format string, args ...interface{}) {
	temp := fmt.Sprintf(format, args...)
	p.i(temp)
}

func (p *HTTPContextProxy) i(message string) {
	if p.infoHandle != nil {
		p.infoHandle(message)
	}
}

func (p *HTTPContextProxy) httpGetReadRequest() *HTTPRequest {
	reader := bufio.NewReader(p.client)
	line, err := httpReadLine(reader)
	if err != nil {
		return nil
	}
	request := &HTTPRequest{
		RequestURI:    "",
		Method:        "",
		URI:           "",
		Proto:         "",
		QueryString:   "",
		Heads:         make(map[string]string),
		ContentLength: 0,
		Content:       nil,
		infoHandle:    p.infoHandle,
	}
	request.httpParseRequestFirstLine(line)
	request.i("获取客户端请求  开始")
	request.i(line)
	request.ContentLength = 0
	for {
		line, err = httpReadLine(reader)
		if err != nil {
			//接收数据出现异常
			return nil
		}
		request.i(line)
		if 0 == len(line) {
			break
		}
		name, value, f := httpParseHead(line)
		if false == f {
			return nil
		}
		if strings.EqualFold("content-length", name) {
			//ContentLength
			len, _ := strconv.ParseInt(value, 10, 32)
			request.ContentLength = int(len)
			continue
		}
		request.Heads[name] = value
	}
	request.httpParseURI()
	request.fi("request.ContentLength = %d", request.ContentLength)
	if request.ContentLength > 0 {
		request.Content = make([]byte, request.ContentLength)
		len := 0
		for len < request.ContentLength {
			size, err := reader.Read(request.Content)
			if err != nil || size <= 0 {
				request.fi("读取客户端数据出现异常! error=%s", err)
				return nil
			}
			len = len + size
		}
		request.fi("content_length:%d content:%s", request.ContentLength, string(request.Content))
	}
	request.i("获取客户端请求  结束")
	return request
}
func (p *HTTPContextProxy) httpSendRequest(request *HTTPRequest) bool {
	writer := bufio.NewWriter(p.server)
	url := request.URI
	if len(request.QueryString) > 0 {
		url = fmt.Sprintf("%s?%s", url, request.QueryString)
	}
	line := fmt.Sprintf("%s %s %s", request.Method, url, request.Proto)
	httpWriteLine(writer, line)
	for k, v := range request.Heads {
		name := k
		value := v
		if strings.EqualFold("host", name) {
			value = p.serverHost
		}
		httpWriteHead(writer, name, value)
	}
	l := strconv.FormatInt(int64(request.ContentLength), 10)
	httpWriteHead(writer, "Content-Length", l)
	httpWriteLine(writer, "")
	if request.ContentLength > 0 {
		httpWriteBuffer(writer, request.Content, 0, request.ContentLength)
	} else {
		httpWriteLine(writer, "")
	}
	err := writer.Flush()
	return err == nil
}
func (p *HTTPContextProxy) httpReceiveResponse(request *HTTPRequest) *HTTPResponse {
	reader := bufio.NewReader(p.server)
	var line string
	var err error
	request.i("读取响应 开始")

	line, err = httpReadLine(reader)
	if err != nil {
		return nil
	}
	//第一行
	response := &HTTPResponse{
		FirstLine:  line,
		HeadNames:  list.New(),
		HeadValues: list.New(),
		infoHandle: p.infoHandle,
	}
	response.FirstLine = line
	request.i(line)
	response.ContentLength = 0
	for {
		line, err = httpReadLine(reader)
		if err != nil {
			return nil
		}
		request.i(line)
		if 0 == len(line) {
			break
		}
		//parse head
		name, value, f := httpParseHead(line)
		if false == f {
			return nil
		}

		if strings.EqualFold("content-length", name) {
			//请求数据长度
			len, _ := strconv.ParseInt(value, 10, 32)
			response.ContentLength = int(len)
		} else if strings.EqualFold("Transfer-Encoding", name) {
			response.TransferEncoding = value
		}
		response.HeadNames.PushBack(name)
		response.HeadValues.PushBack(value)
	}
	//读取body
	if response.ContentLength > 0 {
		//一般响应
		response.Content = make([]byte, response.ContentLength)
		httpReadBuffer(reader, response.Content, 0, response.ContentLength)
		response.fi("content_length:%d content=%s", response.ContentLength, string(response.Content))
	} else if strings.EqualFold("chunked", response.TransferEncoding) {
		//分块编码
		response.TransferEncodingContents = list.New()
		for {
			line, err = httpReadLine(reader)
			if err != nil {
				return nil
			}
			if len(line) == 0 {
				continue
			}
			l, _ := strconv.ParseInt(line, 16, 32)
			if l == 0 {
				break
			}
			temp := make([]byte, l)
			//拷贝已有数据
			httpReadBuffer(reader, temp, 0, int(l))
			response.TransferEncodingContents.PushBack(temp)
			response.i(string(temp))
		}
	}
	request.i("读取响应 结束")
	return response
}
func (p *HTTPContextProxy) httpSendResponse(response *HTTPResponse) {
	writer := bufio.NewWriter(p.client)
	var nameNode *list.Element
	var valueNode *list.Element
	httpWriteLine(writer, response.FirstLine)
	for {
		if nil == nameNode {
			nameNode = response.HeadNames.Front()
			valueNode = response.HeadValues.Front()
		} else {
			nameNode = nameNode.Next()
			valueNode = valueNode.Next()
		}
		if nil == nameNode {
			break
		}
		name := nameNode.Value.(string)
		value := valueNode.Value.(string)
		f := httpWriteHead(writer, name, value)
		if false == f {
			return
		}
	}
	httpWriteLine(writer, "")
	if nil != response.TransferEncodingContents {
		//分块编码
		var node *list.Element
		for {
			if nil == node {
				node = response.TransferEncodingContents.Front()
			} else {
				node = node.Next()
			}
			if nil == node {
				break
			}
			buffer := node.Value.([]byte)
			len := len(buffer)
			lenStr := strconv.FormatInt(int64(len), 16)
			httpWriteLine(writer, lenStr)
			httpWriteBuffer(writer, buffer, 0, len)
			httpWriteLine(writer, "")
			writer.Flush()
		}
		httpWriteLine(writer, "0")
		httpWriteLine(writer, "")
	} else if response.ContentLength > 0 {
		//一般响应信息结构
		httpWriteBuffer(writer, response.Content, 0, response.ContentLength)
	}
	writer.Flush()
}
func httpReadBuffer(reader *bufio.Reader, buffer []byte, start, end int) error {
	for start < end {
		size, err := reader.Read(buffer[start:end])
		if err != nil {
			return err
		}
		if size <= 0 {
			return errors.New("READ_DATA_LENGTH_ERROR")
		}
		start = start + size
	}
	return nil
}

func httpReadLine(reader *bufio.Reader) (string, error) {
	var line string
	buf, _, err := reader.ReadLine()
	if err != nil {
		//接收数据出现异常
		return line, err
	}
	line = string(buf)
	return line, nil
}

func httpParseHead(line string) (key, value string, flag bool) {
	pos := strings.Index(line, ":")
	if pos < 0 {
		return "", "", false
	}
	key = strings.Trim(line[0:pos], " ")
	value = strings.Trim(line[pos+1:], " ")
	return key, value, true
}

func httpWriteLine(writer *bufio.Writer, str string) bool {
	var line string
	if len(str) == 0 {
		line = "\r\n"
	} else {
		line = fmt.Sprintf("%s\r\n", str)
	}
	var count int
	var l int
	var err error
	for count < len(line) {
		l, err = writer.WriteString(line[count:])
		if err != nil || l <= 0 {
			return false
		}
		count = count + l
	}
	return true
}

func httpWriteBuffer(writer *bufio.Writer, buffer []byte, start, end int) bool {
	var l int
	var err error
	for start < end {
		l, err = writer.Write(buffer[start:])
		if err != nil || l <= 0 {
			return false
		}
		start = start + l
	}
	return true
}
func httpWriteHead(writer *bufio.Writer, k, v string) bool {
	s := fmt.Sprintf("%s:%s", k, v)
	return httpWriteLine(writer, s)
}
