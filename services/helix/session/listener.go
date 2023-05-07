package session

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"syscall"
)

func StartListener(ctx context.Context, addr string) {
	setFDLimit()

	lc := net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Errorf("accept err: %v", err)
			return
		}

		session := newSession(conn)
		go session.Run(ctx)
	}
}

func setFDLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	logrus.Infof("Set File Descriptor Limit: %d", rLimit.Cur)
}
