package websocket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"im"

	"github.com/gobwas/ws"
	"github.com/sirupsen/logrus"
)

var _ im.Server = (*server)(nil)

type ServerOptions struct {
	im.Acceptor
	im.ChannelMap
	loginWait time.Duration
	readWait  time.Duration
	writeWait time.Duration
}

type ServerOption func(*ServerOptions)

func WithAcceptor(acceptor im.Acceptor) ServerOption {
	return func(o *ServerOptions) {
		o.Acceptor = acceptor
	}
}

func WithChannelMap(channelMap im.ChannelMap) ServerOption {
	return func(o *ServerOptions) {
		o.ChannelMap = channelMap
	}
}

type server struct {
	im.StateListener
	readWait time.Duration
	im.MessageListener
	addr string
	opts *ServerOptions
}

func (s server) Start() error {
	mux := http.NewServeMux()
	log := logrus.WithFields(logrus.Fields{
		"module": "ws.server",
		"listen": s.addr,
		//"id":     s.ServiceID(),
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// step 1
		rawconn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			resp(w, http.StatusBadRequest, err.Error())
			log.Error(err)
			return
		}

		// step 2 包装conn
		conn := NewConnection(rawconn)

		// step 3
		id, err := s.opts.Accept(conn, s.opts.loginWait)
		if err != nil {
			_ = conn.WriteFrame(im.OpClose, []byte(err.Error()))
			conn.Close()
			return
		}
		if _, ok := s.opts.ChannelMap.Get(id); ok {
			log.Warnf("channel %s existed", id)
			_ = conn.WriteFrame(im.OpClose, []byte("channelId is repeated"))
			conn.Close()
			return
		}

		// step 4
		channel := im.NewChannel(id, conn)
		channel.SetWriteWait(s.opts.writeWait)
		channel.SetReadWait(s.opts.readWait)
		s.opts.ChannelMap.Add(channel)

		go func(ch im.Channel) {
			// step 5
			err := ch.Readloop(s.MessageListener)
			if err != nil {
				log.Info(err)
			}
			// step 6
			s.opts.ChannelMap.Remove(ch.ID())
			err = s.Disconnect(ch.ID())
			if err != nil {
				log.Warn(err)
			}
			ch.Close()
		}(channel)
	})
	log.Infoln("started")
	return http.ListenAndServe(s.addr, mux)
}

func (s server) Shutdown(ctx context.Context) error {
	fmt.Println("ws server stop!")
	return nil
}

func NewServer(addr string, stateListener im.StateListener, readWait time.Duration, listener im.MessageListener, opt ...ServerOption) im.Server {
	srv := &server{
		StateListener:   stateListener,
		readWait:        readWait,
		MessageListener: listener,
		addr:            addr,
		opts: &ServerOptions{
			readWait:   im.DefaultReadWait,
			writeWait:  im.DefaultWriteWait,
			loginWait:  im.DefaultLoginWait,
			ChannelMap: im.DefaultChannelMap,
		},
	}
	for _, o := range opt {
		o(srv.opts)
	}

	if srv.opts.Acceptor == nil {
		//srv.opts.Acceptor =
		panic("acceptor can't be nil")
	}

	return srv
}

func resp(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	if body != "" {
		_, _ = w.Write([]byte(body))
	}
	logrus.Warnf("response with code:%d %s", code, body)
}
