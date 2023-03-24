package main

import (
	"context"
	"errors"
	"im"
	"im/websocket"
	"time"
)

// ServerHandler ServerHandler
type ServerHandler struct {
}

// Accept this connection
func (h *ServerHandler) Accept(conn im.Connection, timeout time.Duration) (string, error) {
	// 1. 读取：客户端发送的鉴权数据包
	frame, err := conn.ReadFrame()
	if err != nil {
		return "", err
	}
	// 2. 解析：数据包内容就是userId
	userID := string(frame.GetPayload())
	// 3. 鉴权：这里只是为了示例做一个fake验证，非空
	if userID == "" {
		return "", errors.New("user id is invalid")
	}
	return userID, nil
}

// Receive default listener
func (h *ServerHandler) Receive(ag im.Agent, payload []byte) {
	ack := string(payload) + " from server "
	_ = ag.Push([]byte(ack))
}

// Disconnect default listener
func (h *ServerHandler) Disconnect(id string) error {
	//logger.Warnf("disconnect %s", id)
	return nil
}

func main() {
	handler := &ServerHandler{}
	service := websocket.NewServer(":8080", handler, time.Minute, handler, websocket.WithAcceptor(handler))
	service.Start()
	defer service.Shutdown(context.Background())
}
