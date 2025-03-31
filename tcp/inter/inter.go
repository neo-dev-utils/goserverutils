package inter

import (
	"net"
)

// 定义一个协程接口对象
type Work interface {
	Work(pid int)
}

// 定义一个服务接口对象
type Server interface {
	Accept(net.Conn) //接受到链接
	Error(error)     //出错
	Close(conn net.Conn)
	Exit() //服务退出需要做的处理
}
