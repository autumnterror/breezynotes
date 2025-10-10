package main

import (
	"fmt"
	"github.com/autumnterror/breezynotes/internal/redis/config"
	"github.com/autumnterror/breezynotes/internal/redis/grpc"
	"github.com/autumnterror/breezynotes/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.redis"
	cfg := config.MustSetup()

	a := grpc.New(cfg)
	go a.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	a.Stop()

	log.Success(op, "stop signal "+fmt.Sprint(sign))
}
