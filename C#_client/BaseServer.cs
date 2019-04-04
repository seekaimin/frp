using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Util.Common;

namespace PenetrateClientApp
{
    public abstract class BaseServer
    {
        public bool DEBUG => true;
        public string ServerName { get; set; } = "system server";
        protected bool? Running { private get; set; } = null;
        /// <summary>
        /// 处理消息
        /// </summary>
        public event Action<string> MessageHandle;
        /// <summary>
        /// 退出服务
        /// </summary>
        public event Action ExitServiceHandle;
        /// <summary>
        /// 退出服务
        /// </summary>
        protected void DoExitService()
        {
            if (this.ExitServiceHandle != null) { try { this.ExitServiceHandle(); } catch (Exception ex) { } }
        }
        public virtual void info(string message)
        {
            if (this.MessageHandle != null)
            {
                this.MessageHandle(string.Format("[{0}]-[{1}]", ThreadShell.Name, message));
            }
        }
        public virtual void info(string format, params object[] args)
        {
            if (this.MessageHandle != null)
            {
                this.MessageHandle(string.Format(format, args));
            }
        }
        public virtual void debug(string message)
        {
            if (this.MessageHandle != null && this.DEBUG)
            {
                this.MessageHandle(string.Format("[{0}]-[{1}]", this.ServerName, message));
            }
        }
        public virtual void debug(string format, params object[] args)
        {
            if (this.MessageHandle != null && this.DEBUG)
            {
                this.MessageHandle(string.Format(format, args));
            }
        }


        private QueuePool<ThreadShell> threadpool = new QueuePool<ThreadShell>();
        private QueuePool<Delegate> delegatepool = new QueuePool<Delegate>();
        protected void AddThread(ThreadShell item)
        {
            this.threadpool.Set(item);
        }
        protected void AddDelegate(Delegate item)
        {
            this.delegatepool.Set(item);
        }

        public void Start()
        {
            if (this.Running != null)
            {
                //服务已经启动
                return;
            }
            //设置服务为启动状态
            this.Running = false;
            if (this.MessageHandle != null)
            {
                this.AddDelegate(this.MessageHandle);
            }
            if (ExitServiceHandle != null)
            {
                this.AddDelegate(ExitServiceHandle);
            }
            if (false == this.Validate())
            {
                return;
            }
            //验证成功
            //清空所有线程 事件
            this.threadpool = new QueuePool<ThreadShell>();
            this.delegatepool = new QueuePool<Delegate>();
            //准备启动任务
            Begin();
            if (this.Running.Value)
            {
                this.info("服务启动成功!");
            }
            else
            {
                this.DoExitService();
            }
        }


        public bool IsRunning { get { return this.Running != null && this.Running.Value; } }



        /// <summary>
        /// 
        /// </summary>
        public void Stop()
        {
            if (this.Running == null)
            {
                //尚未运行
                return;
            }
            this.Running = null;
            try { this.End(); } catch (Exception e) { }
            this.DoExitService();
            lock (this.threadpool)
            {
                while (this.threadpool.Count() > 0)
                {
                    ThreadShell item = this.threadpool.Get();
                    item.DoDispose();
                }
            }
            this.threadpool.DoDispose();
            this.delegatepool.DoDispose();
        }
        protected abstract bool Validate();
        protected abstract void Begin();
        protected abstract void End();

    }





}
