package main

import (
	"fmt"
	"github.com/autumnterror/breezynotes/internal/redis/api"
	"github.com/autumnterror/breezynotes/internal/redis/config"
	"github.com/autumnterror/breezynotes/internal/redis/repository"
	"github.com/autumnterror/utils_go/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.redis"

	cfg := config.MustSetup()

	a := api.New(cfg, repository.New(cfg))
	go a.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	a.Stop()

	log.Success(op, "stop signal "+fmt.Sprint(sign))
}
