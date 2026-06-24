package table

import "github.com/Rafael24595/go-reacterm-core/engine/app/store"

const (
	KeyState store.Key[State] = "table_state"
	KeySync  store.Key[Sync]  = "table_sync"
)

type Sync struct {
	Row *uint16
	Col *uint16
}

type State struct {
	Row uint16
	Col uint16
}
