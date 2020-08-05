package loghelper

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	l "utils/common"
)

/*
	var l log.Logger
	var err error
	err = verify("timeservice.ini")
	if err != nil {
		l = dxlogs.New("", dxlogs.INFO)
		log.SetLogger(l)
		log.Infof("配置文件读取失败", err)
		return
	}
	l = dxlogs.New(conf.LogRoot, conf.LogLevel)
	log.SetLogger(l)
	log.Infof("配置文件读取成功")
*/

//	FATAL   = 0	INFO    = 1	WARNING = 2	ERROR   = 3	DEBUG   = 4
const (
	FATAL   = 0
	INFO    = 1
	WARNING = 2
	ERROR   = 3
	DEBUG   = 4
)

type DxLogger struct {
	RootPath string
	LogLevel int
}

func New(root string, level int) *DxLogger {
	var temp DxLogger
	temp = DxLogger{LogLevel: level}
	if len(root) == 0 {
		temp.RootPath = "logs"
	} else {
		root = strings.TrimRight(root, l.GetSystemSeparator())
		temp.RootPath = root + l.GetSystemSeparator() + "logs"
	}
	return &temp
}
func (p *DxLogger) Debugf(format string, v ...interface{}) {
	if p.LogLevel < DEBUG {
		return
	}
	if v != nil && len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	p.write("DEBUG", format)
}
func (p *DxLogger) Infof(format string, v ...interface{}) {
	if p.LogLevel < INFO {
		return
	}
	if v != nil && len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	p.write("INFO", format)
}
func (p *DxLogger) Warnf(format string, v ...interface{}) {
	if p.LogLevel < WARNING {
		return
	}
	if v != nil && len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	p.write("WARNING", format)
}
func (p *DxLogger) Errorf(format string, v ...interface{}) {
	if p.LogLevel < ERROR {
		return
	}
	if v != nil && len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	p.write("ERROR", format)
}
func (p *DxLogger) Fatalf(format string, v ...interface{}) {
	if p.LogLevel < FATAL {
		return
	}
	if v != nil && len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	p.write("FATAL", format)
}

func (p *DxLogger) getPath() string {
	//判断目录是否存在
	t := time.Now()
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	sp := "\\"
	if !os.IsPathSeparator('\\') {
		sp = "/"
	}
	dir := fmt.Sprintf("%s%s%d%s%d", p.RootPath, sp, year, sp, month)
	fp := fmt.Sprintf("%s%s%d%s%d%s%d.log", p.RootPath, sp, year, sp, month, sp, day)
	//var f *os.File
	var err error
	if !checkFileIsExist(fp) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return ""
		}
		//f, err = os.Create(fp)
		//if err != nil {
		// 	fmt.Println("文件不存在")
		// 	return ""
		// }
	}
	return fp
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//同步锁
var loglock sync.Mutex

//写入日志
func (p *DxLogger) write(level string, message string) {
	if p.LogLevel >= DEBUG {
		fmt.Println(message)
	}
	path := p.getPath()
	if len(path) == 0 {
		return
	}
	loglock.Lock()
	defer loglock.Unlock()
	s := fmt.Sprintf("%s : [%s] - %s\r\n", time.Now().Format("2006-01-02 15:04:05"), level, message)
	fl, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644) //os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer fl.Close()
	fl.WriteString(s)
}
