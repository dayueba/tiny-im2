package im

import (
	"context"
	"time"
)

const (
	DefaultReadWait  = time.Minute * 3
	DefaultWriteWait = time.Second * 10
	DefaultLoginWait = time.Second * 10
	DefaultHeartbeat = time.Second * 55
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
