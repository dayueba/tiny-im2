package main

import (
	"context"
	"im"
	"im/websocket"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
)

type WebsocketDialer struct {
}

// DialAndHandshake DialAndHandshake
func (d *WebsocketDialer) DialAndHandshake(ctx im.DialerContext) (net.Conn, error) {
	// 1 调用ws.Dial拨号
	conn, _, _, err := ws.Dial(context.TODO(), ctx.Address)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	// 2. 发送用户认证信息，示例就是userid
	err = wsutil.WriteClientBinary(conn, []byte(ctx.Id))
	if err != nil {
		return nil, err
	}
	// 3. return conn
	return conn, nil
}

func main() {
	cli := websocket.NewClient("001", "client1", websocket.ClientOptions{})
	cli.SetDialer(&WebsocketDialer{})

	// step2: 建立连接
	err := cli.Connect("ws://127.0.0.1:8080")
	if err != nil {
		logrus.Error(err)
	}
	count := 5
	go func() {
		// step3: 发送消息然后退出
		for i := 0; i < count; i++ {
			err := cli.Send([]byte("hello"))
			if err != nil {
				logrus.Error(err)
				return
			}
			time.Sleep(time.Millisecond * 10)
		}
	}()

	// step4: 接收消息
	recv := 0
	for {
		frame, err := cli.Read()
		if err != nil {
			logrus.Info(err)
			break
		}
		if frame.GetOpCode() != im.OpBinary {
			continue
		}
		recv++
		logrus.Infof("%s receive message [%s]", cli.ID(), frame.GetPayload())
		if recv == count { // 接收完消息
			break
		}
	}
	//退出
	cli.Close()
}
