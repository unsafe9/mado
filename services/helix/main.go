package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/unsafe9/mado/services/helix/session"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	ctx := context.Background()
	logrus.SetLevel(logrus.DebugLevel)

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	beacon := session.NewBeacon()
	go beacon.Run(ctx)
	l := &session.Listener{
		Beacon: beacon,
	}
	l.SetFDLimit(65535)
	l.Run(ctx, ":3000")
}
