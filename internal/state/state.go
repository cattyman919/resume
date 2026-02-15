package state

import "sync"

type AppData struct {
}

type State struct {
	mu   sync.RWMutex
	data AppData
}

func NewState() *State {
	return &State{}
}
