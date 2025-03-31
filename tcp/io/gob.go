/*
*作者：唐杰
*说明：通过gob读取数据
 */
package io

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
	// "goserverutils/log"
)

func RegisterType(value interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("register error")
		}
	}()
	gob.Register(value)
	return

}

func GobEncode(inter interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(inter); err != nil {
		// log.Debug(err.Error())
		return nil, err
	}
	return buf.Bytes(), nil
}

/*
*打包
 */
func WriteGobByTime(conn net.Conn, inter interface{}, t time.Time) error {
	buf, err := GobEncode(inter)
	if err != nil {
		// log.Debug(err.Error())
		return err
	}
	_, err = WriteByTime(conn, buf, t)
	return err
}

func GobDecode(data []byte, inter interface{}) error {
	var buf = bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(inter); err != nil {
		// log.Debug(err.Error())
		return err
	}
	return nil
}

/*
*解包
 */
func ReadGobByTime(conn net.Conn, inter interface{}, t time.Time) error {
	_, data, err := ReadByTime(conn, t)
	if err != nil {
		return err
	}
	return GobDecode(data, inter)
}

/*
*解包
 */
func ReadGob(conn net.Conn, inter interface{}) error {
	_, data, err := Read(conn)
	if err != nil {
		return err
	}
	var buf = bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(inter); err != nil {
		// log.Debug(err.Error())
		return err
	}
	return nil
}
