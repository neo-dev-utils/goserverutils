/*
*作者：唐杰
*说明：通过gob读取数据
 */
package io

import (
	"encoding/json"
	"net"
	"time"
)

/*
*打包
 */
func WriteJsonByTime(conn net.Conn, inter interface{}, t time.Time) error {
	buf, err := json.Marshal(inter)
	if err != nil {
		return err
	}
	_, err = WriteByTime(conn, buf, t)
	if err != nil {
		return err
	}

	return nil
}

/*
*解包
 */
func ReadJsonByTime(conn net.Conn, inter interface{}, t time.Time) error {
	_, buf, err := ReadByTime(conn, t)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, inter)
	if err != nil {
		return err
	}
	return nil
}

/*
*解包
 */
func ReadJson(conn net.Conn, inter interface{}) error {
	_, buf, err := Read(conn)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, inter)
	if err != nil {
		return err
	}
	return nil
}
