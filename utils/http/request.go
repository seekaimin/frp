package ahttp

import (
	"fmt"
	"strings"
)

//HTTPRequest request
type HTTPRequest struct {
	RequestURI    string
	Method        string
	URI           string
	Proto         string
	QueryString   string
	Heads         map[string]string
	ContentLength int
	Content       []byte
	infoHandle    func(string) //信息输出句柄
}

func (request *HTTPRequest) httpParseURI() {
	//GET 方式需要区分uri和参数
	index := strings.Index(request.RequestURI, "?")
	if index < 0 {
		request.URI = request.RequestURI
	} else {
		request.URI = request.RequestURI[0:index]
		request.QueryString = request.RequestURI[index+1 : len(request.RequestURI)-index-1]
	}
}

func (request *HTTPRequest) httpParseRequestFirstLine(line string) bool {
	items := strings.Split(line, " ")
	l := len(items)
	if l < 3 {
		return false
	}
	method := strings.Trim(items[0], " ")
	if strings.EqualFold(METHODGET, method) {
		request.Method = METHODGET
	} else if strings.EqualFold(METHODPOST, method) {
		request.Method = METHODPOST
	} else if strings.EqualFold(MMETHODCONNECTION, method) {
		request.Method = MMETHODCONNECTION
	} else {
		request.Method = method
	}
	uri := strings.Trim(items[1], " ")
	request.RequestURI = uri
	proto := strings.Trim(items[2], " ")
	request.Proto = proto
	return true
}
func (request *HTTPRequest) fi(format string, args ...interface{}) {
	temp := fmt.Sprintf(format, args...)
	request.i(temp)
}

func (request *HTTPRequest) i(message string) {
	if request.infoHandle != nil {
		request.infoHandle(message)
	}
}
