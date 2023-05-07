package session

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"runtime/debug"
)

type Session struct {
	conn net.Conn
}

func newSession(conn net.Conn) *Session {
	return &Session{
		conn: conn,
	}
}

func (s *Session) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("panic: %v", r)
			debug.PrintStack()
			s.conn.Close()
		}
	}()

}
