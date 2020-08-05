package commontools

import (
	"bytes"
	"fmt"
	"io"
)

type StringBuffer struct {
	pool bytes.Buffer
}

func (p *StringBuffer) Bytes() []byte {
	return p.pool.Bytes()
}
func (p *StringBuffer) String() string {
	return string(p.pool.Bytes())
}
func (p *StringBuffer) WriteString(value string) (int, error) {
	return p.pool.WriteString(value)
}
func (p *StringBuffer) Write(value []byte) (int, error) {
	return p.pool.Write(value)
}
func (p *StringBuffer) WriteByte(value byte) error {
	return p.pool.WriteByte(value)
}
func (p *StringBuffer) WriteRune(value rune) (int, error) {
	return p.pool.WriteRune(value)
}
func (p *StringBuffer) WriteTo(w io.Writer) (int64, error) {
	return p.pool.WriteTo(w)
}

//字符串追加
func (p *StringBuffer) WriteJSONString(name string, value string, isEnd bool) {
	if isEnd {
		p.pool.WriteString(fmt.Sprintf("\"%s\":\"%s\"", name, value))
	} else {
		p.pool.WriteString(fmt.Sprintf("\"%s\":\"%s\",", name, value))
	}
}

//字符串追加
func (p *StringBuffer) WriteJSONInt32(name string, value int64, isEnd bool) {

	if isEnd {
		p.pool.WriteString(fmt.Sprintf("\"%s\":%d", name, value))
	} else {
		p.pool.WriteString(fmt.Sprintf("\"%s\":%d,", name, value))
	}
}

//字符串追加
func (p *StringBuffer) WriteJSONBoolean(name string, value bool, isEnd bool) {
	if isEnd {
		p.pool.WriteString(fmt.Sprintf("\"%s\":%v", name, value))
	} else {
		p.pool.WriteString(fmt.Sprintf("\"%s\":%v,", name, value))
	}
}

func (p *StringBuffer) WriteJSONStringBuffer(name string, value StringBuffer, isEnd bool) {
	s := value.String()
	if isEnd {
		p.pool.WriteString(fmt.Sprintf("\"%s\":%s", name, s))
	} else {
		p.pool.WriteString(fmt.Sprintf("\"%s\":%s,", name, s))
	}
}
