package main

import (
	"github.com/autumnterror/breezynotes/internal/auth/config"
	"github.com/autumnterror/breezynotes/internal/auth/grpc"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/internal/auth/psql"
	"github.com/autumnterror/breezynotes/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.auth"
	cfg := config.MustSetup()

	j := jwt.NewWithConfig(cfg)

	db := psql.MustConnect(cfg)

	a := grpc.New(cfg, psql.NewDriver(db.Driver), j)
	go a.MustRun()

	sign := wait()

	a.Stop()
	if err := db.Disconnect(); err != nil {
		log.Error(op, "db disconnect", err)
	}

	log.Success(op, "stop signal "+sign)
}

func wait() string {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	return sign.String()
}
