package service

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/notes"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"

	"github.com/autumnterror/breezynotes/internal/blocknote/config"
)

type BN struct {
	tx  TxRunner
	nts notes.Repo
	tgs tags.Repo
	blk blocks.Repo
	cfg *config.Config
}

func NewNoteService(
	cfg *config.Config,
	tx TxRunner,
	nts notes.Repo,
	blk blocks.Repo,
	tgs tags.Repo,
) *BN {
	return &BN{
		tx:  tx,
		nts: nts,
		cfg: cfg,
		blk: blk,
		tgs: tgs,
	}
}

func (s *BN) Healthz(ctx context.Context) error {
	return s.tx.Healthz(ctx)
}
