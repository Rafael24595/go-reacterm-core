package table

import "github.com/Rafael24595/go-reacterm-core/engine/app/store"

const KeyState store.Key[State] = "table_state"

type State struct {
	Row uint16
	Col uint16
}
