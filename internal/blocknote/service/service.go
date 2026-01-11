package service

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/blocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/notes"
	"github.com/autumnterror/breezynotes/internal/blocknote/repository/tags"

	"github.com/autumnterror/breezynotes/internal/blocknote/config"
)

type TagsService struct {
	tx  TxRunner
	tgs tags.Repo
	cfg *config.Config
}

func NewTagsService(
	tx TxRunner,
	tgs tags.Repo,
	cfg *config.Config,
) *TagsService {
	return &TagsService{
		tx:  tx,
		tgs: tgs,
		cfg: cfg,
	}
}

type BlocksService struct {
	tx  TxRunner
	bls blocks.Repo
	cfg *config.Config
}

func NewBlocksService(
	tx TxRunner,
	bls blocks.Repo,
	cfg *config.Config,
) *BlocksService {
	return &BlocksService{
		tx:  tx,
		bls: bls,
		cfg: cfg,
	}
}

type NotesService struct {
	tx  TxRunner
	nts notes.Repo
	tgs *TagsService
	blk *BlocksService
	cfg *config.Config
}

func NewNoteService(
	tx TxRunner,
	nts notes.Repo,
	cfg *config.Config,
	blk *BlocksService,
	tgs *TagsService,
) *NotesService {
	return &NotesService{
		tx:  tx,
		nts: nts,
		cfg: cfg,
		blk: blk,
		tgs: tgs,
	}
}

func (s *NotesService) Healthz(ctx context.Context) error {
	return s.tx.Healthz(ctx)
}
