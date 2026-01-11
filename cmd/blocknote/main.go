package main

import (
	"fmt"
	"github.com/autumnterror/breezynotes/internal/blocknote/api"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/notes"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/pkgs"
	"github.com/autumnterror/breezynotes/pkg/pkgs/default/textblock"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.blocknote"

	//------------REG-----------
	pkgs.RegisterBlock("text", &textblock.Driver{})
	//------------REG-----------

	cfg := config.MustSetup()

	m := mongo.MustConnect(cfg)
	b := blocks.NewApi(m)
	n := notes.NewApi(m, b)
	g := api.New(cfg, tags.NewApi(m), n, b)
	go g.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	g.Stop()

	log.Success(op, "stop signal "+fmt.Sprint(sign))
}
