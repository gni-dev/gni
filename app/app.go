package app

import (
	"gni.dev/gni"
	"gni.dev/gni/internal/backend"
)

var (
	running bool
)

func Exec(w gni.Widget) {
	if running {
		panic("App already running")
	}
	running = true

	backend.Exec(w)
}
