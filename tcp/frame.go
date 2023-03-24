package tcp

import "im"

var _ im.Frame = (*FrameImpl)(nil)

type FrameImpl struct {
	OpCode  im.OpCode
	Payload []byte
}

func (f *FrameImpl) SetOpCode(code im.OpCode) {
	//TODO implement me
	panic("implement me")
}

func (f *FrameImpl) GetOpCode() im.OpCode {
	//TODO implement me
	panic("implement me")
}

func (f *FrameImpl) SetPayload(payload []byte) {
	//TODO implement me
	panic("implement me")
}

func (f *FrameImpl) GetPayload() []byte {
	//TODO implement me
	panic("implement me")
}
