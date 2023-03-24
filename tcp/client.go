package tcp

import (
	"fmt"
	"im"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type ClientOptions struct {
	Heartbeat time.Duration
	WriteWait time.Duration
}

var _ im.Client = (*Client)(nil)

type Client struct {
	sync.Mutex
	im.Dialer
	once  sync.Once
	id    string
	name  string
	conn  im.Connection
	state int32
	opts  ClientOptions
}

func (c *Client) ID() string {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Name() string {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Connect(addr string) error {
	_, err := url.Parse(addr)
	if err != nil {
		return err
	}
	// 这里是一个CAS原子操作，对比并设置值，是并发安全的。
	if !atomic.CompareAndSwapInt32(&c.state, 0, 1) {
		return fmt.Errorf("client has connected")
	}

	rawconn, err := c.Dialer.DialAndHandshake(im.DialerContext{
		Id:      c.id,
		Name:    c.name,
		Address: addr,
		Timeout: im.DefaultLoginWait,
	})
	if err != nil {
		atomic.CompareAndSwapInt32(&c.state, 1, 0)
		return err
	}
	if rawconn == nil {
		return fmt.Errorf("conn is nil")
	}
	c.conn = NewConnection(rawconn)

	if c.opts.Heartbeat > 0 {
		go func() {
			err := c.heartbealoop()
			if err != nil {
				logrus.WithField("module", "tcp.client").Warn("heartbealoop stopped - ", err)
			}
		}()
	}
	return nil
}

func (c *Client) SetDialer(dialer im.Dialer) {
	c.Dialer = dialer
}

func (c *Client) Send(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Read() (im.Frame, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Close() {
	//TODO implement me
	panic("implement me")
}

func (c *Client) heartbealoop() error {
	tick := time.NewTicker(c.opts.Heartbeat)
	for range tick.C {
		// 发送一个ping的心跳包给服务端
		if err := c.ping(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ping() error {
	logrus.WithField("module", "tcp.client").Tracef("%s send ping to server", c.id)

	err := c.conn.SetWriteDeadline(time.Now().Add(c.opts.WriteWait))
	if err != nil {
		return err
	}
	return c.conn.WriteFrame(im.OpPing, nil)
}
