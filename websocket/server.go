package websocket

import (
	"context"
	"fmt"
	"time"

	"im"
)

type server struct {
	im.Acceptor
	im.StateListener
	readWait time.Duration
	im.ChannelMap
	im.MessageListener
}

func (s server) Start() error {
	fmt.Println("ws server start!")
	return nil
}

func (s server) Shutdown(ctx context.Context) error {
	fmt.Println("ws server stop!")
	return nil
}

func NewServer(acceptor im.Acceptor, stateListener im.StateListener, readWait time.Duration, listener im.MessageListener) im.Server {
	return &server{
		Acceptor:      acceptor,
		StateListener: stateListener,
		readWait:      readWait,
		//ChannelMap:      channelMap,
		MessageListener: listener,
	}
}
