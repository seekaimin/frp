package servicehelper

const (
	//尚未启动
	UnStarted = 0
	//正在运行
	Running = 1
	//已经停止
	Stoped = 2
)

//服务接口
type IService interface {
	//启动服务
	Start()
	//停止服务
	Stop()
	//验证
	Validate() bool
}
