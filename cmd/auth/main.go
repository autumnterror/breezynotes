package main

import (
	api "github.com/autumnterror/breezynotes/internal/auth/api"
	"github.com/autumnterror/breezynotes/internal/auth/config"
	"github.com/autumnterror/breezynotes/internal/auth/infra/psql"
	"github.com/autumnterror/breezynotes/internal/auth/infra/psql/psqltx"
	"github.com/autumnterror/breezynotes/internal/auth/jwt"
	"github.com/autumnterror/breezynotes/internal/auth/service"
	"github.com/autumnterror/utils_go/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.auth"

	cfg := config.MustSetup()

	j := jwt.NewWithConfig(cfg)

	db := psql.MustConnect(cfg)

	s := service.NewAuthService(
		psqltx.NewTxRunner(db.Driver),
		psqltx.NewRepoProvider(db.Driver),
		j,
		cfg,
	)

	a := api.New(cfg, s)
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
