package filehelper

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

//文件全部文件读取
func ReadAll(filePth string) ([]byte, error) {
	f, err := ioutil.ReadFile(filePth)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//读取文件块
func ReadBlock(filePth string, bufSize int, processBlock func([]byte) bool) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, bufSize) //一次读取多少个字节
	bfRd := bufio.NewReader(f)
	for {
		n, err := bfRd.Read(buf)
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			return err
		}
		flag := processBlock(buf[:n]) // n 是成功读取字节数
		if flag == false {
			break
		}
	}
	return nil
}

//按读取文件
func ReadLine(filePth string, processLine func([]byte) bool) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		flag := processLine(line) //放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
		if flag == false {
			break
		}
	}
	return nil
}

//写入文件
func WriteFile(filePath string, buffer []byte) {
	ioutil.WriteFile(filePath, buffer, 0644)
}

//追加文件
func AppendFile(filePath string, buffer []byte) {
	fl, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer fl.Close()
	fl.Write(buffer)
}
