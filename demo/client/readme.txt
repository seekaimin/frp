#参数
      远程服务器地址=server运行的主机地址
      服务器通讯端口=server监听端口
      外网映射端口=现仅支持9102-9110
      内网主机地址/端口=内网需要穿透的主机地址	  
      穿透协议类型=1:TCP;2:HTTP;


#config文件配置结构如下:

<?xml version="1.0" encoding="utf-8" ?>
<configuration>
  <appSettings>

    <!-- 启动类型   1:窗体应用程序;其他:控制台应用程序 -->
    <add key="StartType" value="2"/>
    <!--  节点结构
      ServerHost=远程服务器地址;
      ServerPort=服务器通讯端口;
      RemotePort=外网映射端口;
      LocalHost=需要穿透的局域网主机地址;
      LocalPort=需要穿透的局域网主机端口; 
      ProtocolType=穿透协议类型  1:TCP;2:HTTP
    -->
    <!-- 窗体应用程序只会读取  item0 节点 -->
    <!-- 控制台应用程序 读取 item(0-100) -->
    <add key="item0" value="ServerHost=192.168.1.254;ServerPort=9101;RemotePort=9102;LocalHost=127.0.0.1;LocalPort=5005;ProtocolType=1"/>
    <add key="item1" value="ServerHost=192.168.1.254;ServerPort=9101;RemotePort=9104;LocalHost=192.168.1.125;LocalPort=8666;ProtocolType=2"/>
  </appSettings>
  <startup>
    <supportedRuntime version="v4.0" sku=".NETFramework,Version=v4.5" />
  </startup>
</configuration>