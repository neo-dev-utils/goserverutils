package io

import (
	"fmt"
	"gocommutils/timeUtil"
	"net"
	"os"
	"testing"
	"time"
)

func Test_Read(t *testing.T) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 9999))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Waiting for clients")
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		fmt.Println(conn.RemoteAddr().String(), " tcp connect success")
		go func(con net.Conn) {
			for {
				l, body, err := Read(con)
				if err != nil {
					fmt.Println("err:", err.Error())
					conn.Close()
					return
				}
				fmt.Printf("read msg len:%d,body:%s\n", l, string(body))
			}

		}(conn)

	}
}

func Test_Write(t *testing.T) {
	server := "127.0.0.1:9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	defer conn.Close()
	fmt.Println("connect success")
	for {
		body := "Hello Server"
		l, err := Write(conn, []byte(body))
		if err != nil {
			fmt.Println("err:", err.Error())
			return
		}
		fmt.Printf("write msg len:%d,body:%v\n", l, body)
		time.Sleep(time.Microsecond * 30)
	}
}

type Data struct {
	Int    int
	String string
	Inter  []interface{}
}

func Test_ReadByTime(t *testing.T) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", 9999))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Waiting for clients")
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		fmt.Println(conn.RemoteAddr().String(), " tcp connect success")
		go func(con net.Conn) {
			for {
				data := Data{}
				err = ReadGobByTime(con, &data, timeUtil.Time.NowAddSeconds(10))
				if err != nil {
					fmt.Println("err:", err.Error())
					conn.Close()
					return
				}
				fmt.Printf("read msg body:%v\n", data)

			}

		}(conn)

	}
}

func Test_WirteByTime(t *testing.T) {
	server := "127.0.0.1:9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	defer conn.Close()
	fmt.Println("connect success")
	n := 1

	for {
		go func(num int) {
			aaa := make([]interface{}, 1)
			aaa[0] = num
			data := &Data{
				Int:    num,
				String: "Hello Server",
				Inter:  aaa,
			}
			err := WriteGobByTime(conn, data, timeUtil.Time.NowAddSeconds(10))
			if err != nil {
				fmt.Println("err:", err.Error())
				conn.Close()
				return
			}
			fmt.Printf("write msg data:%v\n", data)
		}(n)
		go func(num int) {
			aaa := make([]interface{}, 1)
			aaa[0] = num
			data := &Data{
				Int:    num,
				String: "Hello Server",
				Inter:  aaa,
			}
			err := WriteGobByTime(conn, data, timeUtil.Time.NowAddSeconds(10))
			if err != nil {
				fmt.Println("err:", err.Error())
				conn.Close()
				return
			}
			fmt.Printf("write msg data:%v\n", data)
		}(n + 1)
		n++
		time.Sleep(time.Millisecond * 10)
	}
}
