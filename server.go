package im

import (
	"context"
	"time"
)

type Server interface {
	Start() error
	Shutdown(context.Context) error
}

type Acceptor interface {
	Accept(Connection, time.Duration) (string, error)
}

type StateListener interface {
	Disconnect(string) error
}

type MessageListener interface {
	Receive(Agent, []byte)
}

type Agent interface {
	ID() string
	Push([]byte) error
}

//func NewServer(protocol string, acceptor Acceptor, stateListener StateListener, readWait time.Duration, listener MessageListener) (Server, error) {
//	var server Server
//	if protocol == "ws" {
//		server = websocket.NewServer(acceptor, stateListener, readWait, listener)
//	} else if protocol == "tcp" {
//
//	} else {
//		return nil, errors.New("不支持此种协议")
//	}
//
//	return server, nil
//}
