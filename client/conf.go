package main

import (
	"encoding/json"
	"fmt"
	filehelper "utils/file"
)

//Configuration conf
type Configuration struct {
	Server       string `xml:"Server,attr"`       //远程服务器地址;
	DataPort     int    `xml:"DataPort,attr"`     //外网映射端口;
	LocalServer  string `xml:"LocalServer,attr"`  //局域网主机地址;
	ProtocolType int    `xml:"ProtocolType,attr"` //协议类型 0:TCP 1:HTTP;
}

// LoadConfig 加载配置
func LoadConfig(path string) []Configuration {
	//读取配置文件
	data, err := filehelper.ReadAll(path)
	if err != nil {
		return nil
	}
	var s []Configuration
	json.Unmarshal(data, &s)
	fmt.Println(s)
	return s
}
