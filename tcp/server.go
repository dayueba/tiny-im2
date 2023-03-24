package tcp

import (
	"context"
	"errors"
	"im"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type ServerOptions struct {
	im.Acceptor
	loginWait time.Duration
	readWait  time.Duration
	writeWait time.Duration
	im.ChannelMap
	im.MessageListener
	im.StateListener
}

var _ im.Server = (*server)(nil)

type server struct {
	addr string
	opts *ServerOptions
}

func (s *server) Start() error {
	log := logrus.WithFields(logrus.Fields{
		"module": "tcp.server",
		"listen": s.addr,
		//"id":     s.ServiceID(),
	})

	lst, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := lst.Accept()
		if err != nil {
			conn.Close()
			log.Warn(err)
			continue
		}

		go func(conn net.Conn) {
			connection := NewConnection(conn)
			// step 3
			id, err := s.opts.Accept(connection, s.opts.loginWait)
			if err != nil {
				_ = connection.WriteFrame(im.OpClose, []byte(err.Error()))
				conn.Close()
				return
			}
			if _, ok := s.opts.ChannelMap.Get(id); ok {
				log.Warnf("channel %s existed", id)
				_ = connection.WriteFrame(im.OpClose, []byte("channelId is repeated"))
				connection.Close()
				return
			}
			//step 4
			channel := im.NewChannel(id, connection)
			channel.SetReadWait(s.opts.readWait)
			channel.SetWriteWait(s.opts.writeWait)
			s.opts.ChannelMap.Add(channel)

			log.Info("accept ", channel)
			//step 5
			err = channel.Readloop(s.opts.MessageListener)
			if err != nil {
				log.Info(err)
			}
			// step 6
			s.opts.ChannelMap.Remove(channel.ID())
			_ = s.opts.StateListener.Disconnect(channel.ID())
			channel.Close()
		}(conn)
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *server) Push(id string, data []byte) error {
	ch, ok := s.opts.ChannelMap.Get(id)
	if !ok {
		return errors.New("channel no found")
	}
	return ch.Push(data)
}

func NewServer() im.Server {
	srv := &server{}

	return srv
}
