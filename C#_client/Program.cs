using System;
using System.Collections.Generic;
using System.Configuration;
using System.Linq;
using System.Threading.Tasks;
using System.Windows.Forms;
using Util.Common.ExtensionHelper;

namespace PenetrateClientApp
{
    static class Program
    {
        /// <summary>
        /// 应用程序的主入口点。
        /// </summary>
        [STAThread]
        static void Main()
        {
            string temp = System.Configuration.ConfigurationManager.AppSettings.Get("StartType");
            if ("1" == temp)
            {
                //窗体应用程序
                WinfortStart();
            }
            else
            {
                //控制台应用程序
                ConsoleStart();
            }
        }

        private static void WinfortStart()
        {
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new Form1());
        }
        private static void ConsoleStart()
        {
            //ServerHost=test.evfort.cn;ServerPort=9201;RemotePort=9203;LocalHost=127.0.0.1;LocalPort=5005;
            List<ConnectionServerService> services = new List<ConnectionServerService>();
            for (int i = 0; i < 100; i++)
            {
                string key = string.Format("item{0}", i);
                string temp = System.Configuration.ConfigurationManager.AppSettings.Get("item" + i.ToString());
                if (string.IsNullOrEmpty(temp))
                {
                    break;
                }
                try
                {
                    //ServerHost = 远程服务器地址;
                    //ServerPort = 服务器通讯端口;
                    //RemotePort = 外网映射端口;
                    //LocalHost = 局域网主机地址;
                    //LocalPort = 局域网主机端口;
                    string ServerHost = "";
                    int ServerPort = 0;
                    int RemotePort = 0;
                    string LocalHost = "";
                    int LocalPort = 0;
                    ProtocolTypes ProtocolType = ProtocolTypes.TCP;
                    temp.ToLower().Split(';').ToList().ForEach((str) =>
                    {
                        string[] items = str.Split('=');
                        if (items.Length > 1)
                        {
                            string k = items[0];
                            string v = items[1];
                            switch (k)
                            {
                                case "serverhost":
                                    ServerHost = v;
                                    break;
                                case "serverport":
                                    ServerPort = v.ToInt32();
                                    break;
                                case "remoteport":
                                    RemotePort = v.ToInt32();
                                    break;
                                case "localhost":
                                    LocalHost = v;
                                    break;
                                case "localport":
                                    LocalPort = v.ToInt32();
                                    break;
                                case "protocoltype":
                                    if ("2".Equals(v))
                                    {
                                        ProtocolType = ProtocolTypes.HTTP;
                                    }
                                    break;
                                default:
                                    break;
                            }
                        }
                    });
                    ConnectionServerService service = new ConnectionServerService()
                    {
                        ServerHost = ServerHost,
                        ServerPort = ServerPort,
                        RemotePort = RemotePort,
                        LocalHost = LocalHost,
                        LocalPort = LocalPort,
                        ProtocolType = ProtocolType
                    };
                    services.Add(service);
                    service.Start();
                }
                catch (Exception ex)
                {
                    Console.WriteLine("启动失败,请检查配置信息:{0}", temp);
                }
            }
            Console.WriteLine("服务启动成功,按任意键结束!");
            Console.ReadLine();
            services.ForEach((service) =>
            {
                service.Stop();
            });
        }

        public static void SaveConfig(params KeyValuePair<string, string>[] items)
        {
            if (items == null || items.Length == 0)
            {
                return;
            }
            //Create the object
            Configuration config = ConfigurationManager.OpenExeConfiguration(ConfigurationUserLevel.None);
            //make changes
            items.ToList().ForEach((item) =>
            {
                config.AppSettings.Settings[item.Key].Value = item.Value;
            });
            //save to apply changes
            config.Save(ConfigurationSaveMode.Modified);
            ConfigurationManager.RefreshSection("appSettings");
        }
    }
}
