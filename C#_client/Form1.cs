using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Configuration;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;
using Util.Common;
using Util.Common.ExtensionHelper;
using Util.Common.Nets;
using Util.WinForm.ExtensionHelper;

namespace PenetrateClientApp
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
        }
        ConnectionServerService service = null;
        private void button1_Click(object sender, EventArgs e)
        {
            this.button1.SetPropertry<bool>("Enabled", false);
            this.button2.SetPropertry<bool>("Enabled", true);
            string serverHost = this.txtServerHost.Text.Trim();
            int serverPort = this.txtServerPort.Text.ToInt32();
            int remotePort = this.txtRemotePort.Text.ToInt32();
            string localHost = this.txtLocalHost.Text.Trim();
            int localPort = this.txtLocalPort.Text.ToInt32();
            ProtocolTypes protocolType = ProtocolTypes.TCP;
            if (cboIsHTTP.Checked)
            {
                protocolType = ProtocolTypes.HTTP;
            }
            try
            {
                service = new ConnectionServerService()
                {
                    ServerHost = serverHost,
                    ServerPort = serverPort,
                    RemotePort = remotePort,
                    LocalHost = localHost,
                    LocalPort = localPort,
                    ProtocolType = protocolType,
                };
                service.ExitServiceHandle += Service_ExitServiceHandle;
                service.Start();

                //反写配置信息
                string val = string.Format("ServerHost={0};ServerPort={1};RemotePort={2}; LocalHost={3};LocalPort={4};ProtocolType={5}"
                    , serverHost, serverPort, remotePort, localHost, localPort, (int)protocolType);
                Program.SaveConfig(new KeyValuePair<string, string>("item0", val));
            }
            catch (Exception ex)
            {
                this.button1.SetPropertry<bool>("Enabled", true);
                this.button2.SetPropertry<bool>("Enabled", false);
                Console.WriteLine(ex.Message);
            }
        }
        private void Service_ExitServiceHandle()
        {
            this.button1.SetPropertry<bool>("Enabled", true);
            this.button2.SetPropertry<bool>("Enabled", false);
        }
        private void button2_Click(object sender, EventArgs e)
        {
            this.button1.SetPropertry<bool>("Enabled", true);
            this.button2.SetPropertry<bool>("Enabled", false);
            try
            {
                service.Stop();
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex);
            }
        }
        private void Form1_FormClosing(object sender, FormClosingEventArgs e)
        {
            button2_Click(null, null);
        }

        private void Form1_Load(object sender, EventArgs e)
        {
            //参数读取
            //ServerHost=test.evfort.cn;ServerPort=9201;RemotePort=9203;LocalHost=127.0.0.1;LocalPort=5005;
            string temp = System.Configuration.ConfigurationManager.AppSettings.Get("item0");
            temp.Split(';').ToList().ForEach((str) =>
            {
                string[] items = str.Split('=');
                if (items.Length > 1)
                {
                    string k = items[0].Trim().ToLower();
                    string v = items[1].Trim();
                    switch (k)
                    {
                        case "serverhost":
                            txtServerHost.Text = v;
                            break;
                        case "serverport":
                            txtServerPort.Text = v;
                            break;
                        case "remoteport":
                            txtRemotePort.Text = v;
                            break;
                        case "localhost":
                            txtLocalHost.Text = v;
                            break;
                        case "localport":
                            txtLocalPort.Text = v;
                            break;
                        case "protocoltype":
                            cboIsHTTP.Checked = "2".Equals(v);
                            break;
                        default:
                            Console.WriteLine("{0}-{1}-{2}", "localhost", k, string.Compare(k, "localhost"));
                            break;
                    }
                }
            });
        }

        private void button3_Click(object sender, EventArgs e)
        {
            string s = "Content-Length: 311";
            byte[] d = Encoding.ASCII.GetBytes(s);
            int size = d.Length;
            Console.WriteLine(size);
            Console.WriteLine(d.ShowString());
        }
    }
}
