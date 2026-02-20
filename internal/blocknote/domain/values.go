package domain

import (
	"time"
)

const (
	WaitTime     = 3 * time.Second
	Db           = "blocknotedb"
	TagColl      = "tags"
	NoteColl     = "notes"
	BlockColl    = "blocks"
	TrashColl    = "trash"
	NoteTagsColl = "notetags"

	ReaderRole = "reader"
	EditorRole = "editor"
)
