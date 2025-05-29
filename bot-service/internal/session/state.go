package session

import (
	"github.com/go-telegram/fsm"
)

const (
	StateDefault             fsm.StateID = "default"
	StateAskCreateCollection fsm.StateID = "ask_create_collection"
	StateCreateCollection    fsm.StateID = "create_collection"
	StateAskRenameCollection fsm.StateID = "ask_rename_collection"
	StateRenameCollection    fsm.StateID = "rename_collection"
	StateAskDeleteCollection fsm.StateID = "ask_delete_collection"
	StateDeleteCollection    fsm.StateID = "delete_collection"
)
