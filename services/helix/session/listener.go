package session

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"syscall"
)

type Listener struct {
	Beacon *Beacon
	sid    uint32
}

func (l *Listener) Run(ctx context.Context, addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	logrus.Infof("listening on %s", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Errorf("accept err: %v", err)
			return
		}

		l.sid++
		session := &Session{
			sid:    l.sid,
			conn:   conn,
			Beacon: l.Beacon,
		}
		go session.Run(ctx)
		l.Beacon.Join(session)
	}
}

func (l *Listener) SetFDLimit(max int) {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	if max < 0 {
		rLimit.Cur = rLimit.Max
	} else {
		rLimit.Cur = uint64(max)
	}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	logrus.Infof("Set File Descriptor Limit: %d", rLimit.Cur)
}
