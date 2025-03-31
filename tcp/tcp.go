/*
*作者：唐杰
*说明：启动一个简单的tcp服务
 */
package tcp

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"goserverutils/tcp/inter"
)

/*
*tcp服务
 */
type TcpServer struct {
	ip     string       // 监听的IP地址
	port   int          // 监听的端口号
	count  int          // 在线的链接数
	pid    int          // 进程pid
	close  chan bool    // 服务状态
	server inter.Server // 服务接口
	status bool         //
}

/*
*创建一个TCP客户端
 */
func NewTcpClient(url string) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", url)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	// conn.SetKeepAlive(true)
	// conn.SetKeepAlivePeriod(time.Second * 250)
	return conn, nil
}

/*
*创建一个TCP服务
 */
func NewTcpServer(ip string, port int, server inter.Server) *TcpServer {
	i := net.ParseIP(ip)
	if i == nil {
		panic("ip error")
	}
	if port <= 0 {
		panic("port error")
	}
	srv := &TcpServer{
		ip:     ip,
		port:   port,
		close:  make(chan bool),
		pid:    0,
		server: server,
		status: false,
	}
	return srv
}

/*
*启动一个TCP服务
 */
func (ts *TcpServer) Run() {
	ts.pid = os.Getpid()
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ts.ip, ts.port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func(m chan os.Signal, srv *TcpServer) {
		msg := <-m
		err := fmt.Sprintf("recv signal %v", msg)
		srv.server.Error(errors.New(err))
		srv.Close()
	}(c, ts)

	//监听关闭
	go func(listen net.Listener, srv *TcpServer) {
		<-srv.close
		//此时开始做结束需要的收尾工作
		ts.server.Exit()
		if err = listen.Close(); err != nil {
			//如果出错强行结束程序
			os.Exit(1)
		}
		close(srv.close)
	}(l, ts)
	ts.status = true
	for {
		conn, err := l.Accept()
		if err != nil {
			if ts.status {
				ts.server.Error(err)
				continue
			}
			return
		}
		go ts.server.Accept(conn)
	}
}

/*
*关闭服务
 */
func (ts *TcpServer) Close() {
	ts.status = false
	ts.close <- true
}

/*
*获取服务的PID
 */
func (ts *TcpServer) GetPid() int {
	return ts.pid
}

/*
*获取服务连接数量
 */
func (ts *TcpServer) GetCount() int {
	return ts.count
}
