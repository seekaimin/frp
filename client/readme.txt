客户端
2020-08-05
启动方式
./frpc.exe -path='conf.json' -debug=true
参数描述
	-path='配置文件路径'		默认值=conf.json
	-debug=true	true:调试模式  	默认值=false

配置文件结构
[
  {
      "DataPort": 9104,
      "LocalServer": "127.0.0.1:8080",
      "ProtocolType": 1,
      "Server": "127.0.0.1:9101"
  },
  {
      "DataPort": 9103,
      "LocalServer": "127.0.0.1:5005",
      "ProtocolType": 0,
      "Server": "127.0.0.1:9101"
  }
]
参数描述
DataPort=服务端开房的数据端口(默认9102-9120)
LocalServer=本地需要透传出去的地址
ProtocolType=1	1:TCP;2:HTTP(请设置为1,配合nginx使用实现HTTP透传);
Server=服务器地址  frps.exe 绑定的地址