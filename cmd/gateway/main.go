package main

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/gateway/clients/auth"
	blocknote "github.com/autumnterror/breezynotes/internal/gateway/clients/blocknote"
	redis "github.com/autumnterror/breezynotes/internal/gateway/clients/redis"
	"github.com/autumnterror/breezynotes/internal/gateway/config"
	"github.com/autumnterror/breezynotes/pkg/log"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.gateway"
	cfg := config.MustSetup()

	g, _ := errgroup.WithContext(context.Background())

	var a *auth.Client
	var r *redis.Client
	var b *blocknote.Client
	var err error

	g.Go(func() error {
		a, err = auth.New(cfg)
		return err
	})
	g.Go(func() error {
		r, err = redis.New(cfg)
		return err
	})
	g.Go(func() error {
		b, err = blocknote.New(cfg)
		return err
	})

	ch(g.Wait())

	if _, err := a.API.Healthz(context.Background(), nil); err == nil {
		log.Success(op, "health auth")
	}
	if _, err := r.API.Healthz(context.Background(), nil); err == nil {
		log.Success(op, "health redis")
	}
	if _, err := b.API.Healthz(context.Background(), nil); err == nil {
		log.Success(op, "health blocknote")
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}

func ch(err error) {
	if err != nil {
		log.Panic(err)
	}
}
