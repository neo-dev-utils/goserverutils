package tcp

import (
	"gocommutils/timeUtil"
	"goserverutils/tcp/io"

	"fmt"
	"net"
	"testing"
	"time"
)

type TestServer struct {
	server *TcpServer
}

func Test_Server(t *testing.T) {
	testServer := &TestServer{}
	server := NewTcpServer("0.0.0.0", 9999, testServer)
	testServer.server = server
	fmt.Println("server run")
	server.Run()
}

func (ts *TestServer) Accept(conn net.Conn) {
	for {
		_, recvData, err := io.Read(conn)
		if err != nil {
			conn.Close()
			return
		}
		fmt.Println(string(recvData))
	}
	fmt.Println("recv accept")
}

func (ts *TestServer) Close(conn net.Conn) {
	fmt.Println("close")
	//ts.server.Close()
}
func (ts *TestServer) Error(err error) {
	fmt.Println("recv err ", err)
}

func (ts *TestServer) Exit() {
	fmt.Println("程序结束开始善后处理 ")
}

func Test_Client(t *testing.T) {
	conn, err := NewTcpClient("127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	n := 0
	t1 := timeUtil.Time.NowTime()
	for {
		_, err := io.Write(conn, []byte(fmt.Sprintf("hello%d", n)))
		if err != nil {
			fmt.Println(err)
			return
		}
		n++
		if (time.Since(t1).Nanoseconds()) >= 1000000 {
			return
		}
		//time.Sleep(time.Microsecond)
	}
}
