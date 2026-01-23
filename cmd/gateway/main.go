package main

import (
	"context"
	"fmt"
	_ "github.com/autumnterror/breezynotes/docs"
	"github.com/autumnterror/breezynotes/internal/gateway/clients/auth"
	blocknote "github.com/autumnterror/breezynotes/internal/gateway/clients/blocknote"
	redis "github.com/autumnterror/breezynotes/internal/gateway/clients/redis"
	"github.com/autumnterror/breezynotes/internal/gateway/config"
	"github.com/autumnterror/breezynotes/internal/gateway/net"
	"github.com/autumnterror/utils_go/pkg/log"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Breezy notes gateway REST UserAPI
// @version 0.1-.-infDev
// @description Full API for BreezyNotes.
// @termsOfService https://about.breezynotes.ru/

// @contact.name Alex "bustard" Provor
// @contact.url https://contacts.breezynotes.ru
// @contact.email help@breezynotes.ru

// @host localhost:8080
// @BasePath /
// @schemes http
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

	health(a, r, b)

	e := net.New(cfg, a, b, r)
	go e.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	s := <-stop

	if err := e.Stop(); err != nil {
		log.Error(op, "stop echo", err)
	}

	log.Green(fmt.Sprintf("%s:%s", op, s.String()))
}

func health(
	a *auth.Client,
	r *redis.Client,
	b *blocknote.Client,
) {
	const op = "cmd.gateway.health"
	ctx, done := context.WithTimeout(context.Background(), 6*time.Second)
	defer done()
	if _, err := a.API.Healthz(ctx, nil); err == nil {
		log.Success(op, "health auth")
	} else {
		log.Error(op, "", err)
	}
	if _, err := r.API.Healthz(ctx, nil); err == nil {
		log.Success(op, "health redis")
	}
	if _, err := b.API.Healthz(ctx, nil); err == nil {
		log.Success(op, "health blocknote")
	}
}

func ch(err error) {
	if err != nil {
		log.Panic(err)
	}
}
