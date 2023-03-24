package websocket

import (
	"im"
	"net"

	"github.com/gobwas/ws"
)

var _ im.Connection = (*ConnectionImpl)(nil)

type ConnectionImpl struct {
	net.Conn
}

func (c *ConnectionImpl) ReadFrame() (im.Frame, error) {
	f, err := ws.ReadFrame(c.Conn)
	if err != nil {
		return nil, err
	}
	return &FrameImpl{raw: f}, nil
}

func (c *ConnectionImpl) WriteFrame(code im.OpCode, payload []byte) error {
	f := ws.NewFrame(ws.OpCode(code), true, payload)
	return ws.WriteFrame(c.Conn, f)
}

func (c *ConnectionImpl) Flush() error {
	return nil
}

func NewConnection(conn net.Conn) im.Connection {
	return &ConnectionImpl{conn}
}
