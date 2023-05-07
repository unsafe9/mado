package session

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/unsafe9/mado/services/helix/message"
	"google.golang.org/protobuf/proto"
	"net"
	"runtime/debug"
)

type Session struct {
	sid    uint32
	conn   net.Conn
	Beacon *Beacon
}

func (s *Session) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("panic: %v", r)
			debug.PrintStack()
		}
		s.conn.Close()
	}()

	logrus.Debugf("session %d joined", s.sid)

	for {
		// Read message content
		buf := make([]byte, 1024)
		n, err := s.conn.Read(buf)
		if err != nil {
			logrus.Errorf("read msg err: %v", err)
			return
		}

		var typ message.MessageType
		if err := proto.Unmarshal(buf[:n], &typ); err != nil {
			logrus.Errorf("unmarshal err: %v", err)
			return
		}

		switch typ.GetType() {
		case message.MessageType_C2SUpdatePosition:
			var msg message.C2SUpdatePosition
			if err := proto.Unmarshal(buf[:n], &msg); err != nil {
				logrus.Errorf("unmarshal err: %v", err)
				return
			}
			logrus.Debugf("session %d: %s", s.sid, msg.String())
			s.Beacon.Broadcast(s, &message.S2CUpdatePosition{
				Type: message.MessageType_S2CUpdatePosition,
				ID:   s.sid,
				X:    msg.X,
				Y:    msg.Y,
				Z:    msg.Z,
			})
		}
	}
}

func (s *Session) Send(m proto.Message) error {
	raw, err := proto.Marshal(m)
	if err != nil {
		logrus.Errorf("marshal err: %v", err)
		return err
	}
	_, err = s.conn.Write(raw)
	return err
}
