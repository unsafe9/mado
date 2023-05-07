package session

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type Beacon struct {
	sessions   map[uint32]*Session
	joins      chan *Session
	multicasts chan request
}
type request struct {
	sid uint32
	msg proto.Message
}

func NewBeacon() *Beacon {
	return &Beacon{
		sessions:   make(map[uint32]*Session, 0),
		joins:      make(chan *Session),
		multicasts: make(chan request),
	}
}

func (b *Beacon) Run(ctx context.Context) {
	for {
		select {
		case s := <-b.joins:
			// TODO : 연결 끊어지면 삭제
			b.sessions[s.sid] = s

		case req := <-b.multicasts:
			b.handleMulticasts(req)

		case <-ctx.Done():
			logrus.Debug("beacon done")
			return
		}
	}
}

func (b *Beacon) handleMulticasts(req request) {
	for _, s := range b.sessions {
		if s.sid == req.sid {
			continue
		}
		if err := s.Send(req.msg); err != nil {
			logrus.Warnf("send err: %v", err)
		}
	}
}

func (b *Beacon) Join(s *Session) {
	b.joins <- s
}

func (b *Beacon) Broadcast(s *Session, m proto.Message) {
	b.multicasts <- request{
		sid: s.sid,
		msg: m,
	}
}
