package websocket

import (
	"im"
	"net"

	"github.com/gobwas/ws"
)

var _ im.Frame = (*FrameImpl)(nil)

type FrameImpl struct {
	raw ws.Frame
}

func (f *FrameImpl) SetOpCode(code im.OpCode) {
	f.raw.Header.OpCode = ws.OpCode(code)
}

func (f *FrameImpl) GetOpCode() im.OpCode {
	return im.OpCode(f.raw.Header.OpCode)
}

func (f *FrameImpl) SetPayload(payload []byte) {
	f.raw.Payload = payload
}

func (f *FrameImpl) GetPayload() []byte {
	if f.raw.Header.Masked {
		ws.Cipher(f.raw.Payload, f.raw.Header.Mask, 0)
	}
	f.raw.Header.Masked = false
	return f.raw.Payload
}

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
