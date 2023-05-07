package main

import (
	"github.com/sirupsen/logrus"
	"github.com/unsafe9/mado/services/helix/message"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"net"
	"time"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	var conn net.Conn

	go func() {
		pos := message.C2SUpdatePosition{
			Type: message.MessageType_C2SUpdatePosition,
			X:    rand.Int31n(100),
			Y:    rand.Int31n(100),
			Z:    rand.Int31n(100),
		}

		for {
			if conn == nil {
				time.Sleep(time.Second)
				continue
			}
			raw, err := proto.Marshal(&pos)
			if err != nil {
				logrus.Warnf("marshal err: %v", err)
				continue
			}
			if _, err := conn.Write(raw); err != nil {
				logrus.Warnf("write err: %v", err)
				conn = nil
				continue
			}
			pos.X += 1
			//logrus.Debugf("write: %s", pos.String())
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			if conn == nil {
				time.Sleep(time.Second)
				continue
			}

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				logrus.Warnf("read err: %v", err)
				continue
			}
			var typ message.MessageType
			if err := proto.Unmarshal(buf[:n], &typ); err != nil {
				logrus.Warnf("unmarshal err: %v", err)
				continue
			}

			switch typ.GetType() {
			case message.MessageType_S2CUpdatePosition:
				var msg message.S2CUpdatePosition
				if err := proto.Unmarshal(buf[:n], &msg); err != nil {
					logrus.Warnf("unmarshal err: %v", err)
					continue
				}
				logrus.Debugf("pos: %s", msg.String())
			}
		}
	}()

	for {
		time.Sleep(time.Second)
		if conn == nil {
			var err error
			conn, err = net.Dial("tcp", "localhost:3000")
			if err != nil {
				logrus.Warnf("dial err: %v", err)
				continue
			}
		}
	}

}
