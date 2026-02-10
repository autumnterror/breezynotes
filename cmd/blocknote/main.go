package main

import (
	"fmt"
	"github.com/autumnterror/breezynotes/internal/blocknote/api"
	"github.com/autumnterror/breezynotes/internal/blocknote/config"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo"
	"github.com/autumnterror/breezynotes/internal/blocknote/infra/mongo/mongotx"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/codeblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/fileblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/headerblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/imgblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/linkblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/listblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/quoteblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block/default/textblock"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/notes"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"
	"github.com/autumnterror/breezynotes/internal/blocknote/service"
	"github.com/autumnterror/utils_go/pkg/log"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "cmd.blocknote"

	//------------REG-----------
	block.RegisterBlock("text", &textblock.Driver{})
	block.RegisterBlock("code", &codeblock.Driver{})
	block.RegisterBlock("file", &fileblock.Driver{})
	block.RegisterBlock("header", &headerblock.Driver{})
	block.RegisterBlock("img", &imgblock.Driver{})
	block.RegisterBlock("link", &linkblock.Driver{})
	block.RegisterBlock("list", &listblock.Driver{})
	block.RegisterBlock("quote", &quoteblock.Driver{})
	//------------REG-----------
	log.Green("Types was registered: ", block.GetRegisteredTypes())

	cfg := config.MustSetup()

	m := mongo.MustConnect(cfg)
	b := blocks.NewApi(m.Blocks())
	t := tags.NewApi(m.Tags())
	n := notes.NewApi(m.Notes(), m.Trash(), m.NoteTags(), t, b)
	g := api.New(cfg, service.NewNoteService(cfg, mongotx.NewTxRunner(m.C), n, b, t))
	go g.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	g.Stop()

	log.Success(op, "stop signal "+fmt.Sprint(sign))
}
