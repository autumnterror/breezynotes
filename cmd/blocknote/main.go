package main

import (
	"fmt"
	"github.com/autumnterror/breezynotes/internal/blocknote/api"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo/mongotx"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/notes"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"
	"github.com/autumnterror/breezynotes/internal/blocknote/service"
	"github.com/autumnterror/breezynotes/pkg/pkgs"
	"github.com/autumnterror/breezynotes/pkg/pkgs/default/textblock"
	"github.com/autumnterror/utils_go/pkg/log"

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
	b := blocks.NewApi(m.Blocks())
	t := tags.NewApi(m.Tags())
	n := notes.NewApi(m.Notes(), m.Trash(), b)
	g := api.New(cfg, service.NewNoteService(cfg, mongotx.NewTxRunner(m.C), n, b, t))
	go g.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	g.Stop()

	log.Success(op, "stop signal "+fmt.Sprint(sign))
}
