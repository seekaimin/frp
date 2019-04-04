using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Util.Common;
using Util.Common.ExtensionHelper;
using Util.Common.Nets;

namespace PenetrateClientApp
{
    public class ConnectionServerService : BaseServer
    {
        public const int HEAD = 0x46474849;
        public const short TYPE_REGISTER = 0x0101;
        public const short TYPE_HEART = 0x0102;
        public const short TYPE_CLIENT_CONNECTION = 0x0103;

        public ProtocolTypes ProtocolType { get; set; } = ProtocolTypes.TCP;
        public string ServerHost { get; set; }
        public int ServerPort { get; set; }
        public int RemotePort { get; set; }
        public string LocalHost { get; set; }
        public int LocalPort { get; set; }
        TcpClientDemo serverSocket;
        public ConnectionServerService()
        {

        }
        protected override bool Validate()
        {
            //注册到服务器
            serverSocket = new TcpClientDemo(new System.Net.Sockets.TcpClient(this.ServerHost, this.ServerPort));
            return true;
        }
        protected override void Begin()
        {
            this.Running = true;
            new Thread(CheckConnection).Start();
            new Thread(Run).Start();
        }
        bool registerFlag = false;
        /// <summary>
        /// 检核连接和发送心跳
        /// </summary>
        void CheckConnection()
        {
            while (this.IsRunning)
            {
                if (false == this.registerFlag)
                {
                    //创建连接
                    try
                    {
                        RegisterServer();
                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine(ex.Message);
                        this.registerFlag = false;
                    }
                }
                if (false == registerFlag)
                {
                    Thread.Sleep(1000);
                    continue;
                }
                //发送心跳
                try
                {
                    this.SendHead();
                }
                catch (Exception ex)
                {
                    this.registerFlag = false;
                }
                Thread.Sleep(2000);
            }
        }
        void RegisterServer()
        {
            serverSocket.DoDispose();
            serverSocket = new TcpClientDemo(new System.Net.Sockets.TcpClient(this.ServerHost, this.ServerPort));
            serverSocket.SetTimeOut(10000, 10000);
            int index = 0;
            byte[] buffer = new byte[255];
            //head(4) + type(2) + len(2) + data
            buffer.Copy(ConnectionServerService.HEAD, ref index);
            buffer.Copy(ConnectionServerService.TYPE_REGISTER, ref index);
            buffer.Copy((short)4, ref index);
            buffer.Copy(this.RemotePort, ref index);
            byte[] data = this.Send(buffer.GetBytes(index), true);
            if (data.Length == 4)
            {
                serverSocket.SetTimeOut(1000, 1000);
                registerFlag = true;
                Console.WriteLine("注册到服务器成功:{0}:{1}-外网端口:{2}-本地主机地址:{3}-{4}-协议类型:{5}",
                    this.ServerHost, this.ServerPort, this.RemotePort, this.LocalHost, this.LocalPort, this.ProtocolType);
            }
            else
            {
                registerFlag = false;
            }
        }
        void SendHead()
        {
            int index = 0;
            byte[] buffer = new byte[255];
            //head(4) + type(2) + len(2) + data
            buffer.Copy(ConnectionServerService.HEAD, ref index);
            buffer.Copy(ConnectionServerService.TYPE_HEART, ref index);
            buffer.Copy((short)4, ref index);
            buffer.Copy(this.RemotePort, ref index);
            this.Send(buffer.GetBytes(index), false);
        }
        void Run()
        {
            while (this.IsRunning)
            {
                try
                {
                    if (false == this.registerFlag)
                    {
                        Thread.Sleep(1000);
                        continue;
                    }
                    byte[] buffrer = this.Receive();
                    int index = 0;
                    int head = buffrer.GetInt32(ref index);
                    if (head != ConnectionServerService.HEAD)
                    {
                        return;
                    }
                    int type = buffrer.GetInt16(ref index);
                    Console.WriteLine(type);
                    if (type == ConnectionServerService.TYPE_REGISTER)
                    {
                        //注册
                        int size = buffrer.GetInt16(ref index);
                        byte[] data = buffrer.GetBytes(size, ref index);
                        string str = Encoding.ASCII.GetString(data);
                        Console.WriteLine("新接入客户端:{0}-{1}", this.RemotePort, str);
                        if (this.ProtocolType == ProtocolTypes.TCP)
                        {
                            this.DoTCP();
                        }
                        else if (this.ProtocolType == ProtocolTypes.HTTP)
                        {
                            this.DoHTTP();
                        }
                        Console.WriteLine("{2}  客户端接入:{0}-{1}", this.RemotePort, str, DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss fff"));
                    }
                }
                catch (TimeoutException socket)
                {
                    //读取超时
                    continue;
                }
                catch (Exception e)
                {
                    SocketException ex = e.InnerException as SocketException;
                    if (ex != null && ex.SocketErrorCode == SocketError.TimedOut)
                    {
                        continue;
                    }
                    this.registerFlag = false;
                }
            }
        }
        private TcpClient createRemoteClient()
        {
            TcpClient socket = null;
            try
            {
                socket = new TcpClient(this.ServerHost, this.ServerPort);
                //发送注册信息
                int s = 4 + 2 + 2 + 4;
                byte[] temp = new byte[s];
                int index = 0;
                temp.Copy(ConnectionServerService.HEAD, ref index);
                temp.Copy(ConnectionServerService.TYPE_CLIENT_CONNECTION, ref index);
                temp.Copy((short)4, ref index);
                temp.Copy(this.RemotePort, ref index);
                socket.GetStream().Write(temp, 0, temp.Length);
            }
            catch (Exception ex)
            {
                socket.DoDispose();
                socket = null;
            }
            return socket;
        }
        private TcpClient createLocalClient()
        {
            TcpClient socket = null;
            try
            {
                socket = new TcpClient(this.LocalHost, this.LocalPort);
            }
            catch (Exception ex)
            {
                socket.DoDispose();
                socket = null;
            }
            return socket;
        }
        private void DoTCP()
        {
            TcpClient remote = null;
            TcpClient local = null;
            NetworkStream remoteStream = null;
            NetworkStream localStream = null;
            try
            {
                remote = this.createRemoteClient();
                local = this.createLocalClient();
                remote.ReceiveTimeout = local.ReceiveTimeout = 5000;
                remoteStream = remote.GetStream();
                localStream = local.GetStream();
                this.TCPCopy(remoteStream, localStream);
                this.TCPCopy(localStream, remoteStream);
            }
            catch (Exception ex)
            {
                //Console.WriteLine(ex.Message);
                remoteStream.DoDispose();
                localStream.DoDispose();
                remote.DoDispose();
                local.DoDispose();
            }
        }
        private void DoHTTP()
        {
            try
            {
                using (TcpClient remote = this.createRemoteClient())
                {
                    remote.ReceiveTimeout = 5000;
                    using (NetworkStream remoteStream = remote.GetStream())
                    {
                        this.HTTPCopy(remoteStream);
                    }
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }
        }

        object lockObj = new object();
        public byte[] Send(byte[] buffer, bool receiveFlag)
        {
            lock (lockObj)
            {
                this.serverSocket.Send(buffer);
                if (receiveFlag)
                {
                    return this.serverSocket.Receive();
                }
                return new byte[0];
            }
        }
        public byte[] Receive()
        {
            lock (lockObj)
            {
                return this.serverSocket.Receive();
            }
        }
        protected override void End()
        {
            serverSocket.DoDispose();
            Console.WriteLine("服务器注销成功:{0}:{1}-外网端口:{2}-本地主机地址:{3}:{4}-协议类型:{5}",
                this.ServerHost, this.ServerPort, this.RemotePort, this.LocalHost, this.LocalPort, this.ProtocolType);

        }

        private void TCPCopy(NetworkStream remoteStream, NetworkStream localStream)
        {
            new Thread(() =>
            {

                byte[] temp = new byte[1316];
                while (this.IsRunning)
                {
                    try
                    {
                        int size = localStream.Read(temp, 0, temp.Length);
                        if (size <= 0)
                        {
                            break;
                        }
                        remoteStream.Write(temp, 0, size);
                    }
                    catch (Exception ex)
                    {
                        SocketException socEx = ex.InnerException as SocketException;
                        if (socEx == null)
                        {
                            break;
                        }
                        if (socEx.SocketErrorCode != SocketError.TimedOut)
                        {
                            break;
                        }
                    }
                }
                remoteStream.DoDispose();
                localStream.DoDispose();
            }).Start();
        }
        int count = 0;
        object a = new object();
        int addCount()
        {
            lock (a)
            {
                count++;
                return count;
            }
        }
        private void HTTPCopy1(NetworkStream remote, NetworkStream local, int c)
        {
            //读取http参数
            byte[] lineBuffer = remote.ReadLineAsBytes();
            string line = Encoding.UTF8.GetString(lineBuffer);
            Console.WriteLine("{0}-{1}", c, line);
            //GET方式需要读取  参数
            string[] str = line.Split(' ');
            if (str.Length < 2)
            {
                throw new Exception("GET方式  URL读取失败");
            }
            string uri = str[1].Trim();
            //发送一个http请求到本地服务端
            string url = "http://{0}:{1}{2}".Fmt(this.LocalHost, this.LocalPort, uri);

            string method = "POST";
            if (line.StartsWith("GET"))
            {
                method = "GET";
            }


            int ContentLength = 0;
            //POST 方式需要读取body部分的参数
            while ((lineBuffer = remote.ReadLineAsBytes()).Length > 0)
            {
                line = Encoding.UTF8.GetString(lineBuffer);
                string tempLine = line.ToLower();
                Console.WriteLine("{0}-{1}", c, line);
                if (tempLine.Length == 0 || tempLine == "" || tempLine == "\r\n")
                {
                    break;
                }
                else if (tempLine.StartsWith("content-length:"))
                {
                    string[] items = line.Split(':');
                    if (items.Length >= 2)
                    {
                        ContentLength = items[1].Trim().ToInt32();
                    }
                }
            }
            byte[] data = new byte[ContentLength];
            if (ContentLength > 0)
            {
                int size = remote.Read(data, 0, ContentLength);
                Console.WriteLine("============SEND DATA===============");
                Console.WriteLine(Encoding.UTF8.GetString(data));
            }
            if (method == "GET")
            {
                this.DoGET(url, remote);
            }
            else if (method == "POST")
            {
                this.DoPost(local, url, data);
            }
            else
            {
                throw new Exception("尚未实现的请求方式");
            }
        }

        private void DoGET(string url, NetworkStream remote)
        {
            HttpWebRequest request = (HttpWebRequest)WebRequest.Create(url);
            //GET请求
            request.Method = "GET";
            request.ReadWriteTimeout = 30000;
            request.ContentType = "text/html;charset=UTF-8";
            using (HttpWebResponse response = (HttpWebResponse)request.GetResponse())
            {
                using (Stream responseStream = response.GetResponseStream())
                {
                    //using (StreamReader streamReader = new StreamReader(responseStream, Encoding.UTF8))
                    {
                        //返回内容
                        int ContentLength = 0;
                        this.Write(remote, "HTTP / 1.1 200 OK");
                        foreach (string head in response.Headers)
                        {
                            //Content - Length:4939
                            //Content - Type:text / html; charset = UTF - 8
                            //Date: Fri, 01 Mar 2019 03:36:47 GMT
                            // Server:Microsoft - IIS / 10.0
                            //X - Powered - By:ASP.NET
                            //Referrer Policy: no-referrer-when-downgrade
                            string h = "{0}:{1}".Fmt(head, response.Headers[head]);
                            if (head == "Content-Length")
                            {
                                ContentLength = response.Headers[head].ToInt32();
                            }
                            else if (head == "Server")
                            {
                                this.Write(remote, "Server:AMServer");
                                continue;
                            }
                            else if (head == "Referrer Policy" || head == "ETag")
                            {
                                continue;
                            }
                            Console.WriteLine(h);
                            this.Write(remote, h);
                        }
                        this.Write(remote, "");

                        Console.WriteLine("===========head end  len:{0}================", ContentLength);
                        //this.Write(remote, "Content-Type:{0}".Fmt(response.ContentType));
                        //this.Write(remote, "Server : AMServer");
                        //this.Write(remote, "Date:{0}".Fmt(DateTime.Now));
                        //this.Write(remote, "Content-Length: {0}".Fmt(data.Length));
                        if (ContentLength > 0)
                        {
                            byte[] data = responseStream.ReadAllBytes(1024);
                            Thread.Sleep(100);
                            Console.WriteLine(Encoding.UTF8.GetString(data));
                            remote.Write(data, 0, data.Length);
                            //this.Write(remote, "");
                        }
                        remote.Flush();
                    }
                }
            }
        }
        private void DoPost(NetworkStream remote, string url, byte[] sendData)
        {
            HttpWebRequest request = (HttpWebRequest)WebRequest.Create(url);
            //GET请求
            request.Method = "POST";
            request.ReadWriteTimeout = 30000;
            request.ContentType = "text/html;charset=UTF-8";
            request.GetRequestStream().Write(sendData, 0, sendData.Length);
            //if (headers != null)
            //{
            //    foreach (var v in headers)
            //    {
            //        request.Headers.Add(v.Key, v.Value);
            //    }
            //}

            //HTTP / 1.1 200 OK
            //Content - Type: text / html; charset = UTF - 8
            //Server: Microsoft - IIS / 10.0
            //X - Powered - By: ASP.NET
            //    Date: Fri, 01 Mar 2019 02:35:30 GMT
            //Content - Length: 4939
            using (HttpWebResponse response = (HttpWebResponse)request.GetResponse())
            {
                using (Stream responseStream = response.GetResponseStream())
                {
                    //using (StreamReader streamReader = new StreamReader(responseStream, Encoding.UTF8))
                    {
                        //返回内容
                        int ContentLength = 0;
                        this.Write(remote, "HTTP / 1.1 200 OK");
                        foreach (string head in response.Headers)
                        {
                            //Content - Length:4939
                            //Content - Type:text / html; charset = UTF - 8
                            //Date: Fri, 01 Mar 2019 03:36:47 GMT
                            // Server:Microsoft - IIS / 10.0
                            //X - Powered - By:ASP.NET
                            //Referrer Policy: no-referrer-when-downgrade
                            string h = "{0}:{1}".Fmt(head, response.Headers[head]);
                            if (head == "Content-Length")
                            {
                                ContentLength = response.Headers[head].ToInt32();
                            }
                            else if (head == "Server")
                            {
                                this.Write(remote, "Server:AMServer");
                                continue;
                            }
                            else if (head == "Referrer Policy" || head == "ETag")
                            {
                                continue;
                            }
                            Console.WriteLine(h);
                            this.Write(remote, h);
                        }
                        this.Write(remote, "");

                        Console.WriteLine("===========head end  len:{0}================", ContentLength);
                        //this.Write(remote, "Content-Type:{0}".Fmt(response.ContentType));
                        //this.Write(remote, "Server : AMServer");
                        //this.Write(remote, "Date:{0}".Fmt(DateTime.Now));
                        //this.Write(remote, "Content-Length: {0}".Fmt(data.Length));
                        if (ContentLength > 0)
                        {
                            byte[] data = responseStream.ReadAllBytes(1024);
                            Thread.Sleep(100);
                            Console.WriteLine(Encoding.UTF8.GetString(data));
                            remote.Write(data, 0, data.Length);
                            //this.Write(remote, "");
                        }
                        remote.Flush();
                    }
                }
            }
        }
        private void Write(NetworkStream stream, string line = "")
        {
            Console.WriteLine(line);
            byte[] buffer = Encoding.UTF8.GetBytes(line + "\r\n");
            stream.Write(buffer, 0, buffer.Length);
        }


        private void HTTPCopy(NetworkStream remote)
        {
            string uri;
            HttpWebRequest request = this.Create(remote, out uri, 5000);
            try
            {
                WebResponse response = request.GetResponse();
                this.WriteResponse(response, remote);
            }
            catch (WebException webex)
            {
                HttpWebResponse res = (HttpWebResponse)webex.Response;
                Console.WriteLine("{1}-{0}", res.StatusCode, res.StatusCode.ToInt32());
                Console.WriteLine(webex.Message);
                Console.WriteLine(webex.StackTrace);
                this.WriteResponse(webex.Response, remote);
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
                Console.WriteLine(ex.StackTrace);
                this.Write(remote, "HTTP/1.1 404");
            }
        }
        public HttpWebRequest Create(NetworkStream remote, out string uri, int timeOut = 5000)
        {
            uri = "";
            //读取第一行
            byte[] buffer = remote.ReadLineAsBytes();
            string line = Encoding.UTF8.GetString(buffer);
            Console.WriteLine("{0}", line);
            //method    url   version
            string[] first = line.Split(' ');
            if (first.Length < 2)
            {
                throw new Exception("GET方式  URL读取失败");
            }
            uri = first[1].Trim();
            if (uri == "/favicon.ico")
            {
                //throw new Exception("/favicon.ico  不做处理");
            }
            //method
            string method = first[0].Trim().ToUpper();
            //url
            string url = "http://{0}:{1}{2}".Fmt(this.LocalHost, this.LocalPort, uri);

            int ContentLength = 0;
            string ContentType = "text/html; charset=UTF-8";
            string RefererUrl = "";
            string UserAgent = null;
            string Accpet = null;
            string Connection = null;
            string IfModifiedSince = null;
            //读取headers
            Dictionary<string, string> headers = new Dictionary<string, string>();
            while ((buffer = remote.ReadLineAsBytes()).Length > 0)
            {
                line = Encoding.UTF8.GetString(buffer);
                string tempLine = line.ToLower();
                Console.WriteLine("request : {0}", line);
                if (tempLine.Length == 0 || tempLine == "" || tempLine == "\r\n")
                {
                    break;
                }
                string[] items = line.Split(':');
                if (items.Length < 2)
                {
                    continue;
                }
                string key = items[0].Trim();
                string value = items[1].Trim();
                if (tempLine.StartsWith("content-length:"))
                {
                    ContentLength = value.ToInt32();
                    continue;
                }
                else if (tempLine.StartsWith("host:"))
                {
                    continue;
                }
                else if (tempLine.StartsWith("content-type:"))
                {
                    ContentType = value;
                    continue;
                }
                else if (tempLine.StartsWith("referer:"))
                {
                    RefererUrl = value;
                    continue;
                }
                else if (tempLine.StartsWith("user-agent:"))
                {
                    UserAgent = value;
                    continue;
                }
                else if (tempLine.StartsWith("accept:"))
                {
                    Accpet = value;
                    continue;
                }
                else if (tempLine.StartsWith("connection:"))
                {
                    Connection = value;
                    continue;
                }
                else if (tempLine.StartsWith("if-modified-since:"))
                {
                    IfModifiedSince = value;
                    continue;
                }

                //读取数据头
                headers.Add(key, value);
            }
            byte[] data = new byte[ContentLength];
            if (ContentLength > 0)
            {
                int index = 0;
                while (this.IsRunning && (index < ContentLength))
                {
                    int size = remote.Read(data, index, ContentLength - index);
                    index += size;
                }
                Console.WriteLine("============SEND DATA START===============");
                Console.WriteLine(Encoding.UTF8.GetString(data));
                Console.WriteLine("============SEND DATA END===============");
            }
            HttpWebRequest request = (HttpWebRequest)WebRequest.Create(url);
            request.Method = method;
            request.ProtocolVersion = HttpVersion.Version10;
            request.Timeout = timeOut;
            request.ContentType = ContentType;
            request.Referer = RefererUrl;
            request.UserAgent = UserAgent;
            request.Accept = Accpet;
            request.KeepAlive = "keep-alive" == Connection;
            //request.IfModifiedSince = DateTime.Now;//IfModifiedSince.ToDateTime();
            foreach (var head in headers)
            {
                request.Headers.Set(head.Key, head.Value);
            }
            if (ContentLength > 0)
            {
                request.GetRequestStream().Write(data, 0, data.Length);
            }
            return request;
        }
        public void WriteResponse(WebResponse res, NetworkStream remote)
        {
            try
            {
                HttpWebResponse response = (HttpWebResponse)res;
                List<string> headers = new List<string>();
                int ContentLength = -1;
                foreach (string head in response.Headers)
                {
                    string line = "{0}:{1}".Fmt(head, response.Headers[head]);
                    if (head == "Server")
                    {
                        line = "Server:AMServer";
                        continue;
                    }
                    else if (head == "Referrer Policy" || head == "ETag")
                    {
                        Console.WriteLine(line);
                        line = "";
                        //continue;
                    }
                    else if (head.ToLower() == "content-length")
                    {
                        ContentLength = response.Headers[head].ToInt32();
                        line = "";
                    }
                    if (line.Length > 0)
                    {
                        headers.Add(line);
                    }
                }
                byte[] data = null;
                if (ContentLength < 0)
                {
                    Encoding encoding = string.IsNullOrEmpty(response.CharacterSet) ? Encoding.UTF8 : Encoding.GetEncoding(response.CharacterSet);
                    Console.WriteLine("-------------------READ RESPONSE DATA START--------------------");
                    using (StreamReader sr = new StreamReader(response.GetResponseStream(), encoding))
                    {
                        data = response.GetResponseStream().ReadAllBytes(1024);
                        if (false == string.IsNullOrEmpty(response.CharacterSet))
                        {
                            string str = Encoding.GetEncoding(response.CharacterSet).GetString(data);
                            Console.WriteLine(str);
                        }
                        //result = sr.ReadToEnd();
                        //Console.WriteLine(result);
                        //data = encoding.GetBytes(result);
                        ContentLength = data.Length;
                    }
                    Console.WriteLine("-------------------READ RESPONSE DATA END--------------------");
                }
                else if (ContentLength > 0)
                {
                    data = new byte[ContentLength];
                    response.GetResponseStream().Read(data, 0, data.Length);
                }
                this.Write(remote, "HTTP/1.1 200");
                foreach (string head in headers)
                {
                    if (head.ToLower().StartsWith("content-length:"))
                    {
                        continue;
                    }
                    this.Write(remote, head);
                }
                if (ContentLength > 0)
                {
                    this.Write(remote, "Content-Length:{0}".Fmt(data.Length));
                    this.Write(remote, "");
                    //this.Write(remote, (data.Length).ToString("X8").TrimStart('0'));
                    //this.Write(remote, result);
                    Console.WriteLine(Encoding.UTF8.GetString(data));
                    remote.Write(data, 0, data.Length);
                }
                else
                {
                    this.Write(remote, "");
                }
                remote.Flush();
            }
            catch (Exception ex)
            {
                Console.WriteLine("WriteResponse  Exception---{0}", ex.Message);
                this.Write(remote, "HTTP/1.1 404");
            }
        }

    }
}
