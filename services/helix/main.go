package main

import (
	"context"
	"github.com/unsafe9/mado/helix/session"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	ctx := context.Background()

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	session.StartListener(ctx, ":3000")
}
