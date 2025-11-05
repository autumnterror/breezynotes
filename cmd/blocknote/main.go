package main

import (
	"fmt"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/grpc"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/notes"
	"github.com/autumnterror/breezynotes/internal/blocknote/mongo/tags"
	"github.com/autumnterror/breezynotes/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.blocknote"
	cfg := config.MustSetup()

	m := mongo.MustConnect(cfg)

	g := grpc.New(cfg, tags.NewApi(m), notes.NewApi(m))
	go g.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	g.Stop()

	log.Success(op, "stop signal "+fmt.Sprint(sign))
}
