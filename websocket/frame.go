package websocket

import (
	"im"

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
