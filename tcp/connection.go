package tcp

import (
	"im"
	"net"
)

var _ im.Connection = (*ConnectionImpl)(nil)

type ConnectionImpl struct {
	net.Conn
}

func NewConnection(conn net.Conn) im.Connection {
	return &ConnectionImpl{conn}
}

// ReadFrame 不是一个线程安全的方法
func (c *ConnectionImpl) ReadFrame() (im.Frame, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConnectionImpl) WriteFrame(code im.OpCode, payload []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConnectionImpl) Flush() error {
	return nil
}
