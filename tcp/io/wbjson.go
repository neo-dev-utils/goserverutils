package io

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

func WbReadJson(conn *websocket.Conn, inter interface{}) error {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	err = json.Unmarshal(message, inter)
	if err != nil {
		return err
	}
	return nil
}

func WbReadJsonByTime(conn *websocket.Conn, inter interface{}, t time.Time) error {
	err := conn.SetReadDeadline(t)
	if err != nil {
		return err
	}
	defer conn.SetReadDeadline(time.Time{})

	_, message, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	err = json.Unmarshal(message, inter)
	if err != nil {
		return err
	}
	return nil
}

func WbWriteJson(conn *websocket.Conn, inter interface{}) error {
	message, err := json.Marshal(inter)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func WbWriteJsonByTime(conn *websocket.Conn, inter interface{}, t time.Time) error {
	err := conn.SetWriteDeadline(t)
	if err != nil {
		return err
	}
	defer conn.SetWriteDeadline(time.Time{})

	message, err := json.Marshal(inter)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}
