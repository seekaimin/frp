package ahttp

import (
	"container/list"
	"fmt"
)

//HTTPResponse response
type HTTPResponse struct {
	FirstLine                string
	HeadNames                *list.List //Set-Sookie  可能出现多次
	HeadValues               *list.List
	TransferEncoding         string     //分块编码标志
	TransferEncodingContents *list.List //分块编码文本信息
	ContentLength            int
	Content                  []byte
	infoHandle               func(string) //信息输出句柄
}

func (p *HTTPResponse) fi(format string, args ...interface{}) {
	temp := fmt.Sprintf(format, args...)
	p.i(temp)
}

func (p *HTTPResponse) i(message string) {
	if p.infoHandle != nil {
		p.infoHandle(message)
	}
}
