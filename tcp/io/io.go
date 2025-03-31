package io

import (
	"bytes"
	"encoding/binary"

	//"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	g_HeadLen = 4
	g_Cach    = 1024
)

// tcp黏包，写包
func Write(conn net.Conn, message []byte) (int, error) {
	buf := make([]byte, g_HeadLen+len(message))
	binary.LittleEndian.PutUint32(buf[0:g_HeadLen], uint32(len(message)))
	copy(buf[g_HeadLen:], message)
	n, err := conn.Write(buf)
	if err != nil {
		return 0, err
	} else if n != len(buf) {
		return 0, fmt.Errorf("write %d less than %d", n, len(buf))
	}
	//log.Debug("write len %d", n)
	return n, err
}

// tcp黏包，读取包
func Read(conn net.Conn) (int, []byte, error) {
	l := make([]byte, g_HeadLen)
	_, err := io.ReadFull(conn, l)
	if err != nil {
		return 0, nil, err
	}

	length := binary.LittleEndian.Uint32(l)
	//if length > 1024 * 1024 {
	//	return 0, nil, errors.New("read buffer length too max")
	//}
	//log.Debug("read len %d", length)
	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return 0, nil, err
	} else {
		return 0, data, nil
	}
}

// 只读取1024的buf
func funcRead(conn net.Conn) (int, []byte, error) {
	buf := make([]byte, g_Cach)
	len, err := conn.Read(buf)
	if err != nil {
		return 0, nil, err
	}
	return len, buf, nil
}

// tcp黏包，写包
func WriteByTime(conn net.Conn, message []byte, t time.Time) (int, error) {
	err := conn.SetWriteDeadline(t)
	if err != nil {
		return 0, err
	}
	id, err := Write(conn, message)
	conn.SetWriteDeadline(time.Time{})
	return id, err
}

// tcp黏包，读取包  超时
func ReadByTime(conn net.Conn, t time.Time) (int, []byte, error) {
	err := conn.SetReadDeadline(t)
	if err != nil {
		return 0, nil, err
	}
	defer conn.SetReadDeadline(time.Time{})
	return Read(conn)
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
