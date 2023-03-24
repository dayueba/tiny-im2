package im

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type ChannelMap interface {
	Add(channel Channel)
	Remove(id string)
	Get(id string) (Channel, bool)
	All() []Channel
}

var DefaultChannelMap ChannelMap

func init() {
	DefaultChannelMap = &ChannelMapImpl{
		channels: new(sync.Map),
	}
}

type ChannelMapImpl struct {
	channels *sync.Map
}

func (c *ChannelMapImpl) Add(channel Channel) {
	if channel.ID() == "" {
		logrus.WithFields(logrus.Fields{
			"module": "ChannelMapImpl",
		}).Error("channel id is required")
	}

	c.channels.Store(channel.ID(), channel)
}

func (c *ChannelMapImpl) Remove(id string) {
	c.channels.Delete(id)
}

func (c *ChannelMapImpl) Get(id string) (Channel, bool) {
	if id == "" {
		logrus.WithFields(logrus.Fields{
			"module": "ChannelsImpl",
		}).Error("channel id is required")
	}

	val, ok := c.channels.Load(id)
	if !ok {
		return nil, false
	}
	return val.(Channel), true
}

func (c *ChannelMapImpl) All() []Channel {
	arr := make([]Channel, 0)
	c.channels.Range(func(key, val interface{}) bool {
		arr = append(arr, val.(Channel))
		return true
	})
	return arr
}

type Channel interface {
	Connection
	Agent
	// Close 关闭连接
	Close() error
	Readloop(lst MessageListener) error
	// SetWriteWait 设置写超时
	SetWriteWait(time.Duration)
	SetReadWait(time.Duration)
}

var _ Channel = (*ChannelImpl)(nil)

// ChannelImpl is a websocket implement of channel
// todo why is a websocket implement? where is the tcp implement
type ChannelImpl struct {
	sync.Mutex
	id string
	Connection
	writechan chan []byte
	once      sync.Once
	writeWait time.Duration
	readWait  time.Duration
	//closed    *Event
	state int32 // 0 init 1 start 2 closed
}

func (ch *ChannelImpl) ID() string {
	return ch.id
}

func (ch *ChannelImpl) Push(payload []byte) error {
	if atomic.LoadInt32(&ch.state) != 1 {
		return fmt.Errorf("channel %s has closed", ch.id)
	}
	// 异步写
	ch.writechan <- payload
	return nil
}

func (ch *ChannelImpl) Readloop(lst MessageListener) error {
	if !atomic.CompareAndSwapInt32(&ch.state, 0, 1) {
		return fmt.Errorf("channel has started")
	}
	log := logrus.WithFields(logrus.Fields{
		"struct": "ChannelImpl",
		"func":   "Readloop",
		"id":     ch.id,
	})
	for {
		_ = ch.SetReadDeadline(time.Now().Add(ch.readWait))

		frame, err := ch.ReadFrame()
		if err != nil {
			log.Info(err)
			return err
		}
		if frame.GetOpCode() == OpClose {
			return errors.New("remote side close the channel")
		}
		if frame.GetOpCode() == OpPing {
			log.Trace("recv a ping; resp with a pong")

			_ = ch.WriteFrame(OpPong, nil)
			_ = ch.Flush()
			continue
		}
		payload := frame.GetPayload()
		if len(payload) == 0 {
			continue
		}
		//err = ch.gpool.Submit(func() {
		//	lst.Receive(ch, payload)
		//})
		//if err != nil {
		//	return err
		//}
		go lst.Receive(ch, payload)
	}
}

func (ch *ChannelImpl) SetWriteWait(writeWait time.Duration) {
	if writeWait == 0 {
		return
	}
	ch.writeWait = writeWait
}

func (ch *ChannelImpl) SetReadWait(readWait time.Duration) {
	if readWait == 0 {
		return
	}
	ch.readWait = readWait
}

func NewChannel(id string, conn Connection) Channel {
	log := logrus.WithFields(logrus.Fields{
		"module": "tcp_channel",
		"id":     id,
	})
	ch := &ChannelImpl{
		id:         id,
		Connection: conn,
		writechan:  make(chan []byte, 5),
		//closed:    NewEvent(),
		writeWait: time.Second * 10, //default value
	}
	go func() {
		err := ch.writeloop()
		if err != nil {
			log.Info(err)
		}
	}()
	return ch
}

func (ch *ChannelImpl) writeloop() error {
	log := logrus.WithFields(logrus.Fields{
		"module": "ChannelImpl",
		"func":   "writeloop",
		"id":     ch.id,
	})
	defer func() {
		log.Debugf("channel %s writeloop exited", ch.id)
	}()
	for payload := range ch.writechan {
		err := ch.WriteFrame(OpBinary, payload)
		if err != nil {
			return err
		}
		chanlen := len(ch.writechan)
		for i := 0; i < chanlen; i++ {
			payload = <-ch.writechan
			err := ch.WriteFrame(OpBinary, payload)
			if err != nil {
				return err
			}
		}
		err = ch.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ch *ChannelImpl) WriteFrame(code OpCode, payload []byte) error {
	_ = ch.Connection.SetWriteDeadline(time.Now().Add(ch.writeWait))
	return ch.Connection.WriteFrame(code, payload)
}
