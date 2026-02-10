package domain2

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
	NoteTagsColl = "note-tags"

	ReaderRole = "reader"
	EditorRole = "editor"
)
